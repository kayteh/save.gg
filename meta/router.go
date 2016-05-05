package meta

import (
	"github.com/julienschmidt/httprouter"
)

type route struct {
	Method   string
	Path     string
	Callback httprouter.Handle
}

var rs []route

// Registers a route, intended to be used as a side effect.
func RegisterRoute(method, path string, fn httprouter.Handle) {
	rs = append(rs, route{
		Method:   method,
		Path:     path,
		Callback: fn,
	})
}

// Mount the routes to a router.
func MountRouter(r *httprouter.Router) {
	for _, h := range rs {
		r.Handle(h.Method, h.Path, h.Callback)
	}
}
