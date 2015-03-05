package core

import (
	"bma/servicecall/constv"
	"fmt"
	"strconv"
)

type Value struct {
	typ int8
	val interface{}
}

func NewValue(typ int8, val interface{}) *Value {
	o := new(Value)
	o.typ = typ
	o.val = val
	return o
}

func CreateValue(val interface{}) *Value {
	t, v, err := ConvertValue(val)
	if err != nil {
		return new(Value)
	}
	if t != constv.TYPES_VAR {
		o := new(Value)
		o.typ = t
		o.val = v
		return o
	} else {
		return v.(*Value)
	}
}

type isstring interface {
	String() string
}

func ConvertValue(val interface{}) (int8, interface{}, error) {
	if val == nil {
		return 0, nil, nil
	}
	switch rv := val.(type) {
	case *Value:
		return constv.TYPES_VAR, rv, nil
	case bool:
		return constv.TYPES_BOOL, rv, nil
	case int:
		return constv.TYPES_INT, int32(rv), nil
	case int8:
		return constv.TYPES_INT, int32(rv), nil
	case int16:
		return constv.TYPES_INT, int32(rv), nil
	case int32:
		return constv.TYPES_INT, int32(rv), nil
	case int64:
		return constv.TYPES_LONG, int64(rv), nil
	case uint:
		return constv.TYPES_INT, int32(rv), nil
	case uint8:
		return constv.TYPES_INT, int32(rv), nil
	case uint16:
		return constv.TYPES_INT, int32(rv), nil
	case uint32:
		return constv.TYPES_LONG, int64(rv), nil
	case uint64:
		return constv.TYPES_LONG, int64(rv), nil
	case float32:
		return constv.TYPES_FLOAT, rv, nil
	case float64:
		return constv.TYPES_DOUBLE, rv, nil
	case string:
		return constv.TYPES_STRING, rv, nil
	case *ValueArray:
		return constv.TYPES_ARRAY, rv, nil
	case []*Value:
		return constv.TYPES_ARRAY, NewValueArray(rv), nil
	case []interface{}:
		return constv.TYPES_ARRAY, CreateValueArray(rv), nil
	case *ValueMap:
		return constv.TYPES_MAP, rv, nil
	case map[string]*Value:
		return constv.TYPES_MAP, NewValueMap(rv), nil
	case map[string]interface{}:
		return constv.TYPES_MAP, CreateValueMap(rv), nil
	case []byte:
		return constv.TYPES_BINARY, rv, nil
	}
	if s, ok := val.(isstring); ok {
		return constv.TYPES_STRING, s.String(), nil
	}
	return 0, nil, fmt.Errorf("unknow value(%T)", val)
}

func (this *Value) As(typ int8) interface{} {
	switch typ {
	case constv.TYPES_BOOL:
		return this.AsBool()
	case constv.TYPES_INT:
		return this.AsInt()
	case constv.TYPES_LONG:
		return this.AsLong()
	case constv.TYPES_FLOAT:
		return this.AsFloat()
	case constv.TYPES_DOUBLE:
		return this.AsDouble()
	case constv.TYPES_STRING:
		return this.AsString()
	case constv.TYPES_VAR:
		return this
	case constv.TYPES_ARRAY:
		return this.AsArray()
	case constv.TYPES_MAP:
		return this.AsMap()
	case constv.TYPES_BINARY:
		return this.AsBinary()
	}
	return nil
}

func (this *Value) ToValue() interface{} {
	switch this.typ {
	case constv.TYPES_NULL:
		return nil
	case constv.TYPES_ARRAY:
		o := this.val.(*ValueArray)
		r := make([]interface{}, len(o.data))
		for i, v := range o.data {
			r[i] = v.ToValue()
		}
		return r
	case constv.TYPES_MAP:
		if o, ok := this.val.(*ValueMap); ok {
			r := make(map[string]interface{})
			for k, v := range o.data {
				r[k] = v.ToValue()
			}
			return r
		}
		return nil
	default:
		return this.val
	}
}

func (this *Value) GetType() int8 {
	return this.typ
}

func (this *Value) GetValue() interface{} {
	return this.val
}

func (this *Value) Dump() string {
	o := this.ToValue()
	return fmt.Sprintf("%v", o)
}

