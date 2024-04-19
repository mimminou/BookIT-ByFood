package server

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mimminou/BookIT-ByFood/back/models"
	"testing"
)

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
	err := AddBook(db, book)
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
		if err != nil {
			//This expects nil because DELETE in sqlite is a no op if it can't match an ID
			t.Errorf("Expected nil, got %v", err)
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
		book := models.Book{40, "LOTR", "Test Author", nil, "1951-12-18"}
		err := UpdateBook(db, book)
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
	})
}
