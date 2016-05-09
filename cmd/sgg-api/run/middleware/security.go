package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"save.gg/sgg/meta"
	m "save.gg/sgg/models"
	util "save.gg/sgg/util/httputil"
	"strings"
)

// Optional flags to pass to a middleware constraint.
//
// CSRF, API Key, or Signed can pass if either is true (inclusive or.) This allows API consumers to use their
// signature to validate instead of getting a CSRF value. Once one is found, execution ends.
//
// API Key-based verification should be *heavily* rate-limited.
//
// Use nil instead of an empty initializer unless you intend for the check to always fail.
type SecurityFlags struct {
	// Individual flags
	CSRF     bool
	Signed   bool
	APIKey   bool
	Internal bool // a.k.a. Universe

	// Security package aliases
	// (only set one of these at a time, they don't play well together)
	All   bool // CSRF, Signed, APIKey
	Loose bool // Signed, NoEnforce

	// Enforcement mode. Set to true if you don't care if the endpoint is validated,
	// but still might want, for instance, a signature.
	NoEnforce bool

	// Ratelimit to enforce
	RateLimit int
}

// Resolves aliases such as SecurityFlags.All.
func (s *SecurityFlags) Resolve() {
	if s.All {
		s.CSRF = true
		s.APIKey = true
		s.Signed = true
		s.Internal = true
	}

	if s.Loose {
		s.Signed = true
		s.Internal = true
		s.NoEnforce = true
	}
}

// Writes a header to communicate what methods of validation were tried in case of failure.
func (s *SecurityFlags) WriteTriedHeader(w http.ResponseWriter) {
	var v []string

	if s.CSRF {
		v = append(v, "CSRF-Token")
	}

	if s.Signed {
		v = append(v, "Signature")
	}

	if s.APIKey {
		v = append(v, "API-Key")
	}

	if s.Internal {
		v = append(v, "Universe")
	}

	w.Header().Add("SGG-Validations-Tried", strings.Join(v, ", "))
}

// Checks all forms of authentication available, returns true if it passed, false if not.
// This function will handle calling various errors, the other end just needs to act upon
// it's return value to stop execution.
func securityCheck(sec *SecurityFlags, w http.ResponseWriter, r *http.Request) bool {
	// in this case, everything passes, as nothing is checked.
	if sec == nil {
		return true
	}

	sec.Resolve()

	if sec.Internal {
		valid, err := m.CheckInternalRequest(r)
		if err != nil {
			meta.App.Log.WithError(err).Error("internal signature validation error")
			//util.InternalServerError(w, err)
		}

		if valid {
			return true
		}
	}

	//TODO(kkz): add normal rate-limiting check here (after internal so it's never enforced for ourselves.)

	if sec.Signed {
		valid, err := m.CheckSignedRequest(r)
		if err != nil {
			meta.App.Log.WithError(err).Error("signature validation error")
			//util.InternalServerError(w, err)
		}

		if valid {
			return true
		}
	}

	if sec.CSRF {
		valid, err := m.CheckCSRFRequest(r)
		if err != nil {
			meta.App.Log.WithError(err).Error("csrf validation error")
			//util.InternalServerError(w, err)
		}

		if valid {
			return true
		}
	}

	if sec.APIKey {
		valid, err := m.CheckAPIKeyRequest(r)
		if err != nil {
			meta.App.Log.WithError(err).Error("api key validation error")
			//util.InternalServerError(w, err)
			//return false
		}

		if valid {
			return true
		}
	}

	// If we reached here, and we aren't enforcing, just return true.
	if sec.NoEnforce {
		return true
	}

	// if it reaches here, all have failed, or this is on purpose.
	//TODO(kkz): add failed security ratelimiting here.
	sec.WriteTriedHeader(w)
	util.Forbidden(w)
	return false

}

// Enforces an endpoint to the security constraints.
func SecurityCheck(fn httprouter.Handle, sec *SecurityFlags) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		pass := true
		if sec != nil {
			pass = securityCheck(sec, w, r)
		}

		if !pass {
			return
		}

		fn(w, r, ps)
	}
}

// Shorthand for a "loose" endpoint. (SC = Signature check)
func SC(fn httprouter.Handle) httprouter.Handle {
	return SecurityCheck(fn, &SecurityFlags{Loose: true})
}

// Shorthand for an "all" endpoint. (CA = Check all)
func CA(fn httprouter.Handle) httprouter.Handle {
	return SecurityCheck(fn, &SecurityFlags{All: true})
}

// Shorthane for internal routes.
func I(fn httprouter.Handle) httprouter.Handle {
	return SecurityCheck(fn, &SecurityFlags{Internal: true})
}
