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
	// fmt.Printf("READ -> %v\n", p[0:n])
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
	// fmt.Printf("READ -> %v\n", b)
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

type Message struct {
	Type byte
	Id   int32
	// Ping,Trasaction
	BoolFlag bool
	// Request
	Service string
	Method  string
	Request *sccore.Request
	Context *sccore.Context
	// Answer
	Answer *sccore.Answer
}

func (this *Message) Reset() {
	this.Type = 0
	this.Id = 0
	this.BoolFlag = false
	this.Service = ""
	this.Method = ""
	this.Request = nil
	this.Context = nil
	this.Answer = nil
}

func (this *MessageReader) NextMessage(msg *Message) (byte, error) {
	msg.Reset()
	for {
		mt, err0 := this.Next()
		if err0 != nil {
			sccore.DoLog("conn read fail - %s", err0)
			return 0, err0
		}
		// sccore.DoLog("read line - %d", mt)
		switch mt {
		case MT_END:
			if msg.Type == MT_ANSWER {
				an := msg.Answer
				if msg.Request != nil {
					an.SetResult(&msg.Request.ValueMap)
					msg.Request = nil
				}
				if msg.Context != nil {
					an.SetContext(&msg.Context.ValueMap)
					msg.Context = nil
				}
			}
			sccore.DoLog("read message -> %d, %v", msg.Type, msg)
			return msg.Type, nil
		case MT_MESSAGE_ID:
			v, err1 := this.ReadMessageId()
			if err1 != nil {
				sccore.DoLog("messageId read fail - %s", err1)
				return 0, err1
			}
			msg.Id = v
		case MT_PING, MT_TRANSACTION:
			v, err1 := Coders.Bool.DoDecode(this)
			if err1 != nil {
				sccore.DoLog("boolFlag read fail - %s", err1)
				return 0, err1
			}
			msg.BoolFlag = v
			msg.Type = mt
		case MT_REQUEST:
			msg.Type = mt
			continue
		case MT_ADDRESS:
			s, m, err1 := this.ReadAddress()
			if err1 != nil {
				sccore.DoLog("messageId read fail - %s", err1)
				return 0, err1
			}
			msg.Service = s
			msg.Method = m
		case MT_DATA:
			n, val, err1 := this.ReadData()
			if err1 != nil {
				sccore.DoLog("data read fail - %s", err1)
				return 0, err1
			}
			if msg.Request == nil {
				msg.Request = sccore.NewRequest()
			}
			msg.Request.Set(n, val)
		case MT_CONTEXT:
			n, val, err1 := this.ReadContext()
			if err1 != nil {
				sccore.DoLog("context read fail - %s", err1)
				return 0, err1
			}
			if msg.Context == nil {
				msg.Context = sccore.NewContext()
			}
			msg.Context.Set(n, val)
		case MT_ANSWER:
			st, s, err1 := this.ReadAnswer()
			if err1 != nil {
				sccore.DoLog("answer read fail - %s", err1)
				return 0, err1
			}
			msg.Type = mt
			if msg.Answer == nil {
				msg.Answer = sccore.NewAnswer()
			}
			msg.Answer.SetStatus(int(st))
			msg.Answer.SetMessage(s)
		default:
			err1 := fmt.Errorf("unknow MessageType(%d)", mt)
			sccore.DoLog("message read fail - %s", err1)
			return 0, err1
		}
	}

}
