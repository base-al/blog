package posts

import (
	"database/sql"

	"github.com/gorilla/mux"

	"base/blog/config"
	"base/blog/modules/posts/controllers"
)

func SetupRoutes(r, apiRouter, adminRouter *mux.Router, db *sql.DB, cfg *config.Config) {
	fc := controllers.NewFrontController(db, cfg)
	ac := controllers.NewAdminController(db, cfg)
	api := controllers.NewAPIController(db, cfg)

	// Front-end routes
	r.HandleFunc("/", fc.ListPostsHandler)
	r.HandleFunc("/post/{id:[0-9]+}", fc.GetPostHandler)

	// Admin routes
	adminRouter.HandleFunc("/posts", ac.ListPostsHandler)
	adminRouter.HandleFunc("/posts/edit/{id:[0-9]+}", ac.EditPostHandler)
	adminRouter.HandleFunc("/posts/update", ac.UpdatePostHandler).Methods("POST")
	adminRouter.HandleFunc("/posts/delete/{id:[0-9]+}", ac.DeletePostHandler)

	// API routes
	apiRouter.HandleFunc("/posts", api.ListPosts).Methods("GET")
	apiRouter.HandleFunc("/posts/{id:[0-9]+}", api.GetPost).Methods("GET")
	apiRouter.HandleFunc("/posts", api.CreatePost).Methods("POST")
	apiRouter.HandleFunc("/posts/{id:[0-9]+}", api.UpdatePost).Methods("PUT")
	apiRouter.HandleFunc("/posts/{id:[0-9]+}", api.DeletePost).Methods("DELETE")
}
