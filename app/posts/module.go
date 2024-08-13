// modules/posts/module.go
package posts

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"base/blog/app/posts/controllers"
	"base/blog/app/posts/models"
	"base/blog/config"
	"base/blog/core"
)

type Module struct {
	core.Module
	db  *gorm.DB
	cfg *config.Config
}

func NewModule(db *gorm.DB, cfg *config.Config) *Module {
	return &Module{
		db:  db,
		cfg: cfg,
	}
}

// Migrate applies migrations for the Posts module
func (m *Module) Migrate(db *gorm.DB) error {
	// Perform migration logic here
	return db.AutoMigrate(&models.Post{})
}

func (m *Module) SetupRoutes(r chi.Router, apiRouter, adminRouter chi.Router) error {
	fc := controllers.NewFrontController(m.db, m.cfg)
	ac := controllers.NewAdminController(m.db, m.cfg)
	api := controllers.NewAPIController(m.db, m.cfg)

	// Front-end routes
	r.Get("/", fc.ListPostsHandler)
	r.Get("/post/{id}", fc.GetPostHandler)

	// Admin routes
	if adminRouter != nil {
		adminRouter.Route("/admin/posts", func(r chi.Router) {
			r.Get("/", ac.Index)
			r.Get("/new", ac.New)
			r.Post("/", ac.Create)
			r.Get("/{id}", ac.Show)
			r.Get("/{id}/edit", ac.Edit)
			r.Put("/{id}", ac.Update)
			r.Delete("/{id}", ac.Destroy)
		})
	}

	// API routes
	if apiRouter != nil {
		apiRouter.Route("/api/posts", func(r chi.Router) {
			r.Get("/", api.Index)
			r.Post("/", api.Create)
			r.Get("/{id}", api.Show)
			r.Put("/{id}", api.Update)
			r.Delete("/{id}", api.Destroy)
		})
	}

	return nil
}
