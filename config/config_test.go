// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import "testing"

// Check if default values are set correctly.
func TestDefaultValues(t *testing.T) {

	var conf Config
	cfg := ""

	conf.ReadFromString(cfg)

	if conf.IRC.Nickname != DEFAULT_NICKNAME {
		t.Error("Default nickname not set correctly!")
	}

	if conf.IRC.User != DEFAULT_USER {
		t.Error("Default user not set correctly!")
	}

	if conf.I18N.Lang != DEFAULT_LANG {
		t.Error("Default language not set correctly!")
	}
}

func TestChannels(t *testing.T) {

	var conf Config
	cfg := `[IRC]
	channels = ["channel1", "#channel2"]
	`

	conf.ReadFromString(cfg)

	if conf.IRC.Channels[0] != "#channel1" && conf.IRC.Channels[1] != "#channel2" {
		t.Error("Channels not set correctly!")
	}

}

func TestConfigExample(t *testing.T) {

	var conf Config
	cfg := `
	[Server]
	hostname = "irc.example.org"
	port = 9999
	useTLS = true
	skipVerify = false
	[IRC]
	nickname = "perpetua-test"
	channels = ["test", "#test1"]
	[I18N]
	lang = "it"
	`

	conf.ReadFromString(cfg)

	if conf.Server.Hostname != "irc.example.org" {
		t.Error("Server hostname not set correctly!")
	}

	if conf.Server.Port != 9999 {
		t.Error("Server port not set correctly!")
	}

	if conf.Server.UseTLS != true {
		t.Error("Option useTSL not set correctly!")
	}

	if conf.Server.SkipVerify != false {
		t.Error("Option skipVerify not set correctly!")
	}

	if conf.IRC.Nickname != "perpetua-test" {
		t.Error("Bot nickname not set correctly!")
	}

	if conf.IRC.Channels[0] != "#test" &&
		conf.IRC.Channels[1] != "#test1" {
		t.Error("Channels not set correctly!")
	}

	if conf.I18N.Lang != "it" {
		t.Error("Lang not set correctly!")
	}
}
