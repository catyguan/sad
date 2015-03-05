package main

import (
	"bma/servicecall/core"
	"bma/servicecall/httpserver"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := httpserver.NewServiceCallMux()
	mux.SetServiceMethod("test", "hello", testHello)
	http.HandleFunc("/", mux.ServeHTTP)
	log.Fatal(http.ListenAndServe(":1080", nil))
}

func testHello(peer core.ServicePeer, req *core.Request, ctx *core.Context) error {
	fmt.Println("Request : ", req.Dump())
	fmt.Println("Context : ", ctx.Dump())
	fmt.Println("Deadline: ", ctx.GetLong("Deadline"))
	peer.WriteAnswer(nil, fmt.Errorf("test error"))
	return nil
}
