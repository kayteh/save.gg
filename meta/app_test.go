package meta

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

func TestApplicationSetup(t *testing.T) {

	a, err := SetupApp()
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("testlog")

	if a.Env == "" {
		t.Error("env unset")
	}

	if a.Conf.Self.Env == "" {
		t.Error("config unset")
	}

	db, err := a.GetPq()
	if err != nil {
		t.Error("get postgres error", err)
	}

	err = db.Ping()
	if err != nil {
		t.Error("db ping error", err)
	}

}
