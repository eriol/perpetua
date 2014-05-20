package perpetua

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/thoj/go-ircevent"
)

const version = "perpetua quote bot v0.1a"

var connection *irc.Connection

func connect() {
	connection = irc.IRC(options.IRC.Nickname, options.IRC.User)
	connection.Version = version
	connection.UseTLS = options.Server.UseTLS
	if options.Server.SkipVerify == true {
		connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	err := connection.Connect(fmt.Sprintf("%s:%d",
		options.Server.Hostname,
		options.Server.Port))

	if err != nil {
		log.Fatal(err)
	}
}

func handleEvents() {
	connection.AddCallback("001", doWelcome)
	connection.AddCallback("JOIN", doJoin)
}

func doWelcome(event *irc.Event) {
	connection.Join("#" + options.IRC.Channel)
}

func doJoin(event *irc.Event) {
	connection.Privmsg(event.Arguments[0], "Hello! I'm "+version)
}

func startIRC() {
	connect()
	handleEvents()
	connection.Loop()
}
