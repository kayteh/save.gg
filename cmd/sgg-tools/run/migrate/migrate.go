package migrate

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	m "github.com/mattes/migrate/migrate"
	a "save.gg/sgg/meta"

	_ "github.com/mattes/migrate/driver/postgres"
)

func getConfig() a.Config {
	return a.NewConfig(a.ResolveConfigLocation())
}

func Up(ctx *cli.Context) {
	c := getConfig()

	log.Println("Running migrations...")

	allErrors, ok := m.UpSync(c.Postgres.URL, "./migrations")
	if !ok {
		log.WithField("errors", allErrors).Errorln("Migration failed.")
	}
}

func Down(ctx *cli.Context) {
	c := getConfig()

	log.Println("Rolling back migrations...")

	allErrors, ok := m.DownSync(c.Postgres.URL, "./migrations")
	if !ok {
		log.WithField("errors", allErrors).Errorln("Migration failed.")
	}
}

func Redo(ctx *cli.Context) {
	c := getConfig()

	log.Println("Rolling back one and doing migration...")

	allErrors, ok := m.RedoSync(c.Postgres.URL, "./migrations")
	if !ok {
		log.WithField("errors", allErrors).Errorln("Migration failed.")
	}
}

func Create(ctx *cli.Context) {
	c := getConfig()
	n := ctx.Args().First()

	if n == "" {
		log.Fatalln("Specify a migration name.")
	}

	f, err := m.Create(c.Postgres.URL, "./migrations", n)
	log.Infof("Created migration %s.\n", f.UpFile.Name)

	if err != nil {
		log.WithField("errors", err).Fatal("Migration failed.")
	}

}
