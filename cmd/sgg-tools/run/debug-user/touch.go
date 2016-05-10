package debuguser

import (
	"github.com/codegangsta/cli"
	"save.gg/sgg/cmd/sgg-tools/run"
	"save.gg/sgg/meta"
	"save.gg/sgg/models"
)

func Touch(ctx *cli.Context) {
	toolsrun.SetupDB()

	slug := ctx.Args().First()

	meta.App.Log.Infof("attempting touch of slug %s", slug)
	u, err := models.UserBySlug(slug)
	if err != nil {
		meta.App.Log.WithError(err).Fatal("user fetch error")
	}

	u.Touch()

	meta.App.Log.Infof("user %s has been touched", u.Username)

}
