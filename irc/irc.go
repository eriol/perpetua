package irc

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"

	"github.com/thoj/go-ircevent"

	"hg.mornie.org/perpetua/config"
	"hg.mornie.org/perpetua/db"
)

const version = "perpetua quote bot v0.1a"

var connection *irc.Connection
var options *config.Options
var store *db.Store

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
	connection.AddCallback("PRIVMSG", doPrivmsg)
}

func doWelcome(event *irc.Event) {
	connection.Join(options.IRC.Channel)
}

func doJoin(event *irc.Event) {

	if event.Nick == options.IRC.Nickname {
		connection.Privmsg(event.Arguments[0], "Hello! I'm "+version)
	} else {

		connection.Privmsg(event.Arguments[0],
			fmt.Sprintf("Hello %s! I'm %s. Do you want a quote?",
				event.Nick,
				version))
	}
}

func doPrivmsg(event *irc.Event) {
	channel := event.Arguments[0]

	// Don't speak in private!
	if channel == options.IRC.Nickname {
		return
	}
	command, person := parseMessage(event.Message())

	if command != "" && person != "" {
		connection.Privmsg(event.Arguments[0], store.GetQuote(person))
	}
}

func parseMessage(message string) (command, person string) {

	re := regexp.MustCompile(options.IRC.Nickname +
		`:?` +
		`\s*` +
		`(?P<command>cita|cosa dice|quote|what does it say)` +
		`\s*(?P<person>[\w\s-']+)`)

	res := re.FindStringSubmatch(message)

	names := re.SubexpNames()
	md := map[string]string{}
	for i, n := range res {
		md[names[i]] = n
	}

	return md["command"], md["person"]
}

func Client(opt *config.Options, db *db.Store) {
	options = opt
	store = db
	connect()
	handleEvents()
	connection.Loop()
}
