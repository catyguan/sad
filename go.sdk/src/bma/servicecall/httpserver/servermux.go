package httpserver

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type pollAnswer struct {
	done   bool
	peer   *HttpServicePeer
	answer *sccore.Answer
	err    error
	timer  *time.Timer
}

type ServiceCallMux struct {
	dispather     ServiceDispatch
	lock          sync.RWMutex
	services      map[string]sccore.ServiceObject
	methods       map[string]map[string]sccore.ServiceMethod
	trans         map[string]*HttpServicePeer
	polls         map[string]*pollAnswer
	seed          int64
	seq           uint32
	clientFactory sccore.ClientFactory
}

func NewServiceCallMux(fac sccore.ClientFactory) *ServiceCallMux {
	o := new(ServiceCallMux)
	o.seed = time.Now().UnixNano()
	o.seq = 0
	o.clientFactory = fac
	return o
}

func (this *ServiceCallMux) SetDispatcher(dis ServiceDispatch) {
	this.dispather = dis
}

func (this *ServiceCallMux) SetServiceObject(name string, so sccore.ServiceObject) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.services == nil {
		this.services = make(map[string]sccore.ServiceObject)
	}
	this.services[name] = so
}

func (this *ServiceCallMux) SetServiceMethod(service string, method string, sm sccore.ServiceMethod) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.methods == nil {
		this.methods = make(map[string]map[string]sccore.ServiceMethod)
	}
	s, ok := this.methods[service]
	if !ok {
		s = make(map[string]sccore.ServiceMethod)
		this.methods[service] = s
	}
	s[method] = sm
}

func (this *ServiceCallMux) Find(s, m string) (sccore.ServiceMethod, error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if ms, ok := this.methods[s]; ok {
		r := ms[m]
		if r != nil {
			return r, nil
		}
	}
	if ss, ok := this.services[s]; ok {
		return ss.GetMethod(m), nil
	}
	return nil, nil
}

func (this *ServiceCallMux) createSeq() string {
	seq := atomic.AddUint32(&this.seq, 1)
	s := fmt.Sprintf("%d_%d", this.seed, seq)
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *ServiceCallMux) getTrans(tid string) *HttpServicePeer {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.trans == nil {
		return nil
	}
	return this.trans[tid]
}

func (this *ServiceCallMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err0 := r.ParseForm()
	if err0 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	dis := this.dispather
	if dis == nil {
		dis = DefaultServiceDispatch
	}
	uri := r.RequestURI
	s, m, err1 := dis(r)
	sccore.DoLog("%s -> %s:%s", uri, s, m)
	if err1 != nil {
		http.Error(w, fmt.Sprintf("dispatch service fail - %s", err1), http.StatusInternalServerError)
		return
	}
	servm, err2 := this.Find(s, m)
	if err2 != nil {
		http.Error(w, fmt.Sprintf("find service method fail - %s", err2), http.StatusInternalServerError)
		return
	}
	sccore.DoLog("%s:%s -> %v", s, m, servm)
	if servm == nil {
		http.Error(w, fmt.Sprintf("service(%s:%s) not found", s, m), http.StatusNotFound)
		return
	}

	qv := r.PostFormValue("q")
	qm := make(map[string]interface{})
	qd := json.NewDecoder(strings.NewReader(qv))
	qd.UseNumber()
	err3 := qd.Decode(&qm)
	if err3 != nil {
		http.Error(w, fmt.Sprintf("decode request fail - %s", err3), http.StatusInternalServerError)
		return
	}

	cv := r.PostFormValue("c")
	cm := make(map[string]interface{})
	cd := json.NewDecoder(strings.NewReader(cv))
	cd.UseNumber()
	err4 := cd.Decode(&cm)
	if err4 != nil {
		http.Error(w, fmt.Sprintf("decode context fail - %s", err4), http.StatusInternalServerError)
		return
	}

	sccore.DoLog("Q : %v", qv)
	sccore.DoLog("C : %v", cv)
	sccore.DoLog("QM : %v", qm)
	sccore.DoLog("CM : %v", cm)

	req := sccore.CreateRequest(qm)
	ctx := sccore.CreateContext(cm)

	aid := ctx.GetString(constv.KEY_ASYNC_ID)
	if aid != "" {
		var pa *pollAnswer
		this.lock.RLock()
		if this.polls != nil {
			if pa2, ok := this.polls[aid]; ok {
				if pa2.done {
					pa = pa2
				}
			}
		}
		this.lock.RUnlock()

		var peer *HttpServicePeer
		var aa *sccore.Answer
		var aerr error
		if pa != nil {
			this.lock.Lock()
			delete(this.polls, aid)
			this.lock.Unlock()
			sccore.DoLog("'%s' poll success", aid)
			pa.timer.Stop()
			aa = pa.answer
			aerr = pa.err
			peer = pa.peer
		} else {
			sccore.DoLog("'%s' polling", aid)
			aa = sccore.NewAnswer()
			aa.SetStatus(constv.STATUS_ASYNC)
			aa.SureResult().Put(constv.KEY_ASYNC_ID, aid)
		}
		doAnswer(peer, w, aa, aerr)
		return
	}

	end := make(chan bool)
	transId := ctx.GetString(constv.KEY_TRANSACTION_ID)
	if transId != "" {
		peer := this.getTrans(transId)
		peer.Post(end, w, req, ctx)
	} else {
		peer := new(HttpServicePeer)
		peer.mux = this
		peer.w = w
		peer.end = end
		go func() {
			err5 := servm(peer, req, ctx)
			if err5 != nil {
				http.Error(w, fmt.Sprintf("service fail - %s", err5), http.StatusInternalServerError)
				return
			}
		}()
	}
	<-end
}

