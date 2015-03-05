package core

import (
	"bytes"
	"fmt"
)

type Address struct {
	typ    string
	api    string
	option *ValueMap
}

func NewAddress() *Address {
	o := new(Address)
	return o
}

func CreateAddress(typ string, api string, opts map[string]interface{}) *Address {
	o := NewAddress()
	o.typ = typ
	o.api = api
	o.option = CreateValueMap(opts)
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

func (this *Address) GetOption() *ValueMap {
	return this.option
}

func (this *Address) SetOption(v *ValueMap) {
	this.option = v
}

func (this *Address) SureOption() *ValueMap {
	if this.option == nil {
		this.option = NewValueMap(nil)
	}
	return this.option
}

func (this *Address) String() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("TYPE:%s, API:%s", this.typ, this.api))
	if this.option != nil {
		buf.WriteString(fmt.Sprintf(", %s", this.option.Dump()))
	}
	return buf.String()
}
