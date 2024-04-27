package services

import (
	_ "github.com/mimminou/BookIT-ByFood/back/docs"
	"net/http"
)

// server swagger yaml file
// Serves Swagger Docs

// @Summary 		Serves Swagger Docs
// @Description 	Serves Swagger Docs
// @Tags			docs
// @Success		200
// @Failure		404
// @Router		/docs/ [get]
func ServeYaml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger.yaml")
}
