package httputil

import (
	"encoding/json"
	"net/http"

	//a "save.gg/sgg/meta"
)

func Output(w http.ResponseWriter, d interface{}) (err error) {
	var o []byte

	//if a.App.Env == "local" {
	o, err = json.MarshalIndent(d, "", "  ")
	// } else {
	// 	o, err = json.Marshal(d)
	// }

	if err != nil {
		return err
	}

	w.Write(o)

	return nil
}

func NotFound(w http.ResponseWriter) (err error) {

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"err":"not found"}`))

	return nil

}

func Forbidden(w http.ResponseWriter) (err error) {

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"err":"forbidden"}`))

	return nil

}
