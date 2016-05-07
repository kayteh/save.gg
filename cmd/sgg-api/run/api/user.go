package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"save.gg/sgg/meta"
	m "save.gg/sgg/models"
	"save.gg/sgg/util/errors"
	util "save.gg/sgg/util/httputil"
)

func init() {
	meta.RegisterRoute("GET", "/api/user/:slug", getUser)
}

// GET /api/user/~
func getUserSelf(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := m.SessionFromRequest(r)
	if err == errors.SessionNotFound {
		util.NotFound(w)
		return
	}

	if err != nil {
		meta.App.Log.WithError(err).WithField("uri", r.RequestURI).Error("unknown error")
		util.InternalServerError(w, err)
		return
	}

	err = s.AttachUser()
	if err != nil {
		meta.App.Log.WithError(err).WithField("uri", r.RequestURI).Error("unknown error getting user")
		util.InternalServerError(w, err)
		return
	}

	util.Output(w, s.User.Presentable())

}

// GET /api/user/:slug
func getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	slug := ps.ByName("slug")

	if slug == "~" {
		getUserSelf(w, r, ps)
		return
	}

	d, err := m.UserBySlug(slug)

	if err == errors.UserNotFound {
		util.NotFound(w)
		return
	}

	if err != nil {
		meta.App.Log.WithError(err).WithField("uri", r.RequestURI).Error("unknown error")
		util.InternalServerError(w, err)
		return
	}

	util.Output(w, d.Presentable())

}
