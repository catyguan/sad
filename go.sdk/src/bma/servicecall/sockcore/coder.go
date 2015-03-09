package sockcore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
)

// lenBytes
type LenBytesCoder int32

func (this LenBytesCoder) DoEncode(w EncodeWriter, bs []byte) error {
	l := len(bs)
	Coders.Int32.DoEncode(w, int32(l))
	if l > 0 {
		_, err := w.Write(bs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this LenBytesCoder) Encode(w EncodeWriter, v interface{}) error {
	return this.DoEncode(w, v.([]byte))
}

func (this LenBytesCoder) DoDecode(r DecodeReader, maxlen int32) ([]byte, error) {
	l, err := Coders.Int32.DoDecode(r)
	if err != nil {
		return nil, err
	}
	if maxlen <= 0 {
		maxlen = DATA_MAXLEN
	}
	if l > maxlen {
		return nil, fmt.Errorf("too large bytes block - %d/%d", l, maxlen)
	}
	p := make([]byte, l)
	if l > 0 {
		_, err = r.Read(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (this LenBytesCoder) Decode(r DecodeReader) (interface{}, error) {
	s, err := this.DoDecode(r, int32(this))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// lenString
type LenStringCoder int32

func (this LenStringCoder) DoEncode(w EncodeWriter, v string) error {
	bs := []byte(v)
	err := Coders.Int32.DoEncode(w, int32(len(bs)))
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (this LenStringCoder) Encode(w EncodeWriter, v interface{}) error {
	return this.DoEncode(w, v.(string))
}

func (this LenStringCoder) DoDecode(r DecodeReader, maxlen int32) (string, error) {
	l, err := Coders.Int32.DoDecode(r)
	if err != nil {
		return "", err
	}
	if maxlen <= 0 {
		maxlen = DATA_MAXLEN
	}
	if l > maxlen {
		return "", fmt.Errorf("too large string block - %d/%d", l, maxlen)
	}
	p := make([]byte, l)
	_, err = r.Read(p)
	if err != nil {
		return "", err
	}
	return string(p), nil
}

func (this LenStringCoder) Decode(r DecodeReader) (interface{}, error) {
	s, err := this.DoDecode(r, int32(this))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// bool
type boolCoder bool

func (this boolCoder) DoEncode(w EncodeWriter, v bool) error {
	b := byte(0)
	if v {
		b = 1
	}
	return w.WriteByte(b)
}

func (this boolCoder) Encode(w EncodeWriter, v interface{}) error {
	return this.DoEncode(w, v.(bool))
}

func (this boolCoder) DoDecode(r DecodeReader) (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, err
}

func (this boolCoder) Decode(r DecodeReader) (interface{}, error) {
	v, err := this.DoDecode(r)
	return v, err
}

// intx
type int32Coder int
type int64Coder int
type uint32Coder int
type uint64Coder int

func (O int32Coder) DoEncode(w EncodeWriter, v int32) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O int32Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(int32))
}
func (O int64Coder) DoEncode(w EncodeWriter, v int64) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O int64Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(int64))
}
func (O uint32Coder) DoEncode(w EncodeWriter, v uint32) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O uint32Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint32))
}
func (O uint64Coder) DoEncode(w EncodeWriter, v uint64) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O uint64Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint64))
}

func (O int32Coder) DoDecode(r DecodeReader) (int32, error) {
	rv, err := binary.ReadVarint(r)
	return int32(rv), err
}
func (O int32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int64Coder) DoDecode(r DecodeReader) (int64, error) {
	rv, err := binary.ReadVarint(r)
	return int64(rv), err
}
func (O int64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint32Coder) DoDecode(r DecodeReader) (uint32, error) {
	rv, err := binary.ReadUvarint(r)
	return uint32(rv), err
}
func (O uint32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint64Coder) DoDecode(r DecodeReader) (uint64, error) {
	rv, err := binary.ReadUvarint(r)
	return uint64(rv), err
}
func (O uint64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}

// fixIntxCoder
type fixUint8Coder int
type fixUint16Coder int
type fixUint32Coder int
type fixUint64Coder int

func (O fixUint8Coder) DoEncode(w EncodeWriter, v uint8) error {
	return w.WriteByte(v)
}
func (O fixUint8Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint8))
}

