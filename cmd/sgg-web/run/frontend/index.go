package frontend

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"save.gg/sgg/meta"
	m "save.gg/sgg/models"
)

func init() {
	meta.RegisterRoute("GET", "/", index)
}

// Index either a) shows the homepage or b) redirects to /dashboard
// It currently only does a.
func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
