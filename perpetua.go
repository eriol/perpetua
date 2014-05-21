package main

import (
	"hg.mornie.org/perpetua/config"
	"hg.mornie.org/perpetua/db"
	"hg.mornie.org/perpetua/irc"
)

func main() {

	var options config.Options
	var store db.Store

	options.Read()

	store.Open(config.DATABASE_FILE)
	// defer store.Close()

	irc.Client(&options, &store)

}
