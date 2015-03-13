package usecase

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"testing"
	"time"
)

func T2estAsyncPoll(t *testing.T) {
	initTest()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	cl.BeginTransaction()
	defer cl.EndTransaction()

	url := "http://localhost:1080/test/async"

	addr := sccore.CreateAddress("http", url, nil)
	if true {
		req := sccore.NewRequest()
		ctx := sccore.NewContext()
		ctx.Put(constv.KEY_ASYNC_MODE, "poll")
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			t.Errorf("invoke fail - %s", err)
			return
		}
		sccore.DoLog("Async answer = %s", answer.Dump())

		if !answer.IsAsync() {
			t.Errorf("must answer async")
			return
		}

		answer2, done, err2 := cl.PollAnswer(addr, answer, ctx, time.Now().Add(10*time.Second), 1000*time.Millisecond)
		if err2 != nil {
			t.Errorf("poll fail - %s", err2)
			return
		}

		if !done {
			t.Errorf("poll timeout")
			return
		}

		if !answer2.IsDone() {
			sccore.DoLog("Answer fail - %d", answer.GetStatus())
			return
		}

		rs := answer2.GetResult()
		if rs != nil {
			sccore.DoLog("Result = %v", rs.Dump())
		}
	}
}

func T2estAsyncCallback(t *testing.T) {
	initTest()

	manager := sccore.NewManager("test")
	cl := manager.CreateClient()
	defer cl.Close()

	cl.BeginTransaction()
	defer cl.EndTransaction()

	url := "http://localhost:1080/test/async"
	addr := sccore.CreateAddress("http", url, nil)

	cburl := "http://localhost:1080/test/ok"
	cbaddr := sccore.CreateAddress("http", cburl, nil)

	if true {
		req := sccore.NewRequest()
		ctx := sccore.NewContext()
		ctx.Put(constv.KEY_ASYNC_MODE, "callback")
		ctx.Put(constv.KEY_CALLBACK, cbaddr.ToValueMap())
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			t.Errorf("invoke fail - %s", err)
			return
		}
		sccore.DoLog("Async answer = %s", answer.Dump())

		if !answer.IsAsync() {
			t.Errorf("must answer async")
			return
		}
		sccore.DoLog("end, check callback")
	}
}
