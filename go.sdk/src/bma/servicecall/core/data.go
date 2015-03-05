package core

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

type Answer struct {
	status  int
	message string
	result  *ValueMap
}

func NewAnswer() *Answer {
	o := new(Answer)
	return o
}

func (this *Answer) Dump() string {
	return ""
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
	if this.result != nil {
		this.result = NewValueMap(nil)
	}
	return this.result
}

func (this *Answer) SetResult(v *ValueMap) {
	this.result = v
}
