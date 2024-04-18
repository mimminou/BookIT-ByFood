package core

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Connects to the database, returns a pointer to the db and an error value
func ConnectDb(dbname, path string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", path+dbname)
	if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}
