package controllers

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/mimminou/BookIT-ByFood/back/services"
	"net/http"
)

func BookController(db *sql.DB) http.Handler {
	booksMux := chi.NewRouter()

	dbRequestHandler := &services.DBRequestHandler{Db: db}
	//Register GET routes
	booksMux.Get("/", dbRequestHandler.GetAll)
	booksMux.Get("/{id}", dbRequestHandler.GetBook)

	//Register POST routes
	booksMux.Post("/", dbRequestHandler.Add)
	booksMux.Put("/{id}", dbRequestHandler.Update)
	booksMux.Delete("/{id}", dbRequestHandler.Delete)

	//Register OPTIONS routes
	booksMux.Options("/", dbRequestHandler.SendOptions)
	booksMux.Options("/{id}", dbRequestHandler.SendOptions)

	return booksMux
}
