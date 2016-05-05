package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"

	"save.gg/sgg/meta"

	"save.gg/sgg/cmd/sgg-tools/run/migrate"
	//"save.gg/sgg/cmd/sgg-tools/run/touch"
)

func main() {

	a := cli.NewApp()
	a.Name = "sgg-tools"
	a.Usage = "Save.gg Cluster/Command/Automation Tools"
	a.Commands = []cli.Command{
		{
			Name:   "migrate",
			Usage:  "Updates the database to a new schema.",
			Action: migrate.Up,
			Subcommands: []cli.Command{
				{
					Name:   "up",
					Action: migrate.Up,
				},
				{
					Name:   "create",
					Action: migrate.Create,
				},
				{
					Name:   "redo",
					Action: migrate.Redo,
				},
				{
					Name:   "down",
					Action: migrate.Down,
				},
			},
		},
		{
			Name:  "touch",
			Usage: "Touches a model.",
			//Action: touch.CliStart(),
		},
		{
			Name: "debug-config",
			Action: func(ctx *cli.Context) {
				c := meta.NewConfig(meta.ResolveConfigLocation())

				log.WithField("config", c).Info("config")
			},
		},
	}

	a.Run(os.Args)

}
