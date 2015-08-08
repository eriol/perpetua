// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "eriol.xyz/perpetua"

import (
	"log"
	"os"
	"os/signal"

	ircevent "github.com/thoj/go-ircevent"
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

	isDone := make(chan bool, 1)
	ircChan := make(chan *ircevent.Connection, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

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

	log.Println("Starting...")
	go irc.Client(&conf, &store, ircChan, isDone)

	sig := <-sigChan
	log.Printf("Got signal %v, exiting now.\n", sig)

	conn := <-ircChan
	conn.Quit()

	done := <-isDone

	if done {
		log.Println("Stopping...")
	}

}
