package core

import (
	"bma/servicecall/constv"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type Client struct {
	id      uint32
	reqSeq  uint32
	manager *Manager
}

func (this *Client) CreateReqId() string {
	id := atomic.AddUint32(&this.reqSeq, 1)
	return fmt.Sprintf("%s_%d_%d", this.manager.name, this.id, id)
}

func (this *Client) doInvoke(addr *Address, req *Request, ctx *Context) (*Answer, error) {
	dr, err := this.manager.GetDriver(addr)
	if err != nil {
		return nil, err
	}
	return dr.Invoke(this, addr, req, ctx)
}

func (this *Client) Invoke(addr *Address, req *Request, ctx *Context) (*Answer, error) {
	if ctx == nil {
		ctx = NewContext()
	}
	if !ctx.Has(constv.KEY_DEADLINE) {
		to := ctx.GetInt(constv.KEY_TIMEOUT)
		if to <= 0 {
			to = 30
		} else {
			ctx.Remove(constv.KEY_TIMEOUT)
		}
		dl := time.Now().Add(time.Second * time.Duration(to)).Unix()
		ctx.Put(constv.KEY_DEADLINE, dl)
	}
	if !ctx.Has(constv.KEY_REQ_ID) {
		ctx.Put(constv.KEY_REQ_ID, this.CreateReqId())
	}
	for {
		a, err := this.doInvoke(addr, req, ctx)
		if err != nil {
			return a, err
		}
		switch a.GetStatus() {
		case 200, 100, 202, 204:
			return a, nil
		case 302:
			return a, nil
		default:
			msg := a.message
			if msg == "" {
				msg = "unknow error"
			}
			return a, errors.New(msg)
		}
		// return a, nil
	}
}

func (this *Client) Close() {

}
