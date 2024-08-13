package controllers

import (
	"database/sql"
	"net/http"

	"base/blog/config"
	"base/blog/modules/users/models"
	"base/blog/utils"
)

type FrontController struct {
	DB     *sql.DB
	Config *config.Config
}

func NewFrontController(db *sql.DB, cfg *config.Config) *FrontController {
	return &FrontController{DB: db, Config: cfg}
}

func (fc *FrontController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user := &models.User{Username: username, Email: email, Password: password}
		err := user.Save(fc.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title": "Register",
	}
	err := utils.RenderTemplate(w, "register.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FrontController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := models.AuthenticateUser(fc.DB, email, password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set session or JWT token here
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title": "Login",
	}
	err := utils.RenderTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FrontController) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session or JWT token
	userID := 1 // This should be retrieved from the session

	user, err := models.GetUserByID(fc.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "User Profile",
		"User":  user,
	}
	err = utils.RenderTemplate(w, "profile.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
