package meta

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"os"
	"testing"
)

var tmpDir string = ""

func writeToml(name string, data string) error {
	if tmpDir == "" {
		return errors.New("tmpDir must be set")
	}

	//log.WithField("filename", tmpDir+"/"+name+".toml").Info("writing out toml")

	f, err := os.OpenFile(tmpDir+"/"+name+".toml", os.O_WRONLY|os.O_CREATE, 0644)

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

	if err = os.Mkdir(td+u, 0777); err != nil {
		return d, err
	}

	d = td + u

	//log.WithField("temp dir", d).Info("using temp dir")

	return d, nil
}

func TestBasicConfiguration(t *testing.T) {
	var err error
	if tmpDir, err = makeTmpDir(); err != nil {
		t.Fatal(err)
	}

	def := `
role = "backend"

[webserver]
port = 8441

	[ssl]
	cert = "/dev/null"
	key = "/dev/null"

[network]
bind = "0.0.0.0"
publish = "127.0.1.1"
`

	writeToml("default", def)

	a := &Application{
		configDir: tmpDir,
		Env:       "test-jig",
		Log:       log.New().WithFields(log.Fields{}),
	}

	a.Conf = a.NewConfig()

	if a.Conf.Network.Publish != "127.0.1.1" {
		t.Fail()
	}
}

func TestCascadingConfiguration(t *testing.T) {
	var err error
	if tmpDir, err = makeTmpDir(); err != nil {
		t.Fatal(err)
	}

	def := `
role = "backend"

[webserver]
port = 8441

	[ssl]
	cert = "/dev/null"
	key = "/dev/null"

[network]
bind = "0.0.0.0"
publish = "127.0.1.1"
`
	writeToml("default", def)

	tj := `
[network]
publish = "127.0.2.1"
`

	writeToml("test-jig", tj)

	a := &Application{
		Env:       "test-jig",
		Log:       log.New().WithFields(log.Fields{}),
		configDir: tmpDir,
	}

	a.Conf = a.NewConfig()

	if a.Conf.Network.Publish != "127.0.2.1" {
		t.Fail()
	}

	if t.Failed() {
		if a.Conf.Network.Publish == "127.0.1.1" {
			t.Log("failed because test config didn't cascade")
		}
	}
}

func TestMain(m *testing.M) {
	r := m.Run()
	os.RemoveAll(tmpDir)
	tmpDir = ""
	os.Exit(r)
}
