// Copyright © 2014 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// Insert some data for testing.
func fixture(databaseName string) {
	db, _ := sql.Open("sqlite3", databaseName)
	defer db.Close()

	db.Exec("INSERT INTO people VALUES (1, 'the fox')")
	db.Exec("INSERT INTO alias VALUES (1, 'volpe')")
	db.Exec("INSERT INTO channels VALUES (1, '#test')")
	db.Exec("INSERT INTO channels VALUES (2, '#test2')")
	db.Exec("INSERT INTO channels VALUES (3, '#test3')")
	db.Exec("INSERT INTO channels VALUES (4, '#test4')")
	db.Exec("INSERT INTO quotes VALUES (1, 1, 'Hatee-hatee-hatee-ho!')")
	db.Exec("INSERT INTO quotes VALUES (2, 1, 'Chacha-chacha-chacha-chow!')")
	db.Exec("INSERT INTO quotes VALUES (3, 1, 'A-oo-oo-oo-ooo!')")
	db.Exec("INSERT INTO quotes VALUES (4, 1, 'Wa-pa-pa-pa-pa-pa-pow!')")
	db.Exec("INSERT INTO quotes_acl VALUES (1, 1)")
	db.Exec("INSERT INTO quotes_acl VALUES (2, 2)")
	db.Exec("INSERT INTO quotes_acl VALUES (3, 4)")
	db.Exec("INSERT INTO quotes_acl VALUES (4, 4)")
}

func TestQuote(t *testing.T) {
	var s Store
	var quote string
	var id int

	dirName, err := ioutil.TempDir("", "perpetua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirName)

	databaseName := path.Join(dirName, "perpetua.sqlite3")

	err = s.Open(databaseName)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	fixture(databaseName)

	id = s.getPerson("the fox")
	if id != 1 {
		t.Fatal("Person id does not match!")
	}

	// Person search is case insensitive
	id = s.getPerson("ThE fOx")
	if id != 1 {
		t.Fatal("Person id does not match!")
	}

	id = s.getChannel("#test2")
	if id != 2 {
		t.Fatal("Channel id does not match!")
	}

	// GetQuote tests

	quote = s.GetQuote("the fox", "#test")
	if quote != "Hatee-hatee-hatee-ho!" {
		t.Fatal("Quote does not match!")
	}

	quote = s.GetQuote("volpe", "#test")
	if quote != "Hatee-hatee-hatee-ho!" {
		t.Fatal("Quote does not match!")
	}

	quote = s.GetQuote("the fox", "#test2")
	if quote != "Chacha-chacha-chacha-chow!" {
		t.Fatal("Quote does not match!")
	}

	quote = s.GetQuote("the fox", "#test3")
	if quote != "" {
		t.Fatal("Quote does not match!")
	}

	// GetQuoteAbout tests
	// Argument search is case insensitive
	quote = s.GetQuoteAbout("the fox", "wa", "#test4")
	if quote != "Wa-pa-pa-pa-pa-pa-pow!" {
		t.Fatal("Quote does not match!")
	}

	quote = s.GetQuoteAbout("the fox", "OOO", "#test4")
	if quote != "A-oo-oo-oo-ooo!" {
		t.Fatal("Quote does not match!")
	}

	quote = s.GetQuoteAbout("volpe", "OoO", "#test4")
	if quote != "A-oo-oo-oo-ooo!" {
		t.Fatal("Quote does not match!")
	}

	quote = s.GetQuoteAbout("the fox", "wa", "#test")
	if quote != "" {
		t.Fatal("Quote does not match!")
	}
}
