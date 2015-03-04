package core

type InvokeContext interface {
}

type Driver interface {
	Invoke(ictx InvokeContext, addr *Address, req *Request, ctx *Context) (*Answer, error)
}

type DriverFactory interface {
	GetDriver(typ string, ctx *Context) (Driver, error)
}
