package services

//Server tests here

import (
	"database/sql"
	"encoding/json"
	"github.com/mimminou/BookIT-ByFood/back/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Setup Mock DB

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

// Testing API endpoints :
func TestGet(t *testing.T) {
	// Create a request to pass to the custom handler
	dbRequestHandler := &DBRequestHandler{Db: db}
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
	dbRequestHandler := &DBRequestHandler{Db: db}
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

	t.Run("Testing Add a Book with an invalide date but correct Format (not YYYY-MM-DD)", func(t *testing.T) {
		t.Log("Testing ADD /books")
		//create body and request

		body := `{"title": "1984", "author": "Some Author", "num_pages": 328, "pub_date": "1949-11-36"}`
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
	dbRequestHandler := &DBRequestHandler{Db: db}
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
	dbRequestHandler := &DBRequestHandler{Db: db}
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

	t.Run("Testing Update a book with wrong date format", func(t *testing.T) {
		t.Log("Testing PUT /books/120")
		//create body and request
		body := `{"title": "1984", "author": "George Orwell", "num_pages": 328, "pub_date": "1949-16-08"}`
		t.Log("REQUEST BODY : ", body)
		req, err := http.NewRequest("PUT", "/books/120", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		dbRequestHandler.Update(rr, req)
		t.Log("RESPONSE BODY : ", rr.Body.String())
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

}
