package server

//Server tests here

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mimminou/BookIT-ByFood/back/models"
)

// Testing API endpoints :
func TestGet(t *testing.T) {
	// Create a request to pass to the custom handler
	dbRequestHandler := &Handler{db: db}
	t.Run("Testing Get All Books", func(t *testing.T) {
		t.Log("Testing GET /books")
		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		dbRequestHandler.GetAll(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		//unmarshall json from body
		var books []models.Book
		jsonErr := json.NewDecoder(rr.Body).Decode(&books)
		if jsonErr != nil {
			t.Fatal(jsonErr)
		}
		if len(books) != 10 {
			t.Errorf("returned wrong number of books: got %v want %v",
				len(books), 10)
		}
	})

	t.Run("Testing Get Single Book", func(t *testing.T) {
		t.Log("Testing GET /books/2")
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
		t.Log("RESPONSE BODY : ", rr.Body.String())
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

		t.Log("RESPONSE BODY : ", rr.Body.String())
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})
}

func TestAdd(t *testing.T) {
	// Create a request to pass to the custom handler
	dbRequestHandler := &Handler{db: db}
	t.Run("Testing Add a Book", func(t *testing.T) {
		t.Log("Testing ADD /books")
		//create body and request

		body := `{"title": "1984", "author": "George Orwell", "num_pages": 328, "pub_date": "1949-06-08"}`
		t.Log("REQUEST BODY : ", body)
		req, err := http.NewRequest("POST", "/books", strings.NewReader(body))

		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// write json here
		dbRequestHandler.Add(rr, req)

		t.Log("RESPONSE BODY : ", rr.Body.String())
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}
		//unmarshall json from body
		var book models.Book
		jsonErr := json.NewDecoder(rr.Body).Decode(&book)
		if jsonErr != nil {
			t.Fatal(jsonErr)
		}
		if book.Title != "1984" {
			t.Errorf("returned wrong json: got %v want %v",
				book.Title, "1984")
		}
	})

	t.Run("Testing Add a Book with an empty mandatory field (title)", func(t *testing.T) {
		t.Log("Testing ADD /books")
		//create body and request

		body := `{"title": "", "author": "Some Author", "num_pages": 328, "pub_date": "1949-06-08"}`
		t.Log("REQUEST BODY : ", body)
		req, err := http.NewRequest("POST", "/books", strings.NewReader(body))

		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// write json here
		dbRequestHandler.Add(rr, req)

		t.Log("RESPONSE BODY : ", rr.Body.String())
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("Testing Add a Book with a missing mandatory field (title)", func(t *testing.T) {
		t.Log("Testing ADD /books")
		//create body and request

		body := `{"author": "Some Author", "num_pages": 328, "pub_date": "1949-06-08"}`
		t.Log("Post Request Body : ", body)
		req, err := http.NewRequest("POST", "/books", strings.NewReader(body))

		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// write json here
		dbRequestHandler.Add(rr, req)
		t.Log("RESPONSE BODY : ", rr.Body.String())

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("Testing Add a Book with an invalide date format (not YYYY-MM-DD)", func(t *testing.T) {
		t.Log("Testing ADD /books")
		//create body and request

		body := `{"title": "1984", "author": "Some Author", "num_pages": 328, "pub_date": "1949-some-string"}`
		t.Log("REQUEST BODY : ", body)
		req, err := http.NewRequest("POST", "/books", strings.NewReader(body))

		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// write json here
		dbRequestHandler.Add(rr, req)
		t.Log("RESPONSE BODY : ", rr.Body.String())

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}

func TestDelete(t *testing.T) {
	// Create a request to pass to the custom handler
	dbRequestHandler := &Handler{db: db}
	t.Run("Testing Delete a Book", func(t *testing.T) {
		t.Log("Testing DELETE /books/5")
		req, err := http.NewRequest("DELETE", "/books/5", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		dbRequestHandler.Delete(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("Testing Delete a non existent Book", func(t *testing.T) {
		t.Log("Testing DELETE /books/100")
		req, err := http.NewRequest("DELETE", "/books/100", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		dbRequestHandler.Delete(rr, req)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})
}

func TestUpdate(t *testing.T) {
	// Create a request to pass to the custom handler
	dbRequestHandler := &Handler{db: db}
	t.Run("Testing Update a Book", func(t *testing.T) {
		t.Log("Testing PUT /books/6")
		//create body and request
		body := `{"title": "1984", "author": "George Orwell", "num_pages": 328, "pub_date": "1949-06-08"}`
		t.Log("REQUEST BODY : ", body)
		req, err := http.NewRequest("PUT", "/books/6", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		dbRequestHandler.Update(rr, req)
		t.Log("RESPONSE BODY : ", rr.Body.String())
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("Testing Update a non existing Book", func(t *testing.T) {
		t.Log("Testing PUT /books/120")
		//create body and request
		body := `{"title": "1984", "author": "George Orwell", "num_pages": 328, "pub_date": "1949-06-08"}`
		t.Log("REQUEST BODY : ", body)
		req, err := http.NewRequest("PUT", "/books/120", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		dbRequestHandler.Update(rr, req)
		t.Log("RESPONSE BODY : ", rr.Body.String())
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})
}
