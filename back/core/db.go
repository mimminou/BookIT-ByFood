package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Connects to the database, returns a pointer to the db and an error value
func ConnectDb(dbname, path string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", path+dbname)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	return db, nil
}

func SetupDB(dbname, path string) error {
	schema, err := os.ReadFile(path + "schema.sql")
	if err != nil {
		log.Fatal(err)
		return err
	}

	db, err := sql.Open("sqlite3", path+dbname)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalln(err)
	}

	//dB setup done
	fmt.Println("DB setup done")

	return nil
}
