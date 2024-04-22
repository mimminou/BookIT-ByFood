package main

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"os"
)

// SetupDB Will create the necessary sqlite DB file with the provided schema
func SetupDatabase() error {

	config, err := readConfig("config.json")
	if err != nil {
		if err == os.ErrNotExist {
			log.Fatalln("Error reading config.json: File does not exist")
		}
		log.Fatalln("Error reading config.json:", err)
	}

	path := config.Db.Path
	dbname := config.Db.Name

	//Load the DB Schema
	schema, err := os.ReadFile(path + "schema.sql")
	if err != nil {
		log.Fatal(err)
		return err
	}

	//Connect to Sqlite DB, if it doesn't exist it will be created
	db, err := sql.Open("sqlite", path+dbname)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	//Create DB Structure from Schema

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalln(err)
	}

	//DB Setup Complete
	log.Println("DB Setup Complete")
	os.Exit(0)
	return nil
}
