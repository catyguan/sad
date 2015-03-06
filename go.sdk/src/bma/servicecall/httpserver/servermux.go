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

type ServiceCallMux struct {
	dispather ServiceDispatch
	lock      sync.RWMutex
	services  map[string]ServiceObject
	methods   map[string]map[string]sccore.ServiceMethod
	trans     map[string]*HttpServicePeer
	transSeed int64
	transSeq  uint32
}

func NewServiceCallMux() *ServiceCallMux {
	o := new(ServiceCallMux)
	o.transSeed = time.Now().UnixNano()
	o.transSeq = 0
	return o
}

func (this *ServiceCallMux) SetDispatcher(dis ServiceDispatch) {
	this.dispather = dis
}

func (this *ServiceCallMux) SetServiceObject(name string, so ServiceObject) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.services == nil {
		this.services = make(map[string]ServiceObject)
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

type reqInfo struct {
	req *sccore.Request
	ctx *sccore.Context
}

type HttpServicePeer struct {
	mux     *ServiceCallMux
	w       http.ResponseWriter
	writed  bool
	transId string
	ch      chan *reqInfo
	end     chan bool
}

func (this *HttpServicePeer) BeginTransaction() (string, error) {
	if this.transId != "" {
		return "", fmt.Errorf("already begin transaction")
	}
	this.ch = make(chan *reqInfo, 1)

	mux := this.mux
	seq := atomic.AddUint32(&mux.transSeq, 1)
	s := fmt.Sprintf("%d%d", mux.transSeed, seq)
	h := md5.New()
	io.WriteString(h, s)
	this.transId = fmt.Sprintf("%x", h.Sum(nil))

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

func (this *HttpServicePeer) WriteAnswer(a *sccore.Answer, err error) error {
	if this.writed {
		return fmt.Errorf("HttpServicePeer already answer")
	}
	if this.w == nil {
		return fmt.Errorf("HttpServicePeer break")
	}
	sccore.DoLog("writeAnswer -> %v, %v", err, a)
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
	if this.transId != "" {
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
	this.w.Write(bs)
	this.writed = true
	close(this.end)
	return nil
}

func (this *HttpServicePeer) Post(end chan bool, w http.ResponseWriter, req *sccore.Request, ctx *sccore.Context) {
	if this.ch == nil {
		sccore.DoLog("post fail, chan nil")
		return
	}
	this.writed = false
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
