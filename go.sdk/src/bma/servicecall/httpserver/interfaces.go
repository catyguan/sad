package httpserver

import (
	sccore "bma/servicecall/core"
	"net/http"
	"strings"
)

type ServiceObject interface {
	GetMethod(name string) sccore.ServiceMethod
}

type ServiceDispatch func(r *http.Request) (string, string, error)

func DefaultServiceDispatch(r *http.Request) (string, string, error) {
	uri := r.URL.RequestURI()
	ps := strings.Split(uri, "/")
	l := len(ps)
	switch l {
	case 0:
		return "home", "index", nil
	case 1:
		return ps[0], "index", nil
	case 2:
		return ps[0], ps[1], nil
	}
	return ps[l-2], ps[l-1], nil
}
