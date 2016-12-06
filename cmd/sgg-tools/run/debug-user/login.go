package debuguser

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"save.gg/sgg/cmd/sgg-tools/run"
	m "save.gg/sgg/models"
	"save.gg/sgg/util/errors"
)

func Login(ctx *cli.Context) {

	handle := ctx.String("username")
	password := ctx.String("password")

	toolsrun.SetupDB()

	u, err := m.UserAuth(handle, password)

	if err == errors.UserAuthBadHandle || err == errors.UserAuthBadPassword {
		log.WithError(err).Error("wrong credentials")
		return
	}

	if err != nil {
		log.WithError(err).Fatal("authentication problem")
		return
	}

	log.Printf("hello %s!", u.Username)

	if ctx.Bool("token") {
		s, err := u.CreateSession()
		if err != nil {
			log.WithError(err).Error("session creation failure")
			return
		}

		t, err := s.Token()
		if err != nil {
			log.WithError(err).Error("token creation failure")
			return
		}

		log.Infof("token: %s", t)
	}

	return

}
