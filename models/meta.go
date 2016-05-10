// Save.gg data model package.
//
// Pretty much any and all business logic exists here.
package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/mediocregopher/radix.v2/pool"
)

var db *sqlx.DB
var redis *pool.Pool

// Prepares the models package to use goroutine-safe database connections.
func PrepModels(d *sqlx.DB, r *pool.Pool) {
	db = d
	redis = r
	return
}
