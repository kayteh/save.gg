package errors

import (
	"errors"
)

var (

	// models/user
	UserNotFound         = errors.New("user not found")
	UserPresentableSave  = errors.New("attempted presentable user save")
	UserNoIDSave         = errors.New("user id not set")
	UsernameInvalid      = errors.New("username invalid")
	UsernameTooLong      = errors.New("username too long")
	UsernameDisallowed   = errors.New("username is not allowed")
	UsernameTaken        = errors.New("username in use")
	UserPasswordTooShort = errors.New("password is too short")
)
