package server

import (
	"database/sql"
	"github.com/mimminou/BookIT-ByFood/back/models"
	"log"
	"net/http"
	"os"
	"testing"
)

var db *sql.DB
var dberr error

var numPages []int = []int{42, 352, 180, 0, 234, 310, 208, 0, 0, 200}
var someInt = 420
var pagesPointers = make([]*int, len(numPages)+10) //add a padding of 10, just so we have headroom to test

func setupMockDB() (*sql.DB, error) {
	// remove some of the num pages
	pagesPointers[2] = nil
	pagesPointers[5] = nil
	pagesPointers[9] = nil
	pagesPointers[14] = &someInt

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	schema := `CREATE TABLE IF NOT EXISTS Books (
    book_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    num_pages INTEGER,
    pub_date DATE NOT NULL
);`

	//apply schema
	db.Exec(schema)

	//creates some mock books
	books := []models.Book{
		{0, "To Kill a Mockingbird", "Harper Lee", nil, "1998-08-30T00:00:00Z"},
		{1, "1984", "George Orwell", nil, "1949-06-08"},
		{2, "The Great Gatsby", "F. Scott Fitzgerald", nil, "1925-04-10"},
		{3, "Pride and Prejudice", "Jane Austen", nil, "1813-01-28"},
		{4, "The Catcher in the Rye", "J.D. Salinger", nil, "1951-07-16"},
		{5, "The Hobbit", "J.R.R. Tolkien", nil, "1937-09-21"},
		{6, "To the Lighthouse", "Virginia Woolf", nil, "1927-05-05"},
		{7, "Moby-Dick", "Herman Melville", nil, "1851-10-18"},
		{8, "Frankenstein", "Mary Shelley", nil, "1818-01-01"},
		{9, "The Picture of Dorian Gray", "Oscar Wilde", nil, "1890-07-20"},
	}

	for i, book := range books {
		_, err := db.Exec("INSERT INTO Books (title, author, num_pages, pub_date) VALUES (?, ?, ?, ?)", book.Title, book.Author, pagesPointers[i], book.Pub_Date)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func SetupMockServer(db *sql.DB) *http.Server {
	serverMux := http.NewServeMux()
	dbRequestHandler := &Handler{db: db}
	serverMux.HandleFunc("/books", dbRequestHandler.Router)
	serverMux.HandleFunc("/books/", dbRequestHandler.Router)

	serverMuxMiddlware := ResponseLogging(serverMux)

	//Hardcoded port number, only for testing
	log.Println("Testing mock server on port", 50503)

	server := &http.Server{Addr: ":50503", Handler: serverMuxMiddlware}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	return server
}

func TestMain(m *testing.M) {
	db, dberr = setupMockDB()
	if dberr != nil {
		log.Fatal(dberr)
	}
	server := SetupMockServer(db)
	defer server.Close()
	defer db.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}
