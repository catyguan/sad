package core

import (
	"time"
)

// Common
type ValueMapWalker func(k string, v *Value) (stop bool)

type ValueArrayWalker func(idx int, v *Value) (stop bool)

type DataConverter func(typ int8, val interface{}) interface{}

// Client
type Closable interface {
	Close()
}

type InvokeContext interface {
	SetProperty(n string, val interface{})
	GetProperty(n string) (interface{}, bool)
	RemoveProperty(n string)
}

type Driver interface {
	CreateConn(typ string, api string) (ServiceConn, error)
}

type ServiceConn interface {
	Invoke(ictx InvokeContext, addr *Address, req *Request, ctx *Context) (*Answer, error)
	WaitAnswer(du time.Duration) (*Answer, error)

	Close()
	End()
}

type ClientFactory func() *Client

// Server
type ServicePeer interface {
	GetDriverType() string

	BeginTransaction() error

	ReadRequest(waitTime time.Duration) (*Request, *Context, error)
	WriteAnswer(a *Answer, err error) error

	SendAsync(ctx *Context, result *ValueMap, timeout time.Duration) error
}

type ServiceProvider func(service, method string) (ServiceMethod, error)

type ServiceObject interface {
	GetMethod(name string) ServiceMethod
}

type ServiceMethod func(peer ServicePeer, req *Request, ctx *Context) error

type ServiceHandler func(peer ServicePeer, service, method string, req *Request, ctx *Context) error