func (this *ServiceCallMux) DoCallback(peer *HttpServicePeer, req *sccore.Request, ctx *sccore.Context) (*sccore.Answer, error) {
	if this.clientFactory == nil {
		return nil, fmt.Errorf("clientFactory is nil")
	}
	cl := this.clientFactory()
	defer cl.Close()
	answer, err := cl.Invoke(peer.callback, req, ctx)
	return answer, err
}

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

func (this *HttpServicePeer) BeginTransaction() (string, error) {
	if this.transId != "" {
		return "", fmt.Errorf("already begin transaction")
	}
	this.ch = make(chan *reqInfo, 1)

	mux := this.mux
	this.transId = mux.createSeq()

	mux.lock.Lock()
	defer mux.lock.Unlock()
	if mux.trans == nil {
		mux.trans = make(map[string]*HttpServicePeer)
	}
	mux.trans[this.transId] = this
	return this.transId, nil
}

func (this *HttpServicePeer) EndTransaction() {
	if this.transId == "" {
		return
	}
	mux := this.mux
	if mux.trans == nil {
		return
	}
	mux.lock.Lock()
	defer mux.lock.Unlock()
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
		if a == nil {
			a = sccore.NewAnswer()
		}
		a.SetStatus(500)
		a.SetMessage(err.Error())
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
		mux.lock.RLock()
		pa := mux.polls[this.asyncId]
		mux.lock.RUnlock()
		if pa != nil {
			pa.answer = a
			pa.err = err
			pa.done = true
			sccore.DoLog("async answer '%s'", this.asyncId)
		} else {
			sccore.DoLog("miss async '%s'", this.asyncId)
		}
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
		an, err2 := this.mux.DoCallback(this, req, ctx)
		sccore.DoLog("callback answer -> %v, %v", err2, an)
		return err
	case 1:
		return fmt.Errorf("HttpServicePeer already answer")
	default:
		if this.w == nil {
			return fmt.Errorf("HttpServicePeer break")
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
		aid := mux.createSeq()
		pa := new(pollAnswer)
		pa.done = false
		pa.peer = this
		pa.answer = nil
		pa.err = nil
		pa.timer = time.AfterFunc(timeout, func() {
			mux.lock.Lock()
			defer mux.lock.Unlock()
			if mux.polls != nil {
				if _, ok := mux.polls[aid]; ok {
					delete(mux.polls, aid)
					sccore.DoLog("async poll(%s) wait timeout", aid)
				}
			}
		})
		mux.lock.Lock()
		if mux.polls == nil {
			mux.polls = make(map[string]*pollAnswer)
		}
		mux.polls[aid] = pa
		mux.lock.Unlock()

		a := sccore.NewAnswer()
		a.SetStatus(constv.STATUS_ASYNC)
		if result == nil {
			result = sccore.NewValueMap(nil)
		}
		result.Put(constv.KEY_ASYNC_ID, aid)
		a.SetResult(result)
		this.WriteAnswer(a, nil)
		this.mode = 2
		this.asyncId = aid
		return nil
	case "callback":
		addrm := ctx.GetMap(constv.KEY_CALLBACK)
		if addrm == nil {
			err0 := fmt.Errorf("HttpServicePeer Async callback miss address")
			this.WriteAnswer(nil, err0)
			return err0
		}
		this.callback = sccore.CreateAddressFromValue(addrm)

		a := sccore.NewAnswer()
		a.SetStatus(constv.STATUS_ASYNC)
		a.SetResult(result)
		this.WriteAnswer(a, nil)
		this.mode = 3
		return nil
	default:
		err := fmt.Errorf("HttpServicePeer not support AsyncMode(%s)", async)
		this.WriteAnswer(nil, err)
		return err
	}
}
