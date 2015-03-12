package sockcore

import (
	sccore "bma/servicecall/core"
	"fmt"
	"io"
)

type MessageReader struct {
	r  io.Reader
	dr DecodeReader
	h  bool
	l  int
	sz int
	mt byte
}

func NewMessageReader(r io.Reader) *MessageReader {
	o := new(MessageReader)
	o.r = r
	if dr, ok := r.(DecodeReader); ok {
		o.dr = dr
	}
	return o
}

func (this *MessageReader) Next() (byte, error) {
	this.h = true
	mt, size, err := this.readHeader()
	this.h = false
	if err != nil {
		return 0, err
	}
	this.mt = mt
	this.sz = size
	this.l = size
	return mt, nil
}

func (this *MessageReader) Len() int {
	return this.sz
}

func (this *MessageReader) Read(p []byte) (n int, err error) {
	if !this.h {
		if this.l <= 0 {
			return 0, io.EOF
		}
		if len(p) > this.l {
			p = p[0:this.l]
		}
	}
	n, err = this.r.Read(p)
	if this.h {
		this.l -= int(n)
	}
	return
}

func (this *MessageReader) ReadByte() (b byte, err error) {
	if !this.h {
		if this.l <= 0 {
			return 0, io.EOF
		}
	}
	if this.dr != nil {
		b, err = this.dr.ReadByte()
		if err != nil {
			return
		}
	} else {
		p := []byte{0}
		_, err = this.r.Read(p)
		if err != nil {
			return
		}
		b = p[0]
	}
	if this.h {
		this.l -= 1
	}
	return
}

func (this *MessageReader) readHeader() (byte, int, error) {
	bs := make([]byte, 4)
	_, err := io.ReadFull(this, bs)
	if err != nil {
		return 0, 0, err
	}
	b, i := this.decodeHeader(bs)
	return b, i, nil
}

func (this *MessageReader) decodeHeader(b []byte) (byte, int) {
	mt := byte(b[0])
	sz := int(b[3]) | int(b[2])<<8 | int(b[1])<<16
	return mt, sz
}

////////
func (this *MessageReader) ReadMessageId() (int32, error) {
	if this.mt != MT_MESSAGE_ID {
		return 0, fmt.Errorf("MT(%d) invalid MT_MESSAGE_ID", this.mt)
	}
	v, err := Coders.FixUint32.DoDecode(this)
	if err != nil {
		return 0, err
	}
	return int32(v), err
}

func (this *MessageReader) ReadAddress() (string, string, error) {
	if this.mt != MT_ADDRESS {
		return "", "", fmt.Errorf("MT(%d) invalid MT_ADDRESS", this.mt)
	}
	s, err1 := Coders.LenString.DoDecode(this, 0)
	if err1 != nil {
		return "", "", err1
	}
	m, err2 := Coders.LenString.DoDecode(this, 0)
	if err2 != nil {
		return "", "", err2
	}
	return s, m, nil
}

func (this *MessageReader) ReadData() (string, *sccore.Value, error) {
	if this.mt != MT_DATA {
		return "", nil, fmt.Errorf("MT(%d) invalid MT_DATA", this.mt)
	}
	s, err1 := Coders.LenString.DoDecode(this, 0)
	if err1 != nil {
		return "", nil, err1
	}
	v, err2 := Coders.Varinat.Decode(this)
	if err2 != nil {
		return "", nil, err2
	}
	val := sccore.CreateValue(v)
	return s, val, nil
}

func (this *MessageReader) ReadContext() (string, *sccore.Value, error) {
	if this.mt != MT_CONTEXT {
		return "", nil, fmt.Errorf("MT(%d) invalid MT_CONTEXT", this.mt)
	}
	s, err1 := Coders.LenString.DoDecode(this, 0)
	if err1 != nil {
		return "", nil, err1
	}
	v, err2 := Coders.Varinat.Decode(this)
	if err2 != nil {
		return "", nil, err2
	}
	val := sccore.CreateValue(v)
	return s, val, nil
}

func (this *MessageReader) ReadAnswer() (int32, string, error) {
	if this.mt != MT_ANSWER {
		return 0, "", fmt.Errorf("MT(%d) invalid MT_ANSWER", this.mt)
	}
	st, err1 := Coders.Int32.DoDecode(this)
	if err1 != nil {
		return 0, "", err1
	}
	msg, err2 := Coders.LenString.DoDecode(this, 0)
	if err2 != nil {
		return 0, "", err2
	}
	return st, msg, nil
}

func (this *MessageReader) NextRequest() (mid int32, s, m string, req *sccore.Request, ctx *sccore.Context, rErr error) {
	for {
		mt, err0 := this.Next()
		if err0 != nil {
			sccore.DoLog("conn read fail - %s", err0)
			rErr = err0
			return
		}
		switch mt {
		case MT_END:
			break
		case MT_MESSAGE_ID:
			v, err1 := this.ReadMessageId()
			if err1 != nil {
				sccore.DoLog("messageId read fail - %s", err1)
				rErr = err1
				return
			}
			mid = v
			sccore.DoLog("message id = %d", mid)
		case MT_REQUEST:
			continue
		case MT_ADDRESS:
			var err1 error
			s, m, err1 = this.ReadAddress()
			if err1 != nil {
				sccore.DoLog("messageId read fail - %s", err1)
				rErr = err1
				return
			}
		case MT_DATA:
			n, val, err1 := this.ReadData()
			if err1 != nil {
				sccore.DoLog("data read fail - %s", err1)
				rErr = err1
				return
			}
			if req == nil {
				req = sccore.NewRequest()
			}
			req.Set(n, val)
		case MT_CONTEXT:
			n, val, err1 := this.ReadContext()
			if err1 != nil {
				sccore.DoLog("context read fail - %s", err1)
				rErr = err1
				return
			}
			if ctx == nil {
				ctx = sccore.NewContext()
			}
			ctx.Set(n, val)
		default:
			err1 := fmt.Errorf("unknow MessageType(%d)", mt)
			sccore.DoLog("message read fail - %s", err1)
			rErr = err1
			return
		}
	}
	return
}
