package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mimminou/BookIT-ByFood/back/docs"
	"github.com/mimminou/BookIT-ByFood/back/models"
	"github.com/mimminou/BookIT-ByFood/back/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	db *sql.DB
}

// ErrMessage is the schema for error responses

// @Description	ErrMessage
// @Property		msg string true "Error message"
type ErrMessage struct {
	Msg string `json:"msg"`
}

// These are good enough for requirments of this project, but it's much better to use a routing library like Gorilla or Chi
var BookIDRegex = regexp.MustCompile("^/books/[0-9]+(/)?$")
var BookRegex = regexp.MustCompile("^/books(/)?$")

// routes requests that have ID based on HTTP methods
func (handler *Handler) Router(w http.ResponseWriter, r *http.Request) {

	if !BookIDRegex.MatchString(r.URL.Path) && !BookRegex.MatchString(r.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Url format invalid, should be /books or /books/[id] where id is a positive integer"})
		w.Write(jsonResponse)
		return
	}

	if BookRegex.MatchString(r.URL.Path) {
		switch r.Method {
		case "POST":
			handler.Add(w, r)
		case "GET":
			handler.GetAll(w, r)
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
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
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// get all records in DB
// Get all books
//
//	@Summary		Get all books
//	@Description	Get all books in the DB
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.Book
//	@Failure		500
//	@Failure		500	{object}	ErrMessage
//	@Failure		404	{object}	ErrMessage
//	@Failure		405
//	@Router			/books/ [get]
func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	books, err := GetBooks(handler.db)
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
//
//	@Summary		Get a single Book
//	@Description	Get a book by ID
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		200	{object}	models.Book
//	@Failure		400
//	@Failure		404	{object}	ErrMessage
//	@Failure		405
//	@Router			/books/{id} [get]
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
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Book not found"})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedBook)
}

// Add a new book
//
//	@Summary		Add a new book
//	@Description	Add a new book
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book	body	models.Book	true	"Book"
//	@Success		200
//	@Failure		400
//	@Failure		404	{object}	ErrMessage
//	@Failure		405
//	@Router			/books/ [post]
func (handler *Handler) Add(w http.ResponseWriter, r *http.Request) {

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

	id, err := AddBook(handler.db, book)
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
//
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
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
		w.Write(jsonResponse)
		return
	}

	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	book.Book_Id = id
	UpdateErr := UpdateBook(handler.db, book)
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

// serve
func Serve(port uint16, db *sql.DB) {

	// simple server, no permission checks, no auth, only input Sanitazation and simple CRUD
	serverMux := http.NewServeMux()
	dbRequestHandler := &Handler{db: db}
	middleWareMux := Cors(ResponseLogging(Logging((serverMux))))
	serverMux.HandleFunc("/books", dbRequestHandler.Router)
	serverMux.HandleFunc("/books/", dbRequestHandler.Router)

	//No idea how to make swagger "see" the file, so i'll serve it myself in hope that that the next handleFunc will catch it
	serverMux.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	//add swagger link
	serverMux.HandleFunc("/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.yaml"),
	))

	fmt.Println("Serving on port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), middleWareMux)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
