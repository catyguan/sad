package httpserver

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type reqInfo struct {
	req *sccore.Request
	ctx *sccore.Context
}

type HttpServicePeer struct {
	mux      *ServiceCallMux
	w        http.ResponseWriter
	mode     int // 0-normal, 1-writed, 2-poll, 3-callback
	asyncId  string
	callback *sccore.Address
	transId  string
	ch       chan *reqInfo
	end      chan bool
}

func (this *HttpServicePeer) GetDriverType() string {
	return "http"
}

func (this *HttpServicePeer) BeginTransaction() error {
	if this.transId != "" {
		return fmt.Errorf("already begin transaction")
	}
	this.ch = make(chan *reqInfo, 1)

	mux := this.mux
	this.transId = mux.serv.CreateSeq()

	mux.serv.Lock.Lock()
	defer mux.serv.Lock.Unlock()
	if mux.trans == nil {
		mux.trans = make(map[string]*HttpServicePeer)
	}
	mux.trans[this.transId] = this
	return nil
}

func (this *HttpServicePeer) endTransaction() {
	if this.transId == "" {
		return
	}
	mux := this.mux
	if mux.trans == nil {
		return
	}
	mux.serv.Lock.Lock()
	defer mux.serv.Lock.Unlock()
	delete(mux.trans, this.transId)
	this.transId = ""
	close(this.ch)
	this.ch = nil
}

func (this *HttpServicePeer) ReadRequest(waitTime time.Duration) (*sccore.Request, *sccore.Context, error) {
	if this.ch == nil {
		return nil, nil, errors.New("not begin transaction")
	}
	timer := time.NewTimer(waitTime)
	select {
	case <-timer.C:
		return nil, nil, fmt.Errorf("timeout")
	case ri := <-this.ch:
		if ri == nil {
			return nil, nil, fmt.Errorf("closed")
		}
		return ri.req, ri.ctx, nil
	}
}

func doAnswer(this *HttpServicePeer, w http.ResponseWriter, a *sccore.Answer, err error) error {
	if err != nil {
		a = sccore.Error2Answer(a, err)
	}
	m := make(map[string]interface{})
	sc := a.GetStatus()
	if sc <= 0 {
		a.SetStatus(200)
	}
	if this != nil && this.transId != "" {
		a.SureContext().Put(constv.KEY_TRANSACTION_ID, this.transId)
	}
	m["Status"] = a.GetStatus()
	msg := a.GetMessage()
	if msg != "" {
		m["Message"] = msg
	}
	rs := a.GetResult()
	if rs != nil {
		m["Result"] = rs.ToMap()
	}
	ctx := a.GetContext()
	if ctx != nil {
		m["Context"] = ctx.ToMap()
	}
	bs, err1 := json.Marshal(m)
	if err1 != nil {
		return err1
	}
	_, err2 := w.Write(bs)
	return err2
}

func (this *HttpServicePeer) WriteAnswer(a *sccore.Answer, err error) error {
	switch this.mode {
	case 2:
		if this.asyncId == "" {
			return fmt.Errorf("poll mode, asyncId empty")
		}
		mux := this.mux
		mux.serv.SetPollAnswer(this.asyncId, a, err)
		this.asyncId = ""
		return nil
	case 3:
		var req *sccore.Request
		if err != nil {
			am := make(map[string]interface{})
			am["Status"] = constv.STATUS_ERROR
			am["Message"] = err.Error()
			req = sccore.CreateRequest(am)
		} else {
			am := a.ToMap()
			req = sccore.CreateRequest(am)
		}
		ctx := sccore.NewContext()
		sccore.DoLog("callback invoke -> %v, %v", err, a)
		an, err2 := this.mux.serv.DoCallback(this.callback, req, ctx)
		sccore.DoLog("callback answer -> %v, %v", err2, an)
		return err
	case 1:
		return fmt.Errorf("HttpServicePeer already answer")
	default:
		if this.w == nil {
			return fmt.Errorf("HttpServicePeer break")
		}
		if this.transId != "" && a.GetStatus() != 100 {
			this.endTransaction()
		}
		sccore.DoLog("writeAnswer -> %v, %v", err, a)
		err2 := doAnswer(this, this.w, a, err)
		this.mode = 1
		close(this.end)
		return err2
	}
}

func (this *HttpServicePeer) Post(end chan bool, w http.ResponseWriter, req *sccore.Request, ctx *sccore.Context) {
	if this.ch == nil {
		sccore.DoLog("post fail, chan nil")
		return
	}
	this.mode = 0
	this.end = end
	this.w = w
	defer func() {
		recover()
	}()
	ri := new(reqInfo)
	ri.req = req
	ri.ctx = ctx
	this.ch <- ri
}

func (this *HttpServicePeer) SendAsync(ctx *sccore.Context, result *sccore.ValueMap, timeout time.Duration) error {
	async := ctx.GetString(constv.KEY_ASYNC_MODE)
	switch async {
	case "", "poll":
		mux := this.mux
		aid := mux.serv.CreatePollAnswer(timeout, this)

		this.mode = 2
		this.asyncId = aid

		a := sccore.NewAnswer()
		a.SetStatus(constv.STATUS_ASYNC)
		if result == nil {
			result = sccore.NewValueMap(nil)
		}
		result.Put(constv.KEY_ASYNC_ID, aid)
		a.SetResult(result)
		return doAnswer(this, this.w, a, nil)
	case "callback":
		addrm := ctx.GetMap(constv.KEY_CALLBACK)
		if addrm == nil {
			err0 := fmt.Errorf("HttpServicePeer Async callback miss address")
			this.WriteAnswer(nil, err0)
			return err0
		}
		this.callback = sccore.CreateAddressFromValue(addrm)

		this.mode = 3
		a := sccore.NewAnswer()
		a.SetStatus(constv.STATUS_ASYNC)
		a.SetResult(result)
		return doAnswer(this, this.w, a, nil)
	default:
		err := fmt.Errorf("HttpServicePeer not support AsyncMode(%s)", async)
		this.WriteAnswer(nil, err)
		return err
	}
}
