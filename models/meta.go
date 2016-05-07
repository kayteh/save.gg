package models

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func PrepModels(d *sqlx.DB) {
	db = d
	return
}
