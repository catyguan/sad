package core

type Client struct {
	id      uint32
	manager *Manager
}

func (this *Client) Invoke(addr *Address, req *Request, ctx *Context) (*Answer, error) {
	dr, err := this.manager.GetDriver(addr)
	if err != nil {
		return nil, err
	}
	a, err2 := dr.Invoke(this, addr, req, ctx)
	if err2 != nil {
		return nil, err2
	}
	return a, nil
}

func (this *Client) Close() {

}
