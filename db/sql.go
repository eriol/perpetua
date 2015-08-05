// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db // import "eriol.xyz/perpetua/db"

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func (s *Store) Open(database string) error {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	s.db = db

	if _, err := os.Stat(database); err != nil {
		if os.IsNotExist(err) {
			s.createDatabase()
		}
	}

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) createDatabase() {

	// quotes_acl table connects a quote with a specific channel so
	// that the quote will be quoted only in that specific channel.
	// Quotes not present in this table are public (quoted in every
	// channel where the bot joined)
	tables := `
	CREATE TABLE people (
	    id INTEGER NOT NULL PRIMARY KEY autoincrement,
	    name TEXT COLLATE NOCASE
	);

	CREATE TABLE quotes (
	    id INTEGER NOT NULL PRIMARY KEY autoincrement,
	    person_id INTEGER NOT NULL,
	    quote TEXT COLLATE NOCASE,
	    FOREIGN KEY(person_id) REFERENCES people(id)
	);

	CREATE TABLE alias (
	    person_id INTEGER NOT NULL,
	    alias TEXT COLLATE NOCASE,
	    FOREIGN KEY(person_id) REFERENCES people(id)
    );

	CREATE TABLE channels (
	    id INTEGER NOT NULL PRIMARY KEY autoincrement,
	    name TEXT COLLATE NOCASE
	);

	CREATE TABLE quotes_acl (
	    quote_id INTEGER NOT NULL,
	    channel_id INTEGER NOT NULL,
	    FOREIGN KEY(quote_id) REFERENCES quotes(id),
	    FOREIGN KEY(channel_id) REFERENCES channels(id)
	);`

	_, err := s.db.Exec(tables)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Store) getPerson(name string) (id int) {
	query_people := `SELECT id FROM people WHERE name = ?`
	query_alias := `SELECT person_id FROM alias WHERE alias = ?`

	s.db.QueryRow(query_people, name).Scan(&id)
	if id == 0 {
		s.db.QueryRow(query_alias, name).Scan(&id)
	}

	return
}

func (s *Store) getChannel(name string) (id int) {
	query := `SELECT id FROM channels WHERE name = ?`

	s.db.QueryRow(query, name).Scan(&id)

	return
}

func (s *Store) GetQuote(person, channel string) (quote string) {
	query := `SELECT quote FROM quotes WHERE person_id = ?
		AND id NOT IN (
			SELECT quote_id FROM quotes_acl
			EXCEPT
			SELECT quote_id FROM quotes_acl WHERE channel_id = ?)
		ORDER BY RANDOM() LIMIT 1;`

	s.db.QueryRow(
		query,
		s.getPerson(person),
		s.getChannel(channel)).Scan(&quote)

	return
}

func (s *Store) GetQuoteAbout(person, argument, channel string) (quote string) {
	// A double quote can't be present in argument because of the
	// regex used but removing anyway
	argument = strings.Replace(argument, "\"", "", -1)

	query := "SELECT quote FROM quotes WHERE person_id = ? " +
		"AND quote LIKE \"%%%s%%\" " +
		"AND id NOT IN ( " +
		"SELECT quote_id FROM quotes_acl " +
		"EXCEPT " +
		"SELECT quote_id FROM quotes_acl WHERE channel_id = ?) " +
		"ORDER BY RANDOM() LIMIT 1;"
	query = fmt.Sprintf(query, argument)

	s.db.QueryRow(
		query,
		s.getPerson(person),
		s.getChannel(channel)).Scan(&quote)

	return
}
