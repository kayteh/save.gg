package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	mw "save.gg/sgg/cmd/sgg-api/run/middleware"
	"save.gg/sgg/meta"
	//m "save.gg/sgg/models"
	//util "save.gg/sgg/util/httputil"
)

func init() {

	meta.RegisterRoute("GET", "/~t/versioned", mw.VR(mw.VRMap{
		"default": versioned,
		"v1":      versionedV1,
		"v1a":     versionedV1a,
	}))

	meta.RegisterRoute("GET", "/~t/valid", mw.CA(secCheck))

}

func versioned(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"version":"2","default":true}`))
}

func versionedV1a(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"version":"1a"}`))
}

func versionedV1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"version":"1"}`))
}

func secCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(`{"verified":true}`))
}
