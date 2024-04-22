package database

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"log"
)

// Connects to the database, returns a pointer to the db and an error value
func ConnectDb(dbname, path string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", path+dbname)
	if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}
