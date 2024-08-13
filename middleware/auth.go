package middleware

import (
	"net/http"

	"base/blog/config"
)

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for session cookie or JWT token
			// This is a basic example, you should implement proper authentication
			cookie, err := r.Cookie("session_token")
			if err != nil || cookie.Value != "admin" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
