package core

import "fmt"

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
}

func NewManager() *Manager {
	o := new(Manager)
	return o
}

func (this *Manager) CreateClient() *Client {
	o := new(Client)
	o.manager = this
	return o
}

func (this *Manager) GetDriver(addr *Address) (Driver, error) {
	df := GetDriverFactory(addr.GetType())
	if df == nil {
		return nil, fmt.Errorf("unknow driver(%s)", addr.GetType())
	}
	return df.GetDriver(addr.GetAPI(), addr.GetContext())
}
