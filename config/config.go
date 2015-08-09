// Copyright © 2014-2015 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config // import "eriol.xyz/perpetua/config"

import (
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

// Config is used to store information about the configuration of the bot.
type Config struct {
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
func (c *Config) Read(configFile string) (err error) {

	if configFile == "" {
		configFile = CONFIG_FILE
	}

	if _, err := toml.DecodeFile(configFile, c); err != nil {
		return err
	}

	c.setDefaultValues()

	return nil
}

// Read configuration from string.
func (c *Config) ReadFromString(config string) (err error) {

	if _, err := toml.Decode(config, c); err != nil {
		return err
	}

	c.setDefaultValues()

	return nil
}

// Set default values for not provided entries.
func (c *Config) setDefaultValues() {

	if c.IRC.Nickname == "" {
		c.IRC.Nickname = DEFAULT_NICKNAME
	}
	if c.IRC.User == "" {
		c.IRC.User = DEFAULT_USER
	}

	if c.I18N.Lang == "" {
		c.I18N.Lang = DEFAULT_LANG
	}

	// Add a # at the beginning of the channel name if it's not there yet.
	for i, channel := range c.IRC.Channels {
		if string(channel[0]) == "#" {
			continue
		} else {
			c.IRC.Channels[i] = "#" + channel
		}
	}
}
