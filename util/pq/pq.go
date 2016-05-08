// Helpers for andvanced data structures in database/sql(x).
package pq

import (
	"errors"
)

var (
	// Errors that this package will use
	ErrTypeMismatch = errors.New("save.gg/pq: type mismatch")
)
