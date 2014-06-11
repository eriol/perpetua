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

// Options is used by Gcfg to store data read from CONFIG_FILE or a string.
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

// Read configuration from file specified by configFile and use
// the default config file CONFIG_FILE if configFile is empty.
func (o *Options) Read(configFile string) {

	if configFile == "" {
		configFile = CONFIG_FILE
	}

	err := gcfg.ReadFileInto(o, configFile)

	if err != nil {
		log.Fatal(err)
	}

	o.setDefaultValues()
}

// Read configuration from string.
func (o *Options) ReadFromString(config string) {

	err := gcfg.ReadStringInto(o, config)

	if err != nil {
		log.Fatal(err)
	}

	o.setDefaultValues()
}

// Set default values for not provided entries.
func (o *Options) setDefaultValues() {

	if o.IRC.Nickname == "" {
		o.IRC.Nickname = DEFAULT_NICKNAME
	}
	if o.IRC.User == "" {
		o.IRC.User = DEFAULT_USER
	}

	if o.I18N.Lang == "" {
		o.I18N.Lang = DEFAULT_LANG
	}

	// Add a # at the beginning of the channel name if it's not there yet.
	// gcfg use # for comments so if you want to insert a # you must enclose
	// channel inside double quote marks.
	for i, channel := range o.IRC.Channel {
		if string(channel[0]) == "#" {
			continue
		} else {
			o.IRC.Channel[i] = "#" + channel
		}
	}

}
