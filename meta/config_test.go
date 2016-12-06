package meta

import (
	//"errors"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"os"
	"testing"
)

func writeToml(name string, data string) error {

	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	f.Sync()
	f.Close()

	return nil
}

func makeTmpDir() (d string, err error) {
	td := os.TempDir() + "gg.save.test/"
	u := uuid.NewV4().String()
	d = td + u

	if err = os.MkdirAll(d+"/conf.d", 0777); err != nil {
		return d, err
	}

	//log.WithField("temp dir", d).Info("using temp dir")

	return d, nil
}

func destroyTmpDir(tmpDir string) error {
	return os.RemoveAll(tmpDir)
}

func TestBasicConfiguration(t *testing.T) {
	tmpDir, err := makeTmpDir()
	if err != nil {
		t.Fatal(err)
	}
	def := `
[self]
env = "test"
signing_key = "1234567890abcdefghijklmnopqrstuvwxyz"
`

	writeToml(tmpDir+"/app.toml", def)

	c := NewConfig(tmpDir + "/app.toml")

	if c.Self.Env != "test" {
		log.WithField("config", c).Error("configuration didn't set.")
		t.Fail()
	}

	destroyTmpDir(tmpDir)

}

func TestCascadingConfiguration(t *testing.T) {
	tmpDir, err := makeTmpDir()
	if err != nil {
		t.Fatal(err)
	}

	def := `
[self]
env = "test"

[security.signing_keys]
CSRF = "1234567890abcdefghijklmnopqrstuvwxyz"
`

	writeToml(tmpDir+"/app.toml", def)

	def2 := `
[self]
env = "test-00"
`

	writeToml(tmpDir+"/conf.d/00-test.toml", def2)

	c := NewConfig(tmpDir + "/app.toml")

	if c.Self.Env != "test-00" {
		log.WithField("config", c).Error("configuration didn't set.")
		t.Fail()
	}

	if c.Security.SigningKeys.CSRF != "1234567890abcdefghijklmnopqrstuvwxyz" {
		log.Error("configuration overwrote without superceding")
		t.Fail()
	}

	destroyTmpDir(tmpDir)
}

func TestMain(m *testing.M) {
	r := m.Run()
	os.Exit(r)
}
