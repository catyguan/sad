package usecase

import (
	sccore "bma/servicecall/core"
	"fmt"
	"testing"
)

func T2estTransaction(t *testing.T) {
	initTest()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	cl.BeginTransaction()
	defer cl.EndTransaction()

	url := "http://localhost:1080/test/login"

	addr := sccore.CreateAddress("http", url, nil)
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
