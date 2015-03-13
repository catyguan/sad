package sockcore

import (
	sccore "bma/servicecall/core"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// poolConfig
type poolConfig struct {
	Net       string
	Address   string
	TimeoutMS int
	PoolSize  int
	IdleMS    int
}

func (this *poolConfig) Valid() error {
	if this.Address == "" {
		return errors.New("address empty")
	}
	if this.Net == "" {
		this.Net = "tcp"
	}
	if this.PoolSize < 0 {
		this.PoolSize = 0
	}
	if this.IdleMS <= 0 {
		this.IdleMS = 60 * 1000
	}
	return nil
}

// dialPool
type connItem struct {
	conn     net.Conn
	waitTime time.Time
	pingTime time.Time
}
type dialPool struct {
	config poolConfig
	wait   chan *connItem
	closed chan bool
}

func newDialPool(cfg *poolConfig) *dialPool {
	err := cfg.Valid()
	if err != nil {
		panic(err)
	}

	this := new(dialPool)
	this.config = *cfg
	this.closed = make(chan bool, 1)

	this.wait = make(chan *connItem, this.config.PoolSize)

	return this
}

func (this *dialPool) String() string {
	return fmt.Sprintf("dialPool[%s, %d/%d]", this.config.Address, len(this.wait), this.config.PoolSize)
}

func (this *dialPool) CheckSocket(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(1))
	one := make([]byte, 1)
	n, err := conn.Read(one)
	if n > 0 {
		conn.Close()
		return false
	}
	conn.SetReadDeadline(time.Time{})
	if err == nil {
		return true
	}
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return true
	}
	return false
}

