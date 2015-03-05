package httpserver

import (
	sccore "bma/servicecall/core"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ServiceCallMux struct {
	dispather ServiceDispatch
	lock      sync.RWMutex
	services  map[string]ServiceObject
	methods   map[string]map[string]sccore.ServiceMethod
}

func NewServiceCallMux() *ServiceCallMux {
	o := new(ServiceCallMux)
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
	s, m, err1 := dis(r)
	if err1 != nil {
		http.Error(w, fmt.Sprintf("dispatch service fail - %s", err1), http.StatusInternalServerError)
		return
	}
	servm, err2 := this.Find(s, m)
	if err2 != nil {
		http.Error(w, fmt.Sprintf("find service method fail - %s", err2), http.StatusInternalServerError)
		return
	}
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

	// fmt.Println("Q : ", qv)
	// fmt.Println("C : ", cv)
	// fmt.Println("QM : ", qm)
	// fmt.Println("CM : ", cm)

	req := sccore.CreateRequest(qm)
	ctx := sccore.CreateContext(cm)

	peer := new(HttpServicePeer)
	peer.w = w
	err5 := servm(peer, req, ctx)

	if err5 != nil {
		http.Error(w, fmt.Sprintf("service fail - %s", err5), http.StatusInternalServerError)
		return
	}
}

type HttpServicePeer struct {
	w      http.ResponseWriter
	writed bool
}

func (this *HttpServicePeer) BeginTransaction() {

}

func (this *HttpServicePeer) EndTransaction() {

}

func (this *HttpServicePeer) ReadRequest(waitTime time.Duration) (*sccore.Request, *sccore.Context, error) {
	return nil, nil, nil
}

func (this *HttpServicePeer) WriteAnswer(a *sccore.Answer, err error) error {
	if this.writed {
		return fmt.Errorf("HttpServicePeer already answer")
	}
	if this.w == nil {
		return fmt.Errorf("HttpServicePeer break")
	}
	if err != nil {
		if a == nil {
			a = sccore.NewAnswer()
		}
		a.SetStatus(500)
		a.SetMessage(err.Error())
	}
	m := make(map[string]interface{})
	m["Status"] = a.GetStatus()
	msg := a.GetMessage()
	if msg != "" {
		m["Message"] = msg
	}
	rs := a.GetResult()
	if rs != nil {
		rsm := rs.ToMap()
		m["Result"] = rsm
	}
	bs, err1 := json.Marshal(m)
	if err1 != nil {
		return err1
	}
	this.w.Write(bs)
	this.writed = true
	return nil
}
