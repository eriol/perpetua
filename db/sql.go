package db

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

func (s *Store) Open(database string) {
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
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) createDatabase() {

	// TODO: create all tables in one variable when multiple statements will
	// be supported see https://github.com/mattn/go-sqlite3/issues/60
	sql_people_table := `CREATE TABLE people (
	id INTEGER NOT NULL PRIMARY KEY autoincrement,
	name TEXT COLLATE NOCASE
);`

	sql_quotes_table := `CREATE TABLE quotes (
	id INTEGER NOT NULL PRIMARY KEY autoincrement,
	person_id INTEGER NOT NULL,
	quote TEXT COLLATE NOCASE,
	FOREIGN KEY(person_id) REFERENCES people(id)
);`

	sql_alias_table := `CREATE TABLE alias (
	person_id INTEGER NOT NULL,
	alias TEXT COLLATE NOCASE,
	FOREIGN KEY(person_id) REFERENCES people(id)
);`

	sql_channels_table := `CREATE TABLE channels (
	id INTEGER NOT NULL PRIMARY KEY autoincrement,
	name TEXT COLLATE NOCASE
);`
	// quotes_acl table connects a quote with a specific channel so
	// that the quote will be quoted only in that specific channel.
	// Quotes not present in this table are public (quoted in every
	// channel where the bot joined)
	sql_quotes_acl_table := `CREATE TABLE quotes_acl (
	quote_id INTEGER NOT NULL,
	channel_id INTEGER NOT NULL,
	FOREIGN KEY(quote_id) REFERENCES quotes(id),
	FOREIGN KEY(channel_id) REFERENCES channels(id)
);`

	tables := [5]string{
		sql_people_table,
		sql_quotes_table,
		sql_alias_table,
		sql_channels_table,
		sql_quotes_acl_table,
	}

	for _, table := range tables {
		s.db.Exec(table)
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

func (s *Store) GetQuote(person string) (quote string) {
	query := `SELECT quote FROM quotes WHERE person_id = ?
		ORDER BY RANDOM() LIMIT 1;`
	s.db.QueryRow(query, s.getPerson(person)).Scan(&quote)
	return quote
}

func (s *Store) GetQuoteAbout(person, argument string) (quote string) {
	// A double quote can't be present in argument because of the
	// regex used but removing anyway
	argument = strings.Replace(argument, "\"", "", -1)

	query := "SELECT quote FROM quotes WHERE person_id = ? " +
		"AND quote LIKE \"%%%s%%\" " +
		"ORDER BY RANDOM() LIMIT 1;"
	query = fmt.Sprintf(query, argument)

	s.db.QueryRow(query, s.getPerson(person)).Scan(&quote)
	return quote
}
