package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"save.gg/sgg/meta"
	m "save.gg/sgg/models"
	u "save.gg/sgg/util/httputil"
)

func init() {
	meta.RegisterRoute("GET", "/api/user/:slug", getUser)
	meta.RegisterRoute("GET", "/api/user/:slug/comments", getUserComments)
	meta.RegisterRoute("GET", "/api/user/:slug/saves", getUserSaves)
}

// GET /api/user/:slug
func getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	slug := ps.ByName("slug")

	d := m.UserBySlug(slug)

	if d == nil {
		u.NotFound(w)
	} else {
		u.Output(w, d)
	}

}

// GET /api/user/:slug/comments
func getUserComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	slug := ps.ByName("slug")

	q := r.URL.Query()
	fc := m.FetchConfigFromQuery(q)

	d := m.UserBySlug(slug).FetchComments(fc)

	if d == nil {
		u.NotFound(w)
	} else {
		u.Output(w, d)
	}

}

// GET /api/user/:slug/saves
func getUserSaves(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	slug := ps.ByName("slug")

	q := r.URL.Query()
	fc := m.FetchConfigFromQuery(q)

	d := m.UserBySlug(slug).FetchSaves(fc)

	if d == nil {
		u.NotFound(w)
	} else {
		u.Output(w, d)
	}

}
