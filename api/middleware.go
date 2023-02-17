package api

import (
	"net/http"
	"time"
	"log"
)

func Logging(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start := time.Now()
		log.Printf("Starting Server...")
		log.Printf("Started %v %v", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
		log.Printf("Completed %v in %v", r.URL.Path, time.Since(start))
	})
}