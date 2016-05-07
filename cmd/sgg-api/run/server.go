package run

import (
	"net/http"
	"save.gg/sgg/meta"

	"github.com/julienschmidt/httprouter"
	"save.gg/sgg/models"

	_ "save.gg/sgg/cmd/sgg-api/run/api"
)

func Start() {
	meta.App.Log.Info("Starting api server...")

	r := httprouter.New()
	meta.MountRouter(r)

	pq, err := meta.App.GetPq()
	if err != nil {
		meta.App.Log.WithError(err).Fatalln("Database failed")
	}

	models.PrepModels(pq)

	config := meta.App.Conf

	meta.App.Log.Infof("sgg-api is now serving on https://%s...", config.DevServer.Addr)

	if meta.App.Env == "local" {
		meta.App.Log.Info("Happy coding!~")
	}

	http.ListenAndServeTLS(config.DevServer.Addr, config.Webserver.TLS.Cert, config.Webserver.TLS.Private, r)

}
