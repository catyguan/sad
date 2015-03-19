package httpclient

import (
	sccore "bma/servicecall/core"
	"bma/servicecall/usecase"
	"fmt"
	"os"
	"testing"
	"time"
)

func runTest(t *testing.T, f func(m *sccore.Manager, ab sccore.AddressBuilder) error) {
	time.AfterFunc(20*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
	sccore.SetLogger(sccore.LoggerGo)

	m := sccore.NewManager("test")
	ab := sccore.NewAddressBuilder()
	ab.Type = "http"
	ab.API = "http://localhost:1080/$SNAME$/$MNAME$"
	err := f(m, ab.Build)
	if err != nil {
		t.Error(err)
	}
}

func T2estBaseHello(t *testing.T) {
	runTest(t, usecase.SCIHello)
}

func T2estBaseBinary(t *testing.T) {
	runTest(t, usecase.SCIBinary)
}

func T2estBaseAdd(t *testing.T) {
	runTest(t, usecase.SCIAdd)
}

func T2estBaseRedirect(t *testing.T) {
	runTest(t, usecase.SCIRedirect)
}

func T2estTransLogin(t *testing.T) {
	runTest(t, usecase.SCITrans)
}

func TestAsyncPoll(t *testing.T) {
	runTest(t, usecase.SCIAsyncPoll)
}

func T2estAsyncCallback(t *testing.T) {
	runTest(t, usecase.SCIAsyncCallback)
}

func T2estExportImport(t *testing.T) {
	runTest(t, usecase.SCIExportImport)
}
