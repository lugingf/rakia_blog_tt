package middleware

import (
	"log"
	"net/http"
)

func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalPath := r.URL.Path
		if len(originalPath) > 1 && originalPath[len(originalPath)-1] == '/' {
			r.URL.Path = originalPath[:len(originalPath)-1]
		}
		log.Printf("Original URL path: %s, Updated URL path: %s", originalPath, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
