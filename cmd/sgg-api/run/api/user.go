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
	meta.RegisterRoute("PATCH", "/api/user/:slug", patchUser)
}

// GET /api/user/~
func getUserSelf(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := m.SessionFromRequest(r)
	if err == errors.SessionNotFound || err == errors.SessionTokenInvalid {
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

// PATCH /api/user/:slug
func patchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := m.SessionFromRequest(r)
	if err == errors.SessionNotFound || err == errors.SessionTokenInvalid {
		util.Forbidden(w)
		return
	}

	slug := ps.ByName("slug")

	var u *m.User

	if slug == "~" {
		err = s.AttachUser()
		u = s.User
	} else {
		u, err = m.UserBySlug(slug)
	}

	if err == errors.UserNotFound {
		util.NotFound(w)
		return
	}

	if err != nil {
		meta.App.Log.WithError(err).WithField("uri", r.RequestURI).Error("unknown error")
		util.InternalServerError(w, err)
		return
	}

	s.AttachUser()
	if !u.UserCanModify(s.User) {
		util.Forbidden(w)
		return
	}

	var i map[string]interface{}

	err = util.Input(r, &i)
	if err != nil {
		meta.App.Log.WithError(err).Error("json decode problem")
		util.BadInput(w)
		return
	}

	err = u.Patch(i)
	//TODO(kkz): handle specific problems here, e.g. specify what input is bad
	if err != nil {
		meta.App.Log.WithError(err).Error("user patch problem")
		util.InternalServerError(w, err)
	}

	util.NoContent(w)

}
