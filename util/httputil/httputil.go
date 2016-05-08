// HTTP REST API utilities
package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	a "save.gg/sgg/meta"
	_ "save.gg/sgg/util/errors"
)

// Decodes a JSON payload into the interface
func Input(r *http.Request, i interface{}) (err error) {
	d := json.NewDecoder(r.Body)
	return d.Decode(i)
}

// Outputs d to JSON, and that JSON to the ResponseWriter. This is an express route.
func Output(w http.ResponseWriter, d interface{}) (err error) {
	var o []byte

	if a.App.Env == "local" {
		o, err = json.MarshalIndent(d, "", "    ")
	} else {
		o, err = json.Marshal(d)
	}

	if err != nil {
		return err
	}

	w.Write(o)

	return nil
}

// Returns a generic 204 response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

// Returns a generic 403 response
func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"err":"forbidden"}`))
}

// Returns a generic 404 response
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"err":"not found"}`))
}

// Returns a generic 422 response
func BadInput(w http.ResponseWriter) {
	w.WriteHeader(422)
	w.Write([]byte(`{"err":"input unprocessable"}`))
}

// Returns a generic 500 response
func InternalServerError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"err":"%s"}`, e.Error())))
}