func (O fixUint16Coder) DoEncode(w EncodeWriter, v uint16) error {
	bs := [2]byte{}
	b := bs[:]
	binary.BigEndian.PutUint16(b, uint16(v))
	_, err := w.Write(b)
	return err
}
func (O fixUint16Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint16))
}
func (O fixUint32Coder) DoEncode(w EncodeWriter, v uint32) error {
	bs := [4]byte{}
	b := bs[:]
	binary.BigEndian.PutUint32(b, uint32(v))
	_, err := w.Write(b)
	return err
}
func (O fixUint32Coder) Encode(w EncodeWriter, v interface{}) error {
	O.DoEncode(w, v.(uint32))
	return nil
}
func (O fixUint64Coder) DoEncode(w EncodeWriter, v uint64) error {
	bs := [8]byte{}
	b := bs[:]
	binary.BigEndian.PutUint64(b, uint64(v))
	_, err := w.Write(b)
	return err
}
func (O fixUint64Coder) Encode(w EncodeWriter, v interface{}) error {
	O.DoEncode(w, v.(uint64))
	return nil
}
func (O fixUint8Coder) DoDecode(r DecodeReader) (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}
func (O fixUint8Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}

func (O fixUint16Coder) DoDecode(r DecodeReader) (uint16, error) {
	bs := [2]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}
func (O fixUint16Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O fixUint32Coder) DoDecode(r DecodeReader) (uint32, error) {
	bs := [4]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}
func (O fixUint32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O fixUint64Coder) DoDecode(r DecodeReader) (uint64, error) {
	bs := [8]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}
func (O fixUint64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}

// Float32 Float64 Coder
type float32Coder int
type float64Coder int

func (O float32Coder) DoEncode(w EncodeWriter, v float32) error {
	iv := math.Float32bits(v)
	return Coders.FixUint32.DoEncode(w, iv)
}
func (O float32Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(float32))
}
func (O float64Coder) DoEncode(w EncodeWriter, v float64) error {
	iv := math.Float64bits(v)
	return Coders.FixUint64.DoEncode(w, iv)
}
func (O float64Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(float64))
}

func (O float32Coder) DoDecode(r DecodeReader) (float32, error) {
	bs := [4]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	iv := binary.BigEndian.Uint32(b)
	return math.Float32frombits(iv), nil
}
func (O float32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O float64Coder) DoDecode(r DecodeReader) (float64, error) {
	bs := [8]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	iv := binary.BigEndian.Uint64(b)
	return math.Float64frombits(iv), nil
}
func (O float64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}

// varCoder
type varCoder int

