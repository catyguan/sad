package httpserver

import (
	"net/http"
	"strings"
)

type ServiceDispatch func(r *http.Request) (string, string, error)

func DefaultServiceDispatch(r *http.Request) (string, string, error) {
	uri := r.RequestURI
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
