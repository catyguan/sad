package core

import (
	"bytes"
	"fmt"
)

type Address struct {
	typ string
	api string
	ctx *Context
}

func NewAddress() *Address {
	o := new(Address)
	return o
}

func CreateAddress(typ string, api string, ctx map[string]interface{}) *Address {
	o := NewAddress()
	o.typ = typ
	o.api = api
	o.ctx = CreateContext(ctx)
	return o
}

func (this *Address) GetType() string {
	return this.typ
}

func (this *Address) SetType(v string) {
	this.typ = v
}

func (this *Address) GetAPI() string {
	return this.api
}

func (this *Address) SetAPI(v string) {
	this.api = v
}

func (this *Address) GetContext() *Context {
	return this.ctx
}

func (this *Address) SetContext(v *Context) {
	this.ctx = v
}

func (this *Address) SureContext() *Context {
	if this.ctx == nil {
		this.ctx = NewContext()
	}
	return this.ctx
}

func (this *Address) String() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("TYPE:%s, API:%s", this.typ, this.api))
	if this.ctx != nil {
		buf.WriteString(fmt.Sprintf(", %s", this.ctx))
	}
	return buf.String()
}
