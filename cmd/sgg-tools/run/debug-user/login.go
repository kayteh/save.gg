package debuguser

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/segmentio/go-prompt"
	"save.gg/sgg/cmd/sgg-tools/run"
	m "save.gg/sgg/models"
)

func Login(ctx *cli.Context) {

	username := ctx.String("username")
	if username == "" {
		username = prompt.StringRequired("username >")
	}

	slug := m.TransformUsernameToSlug(username)

	password := prompt.PasswordMasked("password >")

	toolsrun.SetupDB()

	u, err := m.UserBySlug(slug)
	if err != nil {
		log.WithError(err).Fatalln("get user failed")
	}

	pass, err := u.TestSecret(password)
	if err != nil {
		log.WithError(err).Fatalln("get user failed")
	}

	log.WithField("passed?", pass).Info("output")

	return

}
