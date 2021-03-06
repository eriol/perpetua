// Copyright © 2014-2015 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db // import "eriol.xyz/perpetua/db"

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

// Open SQLite3 database specified in database path.
func (s *Store) Open(database string) (err error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}
	s.db = db

	if _, err := os.Stat(database); err != nil {
		if os.IsNotExist(err) {
			if err := s.createDatabase(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Close SQLite3 database used by Store.
func (s *Store) Close() error {
	return s.db.Close()
}

// Create database tables.
func (s *Store) createDatabase() error {

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
		return err
	}

	return nil
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
	argument = fmt.Sprintf("%%%s%%", argument)

	query := `SELECT quote FROM quotes WHERE person_id = ?
		AND quote LIKE ?
		AND id NOT IN (
		SELECT quote_id FROM quotes_acl
		EXCEPT
		SELECT quote_id FROM quotes_acl WHERE channel_id = ?)
		ORDER BY RANDOM() LIMIT 1`

	s.db.QueryRow(
		query,
		s.getPerson(person),
		argument,
		s.getChannel(channel)).Scan(&quote)

	return
}
