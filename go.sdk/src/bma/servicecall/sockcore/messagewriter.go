package sockcore

import (
	sccore "bma/servicecall/core"
	"bytes"
	"io"
)

const (
	HEADER_SIZE = 4
)

var (
	endData  = []byte{0, 0, 0, 0}
	flagData = []byte{2, 0, 0, 0}
)

type MessageWriter struct {
	w    io.Writer
	ew   EncodeWriter
	hbuf []byte
	buf  *bytes.Buffer
}

func NewMessageWriter(w io.Writer) *MessageWriter {
	o := new(MessageWriter)
	o.w = w
	if ew, ok := w.(EncodeWriter); ok {
		o.ew = ew
	}
	o.hbuf = make([]byte, 4)
	return o
}

func (this *MessageWriter) Close() {
	if this.buf != nil {
		this.buf.Reset()
		this.buf = nil
	}
}

func (this *MessageWriter) Write(p []byte) (n int, err error) {
	return this.w.Write(p)
}

func (this *MessageWriter) WriteByte(c byte) error {
	if this.ew != nil {
		return this.ew.WriteByte(c)
	}
	bs := []byte{c}
	_, err := this.w.Write(bs)
	return err
}

func (this *MessageWriter) writeHeader(mt byte, sz int) error {
	this.hbuf[0] = mt
	this.hbuf[1] = byte(sz >> 16)
	this.hbuf[2] = byte(sz >> 8)
	this.hbuf[3] = byte(sz)
	_, err := this.w.Write(this.hbuf)
	return err
}

func (this *MessageWriter) writeEnd() error {
	_, err := this.w.Write(endData)
	return err
}

func (this *MessageWriter) WriteMessageId(mid int32) error {
	this.writeHeader(MT_MESSAGE_ID, 4)
	return Coders.FixUint32.DoEncode(this, uint32(mid))
}

func (this *MessageWriter) WriteFlag() error {
	_, err := this.w.Write(flagData)
	return err
}

func (this *MessageWriter) sbuf() *bytes.Buffer {
	if this.buf == nil {
		this.buf = bytes.NewBuffer([]byte{})
	}
	this.buf.Reset()
	return this.buf
}

func (this *MessageWriter) WriteAddress(s string, m string) error {
	buf := this.sbuf()
	err1 := Coders.LenString.DoEncode(buf, s)
	if err1 != nil {
		return err1
	}
	err2 := Coders.LenString.DoEncode(buf, m)
	if err2 != nil {
		return err2
	}
	err0 := this.writeHeader(MT_ADDRESS, buf.Len())
	if err0 != nil {
		return err0
	}
	_, errW := this.w.Write(buf.Bytes())
	if errW != nil {
		return errW
	}
	return nil
}

func (this *MessageWriter) WriteData(n string, val *sccore.Value) error {
	buf := this.sbuf()
	err1 := Coders.LenString.DoEncode(buf, n)
	if err1 != nil {
		return err1
	}
	var v interface{}
	if val != nil {
		v = val.ToValue()
	}
	err2 := Coders.Varinat.Encode(buf, v)
	if err2 != nil {
		return err2
	}
	err0 := this.writeHeader(MT_DATA, buf.Len())
	if err0 != nil {
		return err0
	}
	_, errW := this.w.Write(buf.Bytes())
	if errW != nil {
		return errW
	}
	return nil
}

func (this *MessageWriter) WriteContext(n string, val *sccore.Value) error {
	buf := this.sbuf()
	err1 := Coders.LenString.DoEncode(buf, n)
	if err1 != nil {
		return err1
	}
	var v interface{}
	if val != nil {
		v = val.ToValue()
	}
	err2 := Coders.Varinat.Encode(buf, v)
	if err2 != nil {
		return err2
	}
	err0 := this.writeHeader(MT_CONTEXT, buf.Len())
	if err0 != nil {
		return err0
	}
	_, errW := this.w.Write(buf.Bytes())
	if errW != nil {
		return errW
	}
	return nil
}

func (this *MessageWriter) WriteAnswer(st int, msg string) error {
	buf := this.sbuf()
	err1 := Coders.Int32.DoEncode(buf, int32(st))
	if err1 != nil {
		return err1
	}
	err2 := Coders.LenString.DoEncode(buf, msg)
	if err2 != nil {
		return err2
	}
	err0 := this.writeHeader(MT_CONTEXT, buf.Len())
	if err0 != nil {
		return err0
	}
	_, errW := this.w.Write(buf.Bytes())
	if errW != nil {
		return errW
	}
	return nil
}
