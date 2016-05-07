package meta

import (
	_ "database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Application struct {
	Conf Config
	Log  *log.Entry

	Env string
}

var App *Application

func SetupApp() (*Application, error) {

	a := Application{}

	a.Log = log.New().WithFields(log.Fields{})

	configLoc := ResolveConfigLocation()
	a.Conf = NewConfig(configLoc)

	a.Env = a.Conf.Self.Env

	return &a, nil
}

func (a Application) GetPq() (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("postgres", a.Conf.Postgres.URL)

	return db, err

}
