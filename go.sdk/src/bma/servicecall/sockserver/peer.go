package sockserver

import (
	sccore "bma/servicecall/core"
	"net"
	"time"
)

type SocketServicePeer struct {
	mux       *ServiceCallMux
	conn      net.Conn
	mode      int // 0-normal, 1-writed, 2-poll, 3-callback, 4-push
	asyncId   string
	callback  *sccore.Address
	transId   string
	messageId int32
}

func (this *SocketServicePeer) doAnswer(conn net.Conn, mid int32, a *sccore.Answer, err error) error {
	if err != nil {
		a = sccore.Error2Answer(a, err)
	}
	// m := make(map[string]interface{})
	// sc := a.GetStatus()
	// if sc <= 0 {
	// 	a.SetStatus(200)
	// }
	// if this.transId != "" {
	// 	a.SureContext().Put(constv.KEY_TRANSACTION_ID, this.transId)
	// }
	// m["Status"] = a.GetStatus()
	// msg := a.GetMessage()
	// if msg != "" {
	// 	m["Message"] = msg
	// }
	// rs := a.GetResult()
	// if rs != nil {
	// 	m["Result"] = rs.ToMap()
	// }
	// ctx := a.GetContext()
	// if ctx != nil {
	// 	m["Context"] = ctx.ToMap()
	// }
	// bs, err1 := json.Marshal(m)
	// if err1 != nil {
	// 	return err1
	// }
	// _, err2 := w.Write(bs)
	// return err2
	return nil
}

func (this *SocketServicePeer) BeginTransaction() (string, error) {
	// if this.transId != "" {
	// 	return "", fmt.Errorf("already begin transaction")
	// }
	// this.ch = make(chan *reqInfo, 1)

	// mux := this.mux
	// this.transId = mux.createSeq()

	// mux.lock.Lock()
	// defer mux.lock.Unlock()
	// if mux.trans == nil {
	// 	mux.trans = make(map[string]*SocketServicePeer)
	// }
	// mux.trans[this.transId] = this
	// return this.transId, nil
	return "", nil
}

func (this *SocketServicePeer) EndTransaction() {
	// if this.transId == "" {
	// 	return
	// }
	// mux := this.mux
	// if mux.trans == nil {
	// 	return
	// }
	// mux.lock.Lock()
	// defer mux.lock.Unlock()
	// delete(mux.trans, this.transId)
	// this.transId = ""
	// close(this.ch)
	// this.ch = nil
}

func (this *SocketServicePeer) ReadRequest(waitTime time.Duration) (*sccore.Request, *sccore.Context, error) {
	// if this.ch == nil {
	// 	return nil, nil, errors.New("not begin transaction")
	// }
	// timer := time.NewTimer(waitTime)
	// select {
	// case <-timer.C:
	// 	return nil, nil, fmt.Errorf("timeout")
	// case ri := <-this.ch:
	// 	if ri == nil {
	// 		return nil, nil, fmt.Errorf("closed")
	// 	}
	// 	return ri.req, ri.ctx, nil
	// }
	return nil, nil, nil
}

func (this *SocketServicePeer) WriteAnswer(a *sccore.Answer, err error) error {
	// switch this.mode {
	// case 2:
	// 	if this.asyncId == "" {
	// 		return fmt.Errorf("poll mode, asyncId empty")
	// 	}
	// 	mux := this.mux
	// 	mux.lock.RLock()
	// 	pa := mux.polls[this.asyncId]
	// 	mux.lock.RUnlock()
	// 	if pa != nil {
	// 		pa.answer = a
	// 		pa.err = err
	// 		pa.done = true
	// 		sccore.DoLog("async answer '%s'", this.asyncId)
	// 	} else {
	// 		sccore.DoLog("miss async '%s'", this.asyncId)
	// 	}
	// 	this.asyncId = ""
	// 	return nil
	// case 3:
	// 	var req *sccore.Request
	// 	if err != nil {
	// 		am := make(map[string]interface{})
	// 		am["Status"] = constv.STATUS_ERROR
	// 		am["Message"] = err.Error()
	// 		req = sccore.CreateRequest(am)
	// 	} else {
	// 		am := a.ToMap()
	// 		req = sccore.CreateRequest(am)
	// 	}
	// 	ctx := sccore.NewContext()
	// 	sccore.DoLog("callback invoke -> %v, %v", err, a)
	// 	an, err2 := this.mux.DoCallback(this, req, ctx)
	// 	sccore.DoLog("callback answer -> %v, %v", err2, an)
	// 	return err
	// case 1:
	// 	return fmt.Errorf("SocketServicePeer already answer")
	// default:
	// 	if this.w == nil {
	// 		return fmt.Errorf("SocketServicePeer break")
	// 	}
	// 	sccore.DoLog("writeAnswer -> %v, %v", err, a)
	// 	err2 := doAnswer(this, this.w, a, err)
	// 	this.mode = 1
	// 	close(this.end)
	// 	return err2
	// }
	return nil
}

func (this *SocketServicePeer) SendAsync(ctx *sccore.Context, result *sccore.ValueMap, timeout time.Duration) error {
	// async := ctx.GetString(constv.KEY_ASYNC_MODE)
	// switch async {
	// case "", "poll":
	// 	mux := this.mux
	// 	aid := mux.createSeq()
	// 	pa := new(pollAnswer)
	// 	pa.done = false
	// 	pa.peer = this
	// 	pa.answer = nil
	// 	pa.err = nil
	// 	pa.timer = time.AfterFunc(timeout, func() {
	// 		mux.lock.Lock()
	// 		defer mux.lock.Unlock()
	// 		if mux.polls != nil {
	// 			if _, ok := mux.polls[aid]; ok {
	// 				delete(mux.polls, aid)
	// 				sccore.DoLog("async poll(%s) wait timeout", aid)
	// 			}
	// 		}
	// 	})
	// 	mux.lock.Lock()
	// 	if mux.polls == nil {
	// 		mux.polls = make(map[string]*pollAnswer)
	// 	}
	// 	mux.polls[aid] = pa
	// 	mux.lock.Unlock()

	// 	a := sccore.NewAnswer()
	// 	a.SetStatus(constv.STATUS_ASYNC)
	// 	if result == nil {
	// 		result = sccore.NewValueMap(nil)
	// 	}
	// 	result.Put(constv.KEY_ASYNC_ID, aid)
	// 	a.SetResult(result)
	// 	this.WriteAnswer(a, nil)
	// 	this.mode = 2
	// 	this.asyncId = aid
	// 	return nil
	// case "callback":
	// 	addrm := ctx.GetMap(constv.KEY_CALLBACK)
	// 	if addrm == nil {
	// 		err0 := fmt.Errorf("SocketServicePeer Async callback miss address")
	// 		this.WriteAnswer(nil, err0)
	// 		return err0
	// 	}
	// 	this.callback = sccore.CreateAddressFromValue(addrm)

	// 	a := sccore.NewAnswer()
	// 	a.SetStatus(constv.STATUS_ASYNC)
	// 	a.SetResult(result)
	// 	this.WriteAnswer(a, nil)
	// 	this.mode = 3
	// 	return nil
	// default:
	// 	err := fmt.Errorf("SocketServicePeer not support AsyncMode(%s)", async)
	// 	this.WriteAnswer(nil, err)
	// 	return err
	// }
	return nil
}
