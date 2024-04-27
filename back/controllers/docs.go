package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mimminou/BookIT-ByFood/back/services"
	"github.com/swaggo/http-swagger"
	"net/http"
)

func DocsController() http.Handler {
	docsMux := chi.NewRouter()
	docsMux.Get("/swagger.yaml", services.ServeYaml)
	httpSwagger.URL("/docs/swagger.yaml")
	docsMux.Get("/*", httpSwagger.Handler(httpSwagger.URL("/docs/swagger.yaml")))
	return docsMux
}
