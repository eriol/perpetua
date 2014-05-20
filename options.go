package perpetua

import (
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/gcfg"
)

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
		Nickname, User, Channel string
	}
}

func (o *Options) read() {

	err := gcfg.ReadFileInto(o, CONFIG_FILE)

	if err != nil {
		log.Fatal(err)
	}

}
