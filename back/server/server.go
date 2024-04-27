package server

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mimminou/BookIT-ByFood/back/controllers"
	"log"
	"net/http"
)

// serve
func Serve(port uint16, db *sql.DB) {
	serverMux := chi.NewRouter()

	serverMux.Use(middleware.StripSlashes)
	serverMux.Use(middleware.Logger)
	serverMux.Use(Cors)

	//Mount Books Controller
	serverMux.Mount("/books", controllers.BookController(db))

	//Mount Docs Controller
	serverMux.Mount("/docs", controllers.DocsController())

	//Mount UrlCleaner Controller
	serverMux.Mount("/url", controllers.UrlCleanerController())

	fmt.Println("Serving on port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), serverMux)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