func string_int64(v string) (int64, bool) {
	o, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return o, true
	}
	return 0, false
}

func string_float64(v string) (float64, bool) {
	o, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return o, true
	}
	return 0, false
}

func (this *Value) AsBool() bool {
	if this.val == nil {
		return false
	}
	switch this.typ {
	case constv.TYPES_BOOL:
		if rv, ok := this.val.(bool); ok {
			return rv
		}
		return false
	case constv.TYPES_INT, constv.TYPES_LONG, constv.TYPES_FLOAT, constv.TYPES_DOUBLE:
		rv := this.AsInt()
		return rv != 0
	case constv.TYPES_STRING:
		if rv, ok := this.val.(string); ok {
			o, err := strconv.ParseBool(rv)
			if err == nil {
				return o
			}
		}
		return false
	case constv.TYPES_ARRAY:
		if rv, ok := this.val.(*ValueArray); ok {
			return rv != nil && rv.Len() > 0
		}
		return false
	case constv.TYPES_MAP:
		if rv, ok := this.val.(*ValueMap); ok {
			return rv != nil && rv.Len() > 0
		}
		return false
	case constv.TYPES_BINARY:
		if rv, ok := this.val.([]byte); ok {
			return rv != nil && len(rv) > 0
		}
		return false
	}
	return false
}

func (this *Value) AsInt() int32 {
	if this.val == nil {
		return 0
	}
	switch this.typ {
	case constv.TYPES_BOOL:
		if rv, ok := this.val.(bool); ok && rv {
			return 1
		}
		return 0
	case constv.TYPES_INT:
		if rv, ok := this.val.(int32); ok {
			return rv
		}
		return 0
	case constv.TYPES_LONG:
		if rv, ok := this.val.(int64); ok {
			return int32(rv)
		}
		return 0
	case constv.TYPES_FLOAT:
		if rv, ok := this.val.(float32); ok {
			return int32(rv)
		}
		return 0
	case constv.TYPES_DOUBLE:
		if rv, ok := this.val.(float64); ok {
			return int32(rv)
		}
		return 0
	case constv.TYPES_STRING:
		if rv, ok := this.val.(string); ok {
			o, sok := string_int64(rv)
			if sok {
				return int32(o)
			}
		}
		return 0
	case constv.TYPES_ARRAY:
		return 0
	case constv.TYPES_MAP:
		return 0
	case constv.TYPES_BINARY:
		return 0
	}
	return 0
}

func (this *Value) AsLong() int64 {
	if this.val == nil {
		return 0
	}
	switch this.typ {
	case constv.TYPES_BOOL:
		if rv, ok := this.val.(bool); ok && rv {
			return 1
		}
		return 0
	case constv.TYPES_INT:
		if rv, ok := this.val.(int32); ok {
			return int64(rv)
		}
		return 0
	case constv.TYPES_LONG:
		if rv, ok := this.val.(int64); ok {
			return int64(rv)
		}
		return 0
	case constv.TYPES_FLOAT:
		if rv, ok := this.val.(float32); ok {
			return int64(rv)
		}
		return 0
	case constv.TYPES_DOUBLE:
		if rv, ok := this.val.(float64); ok {
			return int64(rv)
		}
		return 0
	case constv.TYPES_STRING:
		if rv, ok := this.val.(string); ok {
			o, sok := string_int64(rv)
			if sok {
				return o
			}
		}
		return 0
	case constv.TYPES_ARRAY:
		return 0
	case constv.TYPES_MAP:
		return 0
	case constv.TYPES_BINARY:
		return 0
	}
	return 0
}

func (this *Value) AsFloat() float32 {
	if this.val == nil {
		return 0
	}
	switch this.typ {
	case constv.TYPES_BOOL:
		if rv, ok := this.val.(bool); ok && rv {
			return 1
		}
		return 0
	case constv.TYPES_INT:
		if rv, ok := this.val.(int32); ok {
			return float32(rv)
		}
		return 0
	case constv.TYPES_LONG:
		if rv, ok := this.val.(int64); ok {
			return float32(rv)
		}
		return 0
	case constv.TYPES_FLOAT:
		if rv, ok := this.val.(float32); ok {
			return float32(rv)
		}
		return 0
	case constv.TYPES_DOUBLE:
		if rv, ok := this.val.(float64); ok {
			return float32(rv)
		}
		return 0
	case constv.TYPES_STRING:
		if rv, ok := this.val.(string); ok {
			o, sok := string_float64(rv)
			if sok {
				return float32(o)
			}
		}
		return 0
	case constv.TYPES_ARRAY:
		return 0
	case constv.TYPES_MAP:
		return 0
	case constv.TYPES_BINARY:
		return 0
	}
	return 0
}

