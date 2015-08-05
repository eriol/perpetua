// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "eriol.xyz/perpetua"

import (
	"eriol.xyz/perpetua/config"
	"eriol.xyz/perpetua/db"
	"eriol.xyz/perpetua/irc"
)

func main() {

	var conf config.Config
	var store db.Store

	// TODO add a command line option to specify a config file
	conf.Read("")

	store.Open(config.DATABASE_FILE)
	defer store.Close()

	irc.Client(&conf, &store)

}
