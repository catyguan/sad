package usecase

import (
	sccore "bma/servicecall/core"
	"bma/servicecall/sockcore"
	"fmt"
	"testing"
)

func TestTransaction(t *testing.T) {
	initTest()

	pool := sockcore.SocketPool()
	pool.InitPoolSize(3)
	pool.Start()
	defer pool.Close()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	cl.BeginTransaction()
	defer cl.EndTransaction()

	addr := maddr("test", "login")

	key := ""
	if true {
		req := sccore.NewRequest()
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
				key = rs.GetString("Key")
				fmt.Println("get login key -> ", key)
			}
		} else {
			fmt.Println("Invoke fail", answer.GetStatus())
			return
		}
	}

	if true {
		req := sccore.NewRequest()
		req.Put("User", "test")
		req.Put("Password", key)
		ctx := sccore.NewContext()
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			t.Errorf("invoke fail - %s", err)
			return
		}
		fmt.Println(answer.Dump())
	}
}
