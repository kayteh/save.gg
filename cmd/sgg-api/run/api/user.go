package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"save.gg/sgg/meta"
	m "save.gg/sgg/models"
	util "save.gg/sgg/util/httputil"
)

func init() {
	meta.RegisterRoute("GET", "/api/user/:slug/comments", getUserComments)
	meta.RegisterRoute("GET", "/api/user/:slug/saves", getUserSaves)
	meta.RegisterRoute("GET", "/api/user/:slug", getUser)
}

// GET /api/user/:slug
func getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	slug := ps.ByName("slug")

	d, _ := m.UserBySlug(slug)

	if d == nil {
		util.NotFound(w)
	} else {
		util.Output(w, d)
	}

}

// GET /api/user/:slug/comments
func getUserComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	slug := ps.ByName("slug")

	q := r.URL.Query()
	fc := m.FetchConfigFromQuery(q)

	u, _ := m.UserBySlug(slug)

	d, _ := u.FetchSaves(fc)

	if d == nil {
		util.NotFound(w)
	} else {
		util.Output(w, d)
	}

}

// GET /api/user/:slug/saves
func getUserSaves(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	slug := ps.ByName("slug")

	q := r.URL.Query()
	fc := m.FetchConfigFromQuery(q)

	u, _ := m.UserBySlug(slug)

	d, _ := u.FetchSaves(fc)

	if d == nil {
		util.NotFound(w)
	} else {
		util.Output(w, d)
	}

}
