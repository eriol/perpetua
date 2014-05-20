package perpetua

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func (s *Store) open() {
	db, err := sql.Open("sqlite3", DATABASE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
}

func (s *Store) close() {
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
	person_id INTEGER NOT NULL,
	quote TEXT,
	FOREIGN KEY(person_id) REFERENCES people(id)
);
	`
	s.db.Exec(sql_people_table)
	s.db.Exec(sql_quotes_table)
}
