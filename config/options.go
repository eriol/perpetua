// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

const (
	DEFAULT_LANG     = "en"
	DEFAULT_NICKNAME = "perpetua"
	DEFAULT_USER     = "perpetua"
	Version          = "0.1a"
)

var (
	BASE_DIR      = path.Join(os.ExpandEnv("$HOME"), ".perpetua")
	CONFIG_FILE   = path.Join(BASE_DIR, "perpetua.toml")
	DATABASE_FILE = path.Join(BASE_DIR, "perpetua.sqlite3")
)

// Options is used to store data read from CONFIG_FILE or a string.
type Options struct {
	Server struct {
		Hostname           string
		Port               uint16
		UseTLS, SkipVerify bool
	}
	IRC struct {
		Nickname, User string
		Channels       []string
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

	if _, err := toml.DecodeFile(configFile, o); err != nil {
		log.Fatal(err)
	}

	o.setDefaultValues()
}

// Read configuration from string.
func (o *Options) ReadFromString(config string) {

	if _, err := toml.Decode(config, o); err != nil {
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
	for i, channel := range o.IRC.Channels {
		if string(channel[0]) == "#" {
			continue
		} else {
			o.IRC.Channels[i] = "#" + channel
		}
	}
}
