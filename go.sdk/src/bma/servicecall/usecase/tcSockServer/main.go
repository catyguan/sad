package main

import (
	"bma/servicecall/core"
	_ "bma/servicecall/sockclient"
	"bma/servicecall/sockserver"
	"bma/servicecall/usecase"
	"flag"
	"fmt"
	"log"
	"net"
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

	mux := sockserver.NewServiceCallMux(manager.CreateClient)
	mux.SetServiceMethod("test", "ok", usecase.SM_OK)
	mux.SetServiceMethod("test", "echo", usecase.SM_Echo)
	mux.SetServiceMethod("test", "hello", usecase.SM_Hello)
	mux.SetServiceMethod("test", "add", usecase.SM_Add)
	mux.SetServiceMethod("test", "error", usecase.SM_Error)
	mux.SetServiceMethod("test", "redirect", usecase.SM_Redirect)
	mux.SetServiceMethod("test", "login", usecase.SM_Login)
	mux.SetServiceMethod("test", "async", usecase.SM_Async)

	fmt.Printf("start at %d\n", input_port)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", input_port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		core.DoLog("accept %s", conn.RemoteAddr())
		go func(c net.Conn) {
			mux.Run(conn)
		}(conn)
	}
}
