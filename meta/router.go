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
//
// The best way to utilize this is to call this in a package's `init()` function,
// and importing that in the web server init as a non-capturing import, such as
// `_ "path/to/routes"`
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
