package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mimminou/BookIT-ByFood/back/core"
)

// config struct for config.json
type server_config struct {
	Port uint16 `json:"port"`
}

type db_config struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type config struct {
	Server server_config `json:"server"`
	Db     db_config     `json:"database"`
}

// Print small help message that demonstrates usage
func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("-h : Prints this help message and exits")
	fmt.Println("-s : Setup DB and exits, needs to run once before running the server first time")
}

// reads json config from file, returns pointer to config struct
func readConfig(filename string) (*config, error) {
	var config config
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading config.json: ", err)
		return nil, err
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("Error reading config.json: ", err)
		return nil, err
	}
	return &config, nil
}

// Handles command line args
func handleArgs() {
	if len(os.Args) == 1 {
		// no args provided, run server without setup
		return
	}
	if len(os.Args) > 2 {
		fmt.Println("Error: Too many arguments, please only provide one arguement")
		showHelp()
		os.Exit(3)
	}

	switch os.Args[1] {
	case "-s":

		db := SetupDatabase()
		if db != nil {
			fmt.Println("Error setting up database: ", db)
			os.Exit(5)
		}
	case "-h":
		showHelp()
		os.Exit(0)
	default:
		fmt.Println("Error: Invalid argument")
		showHelp()
		os.Exit(4)
	}
}

func checkDB(name, path string) {
	if _, err := os.Stat(path + name); err == nil {
		return
	} else {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("Error : Db does not exist yet, please run the server with '-s' flag to setup the DB first")
		} else {
			log.Fatal("Error : ", err)
		}
	}

}

func main() {
	// Attempt to read the config.json file
	config, err := readConfig("config.json")
	if err != nil {
		if err == os.ErrNotExist {
			fmt.Println("Error reading config.json: File does not exist")
			os.Exit(1)
		}
		fmt.Println("Error reading config.json:", err)
		os.Exit(2)
	}

	handleArgs()

	checkDB(config.Db.Name, config.Db.Path)

	// connect to DB, then pass DB instance to server
	db, err := core.ConnectDb(config.Db.Name, config.Db.Path)
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		os.Exit(1)
	}
	defer db.Close()

	core.Serve(config.Server.Port, db)
}
