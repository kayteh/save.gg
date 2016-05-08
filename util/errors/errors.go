// Global errors.
//
// This is in it's own package to prevent cyclic imports.
package errors

import (
	"errors"
)

var (

	// models/user
	UserNotFound         = errors.New("sgg: user not found")
	UserPresentableSave  = errors.New("sgg: attempted presentable user save")
	UserNoIDSave         = errors.New("sgg: user id not set")
	UsernameInvalid      = errors.New("sgg: username invalid")
	UsernameTooLong      = errors.New("sgg: username too long")
	UsernameDisallowed   = errors.New("sgg: username is not allowed")
	UsernameTaken        = errors.New("sgg: username in use")
	UserPasswordTooShort = errors.New("sgg: password is too short")
	UserAuthBadHandle    = errors.New("sgg: bad username/email")
	UserAuthBadPassword  = errors.New("sgg: bad password")

	// models/session
	SessionTokenInvalid = errors.New("sgg: token invalid")
	SessionNotFound     = errors.New("sgg: session not found")
)
