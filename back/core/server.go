package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	db *sql.DB
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// get with offset
func (handler *Handler) GetOffset(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	var (
		defaultLimit  = 10
		defaultOffset = 0
	)

	query := r.URL.Query()
	limitStr, offsetStr := query.Get("limit"), query.Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	books, err := GetOffsetBooks(handler.db, limit, offset)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(books)
}

// get all records in DB
func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	books, err := GetBooks(handler.db)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(books)
}

// get a single Book
func (handler *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
		// json error
		w.WriteHeader(http.StatusBadRequest)
	}

	thisBook, err := GetBook(handler.db, book.Book_Id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(thisBook)
}

// add new book
func (handler *Handler) Add(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// it is possible to check if designation already exists, then if it does just call Update() and incremenet the qte by whatever is passed, but that's a design / management decision
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	err := AddBook(handler.db, book)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
}

// delete existing artcile
func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	err := DeleteBook(handler.db, book.Book_Id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

// update existing book
func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	err := UpdateBook(handler.db, book)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusNoContent)
}

// serve
func Serve(port uint16, db *sql.DB) {
	// simple server, no permission checking, no auth, just simple CRUD
	fmt.Println("Serving on port", port)

	handler := Handler{db: db}
	http.HandleFunc("/books/", handler.GetAll)
	http.HandleFunc("/books/add", handler.Add)
	http.HandleFunc("/books/delete", handler.Delete)
	http.HandleFunc("/books/update", handler.Update)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
