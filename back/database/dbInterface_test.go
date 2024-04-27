package database

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"github.com/mimminou/BookIT-ByFood/back/models"
	"log"
	"os"
	"testing"
)

var db *sql.DB
var dberr error

var numPages []int = []int{42, 352, 180, 0, 234, 310, 208, 0, 0, 200}
var someInt = 420
var pagesPointers = make([]*int, len(numPages)+10) //add a padding of 10, just so we have headroom to test

func TestMain(m *testing.M) {
	db, dberr = setupMockDB()
	if dberr != nil {
		log.Fatal(dberr)
	}
	defer db.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockDB() (*sql.DB, error) {
	// remove some of the num pages
	pagesPointers[2] = nil
	pagesPointers[5] = nil
	pagesPointers[9] = nil
	pagesPointers[14] = &someInt

	db, err := sql.Open("sqlite", ":memory:")
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

func TestControllerGet(t *testing.T) {
	if dberr != nil {
		t.Fatal(dberr)
	}

	t.Run("Testing Get All Books", func(t *testing.T) {
		books, err := GetBooks(db)
		if err != nil {
			t.Fatal(err)
		}
		if len(books) != 10 {
			t.Errorf("Expected 10 books, got %d", len(books))
		}

	})

	t.Run("Testing Get Single Book", func(t *testing.T) {
		book, err := GetBook(db, 2) // sqlite Ids start with 1
		if err != nil {
			t.Fatal(err)
		}
		if book.Title != "1984" {
			t.Errorf("Expected '1984', got %s", book.Title)
		}
	})

	t.Run("Testing getting non existing book", func(t *testing.T) {
		book, err := GetBook(db, 100)
		if err == nil {
			t.Errorf("Expected error, got %v", book)
		}
	})
}

func TestControllerAdd(t *testing.T) {
	if dberr != nil {
		t.Fatal(dberr)
	}

	book := models.Book{10, "The Hobbit", "J.R.R. Tolkien", pagesPointers[14], "1937-09-21"}
	_, err := AddBook(db, book)
	if err != nil {
		t.Fatal(err)
	}
}

func TestControllerDelete(t *testing.T) {
	if dberr != nil {
		t.Fatal(dberr)
	}

	t.Run("Testing Delete Book", func(t *testing.T) {
		err := DeleteBook(db, 4)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Testing deleting non existing book", func(t *testing.T) {
		err := DeleteBook(db, 60)
		if err != sql.ErrNoRows {
			t.Errorf("Expected NoRows error, got %v", err)
		}
	})
}

func TestControllerUpdate(t *testing.T) {
	if dberr != nil {
		t.Fatal(dberr)
	}

	t.Run("Testing Update Book", func(t *testing.T) {
		book := models.Book{7, "Moby-Dick", "Test Author", nil, "1851-12-18"}
		err := UpdateBook(db, book)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Testing update non Existing book", func(t *testing.T) {
		book := models.Book{50, "LOTR", "Test Author", nil, "1951-12-18"}
		err := UpdateBook(db, book)
		if err != sql.ErrNoRows {
			t.Errorf("Expected NoRows error, got %v", err)
		}
	})
}
