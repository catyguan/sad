package httpserver

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"crypto/md5"
	"encoding/json"
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
