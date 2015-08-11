// Copyright © 2014-2015 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI18nKeyJoin(t *testing.T) {
	var s string

	s = i18nKeyJoin("en", "quote")
	assert.Equal(t, s, "quote|what does it say", "quote key for en language error")

	s = i18nKeyJoin("en", "about")
	assert.Equal(t, s, "about", "about key for en language error")

	s = i18nKeyJoin("it", "quote")
	assert.Equal(t, s, "cita|che dice|cosa dice|che cosa dice",
		"quote key for it language error")

	s = i18nKeyJoin("it", "about")
	assert.Equal(t, s, "su|sul|sulla|sullo|sui|sugli|sulle",
		"quote key for it language error")
}
