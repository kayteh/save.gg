package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"save.gg/sgg/meta"
	m "save.gg/sgg/models"
	u "save.gg/sgg/util/httputil"
)

func init() {
	meta.RegisterRoute("GET", "/api/saves/:id", getSave)
}

func getSave(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")

	d := m.SaveById(id)

	if d == nil {
		u.NotFound(w)
	} else {
		u.Output(w, d)
	}

}

func getSaveComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")

	q := r.URL.Query()
	fc := m.FetchConfigFromQuery(q)

	d := m.SaveById(id).FetchComments(fc)

	if d == nil {
		u.NotFound(w)
	} else {
		u.Output(w, d)
	}

}
