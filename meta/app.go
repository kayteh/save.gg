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
//
// This package can *never* import any save.gg/sgg packages. Period.
// Doing so will cause build failures unless that package in only imported by
// this package.
package meta

import (
	log "github.com/Sirupsen/logrus"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	redis "github.com/mediocregopher/radix.v2/pool"
	r "gopkg.in/dancannon/gorethink.v2"
)

// Structure for core app data like configuration, a logger, and anything pertinent to any part of the app.
// Data services should be placed here at risk of not existing when you actually want to use them.
type Application struct {
	Conf    Config
	Log     *log.Entry
	Influx  influx.Client
	Rethink *r.Session
	Redis   *redis.Pool

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

// Returns a postgres sqlx connection. This is sort-of goroutine-safe.
// It needs to be mounted manually, but will work find after that.
// See models.PrepModels().
func (a Application) GetPq() (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("postgres", a.Conf.Postgres.URL)

	return db, err

}

// Returns an InfluxDB connection. This is goroutine-safe.
func (a Application) GetInflux() (i influx.Client, err error) {
	i, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     a.Conf.Influx.Addr,
		Username: a.Conf.Influx.User,
		Password: a.Conf.Influx.Pass,
	})

	return i, err
}

// Returns a RethinkDB session. This is goroutine-safe.
func (a Application) GetRethink() (*r.Session, error) {
	var s *r.Session
	var err error

	s, err = r.Connect(r.ConnectOpts{
		Address: a.Conf.Rethink.Addr,
	})

	return s, err
}

// Returns a Redis session. This is goroutine-safe.
func (a Application) GetRedis() (*redis.Pool, error) {

	p, err := redis.New("tcp", a.Conf.Redis.Addr, 10)

	if err != nil {
		return nil, err
	}

	r := p.Cmd("PING")
	if r.Err != nil {
		return nil, err
	}

	return p, nil

}
