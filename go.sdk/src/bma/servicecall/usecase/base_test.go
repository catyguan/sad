package usecase

import (
	sccore "bma/servicecall/core"
	"fmt"
	"testing"
)

func TestBase(t *testing.T) {
	manager := sccore.NewManager()
	cl := manager.CreateClient()
	defer cl.Close()

	addr := sccore.CreateAddress("http", "http://api.myhost.com/test/hello", nil)
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
		fmt.Println("RESULT ===", rs.Dump())
	} else {
		fmt.Println("not done")
	}
}
