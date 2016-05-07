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
	pq, err := app.GetPq()
	if err != nil {
		log.WithError(err).Fatal("DB initialization failed.")
	}
	m.PrepModels(pq)

}
