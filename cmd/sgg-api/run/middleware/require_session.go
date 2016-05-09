package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	//"save.gg/sgg/meta"
	m "save.gg/sgg/models"
	"save.gg/sgg/util/errors"
	util "save.gg/sgg/util/httputil"
)

// httprouter.Handle with a *models.Session argument as well.
type SessionHandler func(http.ResponseWriter, *http.Request, httprouter.Params, *m.Session)

// Validates and attaches a session to a request. The underlying handler will never be called
// if this fails to validate. It attaches the session to prevent a double-query to find the same session.
func RequireSession(fn SessionHandler, sec *SecurityFlags) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		s, err := m.SessionFromRequest(r)

		if err == errors.SessionTokenInvalid || err == errors.SessionNotFound || s == nil {
			util.Forbidden(w)
			return
		}

		if err != nil {
			util.InternalServerError(w, err)
			return
		}

		if sec != nil {
			pass := securityCheck(sec, w, r)
			if !pass {
				return
			}
		}

		fn(w, r, ps, s)
	}
}
