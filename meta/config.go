package meta

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Resolves the configuration location from various sources.
//TODO(kkz): Actually implement this properly
func ResolveConfigLocation() string {
	return "./config/app.toml"
}

// Parses many TOML files from a few different locations to create the app config.
//
// Precedence
//
// command-line,
//
// conf.d in glob order,
//
// root (app.toml),
//
// Using this, you could set two postgres configs, 00-postgres.toml and
// 01-postgres.toml, and the latter will overwrite any values 00 has, and
// 00 will overwrite any values app.toml has. Exploit this for production config.
func NewConfig(path string) (conf Config) {
	good := false

	// get outer config dir
	d := filepath.Dir(path)

	// check if main config exists
	_, err := os.Stat(path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			//log.Print("config: base config not found")
		} else {
			log.Panicf("config: stat err :: %v", err)
		}
	} else {
		_, err = toml.DecodeFile(path, &conf)
		if err != nil {
			log.Panicf("config: decode err :: %v", err)
		} else {
			//log.Printf("config: loaded root: %s", path)
			good = true
		}
	}

	// check conf.d
	cd, err := os.Stat(d)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			//log.Printf("config: %s doesn't exist, skipping")
		} else {
			log.Panicf("config: stat err :: %v", err)
		}
	} else {
		if cd.IsDir() {
			cdf, _ := filepath.Glob(filepath.Join(d, "conf.d", "*.toml"))

			for _, c := range cdf {
				_, err = toml.DecodeFile(c, &conf)
				if err != nil {
					//log.Printf("config: decode err (file: %s) :: %v", c, err)
				} else {
					//log.Printf("config: loaded conf.d: %s", c)
					good = true
				}
			}
		} else {
			//log.Printf("config: %s isn't a directory, skipping")
		}
	}

	if !good {
		log.Fatal("config: i didn't load any sort of config. exiting.")
	}

	return conf
}

// Configuration structure. Keep this in alphabetical order.
//TODO(kkz): consolidate and cleanup
type Config struct {
	Cache      cacheConfig
	DevServer  devserverConfig `toml:"dev-server"`
	NATS       natsConfig
	Postgres   pgConfig
	Self       selfConfig
	Validation validationConfig
	Webserver  webserverConfig
}

type cacheConfig struct {
	Backend string

	TTL struct {
		Comment time.Duration
		User    time.Duration
		Save    time.Duration
		Session time.Duration
	}

	Redis struct {
		Addr string
	}
}

type devserverConfig struct {
	Addr string
}

type natsConfig struct {
	URL string
}

type pgConfig struct {
	URL string
}

type selfConfig struct {
	Env           string
	Revision      string
	SessionCookie string `toml:"session_cookies"`
	SigningKey    string `toml:"signing_key"`
}

type validationConfig struct {
	PasswordLength  int      `toml:"password_length"`
	DisallowedSlugs []string `toml:"disallowed_slugs"`
	UsernameLength  int      `toml:"username_length"`
}

type webserverConfig struct {
	TLS struct {
		Cert    string
		Private string
	}
}
