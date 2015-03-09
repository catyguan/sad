package core

import (
	"bma/servicecall/constv"
	"errors"
	"fmt"
	"time"
)

type Client struct {
	manager   *Manager
	id        uint32
	reqSeq    uint32
	inTrans   bool
	transId   string
	sessionId string
}

func (this *Client) CreateReqId() string {
	this.reqSeq++
	seq := this.reqSeq
	return fmt.Sprintf("%s_%d_%d", this.manager.name, this.id, seq)
}

func (this *Client) GetSessionId() string {
	return this.sessionId
}

func (this *Client) SetSessionId(v string) {
	this.sessionId = v
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
		if this.inTrans && this.transId != "" {
			ctx.Put(constv.KEY_TRANSACTION_ID, this.transId)
		}
		if this.sessionId != "" {
			ctx.Put(constv.KEY_SESSION_ID, this.sessionId)
		}
		a, err := this.doInvoke(addr, req, ctx)
		if err != nil {
			return a, err
		}
		if this.inTrans {
			actx := a.GetContext()
			if actx != nil {
				tid := actx.GetString(constv.KEY_TRANSACTION_ID)
				if tid != "" {
					this.transId = tid
				}
				sid := actx.GetString(constv.KEY_SESSION_ID)
				if sid != "" {
					this.sessionId = sid
				}
			}
		}
		switch a.GetStatus() {
		case 200, 100, 202, 204:
			return a, nil
		case 302:
			rs := a.GetResult()
			if rs == nil {
				return nil, fmt.Errorf("redirect address empty")
			}
			addr = CreateAddressFromValue(rs)
			if errA := addr.Valid(); errA != nil {
				return nil, errA
			}
			DoLog("redirect -> %s", addr.String())
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

func (this *Client) BeginTransaction() bool {
	if !this.inTrans {
		this.inTrans = true
		this.transId = ""
		return true
	}
	return false
}

func (this *Client) EndTransaction() {
	this.inTrans = false
	this.transId = ""
}

func (this *Client) IsTransacion() bool {
	return this.inTrans
}

func (this *Client) Export() map[string]interface{} {
	r := make(map[string]interface{})
	if this.inTrans && this.transId != "" {
		r["TransId"] = this.transId
	}
	if this.sessionId != "" {
		r["SessionId"] = this.sessionId
	}
	return r
}

func (this *Client) Import(data map[string]interface{}) error {
	if data == nil {
		return nil
	}
	if this.inTrans {
		if sv, ok := data["TransId"]; ok {
			if s, ok2 := sv.(string); ok2 {
				this.transId = s
			}
		}
	}
	if sv, ok := data["SessionId"]; ok {
		if s, ok2 := sv.(string); ok2 {
			this.sessionId = s
		}
	}
	return nil
}
