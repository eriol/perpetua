// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package irc

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

var connection *irc.Connection
var options *config.Options
var store *db.Store

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
	for _, channel := range options.IRC.Channels {
		connection.Join(channel)
		connection.Log.Println("Joined to " + channel)
	}
}

func doJoin(event *irc.Event) {
	channel := event.Arguments[0]

	if event.Nick == options.IRC.Nickname {
		connection.Privmsg(channel, "Hello! I'm "+version)
	} else {
		connection.Privmsg(channel,
			fmt.Sprintf("Hello %s! I'm %s. Do you want a quote?",
				event.Nick,
				version))
	}
}

func doPrivmsg(event *irc.Event) {
	channel := event.Arguments[0]
	var quote string

	// Don't speak in private!
	if channel == options.IRC.Nickname {
		return
	}
	command, person, extra, argument := parseMessage(event.Message())

	if command != "" && person != "" {

		quote = store.GetQuote(person, channel)

		if extra != "" && argument != "" {
			quote = store.GetQuoteAbout(person, argument, channel)
		}

		connection.Privmsg(channel, quote)
	}
}

func parseMessage(message string) (command, person, extra, argument string) {
	var names []string
	lang := options.I18N.Lang

	reArgument := regexp.MustCompile(options.IRC.Nickname +
		`:?` +
		`\s+` +
		`(?P<command>` + i18nKeyJoin(lang, "quote") + `)` +
		`\s+` +
		`(?P<person>[\w\s-'\p{Latin}]+)` +
		`(?:\s+)` +
		`(?P<extra>` + i18nKeyJoin(lang, "about") + `)` +
		`(?:\s+)` +
		`(?P<argument>[\w\s-'\p{Latin}]+)`)

	re := regexp.MustCompile(options.IRC.Nickname +
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

func Client(opt *config.Options, db *db.Store) {
	options = opt
	store = db
	connect()
	handleEvents()
	connection.Loop()
}
