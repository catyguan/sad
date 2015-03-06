package usecase

import (
	sccore "bma/servicecall/core"
	"fmt"
	"testing"
)

func T2estRedirect(t *testing.T) {
	initTest()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	url := "http://localhost:1080/test/redirect"

	addr := sccore.CreateAddress("http", url, nil)
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
