// Save.gg Toolbox.
//
// This package is documented at docs/Toolbox.md.
// There's a lot to explain, and code documentation isn't going
// to help much.
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
	"runtime"

	"save.gg/sgg/meta"

	debugUser "save.gg/sgg/cmd/sgg-tools/run/debug-user"
	"save.gg/sgg/cmd/sgg-tools/run/migrate"
	//"save.gg/sgg/cmd/sgg-tools/run/touch"
)

func main() {

	workers := runtime.NumCPU()
	runtime.GOMAXPROCS(workers)

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
				{
					Name:   "rethink",
					Action: migrate.Rethink,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "reset,r"},
					},
				},
				{
					Name:   "influx",
					Action: migrate.Influx,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "reset,r"},
					},
				},
			},
		},
		{
			Name:  "debug-user",
			Usage: "debug user-related things",
			Subcommands: []cli.Command{
				{
					Name:   "register",
					Action: debugUser.Register,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "test, t"},
						cli.StringFlag{Name: "username, u"},
						cli.StringFlag{Name: "password, p"},
						cli.StringFlag{Name: "email, e"},
						cli.BoolFlag{Name: "admin, a"},
					},
				},
				{
					Name:   "login",
					Action: debugUser.Login,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "username, u"},
						cli.BoolFlag{Name: "token, t"},
					},
				},
				{
					Name:   "touch",
					Action: debugUser.Touch,
				},
			},
		},
		{
			Name: "debug-config",
			Action: func(ctx *cli.Context) {
				c := meta.NewConfig(meta.ResolveConfigLocation())

				log.WithField("config", c).Info("config")
			},
		},
		{
			Name: "version",
			Action: func(ctx *cli.Context) {
				log.Infof("Save.gg :: Version " + meta.Version)
			},
		},
	}

	a.Run(os.Args)

}
