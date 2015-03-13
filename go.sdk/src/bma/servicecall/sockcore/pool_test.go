package sockcore

import (
	sccore "bma/servicecall/core"
	"fmt"
	"testing"
	"time"
)

func T2estDialPool(t *testing.T) {

	pool := SocketPool()
	pool.InitPoolSize(10)
	pool.Start()
	defer func() {
		pool.Close()
		time.Sleep(100 * time.Millisecond)
	}()

	time.Sleep(time.Duration(500) * time.Millisecond)

	addr := sccore.NewAddress()
	addr.SetAPI("tcp:127.0.0.1:1080:test:hello")
	f := func() {
		key, conn, err := pool.GetSocket(addr, nil, 0)
		if err != nil {
			fmt.Printf("GetSocket fail - %s\n", err)
		}
		if conn != nil {
			fmt.Printf("GetSocket -> %p\n", conn)
			time.AfterFunc(1*time.Second, func() {
				pool.ReturnSocket(key, conn)
			})
		}
	}
	for i := 0; i < 4; i++ {
		go f()
	}

	time.Sleep(2000 * time.Millisecond)
	fmt.Println("-------- next turn --------")

	for i := 0; i < 5; i++ {
		go f()
	}

	time.Sleep(time.Duration(5) * time.Second)
}

func TestPoolIdlePing(t *testing.T) {

	sccore.SetLogger(sccore.LoggerFmtPrint)

	pool := SocketPool()
	pool.InitPoolSize(10)
	// pool.InitPoolIdleTimeMS(100)
	pool.Start()
	defer func() {
		pool.Close()
		time.Sleep(100 * time.Millisecond)
	}()

	time.Sleep(time.Duration(500) * time.Millisecond)

	addr := sccore.NewAddress()
	addr.SetAPI("tcp:127.0.0.1:1080:test:hello")
	key, conn, err := pool.GetSocket(addr, nil, 0)
	if err != nil {
		fmt.Printf("GetSocket fail - %s\n", err)
	}
	if conn != nil {
		fmt.Printf("GetSocket -> %p\n", conn)
		pool.ReturnSocket(key, conn)
	}

	time.Sleep(time.Duration(6) * time.Second)
}
