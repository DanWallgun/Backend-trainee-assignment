package middleware

import (
	"log"
	"net/http"
	"time"
)

// Log - loggin middleware
func Log(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Method: %s; RemoteAddr: %s; URL: %s", r.Method, r.RemoteAddr, r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Printf("Time: %s\n", time.Since(start))
	})
}
