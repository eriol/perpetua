// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/gcfg"
)

const Version = "0.1a"

const DEFAULT_LANG = "en"

const DEFAULT_NICKNAME = "perpetua"
const DEFAULT_USER = "perpetua"

var BASE_DIR = filepath.Join(os.ExpandEnv("$HOME"), ".perpetua")
var CONFIG_FILE = filepath.Join(BASE_DIR, "perpetua.gcfg")
var DATABASE_FILE = filepath.Join(BASE_DIR, "perpetua.sqlite3")

// Options is used by Gcfg to store data read from CONFIG_FILE.
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
	I18N struct {
		Lang string
	}
}

// Read configuration from default config file specified by
// CONFIG_FILE and set default values for not provided entries.
func (o *Options) Read() {

	err := gcfg.ReadFileInto(o, CONFIG_FILE)

	if err != nil {
		log.Fatal(err)
	}

	if o.IRC.Nickname == "" {
		o.IRC.Nickname = DEFAULT_NICKNAME
	}
	if o.IRC.User == "" {
		o.IRC.User = DEFAULT_USER
	}

	if o.I18N.Lang == "" {
		o.I18N.Lang = DEFAULT_LANG
	}

}
