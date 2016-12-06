package debuguser

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"save.gg/sgg/cmd/sgg-tools/run"
	m "save.gg/sgg/models"
)

func Register(ctx *cli.Context) {
	toolsrun.SetupDB()

	// Get info
	email := ctx.String("email")
	password := ctx.String("password")
	username := ctx.String("username")
	admin := ctx.Bool("admin")

	// Create user
	u := m.NewUser()
	u.Activated = true

	err := u.CreateSecret(password)
	if err != nil {
		log.WithError(err).Fatal("Password error")
	}

	err = u.SetUsername(username)
	if err != nil {
		log.WithError(err).Fatal("Username error")
	}

	err = u.SetEmail(email)
	if err != nil {
		log.WithError(err).Fatal("Email error")
	}

	if admin {
		u.ACL.Append("admin")
	}

	err = u.Insert()
	if err != nil {
		log.WithError(err).Fatal("Save error")
	}

	log.Printf("Successfully created user %s", username)
}
