package debuguser

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/segmentio/go-prompt"
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

	if username == "" {
		username = prompt.StringRequired("Username")
	}

	if email == "" {
		email = prompt.StringRequired("Email")
	}

	if password == "" {
		password = prompt.PasswordMasked("Password")
	}

	// Create user
	u := m.User{
		Email:     email,
		Activated: true,
	}

	err := u.CreateSecret(password)
	if err != nil {
		log.WithError(err).Fatal("Password error")
	}

	err = u.SetUsername(username)
	if err != nil {
		log.WithError(err).Fatal("Username error")
	}

	if admin {
		u.AppendACL("admin")
	}

	err = u.Insert()
	if err != nil {
		log.WithError(err).Fatal("Save error")
	}

	log.Printf("Successfully created user %s", username)
}
