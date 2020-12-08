package middleware

import (
	"io/ioutil"
	"log"
	"net/http"
)

// LoggerHandler is the logger for Handlers
func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		if len(bodyBytes) == 0 {
			log.Printf("End-Point: %s, Method: %s, Data: %s", r.RequestURI, r.Method, "nil")
		} else {
			log.Printf("End-Point: %s, Method: %s, Data: %s", r.RequestURI, r.Method, string(bodyBytes))
		}
		next.ServeHTTP(w, r)
	})
    }