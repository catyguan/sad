package main

import (
	"bma/servicecall/core"
	_ "bma/servicecall/httpclient"
	"bma/servicecall/httpserver"
	"bma/servicecall/usecase"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var input_port int

func Init() {
	flag.IntVar(&input_port, "p", 1080, "server port")
}

func main() {
	Init()
	flag.Parse()

	core.SetLogger(core.LoggerFmtPrint)

	manager := core.NewManager("")

	mux := httpserver.NewServiceCallMux(manager.CreateClient)
	mux.SetServiceMethod("test", "ok", usecase.SM_OK)
	mux.SetServiceMethod("test", "echo", usecase.SM_Echo)
	mux.SetServiceMethod("test", "hello", usecase.SM_Hello)
	mux.SetServiceMethod("test", "add", usecase.SM_Add)
	mux.SetServiceMethod("test", "error", usecase.SM_Error)
	mux.SetServiceMethod("test", "redirect", usecase.SM_Redirect)
	mux.SetServiceMethod("test", "login", usecase.SM_Login)
	mux.SetServiceMethod("test", "async", usecase.SM_Async)

	http.HandleFunc("/", mux.ServeHTTP)
	fmt.Printf("start at %d\n", input_port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", input_port), nil))
}
