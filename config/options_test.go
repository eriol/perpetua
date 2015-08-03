// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import "testing"

// Check if default values are set correctly.
func TestDefaultValues(t *testing.T) {

	var option Options
	cfg := ""

	option.ReadFromString(cfg)

	if option.IRC.Nickname != DEFAULT_NICKNAME {
		t.Error("Default nickname not set correctly!")
	}

	if option.IRC.User != DEFAULT_USER {
		t.Error("Default user not set correctly!")
	}

	if option.I18N.Lang != DEFAULT_LANG {
		t.Error("Default language not set correctly!")
	}
}

func TestChannels(t *testing.T) {

	var option Options
	cfg := `[IRC]
	channels = ["channel1", "#channel2"]
	`

	option.ReadFromString(cfg)

	if option.IRC.Channels[0] != "#channel1" &&
		option.IRC.Channels[1] != "#channel2" {
		t.Error("Channels not set correctly!")
	}

}

func TestConfigExample(t *testing.T) {

	var option Options
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

	option.ReadFromString(cfg)

	if option.Server.Hostname != "irc.example.org" {
		t.Error("Server hostname not set correctly!")
	}

	if option.Server.Port != 9999 {
		t.Error("Server port not set correctly!")
	}

	if option.Server.UseTLS != true {
		t.Error("Option useTSL not set correctly!")
	}

	if option.Server.SkipVerify != false {
		t.Error("Option skipVerify not set correctly!")
	}

	if option.IRC.Nickname != "perpetua-test" {
		t.Error("Bot nickname not set correctly!")
	}

	if option.IRC.Channels[0] != "#test" &&
		option.IRC.Channels[1] != "#test1" {
		t.Error("Channels not set correctly!")
	}

	if option.I18N.Lang != "it" {
		t.Error("Lang not set correctly!")
	}
}
