package run

import (
	"net/http"
	"save.gg/sgg/meta"

	"github.com/julienschmidt/httprouter"

	_ "save.gg/sgg/cmd/sgg-api/run/api"
)

func Start() {
	meta.App.Log.Info("Starting development server...")

	r := httprouter.New()
	meta.MountRouter(r)

	config := meta.App.Conf

	meta.App.Log.Infof("sgg-dev is now serving on https://%s...", config.DevServer.Addr)
	meta.App.Log.Info("Happy coding!~")
	http.ListenAndServeTLS(config.DevServer.Addr, config.Webserver.TLS.Cert, config.Webserver.TLS.Private, r)

}
