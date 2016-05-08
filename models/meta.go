// Save.gg data model package.
//
// Pretty much any and all business logic exists here.
package models

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// Prepares the models package to use goroutine-safe database connections.
func PrepModels(d *sqlx.DB) {
	db = d
	return
}
