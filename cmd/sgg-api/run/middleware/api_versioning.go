package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
)

var (
	versionRegex = regexp.MustCompile(`^application/vnd.svgg.(v[0-9]+(?:[a-z]+)?)\+json$`)
)

type VRMap map[string]httprouter.Handle

func VR(i VRMap) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		v := parseAcceptVersion(r.Header.Get("Accept"))

		route, ok := i[v]
		if !ok {
			w.Header().Add("Sgg-Api-Version", "default")
			i["default"](w, r, ps)
			return
		}

		w.Header().Add("Sgg-Api-Version", v)
		route(w, r, ps)

	}
}

// application/vnd.svgg[.version]+json
func parseAcceptVersion(h string) string {

	s := versionRegex.FindStringSubmatch(h)

	if len(s) == 2 {
		return s[1]
	}

	return "default"
}
