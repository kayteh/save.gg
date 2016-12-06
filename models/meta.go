// Save.gg data model package.
//
// Pretty much any and all business logic exists here.
package models

import (
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mediocregopher/radix.v2/pool"
	"save.gg/sgg/meta"
)

var (
	db      *sqlx.DB
	redis   *pool.Pool
	influx  influxdb.Client
)

// Connector is a helper for models to have goroutine-safe datastore connections.
type Connector struct {
	Pq      *sqlx.DB
	Redis   *pool.Pool
	Influx  influxdb.Client
}

// Prepares the models package to use goroutine-safe database connections.
func PrepModels(c *Connector) {
	db = c.Pq
	redis = c.Redis
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

	in, err := a.GetInflux()
	if err != nil {
		return nil, err
	}

	return &Connector{
		Pq:      pq,
		Redis:   rd,
		Influx:  in,
	}, nil
}