func (this *dialPool) GetSocket(timeout time.Duration) (net.Conn, error) {
	conn := func() net.Conn {
		for {
			select {
			case o := <-this.wait:
				if o == nil {
					return nil
				}
				if this.CheckSocket(o.conn) {
					return o.conn
				}
				o.conn.Close()
			default:
				// fmt.Println("NewConnect")
				return nil
			}
		}
	}()
	if conn != nil {
		return conn, nil
	}
	conn, err := this.doDial(timeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (this *dialPool) ReturnSocket(conn net.Conn) {
	now := time.Now()
	o := new(connItem)
	o.conn = conn
	o.waitTime = now
	o.pingTime = now
	// fmt.Println("ReturnConnect")
	if !this.put(o) {
		// fmt.Println("ReturnFail")
	}
}

func (this *dialPool) CloseSocket(conn net.Conn) {
	conn.Close()
}

func (this *dialPool) doDial(timeout time.Duration) (net.Conn, error) {
	var conn net.Conn
	var err error
	cfg := &this.config
	if cfg.TimeoutMS > 0 {
		timeout = time.Duration(cfg.TimeoutMS) * time.Millisecond
	}
	if timeout == 0 {
		timeout = 5 * time.Millisecond
	}
	conn, err = net.DialTimeout(cfg.Net, cfg.Address, timeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (this *dialPool) close() {
	done := false
	for {
		if done {
			break
		}
		select {
		case o := <-this.wait:
			o.conn.Close()
		default:
			done = true
		}
	}
	close(this.wait)
	close(this.closed)
}

func (this *dialPool) put(item *connItem) bool {
	defer func() {
		if recover() != nil {
			item.conn.Close()
		}
	}()
	select {
	case this.wait <- item: // Put 2 in the channel unless it is full
		return true
	default:
		item.conn.Close()
		return false
	}
}

func (this *dialPool) idlePing() {
	go func() {
		timer := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-this.closed:
				return
			case <-timer.C:
			}

			// check idle
			idleDu := time.Duration(this.config.IdleMS) * time.Millisecond
			pingDu := 15 * time.Second
			l := len(this.wait)
			// fmt.Println("do idleping", idleDu, l)
			for i := 0; i < l; i++ {
				now := time.Now()
				done := false
				select {
				case item := <-this.wait:
					if now.Sub(item.waitTime) > idleDu {
						// close
						sccore.DoLog("'%s' idle break", item.conn.RemoteAddr())
						item.conn.Close()
					} else {
						if now.Sub(item.pingTime) > pingDu {
							// ping
							if this.doPing(item.conn) {
								item.pingTime = now
								this.put(item)
							}
						}
					}
				default:
					done = true
				}
				if done {
					break
				}
			}
		}
	}()
}

func (this *dialPool) doPing(conn net.Conn) bool {
	// fmt.Println("do ping")
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer func() {
		conn.SetDeadline(time.Time{})
	}()
	_, err := conn.Write(pingData)
	if err != nil {
		sccore.DoLog("'%s' ping write fail %s", conn.RemoteAddr(), err)
		conn.Close()
		return false
	}
	bs := make([]byte, 9)
	_, err2 := io.ReadFull(conn, bs)
	if err2 != nil {
		sccore.DoLog("'%s' ping read fail %s", conn.RemoteAddr(), err2)
		conn.Close()
		return false
	}
	for i, v := range bs {
		if v != pingRData[i] {
			sccore.DoLog("'%s' ping invalid response %v", conn.RemoteAddr(), bs)
			conn.Close()
			return false
		}
	}

	return true
}

// Pool
type Pool struct {
	lock   sync.RWMutex
	dials  map[string]*dialPool
	config poolConfig
}

func (this *Pool) InitPoolSize(ps int) {
	this.config.PoolSize = ps
}

func (this *Pool) InitPoolTimeoutMS(ms int) {
	this.config.TimeoutMS = ms
}

func (this *Pool) InitPoolIdleTimeMS(ms int) {
	this.config.IdleMS = ms
}

func (this *Pool) Start() {
}

func (this *Pool) Close() {
	tmpm := make([]*dialPool, 0)
	this.lock.Lock()
	for k, dial := range this.dials {
		tmpm = append(tmpm, dial)
		delete(this.dials, k)
	}
	this.lock.Unlock()

	for _, dial := range tmpm {
		dial.close()
	}
}

func (this *Pool) GetSocket(addr *sccore.Address, api *SocketAPI, timeout time.Duration) (string, net.Conn, error) {
	if api == nil {
		var err0 error
		api, err0 = ParseSocketAPI(addr.GetAPI())
		if err0 != nil {
			return "", nil, err0
		}
	}
	err := api.Valid()
	if err != nil {
		return "", nil, err
	}

	var dial *dialPool
	key := api.Key()
	this.lock.RLock()
	if this.dials != nil {
		dial = this.dials[key]
	}
	this.lock.RUnlock()
	if dial == nil {
		host := api.Host
		if api.Port != 0 {
			host = fmt.Sprintf("%s:%d", api.Host, api.Port)
		}
		// fmt.Printf("host => %v\n", host)

		cfg := new(poolConfig)
		*cfg = this.config
		cfg.Net = api.Type
		cfg.Address = host
		opt := addr.GetOption()
		if opt != nil {
			if opt.Has("PoolSize") {
				cfg.PoolSize = int(opt.GetInt("PoolSize"))
			}
			if opt.Has("Timeout") {
				cfg.TimeoutMS = int(opt.GetInt("Timeout"))
			}
			if opt.Has("Idle") {
				cfg.IdleMS = int(opt.GetInt("Idle"))
			}
		}

		this.lock.Lock()
		if this.dials == nil {
			this.dials = make(map[string]*dialPool)
		}
		dial = this.dials[key]
		if dial == nil {
			dial = newDialPool(cfg)
			this.dials[key] = dial
			dial.idlePing()
		}
		this.lock.Unlock()
	}
	if timeout <= 0 {
		timeout = time.Duration(this.config.TimeoutMS) * time.Millisecond
	}
	conn, err2 := dial.GetSocket(timeout)
	return key, conn, err2
}

func (this *Pool) ReturnSocket(key string, conn net.Conn) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.dials != nil {
		dial := this.dials[key]
		if dial != nil {
			dial.ReturnSocket(conn)
		}
	}
}

func (this *Pool) CloseSocket(key string, conn net.Conn) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.dials != nil {
		dial := this.dials[key]
		if dial != nil {
			dial.CloseSocket(conn)
			return
		}
	}
	conn.Close()
}

var (
	pingData  = []byte{9, 0, 0, 1, 0, 0, 0, 0, 0}
	pingRData = []byte{9, 0, 0, 1, 1, 0, 0, 0, 0}
)

var (
	gPool Pool
)

func SocketPool() *Pool {
	return &gPool
}
