package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"base/blog/config"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type ContextKey string

const (
	SessionCookieName = "session_token"
	UserContextKey    = ContextKey("user")
)

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func SetSession(w http.ResponseWriter, userID int, cfg *config.Config) error {
	token, err := GenerateSessionToken()
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Hour * 24 / time.Second), // 24 hours
	})

	// Here you would typically store the session in a database or cache
	// For simplicity, we're just using the cookie

	return nil
}

func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

func GetUserFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserContextKey).(int)
	return userID, ok
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain-text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Admin session management
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// SetSetAdminSessionSession sets a session value for the admin
func SetAdminSession(w http.ResponseWriter, r *http.Request, key string, value interface{}) {
	session, _ := store.Get(r, "session-name")
	session.Values[key] = value
	session.Save(r, w)
}

// GetAdminSessionValue gets a session value for the admin
func GetAdminSessionValue(r *http.Request, key string) interface{} {
	session, _ := store.Get(r, "session-name")
	return session.Values[key]
}

// ClearAdminSession clears a session value for the admin
func ClearAdminSession(w http.ResponseWriter, r *http.Request, key string) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, key)
	session.Save(r, w)
}
