package core

import (
	"fmt"
	"time"
)
import (
	"sync/atomic"
)

var (
	gDS map[string]DriverFactory
)

func GetDriverFactory(typ string) DriverFactory {
	if gDS == nil {
		return nil
	}
	df := gDS[typ]
	if df == nil {
		return nil
	}
	return df
}

func InitDriverFactory(typ string, df DriverFactory) {
	if gDS == nil {
		gDS = make(map[string]DriverFactory)
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
	o := new(Client)
	o.manager = this
	o.id = atomic.AddUint32(&this.clientSeq, 1)
	return o
}

func (this *Manager) GetDriver(addr *Address) (Driver, error) {
	df := GetDriverFactory(addr.GetType())
	if df == nil {
		return nil, fmt.Errorf("unknow driver(%s)", addr.GetType())
	}
	return df.GetDriver(addr.GetAPI())
}