func (this varCoder) Encode(w EncodeWriter, v interface{}) error {
	if v == nil {
		w.WriteByte(0)
		return nil
	}
	var err error
	var b [binary.MaxVarintLen64]byte
	bs := b[:]

	if rb, ok := v.([]byte); ok {
		err = w.WriteByte(10)
		if err != nil {
			return err
		}
		err = Coders.LenBytes.DoEncode(w, rb)
		if err != nil {
			return err
		}
		return nil
	}

	tv := reflect.ValueOf(v)
	switch tv.Kind() {
	case reflect.Bool:
		err = w.WriteByte(1)
		if err != nil {
			return err
		}
		rv := tv.Bool()
		return Coders.Bool.DoEncode(w, rv)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tvk := byte(3)
		rv := tv.Int()
		if rv <= 2147483647 && rv >= -2147483648 {
			tvk = byte(2)
		}
		err = w.WriteByte(tvk)
		if err != nil {
			return err
		}
		l := binary.PutVarint(bs, rv)
		_, err = w.Write(b[:l])
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv := tv.Uint()
		if rv > 9223372036854775807 {
			err = w.WriteByte(6)
			if err != nil {
				return err
			}
			s := fmt.Sprintf("%d", rv)
			return Coders.LenString.DoEncode(w, s)
		} else {
			tvk := byte(3)
			if rv <= 2147483647 {
				tvk = byte(2)
			}
			err = w.WriteByte(tvk)
			if err != nil {
				return err
			}
			l := binary.PutUvarint(bs, rv)
			_, err = w.Write(b[:l])
			return err
		}
	case reflect.Float32:
		err = w.WriteByte(4)
		if err != nil {
			return err
		}
		return Coders.Float32.DoEncode(w, float32(tv.Float()))
	case reflect.Float64:
		err = w.WriteByte(5)
		if err != nil {
			return err
		}
		return Coders.Float64.DoEncode(w, tv.Float())
	case reflect.String:
		err = w.WriteByte(6)
		if err != nil {
			return err
		}
		return Coders.LenString.DoEncode(w, tv.String())
	case reflect.Map:
		if tv.Type().Key().Kind() != reflect.String {
			return errors.New("onlye encode map[string]value")
		}
		err = w.WriteByte(9)
		if err != nil {
			return err
		}
		sz := tv.Len()
		l := binary.PutUvarint(bs, uint64(sz))
		_, err = w.Write(bs[:l])
		if err != nil {
			return err
		}
		mkeys := tv.MapKeys()
		for _, k := range mkeys {
			sval := tv.MapIndex(k)
			err = Coders.LenString.DoEncode(w, k.String())
			if err != nil {
				return err
			}
			err = this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Ptr:
		return this.Encode(w, tv.Elem())
	case reflect.Slice:
		err = w.WriteByte(8)
		if err != nil {
			return err
		}
		sz := tv.Len()
		l := binary.PutUvarint(bs, uint64(sz))
		_, err = w.Write(bs[:l])
		if err != nil {
			return err
		}
		for i := 0; i < sz; i++ {
			sval := tv.Index(i)
			err = this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		err = w.WriteByte(9)
		if err != nil {
			return err
		}
		vt := tv.Type()
		sz := vt.NumField()
		l := binary.PutUvarint(bs, uint64(sz))
		_, err = w.Write(bs[:l])
		if err != nil {
			return err
		}
		for i := 0; i < sz; i++ {
			tfield := vt.Field(i)
			sval := tv.Field(i)
			err = Coders.LenString.DoEncode(w, tfield.Name)
			if err != nil {
				return err
			}
			err = this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New(fmt.Sprintf("unknow type %T", v))
	}
}

func (this varCoder) Decode(r DecodeReader) (interface{}, error) {
	var b [binary.MaxVarintLen64]byte
	bs := b[:]
	var err0 error
	bs[0], err0 = r.ReadByte()
	if err0 != nil {
		return nil, err0
	}
	k := bs[0]
	switch k {
	case 0:
		return nil, nil
	case 10:
		return Coders.LenBytes.DoDecode(r, 0)
	case 1:
		return Coders.Bool.DoDecode(r)
	case 2, 3:
		rv, err := binary.ReadVarint(r)
		if err != nil {
			return nil, err
		}
		if k == 2 {
			return int32(rv), nil
		}
		return rv, nil
	case 4:
		return Coders.Float32.Decode(r)
	case 5:
		return Coders.Float64.Decode(r)
	case 6:
		return Coders.LenString.Decode(r)
	case 9:
		l, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
		if l == 0 {
			return nil, nil
		}

		rv := make(map[string]interface{})
		for i := 0; i < int(l); i++ {
			kv, err2 := Coders.LenString.DoDecode(r, 0)
			if err2 != nil {
				return nil, err2
			}
			fv, err3 := this.Decode(r)
			if err3 != nil {
				return nil, err3
			}
			rv[kv] = fv
		}
		return rv, nil
	case 8:
		l, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
		if l == 0 {
			return nil, nil
		}
		rv := make([]interface{}, l)
		for i := 0; i < int(l); i++ {
			fv, err2 := this.Decode(r)
			if err2 != nil {
				return nil, err2
			}
			rv[i] = fv
		}
		return rv, nil
	}
	return nil, nil
}

type NULL int

type allCoder struct {
	LenBytes  LenBytesCoder
	LenString LenStringCoder
	Bool      boolCoder
	Int32     int32Coder
	Int64     int64Coder
	Uint32    uint32Coder
	Uint64    uint64Coder
	FixUint8  fixUint8Coder
	FixUint16 fixUint16Coder
	FixUint32 fixUint32Coder
	FixUint64 fixUint64Coder
	Float32   float32Coder
	Float64   float64Coder
	Varinat   varCoder
	NullValue NULL
}

var (
	Coders allCoder
)
