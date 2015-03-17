package sockclient

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"bma/servicecall/sockcore"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

var (
	gMessageId int32
)

type SocketDriver struct {
}

func (this *SocketDriver) CreateConn(typ, api string) (sccore.ServiceConn, error) {
	o := new(SocketServiceConn)
	return o, nil
}

func init() {
	df := new(SocketDriver)
	sccore.InitDriver(NAME_DRIVER, df)
}

type SocketServiceConn struct {
	key  string
	conn net.Conn
}

func (this *SocketServiceConn) Close() {
	if this.conn != nil {
		sockcore.SocketPool().CloseSocket(this.key, this.conn)
		this.conn = nil
	}
}

func (this *SocketServiceConn) ret() {
	if this.conn != nil {
		sccore.DoLog("return conn %s", this.conn.LocalAddr())
		sockcore.SocketPool().ReturnSocket(this.key, this.conn)
		this.conn = nil
	}
}

func (this *SocketServiceConn) End() {
	this.Close()
}

func (this *SocketServiceConn) Invoke(ictx sccore.InvokeContext, addr *sccore.Address, req *sccore.Request, ctx *sccore.Context) (*sccore.Answer, error) {
	// async := ctx.GetString(constv.KEY_ASYNC_MODE)
	// if async == "push" {
	// 	return nil, fmt.Errorf("http not support AsyncMode(%s)", async)
	// }
	pool := sockcore.SocketPool()

	sapi, errP := sockcore.ParseSocketAPI(addr.GetAPI())
	if errP != nil {
		return nil, errP
	}

	var conn net.Conn
	if this.conn != nil {
		conn = this.conn
	}

	dltm := ctx.GetLong(constv.KEY_DEADLINE)
	dl := time.Unix(dltm, 0)
	if conn == nil {
		du := dl.Sub(time.Now())
		if du <= 0 {
			return nil, errors.New("timeout")
		}
		sccore.DoLog("'%s' connect...", sapi)
		var errC error
		this.key, conn, errC = pool.GetSocket(addr, sapi, du)
		if errC != nil {
			return nil, errC
		}
		this.conn = conn
	} else {
		sccore.DoLog("'%s' use trans socket", sapi)
	}

	conn.SetDeadline(dl)
	defer conn.SetDeadline(time.Time{})

	// opt := addr.GetOption()
	mw := sockcore.NewMessageWriter(conn)

	sccore.DoLog("'%s' write request to '%s'", sapi, conn.RemoteAddr())
	mid := atomic.AddInt32(&gMessageId, 1)
	if mid <= 0 {
		atomic.CompareAndSwapInt32(&gMessageId, mid, 0)
		mid = atomic.AddInt32(&gMessageId, 1)
	}
	err2 := mw.SendRequest(mid, sapi.Service, sapi.Method, req, ctx)
	if err2 != nil {
		this.Close()
		sccore.DoLog("'%s' fail '%s'", sapi, err2)
		return nil, err2
	}

	mr := sockcore.NewMessageReader(conn)
	var msg sockcore.Message
	for {
		mt, errN := mr.NextMessage(&msg)
		if errN != nil {
			this.Close()
			return nil, errN
		}
		switch mt {
		case sockcore.MT_ANSWER:
			switch msg.Answer.GetStatus() {
			case 100:
				sccore.DoLog("keep connection for transaction")
			case 202:
				amode := ctx.GetString(constv.KEY_ASYNC_MODE)
				if amode == "callback" {
					this.ret()
				} else {
					sccore.DoLog("keep connection for async")
				}
			case 200, 204, 302:
				this.ret()
			default:
				this.Close()
			}
			return msg.Answer, nil
		default:
			sccore.DoLog("unknow message(%d) - (%v)", mt, &msg)
		}
	}
}

func (this *SocketServiceConn) WaitAnswer(du time.Duration) (*sccore.Answer, error) {
	conn := this.conn
	if this.conn == nil {
		return nil, fmt.Errorf("invalid connection for ServerPush")
	}

	dl := time.Now().Add(du)
	conn.SetDeadline(dl)
	defer conn.SetDeadline(time.Time{})

	mr := sockcore.NewMessageReader(conn)
	var msg sockcore.Message
	for {
		mt, errN := mr.NextMessage(&msg)
		if errN != nil {
			this.Close()
			return nil, errN
		}
		an := msg.Answer
		switch mt {
		case sockcore.MT_ANSWER:
			switch an.GetStatus() {
			case 202:
			case 200, 204:
				this.ret()
			default:
				this.Close()
			}
			return an, nil
		default:
			sccore.DoLog("unknow message(%d) - (%v)", mt, &msg)
		}
	}
}
