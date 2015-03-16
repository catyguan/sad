package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type PollAnswer struct {
	Done   bool
	Peer   ServicePeer
	Answer *Answer
	Err    error
	Timer  *time.Timer
}

type BaseServiceServ struct {
	ClientFactory ClientFactory
	Lock          sync.RWMutex
	Polls         map[string]*PollAnswer
	Seed          int64
	Seq           uint32
}

func (this *BaseServiceServ) InitBaseServiceServ(fac ClientFactory) {
	this.Seed = time.Now().UnixNano()
	this.Seq = 0
	this.ClientFactory = fac
}

func (this *BaseServiceServ) CreateSeq() string {
	seq := atomic.AddUint32(&this.Seq, 1)
	s := fmt.Sprintf("%d_%d", this.Seed, seq)
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *BaseServiceServ) CreatePollAnswer(du time.Duration, peer ServicePeer) string {
	aid := this.CreateSeq()
	pa := new(PollAnswer)
	pa.Peer = peer
	pa.Timer = time.AfterFunc(du, func() {
		this.Lock.Lock()
		defer this.Lock.Unlock()
		if this.Polls != nil {
			if _, ok := this.Polls[aid]; ok {
				DoLog("poll async timeout '%s'", aid)
			}
			delete(this.Polls, aid)
		}
	})
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.Polls == nil {
		this.Polls = make(map[string]*PollAnswer)
	}
	this.Polls[aid] = pa
	return aid
}

func (this *BaseServiceServ) SetPollAnswer(aid string, an *Answer, err error) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.Polls != nil {
		if pa, ok := this.Polls[aid]; ok {
			pa.Done = true
			pa.Answer = an
			pa.Err = err
			DoLog("poll async answer '%s'", aid)
			return
		}
	}
	DoLog("poll async miss '%s'", aid)
}

func (this *BaseServiceServ) PollAsync(aid string) *PollAnswer {
	if aid == "" {
		return nil
	}
	var pa *PollAnswer
	this.Lock.RLock()
	if this.Polls != nil {
		if pa2, ok := this.Polls[aid]; ok {
			if pa2.Done {
				pa = pa2
			}
		}
	}
	this.Lock.RUnlock()

	if pa != nil {
		this.Lock.Lock()
		delete(this.Polls, aid)
		this.Lock.Unlock()
		DoLog("'%s' poll success", aid)
		pa.Timer.Stop()
		return pa
	} else {
		DoLog("'%s' polling", aid)
		return nil
	}
}

func (this *BaseServiceServ) DoCallback(addr *Address, req *Request, ctx *Context) (*Answer, error) {
	if this.ClientFactory == nil {
		return nil, fmt.Errorf("clientFactory is nil")
	}
	cl := this.ClientFactory()
	defer cl.Close()
	answer, err := cl.Invoke(addr, req, ctx)
	return answer, err
}
