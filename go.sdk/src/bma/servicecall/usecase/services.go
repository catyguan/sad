package usecase

import (
	"bma/servicecall/constv"
	"bma/servicecall/core"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func dumpSM(req *core.Request, ctx *core.Context) {
	fmt.Println("Request : ", req.Dump())
	fmt.Println("Context : ", ctx.Dump())
	fmt.Println("Deadline: ", ctx.GetLong("Deadline"))
}

func SM_Echo(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	dumpSM(req, ctx)

	reqm := req.ToMap()
	a := core.NewAnswer()
	a.SureResult().Put("Data", reqm)
	peer.WriteAnswer(a, nil)
	return nil
}

func SM_OK(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	dumpSM(req, ctx)

	a := core.NewAnswer()
	a.SetStatus(constv.STATUS_OK)
	peer.WriteAnswer(a, nil)
	return nil
}

func SM_Hello(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	dumpSM(req, ctx)

	word := req.GetString("word")
	if word == "" {
		word = "<empty>"
	}
	fmt.Println("Hello ", word)

	r := "Hello " + word
	a := core.NewAnswer()
	a.SureResult().Put("Data", r)

	// peer.WriteAnswer(a, fmt.Errorf("test error"))
	peer.WriteAnswer(a, nil)
	return nil
}

func SM_Add(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	dumpSM(req, ctx)

	pa := req.GetInt("a")
	pb := req.GetInt("b")
	pc := pa + pb
	fmt.Printf("a + b = %d + %d = %d\n", pa, pb, pc)

	a := core.NewAnswer()
	a.SureResult().Put("Data", pc)

	peer.WriteAnswer(a, nil)
	return nil
}

func SM_Error(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	dumpSM(req, ctx)

	errorMsg := req.GetString("Error")
	if errorMsg == "" {
		errorMsg = "<test error>"
	}
	peer.WriteAnswer(nil, errors.New(errorMsg))
	return nil
}

func SM_Redirect(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	dumpSM(req, ctx)

	loc := req.GetMap("Location")
	if loc == nil {
		loc = core.NewValueMap(nil)
		loc.Put("Type", "http")
		loc.Put("API", "http://localhost:1080/test/hello")
	}
	a := core.NewAnswer()
	a.SetStatus(constv.STATUS_REDIRECT)
	a.SetResult(loc)
	peer.WriteAnswer(a, nil)
	return nil
}

func SM_Login(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	err0 := peer.BeginTransaction()
	if err0 != nil {
		return err0
	}

	rand.Seed(time.Now().UTC().UnixNano())
	key := ""
	if true {
		dumpSM(req, ctx)

		key = fmt.Sprintf("%d", 10000000+rand.Intn(99999999-10000000))

		a := core.NewAnswer()
		a.SetStatus(constv.STATUS_CONTINUE)
		rs := a.SureResult()
		rs.Put("Key", key)

		err := peer.WriteAnswer(a, nil)
		if err != nil {
			fmt.Println("WriteAnswer fail - %s", err)
			return nil
		}
	}

	if true {
		req2, ctx2, err := peer.ReadRequest(30 * time.Second)
		if err != nil {
			fmt.Println("ReadRequest fail - ", err)
			return nil
		}

		dumpSM(req2, ctx2)

		user := req2.GetString("User")
		pass := req2.GetString("Password")
		fmt.Println("param", user, pass)
		a := core.NewAnswer()
		rs := a.SureResult()
		if user == "test" && pass == key {
			rs.Put("Done", true)
			rs.Put("UID", 12345)
		} else {
			rs.Put("Done", false)
			rs.Put("Why", "user or pass invalid")
		}
		peer.WriteAnswer(a, nil)
	}

	return nil
}

func SM_Async(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	dumpSM(req, ctx)
	sleepTime := req.GetInt("sleep")
	if sleepTime <= 0 {
		sleepTime = 5
	}
	err := peer.SendAsync(ctx, nil, 1*time.Minute)
	if err != nil {
		core.DoLog("SendAsync fail - %s", err)
		return nil
	}
	go func() {
		time.Sleep(time.Duration(sleepTime) * time.Second)

		a := core.NewAnswer()
		a.SureResult().Put("Word", "Hello Kitty")
		peer.WriteAnswer(a, nil)
	}()
	return nil
}