func (this *Value) AsDouble() float64 {
	if this.val == nil {
		return 0
	}
	switch this.typ {
	case constv.TYPES_BOOL:
		if rv, ok := this.val.(bool); ok && rv {
			return 1
		}
		return 0
	case constv.TYPES_INT:
		if rv, ok := this.val.(int32); ok {
			return float64(rv)
		}
		return 0
	case constv.TYPES_LONG:
		if rv, ok := this.val.(int64); ok {
			return float64(rv)
		}
		return 0
	case constv.TYPES_FLOAT:
		if rv, ok := this.val.(float32); ok {
			return float64(rv)
		}
		return 0
	case constv.TYPES_DOUBLE:
		if rv, ok := this.val.(float64); ok {
			return float64(rv)
		}
		return 0
	case constv.TYPES_STRING:
		if rv, ok := this.val.(string); ok {
			o, sok := string_float64(rv)
			if sok {
				return float64(o)
			}
		}
		return 0
	case constv.TYPES_ARRAY:
		return 0
	case constv.TYPES_MAP:
		return 0
	case constv.TYPES_BINARY:
		return 0
	}
	return 0
}

func (this *Value) AsString() string {
	if this.val == nil {
		return ""
	}
	switch this.typ {
	case constv.TYPES_BOOL, constv.TYPES_INT, constv.TYPES_LONG, constv.TYPES_FLOAT, constv.TYPES_DOUBLE:
		return fmt.Sprintf("%v", this.val)
	case constv.TYPES_STRING:
		if rv, ok := this.val.(string); ok {
			return rv
		}
		return ""
	case constv.TYPES_BINARY:
		if rv, ok := this.val.([]byte); ok {
			return string(rv)
		}
		return ""
	}
	return ""
}

func (this *Value) AsBinary() []byte {
	if this.val == nil {
		return nil
	}
	switch this.typ {
	case constv.TYPES_STRING:
		if rv, ok := this.val.(string); ok {
			return []byte(rv)
		}
		return nil
	case constv.TYPES_BINARY:
		if rv, ok := this.val.([]byte); ok {
			return rv
		}
		return nil
	}
	return nil
}

func (this *Value) AsArray() *ValueArray {
	switch this.typ {
	case constv.TYPES_ARRAY:
		if rv, ok := this.val.(*ValueArray); ok {
			return rv
		}
		return nil
	}
	return nil
}

func (this *Value) AsMap() *ValueMap {
	switch this.typ {
	case constv.TYPES_MAP:
		if rv, ok := this.val.(*ValueMap); ok {
			return rv
		}
		return nil
	}
	return nil
}

type ValueMap struct {
	data map[string]*Value
}

func NewValueMap(data map[string]*Value) *ValueMap {
	o := new(ValueMap)
	o.data = data
	return o
}

func initValueMap(o *ValueMap, data map[string]interface{}) {
	if data != nil {
		r := make(map[string]*Value)
		for k, av := range data {
			r[k] = CreateValue(av)
		}
		o.data = r
	}
}

func CreateValueMap(data map[string]interface{}) *ValueMap {
	o := new(ValueMap)
	initValueMap(o, data)
	return o
}

func (this *ValueMap) Dump() string {
	o := this.ToMap()
	return fmt.Sprintf("%v", o)
}

func (this *ValueMap) Len() int {
	return len(this.data)
}

func (this *ValueMap) GetData() map[string]*Value {
	return this.data
}

func (this *ValueMap) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range this.data {
		r[k] = v.ToValue()
	}
	return r
}

func (this *ValueMap) Has(k string) bool {
	if this.data == nil {
		return false
	}
	_, ok := this.data[k]
	return ok
}

func (this *ValueMap) Get(k string) *Value {
	if this.data == nil {
		return nil
	}
	if rv, ok := this.data[k]; ok {
		return rv
	}
	return nil
}

