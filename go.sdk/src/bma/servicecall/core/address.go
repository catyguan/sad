package core

import (
	"bytes"
	"fmt"
	"strings"
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

func CreateAddressFromMap(vm map[string]interface{}) *Address {
	return CreateAddressFromValue(CreateValueMap(vm))
}

func CreateAddressFromValue(vm *ValueMap) *Address {
	o := NewAddress()
	o.typ = vm.GetString("Type")
	o.api = vm.GetString("API")
	o.option = vm.GetMap("Option")
	return o
}

func (this *Address) Valid() error {
	if this.typ == "" {
		return fmt.Errorf("address type empty")
	}
	if this.api == "" {
		return fmt.Errorf("address api empty")
	}
	return nil
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

func (this *Address) ToValueMap() *ValueMap {
	vm := NewValueMap(nil)
	vm.Put("Type", this.typ)
	vm.Put("API", this.api)
	if this.option != nil && this.option.Len() > 0 {
		vm.Put("Option", this.option)
	}
	return vm
}

// Builder
type SimpleAddressBuilder struct {
	Type   string
	API    string
	Option *ValueMap
}

func NewAddressBuilder() *SimpleAddressBuilder {
	o := new(SimpleAddressBuilder)
	return o
}

func (this *SimpleAddressBuilder) Build(service, method string) *Address {
	s := strings.Replace(this.API, "$SNAME$", service, -1)
	s = strings.Replace(s, "$MNAME$", method, -1)
	r := NewAddress()
	r.SetType(this.Type)
	r.SetAPI(s)
	r.SetOption(this.Option)
	return r
}
