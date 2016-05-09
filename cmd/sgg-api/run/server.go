package run

import (
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"save.gg/sgg/meta"
	"save.gg/sgg/models"
	"time"

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

	hw := handlerWrapper{Router: r}

	http.ListenAndServeTLS(config.DevServer.Addr, config.Webserver.TLS.Cert, config.Webserver.TLS.Private, hw)

}

type responseWriterWrapper struct {
	ow   http.ResponseWriter
	code int
}

func (ww *responseWriterWrapper) Write(b []byte) (i int, err error) {
	i, err = ww.ow.Write(b)
	return i, err
}

func (ww *responseWriterWrapper) WriteHeader(i int) {
	ww.code = i
	ww.ow.WriteHeader(i)
}

func (ww *responseWriterWrapper) Header() http.Header {
	return ww.ow.Header()
}

type handlerWrapper struct {
	Router http.Handler
}

func (hw handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts := time.Now()

	ww := &responseWriterWrapper{ow: w, code: 200}

	if meta.App.Env == "production" {
		ww.Header().Add("Content-Security-Policy", "default-src: 'self'; script-src: 'self' x.svgg.xyz")
		ww.Header().Add("Strict-Transport-Security", "max-age=31536000")
	} else {
		ww.Header().Add("SGG-Message", "not production!")
	}

	hw.Router.ServeHTTP(ww, r)

	time := time.Since(ts)

	logEvent(r, ww, time)
}

func logEvent(r *http.Request, ww *responseWriterWrapper, d time.Duration) {

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	meta.App.Log.WithFields(log.Fields{
		"code":  ww.code,
		"time":  d.Seconds(),
		"agent": r.UserAgent(),
		"ip":    ip,
	}).Infof("%d %s", ww.code, r.RequestURI)

	// todo output to influxdb

}
