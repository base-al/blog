package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"base/blog/config"
	"base/blog/database"
	"base/blog/middleware"
	"base/blog/modules/comments"
	"base/blog/modules/posts"
	"base/blog/modules/posts/controllers"
	"base/blog/modules/users"
	"base/blog/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Initialize templates
	if err := utils.InitTemplates(cfg.Theme); err != nil {
		log.Fatalf("Error initializing templates: %v", err)
	}

	// Initialize router
	r := mux.NewRouter()

	// API routes
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(middleware.APIMiddleware)

	// Admin routes
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware(cfg))

	// Setup module routes
	posts.SetupRoutes(r, apiRouter, adminRouter, db, cfg)
	users.SetupRoutes(r, apiRouter, adminRouter, db, cfg)
	comments.SetupRoutes(r, apiRouter, adminRouter, db, cfg)

	// Explicitly set up the homepage route
	postController := controllers.NewFrontController(db, cfg)
	r.HandleFunc("/", postController.ListPostsHandler).Methods("GET")

	// Serve static files
	fs := http.FileServer(http.Dir("./themes/" + cfg.Theme + "/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Setup middleware for all routes
	r.Use(middleware.LoggingMiddleware)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Start server
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
