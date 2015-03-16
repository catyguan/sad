package sockserver

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"bma/servicecall/sockcore"
	"fmt"
	"io"
	"net"
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
	serv sccore.BaseServiceServ
	sccore.ServiceMux
}

func NewServiceCallMux(fac sccore.ClientFactory) *ServiceCallMux {
	o := new(ServiceCallMux)
	o.serv.InitBaseServiceServ(fac)
	return o
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
		if errR == io.EOF {
			sccore.DoLog("connection closed")
		} else {
			sccore.DoLog("read Request fail - %s", errR)
		}
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
		pa := this.serv.PollAsync(aid)

		if pa != nil {
			aa := pa.Answer
			aerr := pa.Err
			peer2 := pa.Peer.(*SocketServicePeer)
			peer2.doAnswer(conn, mid, aa, aerr)
		} else {
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
