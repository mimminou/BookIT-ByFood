package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mimminou/BookIT-ByFood/back/services"
	"net/http"
)

// implement this
func UrlCleanerController() http.Handler {
	urlMux := chi.NewRouter()
	urlMux.Post("/", services.ProcessUrl)
	return urlMux
}
