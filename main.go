package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"base/blog/app"
	"base/blog/config"
	"base/blog/middleware"
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
	if err := utils.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize templates
	if err := utils.InitTemplates(cfg.Theme); err != nil {
		log.Fatalf("Error initializing templates: %v", err)
	}

	// Initialize Chi router
	r := chi.NewRouter()

	// Use Chi middlewares
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.LoggingMiddleware)

	// Create API subrouter and apply middleware before routes are mounted
	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.APIMiddleware)
	r.Mount("/api/v1", apiRouter)

	// Create Admin subrouter and apply middleware before routes are mounted
	adminRouter := chi.NewRouter()
	adminRouter.Use(middleware.AuthMiddleware(cfg))
	r.Mount("/admin", adminRouter)

	// Initialize and setup module routes
	registry := app.InitModules(utils.DB, cfg)

	// Apply middleware and set up routes
	for name, module := range registry.Modules {
		if err := module.SetupRoutes(r, apiRouter, adminRouter); err != nil {
			log.Fatalf("Error setting up routes for module %s: %v", name, err)
		}
	}

	// Serve static files
	filesDir := http.Dir("./themes/" + cfg.Theme + "/static")
	utils.FileServer(r, "/static", filesDir)

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
