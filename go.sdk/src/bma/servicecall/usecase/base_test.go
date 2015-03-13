package usecase

import (
	sccore "bma/servicecall/core"
	_ "bma/servicecall/httpclient"
	_ "bma/servicecall/sockclient"
	"bma/servicecall/sockcore"
	"fmt"
	"os"
	"testing"
	"time"
)

func initTest() {
	time.AfterFunc(20*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
	sccore.SetLogger(sccore.LoggerGo)
}

func maddr(s, m string) *sccore.Address {
	// typ := "http"
	// api := fmt.Sprintf("http://localhost:1080/%s/%s", s, m)
	typ := "socket"
	api := fmt.Sprintf("tcp:localhost:1080:%s:%s", s, m)

	return sccore.CreateAddress(typ, api, nil)
}

func T2estBase(t *testing.T) {
	initTest()

	pool := sockcore.SocketPool()
	pool.InitPoolSize(3)
	pool.Start()
	defer pool.Close()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	addr := maddr("test", "hello")
	req := sccore.NewRequest()
	req.Put("word", "Kitty")
	ctx := sccore.NewContext()

	answer, err := cl.Invoke(addr, req, ctx)
	if err != nil {
		t.Errorf("invoke fail - %s", err)
		return
	}
	fmt.Println(answer.Dump())

	if answer.IsDone() {
		rs := answer.GetResult()
		if rs != nil {
			fmt.Println("RESULT ===", rs.Dump())
		} else {
			fmt.Println("RESULT NULL")
		}
	} else {
		fmt.Println("not done")
	}
}

func T2estBinary(t *testing.T) {
	initTest()

	pool := sockcore.SocketPool()
	pool.InitPoolSize(3)
	pool.Start()
	defer pool.Close()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	addr := maddr("test", "echo")
	req := sccore.NewRequest()
	req.Put("binary", []byte("Kitty"))
	ctx := sccore.NewContext()

	answer, err := cl.Invoke(addr, req, ctx)
	if err != nil {
		t.Errorf("invoke fail - %s", err)
		return
	}
	fmt.Println(answer.Dump())

	if answer.IsDone() {
		rs := answer.GetResult()
		if rs != nil {
			dat := rs.GetMap("Data")
			fmt.Println("RESULT ===", dat.GetBinary("binary"))
		} else {
			fmt.Println("RESULT NULL")
		}
	} else {
		fmt.Println("not done")
	}
}

func T2estAdd(t *testing.T) {
	initTest()

	pool := sockcore.SocketPool()
	pool.InitPoolSize(3)
	pool.Start()
	defer pool.Close()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	c := int32(0)
	if true {
		addr := maddr("test", "add")
		req := sccore.NewRequest()
		req.Put("a", 1)
		req.Put("b", 2)
		ctx := sccore.NewContext()

		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			t.Errorf("invoke fail - %s", err)
			return
		}
		fmt.Println(answer.Dump())

		if answer.IsDone() {
			rs := answer.SureResult()
			c = rs.GetInt("Data")
		} else {
			fmt.Println("not done")
			return
		}
	}

	if true {
		addr := maddr("test", "add")
		req := sccore.NewRequest()
		req.Put("a", c)
		req.Put("b", 3)
		ctx := sccore.NewContext()

		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			t.Errorf("invoke fail - %s", err)
			return
		}
		fmt.Println(answer.Dump())
	}
}
