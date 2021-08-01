package main

import (
	"fmt"
	"net/http"
)

// enableCORS a middleware that we use for enabling the CORS for our client
func (app *Application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// requestWriter a middleware that we use for writing the request information
func (app *Application) requestWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		next.ServeHTTP(w, r)
	})
}