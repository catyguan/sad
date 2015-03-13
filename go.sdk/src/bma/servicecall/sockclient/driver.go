package sockclient

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"bma/servicecall/sockcore"
	"errors"
	"net"
	"sync/atomic"
	"time"
)

var (
	gMessageId int32
)

type SocketDriverFactory struct {
}

func (this *SocketDriverFactory) GetDriver(typ string) (sccore.Driver, error) {
	o := new(SocketDriver)
	return o, nil
}

func init() {
	df := new(SocketDriverFactory)
	sccore.InitDriverFactory(NAME_DRIVER, df)
}

type SocketDriver struct {
}

type useSocket struct {
	trans map[string]net.Conn
}

func (this *useSocket) Close() {
	pool := sockcore.SocketPool()
	for k, conn := range this.trans {
		pool.CloseSocket(k, conn)
	}
}

func (this *SocketDriver) Invoke(ictx sccore.InvokeContext, addr *sccore.Address, req *sccore.Request, ctx *sccore.Context) (*sccore.Answer, error) {
	// async := ctx.GetString(constv.KEY_ASYNC_MODE)
	// if async == "push" {
	// 	return nil, fmt.Errorf("http not support AsyncMode(%s)", async)
	// }
	pool := sockcore.SocketPool()

	sapi, errP := sockcore.ParseSocketAPI(addr.GetAPI())
	if errP != nil {
		return nil, errP
	}

	key := sapi.Key()
	var conn net.Conn
	var us *useSocket
	so, ok := ictx.GetProperty("socket")
	if ok {
		us, ok = so.(*useSocket)
		if ok {
			conn = us.trans[key]
		}
	}

	dltm := ctx.GetLong(constv.KEY_DEADLINE)
	dl := time.Unix(dltm, 0)
	closemode := 0
	if conn == nil {
		du := dl.Sub(time.Now())
		if du <= 0 {
			return nil, errors.New("timeout")
		}
		sccore.DoLog("'%s' connect...", sapi)
		var errC error
		key, conn, errC = pool.GetSocket(addr, sapi, du)
		if errC != nil {
			return nil, errC
		}
	} else {
		sccore.DoLog("'%s' use trans socket", sapi)
		closemode = 2
	}

	defer func() {
		switch closemode {
		case 0:
			pool.ReturnSocket(key, conn)
		case 1:
			if us != nil {
				delete(us.trans, key)
			}
			pool.CloseSocket(key, conn)
		case 2:
		}
	}()

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
		sccore.DoLog("'%s' fail '%s'", sapi, err2)
		closemode = 1
		return nil, err2
	}

	mr := sockcore.NewMessageReader(conn)
	var msg sockcore.Message
	for {
		mt, errN := mr.NextMessage(&msg)
		if errN != nil {
			closemode = 1
			return nil, errN
		}
		switch mt {
		case sockcore.MT_PING:
			// skip
			sccore.DoLog("ping response(%v)", msg.BoolFlag)
		case sockcore.MT_TRANSACTION:
			// skip
			sccore.DoLog("'%s' transaction(%v)", conn.RemoteAddr(), msg.BoolFlag)
			if msg.BoolFlag {
				closemode = 2
				if us == nil {
					us = new(useSocket)
					us.trans = make(map[string]net.Conn)
					ictx.SetProperty("socket", us)
				}
				us.trans[key] = conn
			} else {
				closemode = 0
				if us != nil {
					delete(us.trans, key)
				}
			}
		case sockcore.MT_ANSWER:
			return msg.Answer, nil
		default:
			sccore.DoLog("unknow message(%d) - (%v)", mt, &msg)
		}
	}
}
