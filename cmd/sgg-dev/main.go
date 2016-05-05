// Dev Server skips NATS-based routing and hooks everything directly.
package main

import (
	"save.gg/sgg/cmd/sgg-dev/run"
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
