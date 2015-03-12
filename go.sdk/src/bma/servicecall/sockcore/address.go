package sockcore

import (
	"fmt"
	"strconv"
	"strings"
)

type SocketAPI struct {
	Type    string
	Host    string
	Port    int
	Service string
	Method  string
}

func NewSocketAPI() *SocketAPI {
	o := new(SocketAPI)
	o.Type = "tcp"
	return o
}

func (this *SocketAPI) Valid() error {
	if this.Type == "" {
		return fmt.Errorf("Type empty")
	}
	if this.Port < 0 {
		return fmt.Errorf("Port(%s) invalid", this.Port)
	}
	if this.Service == "" {
		return fmt.Errorf("Service empty")
	}
	if this.Method == "" {
		return fmt.Errorf("Method empty")
	}
	return nil
}

func ParseSocketAPI(s string) (*SocketAPI, error) {
	ps := strings.SplitN(s, ":", 5)
	if len(ps) != 5 {
		return nil, fmt.Errorf("invalid SocketAPI - %s", s)
	}
	o := NewSocketAPI()
	o.Type = ps[0]
	o.Host = ps[1]
	if ps[2] != "" {
		pv, err1 := strconv.ParseInt(ps[2], 10, 32)
		if err1 != nil {
			return nil, err1
		}
		o.Port = int(pv)
	}
	o.Service = ps[3]
	o.Method = ps[4]
	err2 := o.Valid()
	if err2 != nil {
		return nil, err2
	}
	return o, nil
}

func (this *SocketAPI) String() string {
	return fmt.Sprintf("%s:%s:%d:%s:%s", this.Type, this.Host, this.Port, this.Service, this.Method)
}
