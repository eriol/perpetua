package config

import (
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/gcfg"
)

const Version = "0.1a"

const DEFAULT_NICKNAME = "perpetua"
const DEFAULT_USER = "perpetua"

var BASE_DIR = filepath.Join(os.ExpandEnv("$HOME"), ".perpetua")
var CONFIG_FILE = filepath.Join(BASE_DIR, "perpetua.gcfg")
var DATABASE_FILE = filepath.Join(BASE_DIR, "perpetua.sqlite3")

type Options struct {
	Server struct {
		Hostname           string
		Port               uint16
		UseTLS, SkipVerify bool
	}
	IRC struct {
		Nickname, User string
		Channel        []string
	}
}

func (o *Options) Read() {

	err := gcfg.ReadFileInto(o, CONFIG_FILE)

	if o.IRC.Nickname == "" {
		o.IRC.Nickname = DEFAULT_NICKNAME
	}
	if o.IRC.User == "" {
		o.IRC.User = DEFAULT_USER
	}

	if err != nil {
		log.Fatal(err)
	}

}
