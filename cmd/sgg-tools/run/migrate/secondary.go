package migrate

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	influx "github.com/influxdata/influxdb/client/v2"
	r "gopkg.in/dancannon/gorethink.v2"
	//a "save.gg/sgg/meta"
)

// Setup RethinkDB. No errors are caught here in order to let this grow incrementally.
// It's recommended to never remove anything from here unless it's superfluous.
func Rethink(ctx *cli.Context) {
	conf := getConfig()

	var s *r.Session

	//safeMode := ctx.Bool("safe")

	s, err := r.Connect(r.ConnectOpts{
		Address: conf.Rethink.Addr,
	})
	if err != nil {
		log.WithError(err).Fatal("influxdb connection error")
	}

	r.DBDrop("test").RunWrite(s)
	log.Info("dropped default database `test`")

	if ctx.Bool("reset") {
		r.DBDrop("sgg").RunWrite(s)
		log.Info("reset: dropped database `sgg`")
	}

	r.DBCreate("sgg").RunWrite(s)
	log.Info("created database `sgg`")

	s.Use("sgg")

	r.TableCreate("user_known_ips", r.TableCreateOpts{
		PrimaryKey: "user_id",
	}).RunWrite(s)
	log.Info("create table `user_known_ips`")

	r.TableCreate("user_old_secrets", r.TableCreateOpts{
		PrimaryKey: "user_id",
	}).RunWrite(s)
	log.Info("create table `user_old_secrets`")

	r.TableCreate("rate_limits").RunWrite(s)
	log.Info("create table `rate_limits`")

	r.Table("rate_limits").IndexCreate("key").RunWrite(s)
	log.Info("create index for `key` on table `rate_limits`")

}

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
