package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

//Middleware comes here

// EnableCors Middlware, wildcard * allows any remote to access the API
// handy when testing, never use * in prod
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// Logging Middlware, writes requests to console and to log.txt
// (file size can grow a LOT, checking file size not implemented, not needed for this test project)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logMsg := fmt.Sprintf("Request : %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

		log.Println(logMsg)
		file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening log.txt: ", err)
		}

		_, writeErr := file.WriteString(logMsg)
		if writeErr != nil {
			log.Println("Error writing to log.txt: ", writeErr)
		}

		// Closing after each request will introduce overhead when many requests are sent, but not needed for this test project
		defer file.Close()
		next.ServeHTTP(w, r)
	})
}

// ResponseLogging Middlware, writes responses to console
func ResponseLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := &myResponseWriter{w, bytes.Buffer{}, http.StatusOK} //Status ok as default
		next.ServeHTTP(wr, r)
		log.Printf("Response : %s %d %s\n", r.RemoteAddr, wr.statusCode, wr.body.String())
	})
}

type myResponseWriter struct {
	http.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func (resp *myResponseWriter) Write(b []byte) (int, error) {
	resp.body.Write(b)
	return resp.ResponseWriter.Write(b)
}

func (resp *myResponseWriter) WriteHeader(code int) {
	resp.statusCode = code
	resp.ResponseWriter.WriteHeader(code)
}
