package core

import "fmt"
import (
	"bma/servicecall/constv"
)

type Request struct {
	ValueMap
}

func NewRequest() *Request {
	o := new(Request)
	return o
}

func CreateRequest(data map[string]interface{}) *Request {
	o := new(Request)
	initValueMap(&o.ValueMap, data)
	return o
}

type Context struct {
	ValueMap
}

func NewContext() *Context {
	o := new(Context)
	return o
}

func CreateContext(data map[string]interface{}) *Context {
	o := NewContext()
	initValueMap(&o.ValueMap, data)
	return o
}

func (this *Context) GetTransactioId() string {
	return this.GetString(constv.KEY_TRANSACTION_ID)
}

func (this *Context) GetSessionId() string {
	return this.GetString(constv.KEY_SESSION_ID)
}

type Answer struct {
	status  int
	message string
	result  *ValueMap
	context *ValueMap
}

func NewAnswer() *Answer {
	o := new(Answer)
	return o
}

func (this *Answer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["Status"] = this.status
	if this.message != "" {
		m["Message"] = this.message
	}
	if this.result != nil {
		m["Result"] = this.result.ToMap()
	}
	if this.context != nil {
		m["Context"] = this.context.ToMap()
	}
	return m
}

func (this *Answer) Dump() string {
	return fmt.Sprintf("%v", this.ToMap())
}

func (this *Answer) String() string {
	return this.Dump()
}

func (this *Answer) IsProcessing() bool {
	st := this.GetStatus()
	switch st {
	case 102:
		return true
	}
	return false
}

func (this *Answer) IsDone() bool {
	st := this.GetStatus()
	switch st {
	case 100, 200, 204:
		return true
	}
	return false
}

func (this *Answer) GetStatus() int {
	return this.status
}

func (this *Answer) SetStatus(v int) {
	this.status = v
}

func (this *Answer) GetMessage() string {
	return this.message
}

func (this *Answer) SetMessage(v string) {
	this.message = v
}

func (this *Answer) GetResult() *ValueMap {
	return this.result
}

func (this *Answer) SureResult() *ValueMap {
	if this.result == nil {
		this.result = NewValueMap(nil)
	}
	return this.result
}

func (this *Answer) SetResult(v *ValueMap) {
	this.result = v
}

func (this *Answer) GetContext() *ValueMap {
	return this.context
}

func (this *Answer) SureContext() *ValueMap {
	if this.context == nil {
		this.context = NewValueMap(nil)
	}
	return this.context
}

func (this *Answer) SetContext(v *ValueMap) {
	this.context = v
}

func (this *Answer) SetSessionId(v string) {
	this.SureContext().Put(constv.KEY_SESSION_ID, v)
}
