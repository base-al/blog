package users

import (
	"database/sql"

	"github.com/gorilla/mux"

	"base/blog/config"
	"base/blog/modules/users/controllers"
)

func SetupRoutes(r, apiRouter, adminRouter *mux.Router, db *sql.DB, cfg *config.Config) {
	ac := controllers.NewAdminController(db, cfg)
	fc := controllers.NewFrontController(db, cfg)

	// Admin routes
	admin := r.PathPrefix("/admin/users").Subrouter()
	admin.HandleFunc("", ac.ListUsersHandler).Methods("GET")
	admin.HandleFunc("/edit/{id:[0-9]+}", ac.EditUserHandler).Methods("GET")
	admin.HandleFunc("/update", ac.UpdateUserHandler).Methods("POST")
	admin.HandleFunc("/delete/{id:[0-9]+}", ac.DeleteUserHandler).Methods("POST")

	// Front-end routes
	r.HandleFunc("/register", fc.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/login", fc.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/profile", fc.ProfileHandler).Methods("GET")
}
