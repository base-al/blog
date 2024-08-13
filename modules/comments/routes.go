package comments

import (
	"database/sql"

	"github.com/gorilla/mux"

	"base/blog/config"
	"base/blog/modules/comments/controllers"
)

func SetupRoutes(r, apiRouter, adminRouter *mux.Router, db *sql.DB, cfg *config.Config) {
	fc := controllers.NewFrontController(db, cfg)
	ac := controllers.NewAdminController(db, cfg)

	// Front-end routes
	r.HandleFunc("/comments", fc.GetCommentsHandler).Methods("GET")
	r.HandleFunc("/comments/add", fc.AddCommentHandler).Methods("POST")

	// Admin routes
	admin := r.PathPrefix("/admin/comments").Subrouter()
	admin.HandleFunc("/manage", ac.ManageCommentsHandler).Methods("GET")
	admin.HandleFunc("/delete", ac.DeleteCommentHandler).Methods("DELETE")
}
