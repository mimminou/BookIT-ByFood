package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

// Logging Middlware, writes requests to console
// write to in.log
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logMsg := fmt.Sprintf("%s Request : %s %s %s", time.Now().Format("2006-01-02 15:04:05"), r.RemoteAddr, r.Method, r.URL.Path)

		//write on disk for persistant logs
		f, err := os.OpenFile("in.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(logMsg + "\n"); err != nil {
			log.Println(err)
		}

		fmt.Println(logMsg)
		next.ServeHTTP(w, r)
	})
}

// ResponseLogging Middlware, writes responses to console
// write to out.log
func ResponseLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := &myResponseWriter{w, bytes.Buffer{}, http.StatusOK} //Status ok as default
		logMsg := fmt.Sprintf("%s Response : %s %d %s\n", time.Now().Format("2006-01-02 15:04:05"), r.RemoteAddr, wr.statusCode, wr.body.String())

		//write on disk for persistant logs
		f, err := os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(logMsg); err != nil {
			log.Println(err)
		}

		fmt.Println(logMsg)
		next.ServeHTTP(wr, r)
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
