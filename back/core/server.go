package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Handler struct {
	db *sql.DB
}

type ErrMessage struct {
	Msg string `json:"msg"`
}

// These are good enough for requirments of this project, but it's much better to use a routing Library like Gorilla or Chi
var BookIDRegex = regexp.MustCompile("^/books/[0-9]+(/)?$")
var BookRegex = regexp.MustCompile("^/books(/)?$")

// routes requests that have ID based on HTTP methods
func (handler *Handler) Router(w http.ResponseWriter, r *http.Request) {

	if !BookIDRegex.MatchString(r.URL.Path) && !BookRegex.MatchString(r.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if BookRegex.MatchString(r.URL.Path) {
		switch r.Method {
		case "POST":
			handler.Add(w, r)
		case "GET":
			handler.GetAll(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	//Here we are sure it matches BookID Regex
	switch r.Method {
	case "GET":
		handler.GetBook(w, r)
	case "PUT":
		handler.Update(w, r)
	case "DELETE":
		handler.Delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// get all records in DB
func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	books, err := GetBooks(handler.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrMessage{Msg: err.Error()})
		return
	}
	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrMessage{Msg: "No books found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

// get a single Book
func (handler *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Extract the ID from URI path
	parts := strings.Split(r.URL.Path, "/")
	stringID := parts[len(parts)-1]

	//this err check should be redudant, since the regex already filters anything that can't be converted to an int
	id, err := strconv.Atoi(stringID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fetchedBook, err := GetBook(handler.db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrMessage{Msg: "Book not found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedBook)

}

// add new book
func (handler *Handler) Add(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var book Book
	decodeErr := json.NewDecoder(r.Body).Decode(&book)

	if decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrMessage{Msg: decodeErr.Error()})
		return
	}

	emptyFields := CheckEmptyFields(book)
	if len(emptyFields) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrMessage{Msg: "Bad Request, missing value for fields: " + strings.Join(emptyFields, ", ")})
		return
	}

	err := AddBook(handler.db, book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrMessage{Msg: err.Error()})
		return

	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// delete existing artcile
func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	//Extract the ID from URI path
	parts := strings.Split(r.URL.Path, "/")
	stringID := parts[len(parts)-1]

	id, err := strconv.Atoi(stringID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delErr := DeleteBook(handler.db, id)
	if delErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrMessage{Msg: delErr.Error()})
		return

	}
	w.WriteHeader(http.StatusOK)
}

// update existing book
func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	//Extract the ID from URI path
	parts := strings.Split(r.URL.Path, "/")
	stringID := parts[len(parts)-1]

	id, err := strconv.Atoi(stringID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	book.Book_Id = id
	UpdateErr := UpdateBook(handler.db, book)
	if UpdateErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrMessage{Msg: UpdateErr.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// serve
func Serve(port uint16, db *sql.DB) {

	// simple server, no permission checks, no auth, only input Sanitazation and simple CRUD
	serverMux := http.NewServeMux()
	dbRequestHandler := &Handler{db: db}
	middleWareMux := Cors(Logging(serverMux))
	serverMux.HandleFunc("/books", dbRequestHandler.Router)
	serverMux.HandleFunc("/books/", dbRequestHandler.Router)

	fmt.Println("Serving on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), middleWareMux)
}
