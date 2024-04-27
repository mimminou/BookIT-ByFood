package services

import (
	"database/sql"
	"encoding/json"
	"github.com/mimminou/BookIT-ByFood/back/database"
	"github.com/mimminou/BookIT-ByFood/back/models"
	"github.com/mimminou/BookIT-ByFood/back/utils"

	"log"
	"net/http"
	"strconv"
	"strings"
)

type DBRequestHandler struct {
	Db *sql.DB
}

// ErrMessage is the schema for error responses

// @Description	ErrMessage
// @Property		msg string true "Error message"
type ErrMessage struct {
	Msg string `json:"msg"`
}

// get all records in DB
// Get all books

// @Summary		Get all books
// @Description	Get all books in the DB
// @Tags			books
// @Accept			json
// @Produce		json
// @Success		200	{array}	models.Book
// @Failure		500
// @Failure		500	{object}	ErrMessage
// @Failure		404	{object}	ErrMessage
// @Failure		405
// @Router			/books/ [get]
func (handler *DBRequestHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	books, err := database.GetBooks(handler.Db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
		w.Write(jsonResponse)
		return
	}
	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "No books found"})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

// Get a single Book

// @Summary		Get a single Book
// @Description	Get a book by ID
// @Tags			books
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Book ID"
// @Success		200	{object}	models.Book
// @Failure		400
// @Failure		404	{object}	ErrMessage
// @Failure		405
// @Router			/books/{id} [get]
func (handler *DBRequestHandler) GetBook(w http.ResponseWriter, r *http.Request) {
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

	fetchedBook, err := database.GetBook(handler.Db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Book not found"})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedBook)
}

// Add a new book

// @Summary		Add a new book
// @Description	Add a new book
// @Tags			books
// @Accept			json
// @Produce		json
// @Param			book	body	models.Book	true	"Book"
// @Success		200
// @Failure		400
// @Failure		404	{object}	ErrMessage
// @Failure		405
// @Router			/books/ [post]
func (handler *DBRequestHandler) Add(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var book models.Book
	defer r.Body.Close()

	decodeErr := json.NewDecoder(r.Body).Decode(&book)

	if decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: decodeErr.Error()})
		w.Write(jsonResponse)
		return
	}

	emptyFields := utils.CheckEmptyFields(book)
	if len(emptyFields) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "The following fields are empty: " + strings.Join(emptyFields, ", ")})
		w.Write(jsonResponse)
		return
	}

	if utils.ValidateDate(book.Pub_Date) == false {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Invalid date format. Should be YYYY-MM-DD"})
		w.Write(jsonResponse)
		return
	}

	id, err := database.AddBook(handler.Db, book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
		w.Write(jsonResponse)
		return
	}

	book.Book_Id = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// delete an existing book

//	 @Summary		Delete a book
//		@Description	Delete a book by ID
//		@Tags			books
//		@Accept			json
//		@Produce		json
//		@Param			id	path	int	true	"Book ID"
//		@Success		200
//		@Failure		400
//		@Failure		404	{object}	ErrMessage
//		@Failure		405
//		@Router			/books/{id} [delete]
func (handler *DBRequestHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	delErr := database.DeleteBook(handler.Db, id)
	if delErr != nil {
		if delErr == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: delErr.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: delErr.Error()})
		w.Write(jsonResponse)
		return

	}
	w.WriteHeader(http.StatusOK)
}

// update existing book

// @Summary 		Update a book
// @Description	Update a book by ID
// @Tags			books
// @Accept			json
// @Produce		json
// @Param			id		path		int			true	"Book ID"
// @Param			book	body		models.Book	true	"Book"
// @Success		200		{object}	models.Book
// @Failure		400		{object}	ErrMessage
// @Failure		500		{object}	ErrMessage
// @Failure		404		{object}	ErrMessage
// @Failure		405
// @Router			/books/{id} [put]
func (handler *DBRequestHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	//Extract the ID from URI path
	parts := strings.Split(r.URL.Path, "/")
	stringID := parts[len(parts)-1]

	id, err := strconv.Atoi(stringID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
		w.Write(jsonResponse)
		return
	}

	var book models.Book
	decodeErr := json.NewDecoder(r.Body).Decode(&book)

	if decodeErr != nil {
		log.Println(decodeErr)
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: decodeErr.Error()})
		w.Write(jsonResponse)
		return
	}
	book.Book_Id = id

	if utils.ValidateDate(book.Pub_Date) == false {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Invalid date format. Should be YYYY-MM-DD"})
		w.Write(jsonResponse)
		return
	}

	UpdateErr := database.UpdateBook(handler.Db, book)
	if UpdateErr != nil {
		if UpdateErr == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: UpdateErr.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: UpdateErr.Error()})
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (handler *DBRequestHandler) SendOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
