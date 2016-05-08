// Core systems required by every section of the app.
//
// Most things you'll do here might not be goroutine-safe.
// Make sure to re-define at runtime whenever possible,
// otherwise mister "nil pointer dereference" will ruin your day.
//
// App.Conf, App.Env, and App.Log should always be goroutine safe.
// If they are not, something completely horrible has gone wrong.
//
// This package also manages the side-effect HTTP router in `MountRouter` and
// `RegisterRoute`. Usage on this can be found at their specific documentation.
package meta

import (
	_ "database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Structure for core app data like configuration, a logger, and anything pertinent to any part of the app.
// Data services should be placed here at risk of not existing when you actually want to use them.
type Application struct {
	Conf Config
	Log  *log.Entry

	Env string
}

// This is one of those ugly global variables. It works great :^)
var App *Application

// Build the core systems of the app, and returns it. This should never set App iself,
// that duty is up to the system calling it.
func SetupApp() (*Application, error) {

	a := Application{}

	a.Log = log.New().WithFields(log.Fields{})

	configLoc := ResolveConfigLocation()
	a.Conf = NewConfig(configLoc)

	a.Env = a.Conf.Self.Env

	return &a, nil
}

// Returns a postgres sqlx connection. This is goroutine-safe, but make sure you mount this
// to models.PrepModels at app setup time to be very sure.
func (a Application) GetPq() (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("postgres", a.Conf.Postgres.URL)

	return db, err

}
