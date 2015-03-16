package sockserver

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"bma/servicecall/sockcore"
	"fmt"
	"net"
	"time"
)

type SocketServicePeer struct {
	mux       *ServiceCallMux
	conn      net.Conn
	mode      int // 0-normal, 1-writed, 2-poll, 3-callback, 4-push
	asyncId   string
	callback  *sccore.Address
	messageId int32
}

func (this *SocketServicePeer) GetDriverType() string {
	return "socket"
}

func (this *SocketServicePeer) doAnswer(conn net.Conn, mid int32, a *sccore.Answer, err error) error {
	if err != nil {
		a = sccore.Error2Answer(a, err)
	}
	sc := a.GetStatus()
	if sc <= 0 {
		a.SetStatus(200)
	}
	mr := sockcore.NewMessageWriter(conn)
	err2 := mr.SendAnswer(mid, a)
	return err2
}

func (this *SocketServicePeer) BeginTransaction() error {
	return nil
}

func (this *SocketServicePeer) ReadRequest(waitTime time.Duration) (*sccore.Request, *sccore.Context, error) {
	mux := this.mux
	mr := sockcore.NewMessageReader(this.conn)
	var msg sockcore.Message
	this.conn.SetDeadline(time.Now().Add(waitTime))
	defer func() {
		this.conn.SetDeadline(time.Time{})
	}()
	err := mux.nextMessage(this.conn, mr, &msg)
	if err != nil {
		return nil, nil, err
	}
	this.mode = 0
	this.messageId = msg.Id
	return msg.Request, msg.Context, nil
}

func (this *SocketServicePeer) WriteAnswer(a *sccore.Answer, err error) error {
	switch this.mode {
	case 2:
		if this.asyncId == "" {
			return fmt.Errorf("poll mode, asyncId empty")
		}
		mux := this.mux
		mux.serv.SetPollAnswer(this.asyncId, a, err)
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
		an, err2 := this.mux.serv.DoCallback(this.callback, req, ctx)
		sccore.DoLog("callback answer -> %v, %v", err2, an)
		return err
	case 4:
		sccore.DoLog("pushAnswer -> %v, %v", err, a)
		err2 := this.doAnswer(this.conn, this.messageId, a, err)
		this.mode = 1
		return err2
	case 1:
		return fmt.Errorf("SocketServicePeer already answer")
	default:
		sccore.DoLog("writeAnswer -> %v, %v", err, a)
		err2 := this.doAnswer(this.conn, this.messageId, a, err)
		this.mode = 1
		return err2
	}
}

func (this *SocketServicePeer) SendAsync(ctx *sccore.Context, result *sccore.ValueMap, timeout time.Duration) error {
	async := ctx.GetString(constv.KEY_ASYNC_MODE)
	switch async {
	case "", "poll":
		mux := this.mux
		aid := mux.serv.CreatePollAnswer(timeout, this)
		this.mode = 2
		this.asyncId = aid

		a := sccore.NewAnswer()
		a.SetStatus(constv.STATUS_ASYNC)
		if result == nil {
			result = sccore.NewValueMap(nil)
		}
		result.Put(constv.KEY_ASYNC_ID, aid)
		a.SetResult(result)
		return this.doAnswer(this.conn, this.messageId, a, nil)
	case "push":
		this.mode = 4
		a := sccore.NewAnswer()
		a.SetStatus(constv.STATUS_ASYNC)
		a.SetResult(result)
		return this.doAnswer(this.conn, this.messageId, a, nil)
	case "callback":
		addrm := ctx.GetMap(constv.KEY_CALLBACK)
		if addrm == nil {
			err0 := fmt.Errorf("SocketServicePeer Async callback miss address")
			this.WriteAnswer(nil, err0)
			return err0
		}
		this.callback = sccore.CreateAddressFromValue(addrm)

		this.mode = 3
		a := sccore.NewAnswer()
		a.SetStatus(constv.STATUS_ASYNC)
		a.SetResult(result)
		return this.doAnswer(this.conn, this.messageId, a, nil)
	default:
		err := fmt.Errorf("SocketServicePeer not support AsyncMode(%s)", async)
		this.WriteAnswer(nil, err)
		return err
	}
	return nil
}
