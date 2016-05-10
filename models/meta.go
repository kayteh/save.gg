// Save.gg data model package.
//
// Pretty much any and all business logic exists here.
package models

import (
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mediocregopher/radix.v2/pool"
	r "gopkg.in/dancannon/gorethink.v2"
	"save.gg/sgg/meta"
)

var (
	db      *sqlx.DB
	redis   *pool.Pool
	rethink *r.Session
	influx  influxdb.Client
)

// Connector is a helper for models to have goroutine-safe datastore connections.
type Connector struct {
	Pq      *sqlx.DB
	Redis   *pool.Pool
	Rethink *r.Session
	Influx  influxdb.Client
}

// Prepares the models package to use goroutine-safe database connections.
func PrepModels(c *Connector) {
	db = c.Pq
	redis = c.Redis
	rethink = c.Rethink
	influx = c.Influx
	return
}

// Outputs a connector via a meta.Application, since you'll probably be doing this anyway.
func ConnectorFromApp(a *meta.Application) (*Connector, error) {
	pq, err := a.GetPq()
	if err != nil {
		return nil, err
	}

	rd, err := a.GetRedis()
	if err != nil {
		return nil, err
	}

	re, err := a.GetRethink()
	if err != nil {
		return nil, err
	}

	in, err := a.GetInflux()
	if err != nil {
		return nil, err
	}

	return &Connector{
		Pq:      pq,
		Redis:   rd,
		Rethink: re,
		Influx:  in,
	}, nil
}
