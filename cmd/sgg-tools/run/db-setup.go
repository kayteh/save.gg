package toolsrun

import (
	log "github.com/Sirupsen/logrus"
	a "save.gg/sgg/meta"
	m "save.gg/sgg/models"
)

func SetupDB() {
	app, err := a.SetupApp()
	if err != nil {
		log.WithError(err).Fatal("App initialization failed.")
	}

	a.App = app

	c, err := m.ConnectorFromApp(app)
	if err != nil {
		log.WithError(err).Fatal("connector creation error")
	}

	m.PrepModels(c)

}
