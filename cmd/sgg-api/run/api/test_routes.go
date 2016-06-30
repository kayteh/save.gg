package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	// mw "save.gg/sgg/cmd/sgg-api/run/middleware"
	// "save.gg/sgg/meta"
	//m "save.gg/sgg/models"
	//util "save.gg/sgg/util/httputil"
)

// GET /~t/versioned v2 default
func versioned(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"version":"2","default":true}`))
}

// GET /~t/versioned v1a
func versionedV1a(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"version":"1a"}`))
}

// GET /~t/versioned v1
func versionedV1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"version":"1"}`))
}

// GET /~t/valid
func secCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"verified":true}`))
}
