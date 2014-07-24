// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package irc

import "testing"

func TestI18nKeyJoin(t *testing.T) {
	var s string

	s = i18nKeyJoin("en", "quote")
	if s != "quote|what does it say" {
		t.Fatal("quote key for en language error")
	}
	s = i18nKeyJoin("en", "about")
	if s != "about" {
		t.Fatal("about key for en language error")
	}

	s = i18nKeyJoin("it", "quote")
	if s != "cita|che dice|cosa dice|che cosa dice" {
		t.Fatal("quote key for it language error")
	}
	s = i18nKeyJoin("it", "about")
	if s != "su|sul|sulla|sullo|sui|sugli|sulle" {
		t.Fatal("about key for it language error")
	}
}