func (this *ValueMap) GetBool(k string) bool {
	v := this.Get(k)
	if v == nil {
		return false
	}
	return v.AsBool()
}

func (this *ValueMap) GetInt(k string) int32 {
	v := this.Get(k)
	if v == nil {
		return 0
	}
	return v.AsInt()
}

func (this *ValueMap) GetLong(k string) int64 {
	v := this.Get(k)
	if v == nil {
		return 0
	}
	return v.AsLong()
}

func (this *ValueMap) GetFloat(k string) float32 {
	v := this.Get(k)
	if v == nil {
		return 0
	}
	return v.AsFloat()
}

func (this *ValueMap) GetDouble(k string) float64 {
	v := this.Get(k)
	if v == nil {
		return 0
	}
	return v.AsDouble()
}

func (this *ValueMap) GetString(k string) string {
	v := this.Get(k)
	if v == nil {
		return ""
	}
	return v.AsString()
}

func (this *ValueMap) GetArray(k string) *ValueArray {
	v := this.Get(k)
	if v == nil {
		return nil
	}
	return v.AsArray()
}

func (this *ValueMap) GetMap(k string) *ValueMap {
	v := this.Get(k)
	if v == nil {
		return nil
	}
	return v.AsMap()
}

func (this *ValueMap) GetBinary(k string) []byte {
	v := this.Get(k)
	if v == nil {
		return nil
	}
	return v.AsBinary()
}

func (this *ValueMap) Set(k string, v *Value) {
	if this.data == nil {
		this.data = make(map[string]*Value)
	}
	this.data[k] = v
}

func (this *ValueMap) Put(k string, val interface{}) {
	v := CreateValue(val)
	this.Set(k, v)
}

func (this *ValueMap) CreateMap(k string) *ValueMap {
	v := this.GetMap(k)
	if v != nil {
		return v
	}
	r := NewValueMap(nil)
	this.Set(k, NewValue(constv.TYPES_MAP, r))
	return r
}

func (this *ValueMap) CreateArray(k string) *ValueArray {
	v := this.GetArray(k)
	if v != nil {
		return v
	}
	r := NewValueArray(nil)
	this.Set(k, NewValue(constv.TYPES_ARRAY, r))
	return r
}

func (this *ValueMap) Remove(k string) {
	if this.data == nil {
		return
	}
	delete(this.data, k)
}

func (this *ValueMap) Walk(w ValueMapWalker) {
	for k, v := range this.data {
		if w(k, v) {
			return
		}
	}
}

type ValueArray struct {
	data []*Value
}

func NewValueArray(data []*Value) *ValueArray {
	o := new(ValueArray)
	o.data = data
	return o
}

func CreateValueArray(data []interface{}) *ValueArray {
	o := new(ValueArray)
	if data != nil {
		r := make([]*Value, len(data))
		for i, av := range data {
			r[i] = CreateValue(av)
		}
		o.data = r
	}
	return o
}

func (this *ValueArray) Dump() string {
	o := this.ToArray()
	return fmt.Sprintf("%v", o)
}

func (this *ValueArray) Len() int {
	return len(this.data)
}

func (this *ValueArray) GetData() []*Value {
	return this.data
}

func (this *ValueArray) ToArray() []interface{} {
	r := make([]interface{}, len(this.data))
	for i, v := range this.data {
		r[i] = v.ToValue()
	}
	return r
}

func (this *ValueArray) Get(idx int) *Value {
	if this.data == nil {
		return nil
	}
	if idx < len(this.data) {
		return this.data[idx]
	}
	return nil
}

func (this *ValueArray) Set(idx int, v *Value) bool {
	if this.data == nil {
		return false
	}
	if idx < len(this.data) {
		this.data[idx] = v
		return true
	}
	return false
}

func (this *ValueArray) Add(v *Value) {
	this.data = append(this.data, v)
}

func (this *ValueArray) Remove(idx int) {
	if this.data == nil {
		return
	}
	if idx < 0 || idx >= len(this.data) {
		return
	}
	s := this.data
	this.data = append(s[:idx], s[idx+1:]...)
}

func (this *ValueArray) Walk(w ValueArrayWalker) {
	for i, v := range this.data {
		if w(i, v) {
			return
		}
	}
}
