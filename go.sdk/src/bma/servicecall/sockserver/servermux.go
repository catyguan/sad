package sockserver

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"bma/servicecall/sockcore"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type pollAnswer struct {
	done   bool
	peer   *SocketServicePeer
	answer *sccore.Answer
	err    error
	timer  *time.Timer
}

type ServiceCallMux struct {
	lock          sync.RWMutex
	services      map[string]sccore.ServiceObject
	methods       map[string]map[string]sccore.ServiceMethod
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

func (this *ServiceCallMux) Run(conn net.Conn) {
	defer func() {
		recover()
		conn.Close()
	}()
	for {
		err := this.ServeSocket(conn)
		if err != nil {
			return
		}
	}
}

var pingRData = []byte{9, 0, 0, 1, 1, 0, 0, 0, 0}

func (this *ServiceCallMux) nextMessage(conn net.Conn, mr *sockcore.MessageReader, msg *sockcore.Message) error {
	for {
		mt, errR := mr.NextMessage(msg)
		if errR != nil {
			return errR
		}
		switch mt {
		case sockcore.MT_REQUEST:
			return nil
		case sockcore.MT_PING:
			if !msg.BoolFlag {
				_, err := conn.Write(pingRData)
				if err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("not support messge type(%d)", mt)
		}
	}
}

func (this *ServiceCallMux) ServeSocket(conn net.Conn) error {
	mr := sockcore.NewMessageReader(conn)
	var msg sockcore.Message
	errR := this.nextMessage(conn, mr, &msg)
	if errR != nil {
		sccore.DoLog("next Request fail - %s", errR)
		return errR
	}
	s, m := msg.Service, msg.Method
	if s == "" || m == "" {
		err := fmt.Errorf("address(%s:%s) empty", s, m)
		sccore.DoLog("read Request fail - %s", err)
		return err
	}
	mid := msg.Id
	req := msg.Request
	ctx := msg.Context
	if req == nil {
		req = sccore.NewRequest()
	}
	if ctx == nil {
		ctx = sccore.NewContext()
	}

	peer := new(SocketServicePeer)
	peer.mux = this
	peer.conn = conn
	sccore.DoLog("call -> %s:%s", s, m)

	servm, err2 := this.Find(s, m)
	if err2 != nil {
		err3 := fmt.Errorf("find service method fail - %s", err2)
		peer.doAnswer(conn, mid, nil, err3)
		return err3
	}
	sccore.DoLog("%s:%s -> %v", s, m, servm)
	if servm == nil {
		err3 := fmt.Errorf("service(%s:%s) not found", s, m)
		peer.doAnswer(conn, mid, nil, err3)
		return err3
	}

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

		if pa != nil {
			this.lock.Lock()
			delete(this.polls, aid)
			this.lock.Unlock()
			sccore.DoLog("'%s' poll success", aid)
			pa.timer.Stop()
			aa := pa.answer
			aerr := pa.err
			peer2 := pa.peer
			peer2.doAnswer(conn, mid, aa, aerr)
		} else {
			sccore.DoLog("'%s' polling", aid)
			aa := sccore.NewAnswer()
			aa.SetStatus(constv.STATUS_ASYNC)
			aa.SureResult().Put(constv.KEY_ASYNC_ID, aid)
			peer.doAnswer(conn, mid, aa, nil)
		}
		return nil
	}
	peer.messageId = mid

	err5 := servm(peer, req, ctx)
	if err5 != nil {
		sccore.DoLog("service fail - %s", err5)
		peer.doAnswer(conn, mid, nil, err5)
		return err5
	}
	return nil
}

func (this *ServiceCallMux) DoCallback(peer *SocketServicePeer, req *sccore.Request, ctx *sccore.Context) (*sccore.Answer, error) {
	if this.clientFactory == nil {
		return nil, fmt.Errorf("clientFactory is nil")
	}
	cl := this.clientFactory()
	defer cl.Close()
	answer, err := cl.Invoke(peer.callback, req, ctx)
	return answer, err
}
