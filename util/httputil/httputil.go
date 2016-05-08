// HTTP REST API utilities
package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	a "save.gg/sgg/meta"
)

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

// Returns a generic 404 response
func NotFound(w http.ResponseWriter) (err error) {

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"err":"not found"}`))

	return nil

}

// Returns a generic 403 response
func Forbidden(w http.ResponseWriter) (err error) {

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"err":"forbidden"}`))

	return nil

}

// Returns a generic 500 response
func InternalServerError(w http.ResponseWriter, e error) (err error) {

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"err":"%s"}`, e.Error())))

	return nil
}
