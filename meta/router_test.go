package meta

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"testing"
)

func TestSideeffectRouter(t *testing.T) {
	RegisterRoute("GET", "/test", testRoute)

	r := httprouter.New()

	MountRouter(r)

	h, _, _ := r.Lookup("GET", "/test")

	if h == nil {
		t.Log("/test didn't get mounted")
		t.Fail()
	}

	h2, _, _ := r.Lookup("GET", "/test2")

	if h2 != nil {
		t.Log("/test2 did get mounted")
		t.Fail()
	}

}

func testRoute(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	return
}
