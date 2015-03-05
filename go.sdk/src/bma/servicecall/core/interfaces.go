package core

import (
	"time"
)

type InvokeContext interface {
}

type Driver interface {
	Invoke(ictx InvokeContext, addr *Address, req *Request, ctx *Context) (*Answer, error)
}

type DriverFactory interface {
	GetDriver(typ string) (Driver, error)
}

type ValueMapWalker func(k string, v *Value) (stop bool)

type ValueArrayWalker func(idx int, v *Value) (stop bool)

type ServicePeer interface {
	BeginTransaction()
	EndTransaction()

	ReadRequest(waitTime time.Duration) (*Request, *Context, error)
	WriteAnswer(a *Answer, err error) error
}

type ServiceMethod func(peer ServicePeer, req *Request, ctx *Context) error
