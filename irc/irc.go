// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package irc // import "eriol.xyz/perpetua/irc"

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/thoj/go-ircevent"

	"eriol.xyz/perpetua/config"
	"eriol.xyz/perpetua/db"
)

const version = "perpetua quote bot " + config.Version

var (
	conf  *config.Config
	store *db.Store
)

// Localizated quote and about tokens used to detect the kind of query for
// the bot.
var i18n = map[string]map[string][]string{
	"en": map[string][]string{
		"quote": []string{"quote", "what does it say"},
		"about": []string{"about"},
	},
	"it": map[string][]string{
		"quote": []string{
			"cita",
			"che dice",
			"cosa dice",
			"che cosa dice"},
		"about": []string{
			"su",
			"sul",
			"sulla",
			"sullo",
			"sui",
			"sugli",
			"sulle"},
	},
}

// Join keys from i18n using "|": used inside the regex to perform an
// OR of all keys.
func i18nKeyJoin(lang, key string) string {
	return strings.Join(i18n[lang][key], "|")
}

func connect() (connection *irc.Connection, err error) {
	connection = irc.IRC(conf.IRC.Nickname, conf.IRC.User)
	connection.Version = version
	connection.UseTLS = conf.Server.UseTLS
	if conf.Server.SkipVerify == true {
		connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := connection.Connect(fmt.Sprintf("%s:%d",
		conf.Server.Hostname,
		conf.Server.Port)); err != nil {
		return nil, err
	}

	return connection, nil
}

func doWelcome(event *irc.Event) {
	for _, channel := range conf.IRC.Channels {
		event.Connection.Join(channel)
		event.Connection.Log.Println("Joined to " + channel)
	}
}

func doJoin(event *irc.Event) {
	channel := event.Arguments[0]

	if event.Nick == conf.IRC.Nickname {
		event.Connection.Privmsg(channel, "Hello! I'm "+version)
	} else {
		event.Connection.Privmsg(channel,
			fmt.Sprintf("Hello %s! I'm %s. Do you want a quote?",
				event.Nick,
				version))
	}
}

func doPrivmsg(event *irc.Event) {
	channel := event.Arguments[0]
	var quote string

	// Don't speak in private!
	if channel == conf.IRC.Nickname {
		return
	}
	command, person, extra, argument := parseMessage(event.Message())

	if command != "" && person != "" {

		quote = store.GetQuote(person, channel)

		if extra != "" && argument != "" {
			quote = store.GetQuoteAbout(person, argument, channel)
		}

		event.Connection.Privmsg(channel, quote)
	}
}

func parseMessage(message string) (command, person, extra, argument string) {
	var names []string
	lang := conf.I18N.Lang

	reArgument := regexp.MustCompile(conf.IRC.Nickname +
		`:?` +
		`\s+` +
		`(?P<command>` + i18nKeyJoin(lang, "quote") + `)` +
		`\s+` +
		`(?P<person>[\w\s-'\p{Latin}]+)` +
		`(?:\s+)` +
		`(?P<extra>` + i18nKeyJoin(lang, "about") + `)` +
		`(?:\s+)` +
		`(?P<argument>[\w\s-'\p{Latin}]+)`)

	re := regexp.MustCompile(conf.IRC.Nickname +
		`:?` +
		`\s+` +
		`(?P<command>` + i18nKeyJoin(lang, "quote") + `)` +
		`\s+` +
		`(?P<person>[\w\s-'\p{Latin}]+)`)

	res := reArgument.FindStringSubmatch(message)

	if res == nil {
		res = re.FindStringSubmatch(message)
		names = re.SubexpNames()
	} else {
		names = reArgument.SubexpNames()
	}

	m := map[string]string{}
	for i, n := range res {
		m[names[i]] = n
	}

	return m["command"], m["person"], m["extra"], m["argument"]
}

func NewClient(c *config.Config, db *db.Store, ircChan chan *irc.Connection, done chan bool) {
	conf = c
	store = db

	connection, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	ircChan <- connection

	connection.AddCallback("001", doWelcome)
	connection.AddCallback("JOIN", doJoin)
	connection.AddCallback("PRIVMSG", doPrivmsg)

	connection.Loop()

	done <- true
}
