package server

import (
	"log"
	"net/http"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, "executing logMiddleware")
		next.ServeHTTP(w, r)
		log.Println(r.URL.Path, "executing logMiddleware again")
	})
}
