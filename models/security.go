package models

import (
	"gopkg.in/dgrijalva/jwt-go.v2"
	"net/http"
	"save.gg/sgg/meta"
	"save.gg/sgg/util/errors"
	//"time"
)

//TODO(kkz): implement API keys
func CheckAPIKeyRequest(r *http.Request) (ok bool, err error) {
	if r.Header.Get("API-Key") == "testkey" {
		return true, nil
	}

	return false, nil
}

func CheckCSRFRequest(r *http.Request) (ok bool, err error) {

	ts := r.Header.Get("CSRF-Token")

	if ts == "" {
		return false, nil
	}

	origin := r.Header.Get("Origin")

	if origin == "" {
		origin = "unset"
	}

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		return []byte(meta.App.Conf.CSRF.SigningKey), nil
	})

	if err != nil {
		meta.App.Log.WithError(err).Error("jwt errored")
		return false, errors.CSRFTokenInvalid
	}

	if !token.Valid {
		meta.App.Log.Warn("token invalid")
		return false, errors.CSRFTokenInvalid
	}

	to := token.Claims["origin"].(string)

	if origin != to {
		return false, errors.CSRFOriginMismatch
	}

	sid := token.Claims["session_id"].(string)
	sexists, err := SessionExists(sid)
	if err != nil {
		meta.App.Log.WithError(err).Error("session find errored")
		return false, err
	}

	return sexists, nil

}

func CheckSignedRequest(r *http.Request) (ok bool, err error) {
	return false, nil
}
