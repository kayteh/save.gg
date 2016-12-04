package migrate

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	influx "github.com/influxdata/influxdb/client/v2"
)

// Setup InfluxDB. This implies you're running influxdb in dev, so it'll connect with root:root.
// In production, influxDB should be set up with orchestration systems instead.
func Influx(ctx *cli.Context) {
	conf := getConfig()

	ic := influx.HTTPConfig{
		Addr:     conf.Influx.Addr,
		Username: "root",
		Password: "root",
	}

	i, err := influx.NewHTTPClient(ic)
	if err != nil {
		log.WithError(err).Fatal("influxdb connection error")
	}

	if ctx.Bool("reset") {
		i.Query(influx.NewQuery(fmt.Sprintf(`
			DROP DATABASE sgg;
		`), "", ""))
		log.Info("reset: dropped database `sgg`")
	}

	i.Query(influx.NewQuery(fmt.Sprintf(`
		CREATE USER %s WITH PASSWORD '%s';
	`, conf.Influx.User, conf.Influx.Pass), "", ""))
	log.Infof("created user `%s`", conf.Influx.User)

	i.Query(influx.NewQuery(fmt.Sprintf(`
		CREATE DATABASE sgg;
	`), "", ""))
	log.Info("created database `sgg`")

	i.Query(influx.NewQuery(fmt.Sprintf(`
		GRANT ALL ON sgg TO %s;
	`, conf.Influx.User), "", ""))
	log.Infof("granted all on `sgg` to %s", conf.Influx.User)

}
