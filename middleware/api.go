package middleware

import (
	"net/http"
)

func APIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type to JSON for all API responses
		w.Header().Set("Content-Type", "application/json")

		// You can add more API-specific middleware logic here
		// For example, API key validation, rate limiting, etc.

		next.ServeHTTP(w, r)
	})
}
