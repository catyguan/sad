package core

import (
	"fmt"
	"time"
)
import (
	"sync/atomic"
)

var (
	gDS map[string]Driver
)

func GetDriver(typ string) Driver {
	if gDS == nil {
		return nil
	}
	df := gDS[typ]
	if df == nil {
		return nil
	}
	return df
}

func InitDriver(typ string, df Driver) {
	if gDS == nil {
		gDS = make(map[string]Driver)
	}
	gDS[typ] = df
}

type Manager struct {
	name      string
	clientSeq uint32
}

func NewManager(n string) *Manager {
	if n == "" {
		n = fmt.Sprintf("goscm%d", time.Now().UnixNano())
	}
	o := new(Manager)
	o.name = n
	return o
}

func (this *Manager) CreateClient() *Client {
	return newClient(this, atomic.AddUint32(&this.clientSeq, 1))
}

func (this *Manager) createConn(typ, api string) (ServiceConn, error) {
	df := GetDriver(typ)
	if df == nil {
		return nil, fmt.Errorf("unknow driver(%s)", typ)
	}
	return df.CreateConn(typ, api)
}
