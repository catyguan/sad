package usecase

import (
	"bma/servicecall/constv"
	sccore "bma/servicecall/core"
	"fmt"
	"time"
)

func SCIAsyncPoll(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	addr := ab("test", "async")
	if true {
		req := sccore.NewRequest()
		ctx := sccore.NewContext()
		ctx.Put(constv.KEY_ASYNC_MODE, "poll")
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			return fmt.Errorf("invoke fail - %s", err)
		}
		sccore.DoLog("Async answer = %s", answer.Dump())

		if !answer.IsAsync() {
			return fmt.Errorf("must answer async")
		}

		answer2, done, err2 := cl.PollAnswer(addr, answer, ctx, time.Now().Add(10*time.Second), 1000*time.Millisecond)
		if err2 != nil {
			return fmt.Errorf("poll fail - %s", err2)
		}

		if !done {
			return fmt.Errorf("poll timeout")
		}

		if !answer2.IsDone() {
			sccore.DoLog("Answer fail - %d", answer.GetStatus())
			return nil
		}

		rs := answer2.GetResult()
		if rs != nil {
			sccore.DoLog("Result = %v", rs.Dump())
		}
	}
	return nil
}

func SCIAsyncCallback(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	addr := ab("test", "async")
	cbaddr := ab("test", "ok")

	if true {
		req := sccore.NewRequest()
		ctx := sccore.NewContext()
		ctx.Put(constv.KEY_ASYNC_MODE, "callback")
		ctx.Put(constv.KEY_CALLBACK, cbaddr.ToValueMap())
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			return fmt.Errorf("invoke fail - %s", err)
		}
		sccore.DoLog("Async answer = %s", answer.Dump())

		if !answer.IsAsync() {
			return fmt.Errorf("must answer async")
		}
		sccore.DoLog("end, check callback")
	}
	return nil
}

func SCIAsyncPush(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	addr := ab("test", "async")
	if true {
		req := sccore.NewRequest()
		ctx := sccore.NewContext()
		ctx.Put(constv.KEY_ASYNC_MODE, "push")
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			return fmt.Errorf("invoke fail - %s", err)
		}
		sccore.DoLog("Async answer = %s", answer.Dump())

		if !answer.IsAsync() {
			return fmt.Errorf("must answer async")
		}

		answer2, err2 := cl.WaitAnswer(addr, 5500*time.Millisecond)
		if err2 != nil {
			return fmt.Errorf("waitAnswer fail - %s", err2)
		}

		if !answer2.IsDone() {
			sccore.DoLog("Answer fail - %d", answer.GetStatus())
			return nil
		}

		rs := answer2.GetResult()
		if rs != nil {
			sccore.DoLog("Result = %v", rs.Dump())
		}
	}
	return nil
}
