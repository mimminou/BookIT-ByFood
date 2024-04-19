package server

//Server tests here

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mimminou/BookIT-ByFood/back/models"
)

// Testing API endpoints :
func TestGet(t *testing.T) {
	// Create a request to pass to the custom handler
	dbRequestHandler := &Handler{db: db}
	t.Run("Testing Get All Books", func(t *testing.T) {
		log.Print("Testing GET /books")
		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		dbRequestHandler.GetAll(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		//unmarshall json from body
		var books []models.Book
		jsonErr := json.NewDecoder(rr.Body).Decode(&books)
		if jsonErr != nil {
			t.Fatal(jsonErr)
		}
		if len(books) != 10 {
			t.Errorf("handler returned wrong number of books: got %v want %v",
				len(books), 10)
		}
	})

	t.Run("Testing Get Single Book", func(t *testing.T) {
		log.Print("Testing GET /books/2")
		req, err := http.NewRequest("GET", "/books/2", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		dbRequestHandler.GetBook(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		//unmarshall json from body
		var book models.Book
		jsonErr := json.NewDecoder(rr.Body).Decode(&book)
		t.Log("Response: ", rr.Body.String())
		if jsonErr != nil {
			t.Fatal(jsonErr)
		}
		if book.Title != "1984" {
			t.Errorf("returned wrong book: got %v want %v",
				book.Title, "1984")
		}
	})

	//Testing GetBook with an ID that does not exist
	t.Run("Testing Get Book with non existent ID", func(t *testing.T) {
		t.Log("Testing GET /books/100")
		req, err := http.NewRequest("GET", "/books/100", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		dbRequestHandler.GetBook(rr, req)
		t.Log("Response: ", rr.Body.String())
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})
}

// TODO: Add tests for add, update, delete
