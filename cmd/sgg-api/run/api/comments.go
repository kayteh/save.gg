package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"save.gg/sgg/meta"
	m "save.gg/sgg/models"
	u "save.gg/sgg/util/httputil"
)

func init() {
	meta.RegisterRoute("GET", "/api/comments/:id", getComment)
}

func getComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")

	d := m.CommentById(id)

	if d == nil {
		u.NotFound(w)
	} else {
		u.Output(w, d)
	}

}
