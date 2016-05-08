// API server
package main

import (
	"save.gg/sgg/cmd/sgg-api/run"
	"save.gg/sgg/meta"
)

func main() {
	app, err := meta.SetupApp()

	if err != nil {
		app.Log.WithError(err).Fatal("app couldn't start")
	}

	meta.App = app

	run.Start()
}
