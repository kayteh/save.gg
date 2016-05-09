package models

import (
	"gopkg.in/dgrijalva/jwt-go.v2"
	"net/http"
	"save.gg/sgg/meta"
	"save.gg/sgg/util/errors"
	"time"
)

// Check API keys in the store for equality to the API-Key header, and applies a rate-limiting function.
//
// This is *heavily* rate-limited, and only suitable for light usage and development.
// Ideally, 3 should be allowed per minute, and this should never be accepted for authentication.
// Internal keys are **never** permitted to use this authentication.
//TODO(kkz): implement API keys
func CheckAPIKeyRequest(r *http.Request) (ok bool, err error) {
	key := r.Header.Get("API-Key")

	if key == "" {
		return false, nil
	}

	if key == "testkey" {
		return true, nil
	}

	if !consumerApiKeyLimitCheck(key, "key-verify") {
		return false, errors.APIRateLimited
	}

	c, err := ConsumerByAPIKey(key)
	if err == errors.ConsumerAPIKeyNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if c.Internal {
		return false, nil
	}

	c.RLTick("key-verify")

	return false, nil
}

// Check CSRF header for a Session ID and Origin, and tests if they are valid.
//
// This is computed on the frontend renderer,
// although *Session.GenerateCSRFToken can do it on this side of the app for testing reasons.
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
		return []byte(meta.App.Conf.Security.SigningKeys.CSRF), nil
	})

	if err != nil {
		meta.App.Log.WithError(err).Error("jwt errored")
		return false, errors.SecurityTokenInvalid
	}

	if !token.Valid {
		meta.App.Log.Warn("token invalid")
		return false, errors.SecurityTokenInvalid
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

// Verifies an asymmetric JWT signature. We only keep the public keys.
//
// This is always suggested to be present for API requests,
// and required for all API requests that change data,
// unless API-Key or CSRF is supplied instead.
func CheckSignedRequest(r *http.Request) (ok bool, err error) {

	ts := r.Header.Get("Signature")

	if ts == "" {
		return false, nil
	}

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		k, err := LookupPublicKey(token.Header["api_key"].(string))
		return k, err
	})

	if err != nil {
		meta.App.Log.WithError(err).Error("jwt errored")
		return false, errors.SecurityTokenInvalid
	}

	if !token.Valid {
		meta.App.Log.Warn("token invalid")
		return false, errors.SecurityTokenInvalid
	}

	if checkSignedClaims(r, token) {
		return true, nil
	}

	return false, nil
}

// Verifies an internal (usually web->api) request from the Universe header.
//
// Since if this is a valid signature, and we trust it a fair bit more,
// this should never be rate limited. If it fails, the security middleware
// should ratelimit it as normal instead AND log it in the security log.
// This is clever engineering. Yes.
//
// We check the same claims as a signed request here.
func CheckInternalRequest(r *http.Request) (ok bool, err error) {

	ts := r.Header.Get("Universe")

	if ts == "" {
		return false, nil
	}

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		return []byte(meta.App.Conf.Security.SigningKeys.Internal), nil
	})

	if err != nil {
		meta.App.Log.WithError(err).Error("jwt errored")
		return false, errors.SecurityTokenInvalid
	}

	if !token.Valid {
		meta.App.Log.Warn("token invalid")
		return false, errors.SecurityTokenInvalid
	}

	if checkSignedClaims(r, token) {
		return true, nil
	}

	return false, nil
}

// Check claims against the request URI, http method, agent string, and creation time.
// The time claim has a maximum tolerance of 10 seconds.
func checkSignedClaims(r *http.Request, t *jwt.Token) bool {

	ct := time.Unix(t.Claims["time"].(int64), 0)

	timeOk := time.Since(ct) < 10*time.Second

	return timeOk &&
		t.Claims["uri"].(string) == r.RequestURI &&
		t.Claims["method"].(string) == r.Method &&
		t.Claims["agent"].(string) == r.UserAgent()

}
