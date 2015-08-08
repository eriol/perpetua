// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "eriol.xyz/perpetua"

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	"eriol.xyz/perpetua/config"
	"eriol.xyz/perpetua/db"
	"eriol.xyz/perpetua/irc"
)

func main() {

	var (
		conf  config.Config
		store db.Store
	)
	var (
		configFile = kingpin.Flag("config", "Configuration file path.").Short('c').Default("").String()
	)

	kingpin.Version(config.Version)
	kingpin.CommandLine.Help = "Quote bot for IRC."
	kingpin.Parse()

	if err := conf.Read(*configFile); err != nil {
		log.Fatal(err)
	}

	if err := store.Open(config.DATABASE_FILE); err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	irc.Client(&conf, &store)

}
