package httpserver

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ServiceCallMux struct {
	dispather ServiceDispatch
	serv      sccore.BaseServiceServ
	sccore.ServiceMux
	trans map[string]*HttpServicePeer
}

func NewServiceCallMux(fac sccore.ClientFactory) *ServiceCallMux {
	o := new(ServiceCallMux)
	o.serv.InitBaseServiceServ(fac)
	return o
}

func (this *ServiceCallMux) getTrans(tid string) *HttpServicePeer {
	this.serv.Lock.RLock()
	defer this.serv.Lock.RUnlock()
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
	if qv != "" {
		qd := json.NewDecoder(strings.NewReader(qv))
		qd.UseNumber()
		err3 := qd.Decode(&qm)
		if err3 != nil {
			http.Error(w, fmt.Sprintf("decode request fail - %s", err3), http.StatusInternalServerError)
			return
		}
	}

	cv := r.PostFormValue("c")
	cm := make(map[string]interface{})
	if cv != "" {
		cd := json.NewDecoder(strings.NewReader(cv))
		cd.UseNumber()
		err4 := cd.Decode(&cm)
		if err4 != nil {
			http.Error(w, fmt.Sprintf("decode context fail - %s", err4), http.StatusInternalServerError)
			return
		}
	}

	sccore.DoLog("Q : %v", qv)
	sccore.DoLog("C : %v", cv)
	sccore.DoLog("QM : %v", qm)
	sccore.DoLog("CM : %v", cm)

	req := sccore.CreateRequest(qm)
	ctx := sccore.CreateContext(cm)

	aid := ctx.GetString(constv.KEY_ASYNC_ID)
	if aid != "" {
		pa := this.serv.PollAsync(aid)
		var peer *HttpServicePeer
		var aa *sccore.Answer
		var aerr error
		if pa != nil {
			aa = pa.Answer
			aerr = pa.Err
			peer = pa.Peer.(*HttpServicePeer)
		} else {
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
