package middleware

import (
	"base/blog/config"
	"base/blog/utils"
	"context"
	"net/http"
)

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(utils.SessionCookieName)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Here you would typically validate the session token against a database or cache
			// For simplicity, we're just checking if it exists
			if cookie.Value != "" {
				// In a real application, you'd fetch the user ID associated with this session
				userID := 1 // This should come from your session store
				ctx := context.WithValue(r.Context(), utils.UserContextKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		})
	}
}
