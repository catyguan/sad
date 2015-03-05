package usecase

import (
	sccore "bma/servicecall/core"
	_ "bma/servicecall/httpclient"
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func TestBase(t *testing.T) {
	safeCall()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	// url := "http://api.myhost.com/test/hello"
	// url := "http://cn.bing.com/"
	url := "http://localhost:1080/test/hello"

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
