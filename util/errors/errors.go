// Global errors.
//
// This is in it's own package to prevent cyclic imports.
package errors

import (
	"errors"
)

var (

	// models/user
	UserAuthBadHandle    = errors.New("sgg: bad username/email")
	UserAuthBadPassword  = errors.New("sgg: bad password")
	UserEmailInvalid     = errors.New("sgg: email invalid")
	UserNoIDSave         = errors.New("sgg: user id not set")
	UserNotFound         = errors.New("sgg: user not found")
	UserPresentableSave  = errors.New("sgg: attempted presentable user save")
	UserPasswordTooShort = errors.New("sgg: password is too short")
	UsernameDisallowed   = errors.New("sgg: username is not allowed")
	UsernameInvalid      = errors.New("sgg: username invalid")
	UsernameTaken        = errors.New("sgg: username in use")
	UsernameTooLong      = errors.New("sgg: username too long")

	// models/session
	SessionTokenInvalid = errors.New("sgg: token invalid")
	SessionNotFound     = errors.New("sgg: session not found")
)
