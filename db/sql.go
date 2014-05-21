package db

import (
	"database/sql"
	"log"
	"os"

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
	name TEXT
);
	`
	sql_quotes_table := `CREATE TABLE quotes (
	id INTEGER NOT NULL PRIMARY KEY autoincrement,
	person_id INTEGER NOT NULL,
	quote TEXT,
	FOREIGN KEY(person_id) REFERENCES people(id)
);
	`
	s.db.Exec(sql_people_table)
	s.db.Exec(sql_quotes_table)
}

func (s *Store) GetQuote(person string) (quote string) {
	var q string
	query := `SELECT quote FROM quotes WHERE person_id =
		(SELECT id FROM people WHERE name = ?) ORDER BY RANDOM() LIMIT 1;`
	s.db.QueryRow(query, person).Scan(&q)
	return q
}
