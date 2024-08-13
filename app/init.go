package app

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"base/blog/config"
	"base/blog/core"
	"base/blog/middleware"

	"base/blog/app/comments"
	"base/blog/app/posts"
)

func InitModules(db *gorm.DB, cfg *config.Config) *core.Registry {
	registry := core.NewRegistry()

	// Register modules
	registry.RegisterModule("posts", posts.NewModule(db, cfg))
	registry.RegisterModule("comments", comments.NewModule(db, cfg))

	// Migrate modules
	registry.MigrateModules(db)

	return registry
}

func SetupRoutes(r chi.Router, registry *core.Registry, cfg *config.Config) {
	// Apply logging middleware to all routes
	r.Use(middleware.LoggingMiddleware)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.APIMiddleware)
		for _, module := range registry.Modules {
			module.SetupRoutes(r, r, nil)
		}
	})

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg))
		for _, module := range registry.Modules {
			module.SetupRoutes(r, nil, r)
		}
	})

	// Front-end routes
	for _, module := range registry.Modules {
		module.SetupRoutes(r, nil, nil)
	}
}
