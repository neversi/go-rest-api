package middleware

import (
	"log"
	"net/http"
)

// LoggerHandler is the logger for Handlers
func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("End-Point: %s, Method: %s", r.RequestURI, r.Method)
	})
    }