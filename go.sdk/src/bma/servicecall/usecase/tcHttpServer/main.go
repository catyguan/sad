package main

import (
	"bma/servicecall/core"
	"bma/servicecall/httpserver"
	"bma/servicecall/usecase"
	"log"
	"net/http"
)

func main() {
	core.SetLogger(core.LoggerFmtPrint)

	mux := httpserver.NewServiceCallMux()
	mux.SetServiceMethod("test", "ok", usecase.SM_OK)
	mux.SetServiceMethod("test", "echo", usecase.SM_Echo)
	mux.SetServiceMethod("test", "hello", usecase.SM_Hello)
	mux.SetServiceMethod("test", "add", usecase.SM_Add)
	mux.SetServiceMethod("test", "error", usecase.SM_Error)
	mux.SetServiceMethod("test", "redirect", usecase.SM_Redirect)
	mux.SetServiceMethod("test", "login", usecase.SM_Login)
	mux.SetServiceMethod("test", "async", usecase.SM_Async)

	http.HandleFunc("/", mux.ServeHTTP)
	log.Fatal(http.ListenAndServe(":1080", nil))
}
