package core

import (
	"bma/servicecall/constv"
	"errors"
	"fmt"
	"time"
)

type Client struct {
	manager *Manager
	id      uint32
	reqSeq  uint32
	conns   map[string]ServiceConn
	// status
	sessionId string
	props     map[string]interface{}
}

func newClient(m *Manager, id uint32) *Client {
	o := new(Client)
	o.manager = m
	o.id = id
	o.conns = make(map[string]ServiceConn)
	return o
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

func (this *Client) getConn(addr *Address) (ServiceConn, error) {
	api := addr.GetAPI()
	if conn, ok := this.conns[api]; ok {
		return conn, nil
	}
	typ := addr.GetType()
	conn, err := this.manager.createConn(typ, api)
	if err != nil {
		return nil, err
	}
	this.conns[api] = conn
	return conn, nil
}

func (this *Client) closeConn(addr *Address) {
	api := addr.GetAPI()
	conn := this.conns[api]
	delete(this.conns, api)
	if conn != nil {
		conn.Close()
	}
}

func (this *Client) doInvoke(addr *Address, req *Request, ctx *Context) (*Answer, error) {
	conn, err := this.getConn(addr)
	if err != nil {
		return nil, err
	}
	a, err2 := conn.Invoke(this, addr, req, ctx)
	if err2 != nil {
		this.closeConn(addr)
		return nil, err2
	}
	actx := a.GetContext()
	if actx != nil {
		sid := actx.GetString(constv.KEY_SESSION_ID)
		if sid != "" {
			this.sessionId = sid
		}
	}
	st := a.GetStatus()
	switch st {
	case 100, 200, 202, 204, 302:
	default:
		this.closeConn(addr)
	}
	return a, nil
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
		if this.sessionId == "" {
			ctx.Remove(constv.KEY_SESSION_ID)
		} else {
			ctx.Put(constv.KEY_SESSION_ID, this.sessionId)
		}
		a, err := this.doInvoke(addr, req, ctx)
		if err != nil {
			return a, err
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
	for _, v := range this.props {
		if c, ok := v.(Closable); ok {
			c.Close()
		}
	}
	for k, conn := range this.conns {
		delete(this.conns, k)
		conn.End()
	}
}

func copyv(v interface{}) interface{} {
	if v == nil {
		return v
	}
	switch o := v.(type) {
	case map[string]interface{}:
		r := make(map[string]interface{})
		for k, vv := range o {
			r[k] = copyv(vv)
		}
		return r
	case []interface{}:
		r := make([]interface{}, len(o))
		for i, vv := range o {
			r[i] = copyv(vv)
		}
		return r
	default:
		return v
	}
}

func (this *Client) Export() map[string]interface{} {
	r := make(map[string]interface{})
	if this.sessionId != "" {
		r["SessionId"] = this.sessionId
	}
	if this.props != nil {
		m := make(map[string]interface{})
		for k, v := range this.props {
			m[k] = copyv(v)
		}
		r["Props"] = m
	}
	return r
}

func vtos(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (this *Client) Import(data map[string]interface{}) error {
	if data == nil {
		return nil
	}
	if sv, ok := data["SessionId"]; ok {
		if s, ok2 := sv.(string); ok2 {
			this.sessionId = s
		}
	}
	if mv, ok := data["Props"]; ok {
		if m, ok2 := mv.(map[string]interface{}); ok2 {
			for k, v := range m {
				if this.props == nil {
					this.props = make(map[string]interface{})
				}
				this.props[k] = copyv(v)
			}
		}
	}
	return nil
}

func (this *Client) PollAnswer(addr *Address, an *Answer, ctx *Context, endTime time.Time, sleepDur time.Duration) (*Answer, bool, error) {
	aid := an.GetAsyncId()
	if aid == "" {
		return nil, true, fmt.Errorf("miss AsyncId")
	}
	req := NewRequest()
	ctx.Put(constv.KEY_ASYNC_ID, aid)
	for {
		if time.Now().After(endTime) {
			return nil, false, nil
		}
		an2, err := this.Invoke(addr, req, ctx)
		if err != nil {
			return nil, true, err
		}
		if !an2.IsAsync() {
			return an2, true, nil
		}
		if time.Now().After(endTime) {
			return nil, false, nil
		}
		if sleepDur <= 0 {
			return nil, false, nil
		}
		time.Sleep(sleepDur)
	}
}

func (this *Client) WaitAnswer(du time.Duration) (*Answer, bool, error) {
	aid := an.GetAsyncId()
	if aid == "" {
		return nil, true, fmt.Errorf("miss AsyncId")
	}
	req := NewRequest()
	ctx.Put(constv.KEY_ASYNC_ID, aid)
	for {
		if time.Now().After(endTime) {
			return nil, false, nil
		}
		an2, err := this.Invoke(addr, req, ctx)
		if err != nil {
			return nil, true, err
		}
		if !an2.IsAsync() {
			return an2, true, nil
		}
		if time.Now().After(endTime) {
			return nil, false, nil
		}
		if sleepDur <= 0 {
			return nil, false, nil
		}
		time.Sleep(sleepDur)
	}
}

func (this *Client) SetProperty(n string, val interface{}) {
	if this.props == nil {
		this.props = make(map[string]interface{})
	}
	this.props[n] = val
}

func (this *Client) GetProperty(n string) (interface{}, bool) {
	if this.props == nil {
		return nil, false
	}
	r, ok := this.props[n]
	return r, ok
}

func (this *Client) RemoveProperty(n string) {
	if this.props == nil {
		return
	}
	delete(this.props, n)
}
