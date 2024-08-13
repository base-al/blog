// modules/comments/module.go
package comments

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"base/blog/app/comments/controllers"
	"base/blog/app/comments/models"
	"base/blog/config"
	"base/blog/core"
)

type Module struct {
	core.Module
}

func NewModule(db *gorm.DB, cfg *config.Config) *Module {
	m := &Module{
		Module: core.Module{
			DB:          db,
			Cfg:         cfg,
			Controllers: make(map[string]core.Controller),
		},
	}
	m.Controllers["front"] = controllers.NewFrontController(db, cfg)
	m.Controllers["admin"] = controllers.NewAdminController(db, cfg)
	// If you have an API controller, uncomment the following line:
	// m.Controllers["api"] = controllers.NewAPIController(db, cfg)
	return m
}

// Migrate applies migrations for the Posts module
func (m *Module) Migrate(db *gorm.DB) error {
	// Perform migration logic here
	return db.AutoMigrate(&models.Comment{})
}
func (m *Module) SetupRoutes(r chi.Router, apiRouter, adminRouter chi.Router) error {
	fc := m.Controllers["front"].(*controllers.FrontController)
	ac := m.Controllers["admin"].(*controllers.AdminController)
	// If you have an API controller, uncomment the following line:
	// api := m.Controllers["api"].(*controllers.APIController)

	// Front-end routes
	r.Route("/comments", func(r chi.Router) {
		r.Get("/", fc.GetCommentsHandler)
		r.Post("/add", fc.AddCommentHandler)
	})

	// Admin routes
	if adminRouter != nil {
		adminRouter.Route("/comments", func(r chi.Router) {
			r.Get("/manage", ac.ManageCommentsHandler)
			r.Delete("/delete", ac.DeleteCommentHandler)
		})
	}

	// API routes
	// If you have API routes, you can add them here:
	// if apiRouter != nil {
	//     apiRouter.Route("/comments", func(r chi.Router) {
	//         r.Get("/", api.ListComments)
	//         r.Post("/", api.AddComment)
	//         // Add other API routes as needed
	//     })
	// }

	return nil
}
