package admin

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"base/blog/config"
	"base/blog/core"
	"base/blog/core/admin/controllers"
	"base/blog/core/admin/models"
)

type Module struct {
	core.Module
}

// NewModule initializes the Admin module with its controllers
func NewModule(db *gorm.DB, cfg *config.Config) *Module {
	m := &Module{
		Module: core.Module{
			DB:          db,
			Cfg:         cfg,
			Controllers: make(map[string]core.Controller),
		},
	}

	// Initialize and map controllers
	m.Controllers["admin"] = controllers.NewAdminController(db, cfg)
	// Uncomment if you have an API controller
	// m.Controllers["api"] = controllers.NewAPIController(db, cfg)

	return m
}

// Migrate applies migrations for the Admin module, ensuring relevant tables are created
func (m *Module) Migrate(db *gorm.DB) error {
	// Perform migration for admin-related models, such as administrators, settings, etc.
	return db.AutoMigrate(&models.Administrator{}, &models.Setting{})
}

// SetupRoutes sets up the routes for the Admin module, especially under /admin
func (m *Module) SetupRoutes(r chi.Router, apiRouter, adminRouter chi.Router) error {
	ac := m.Controllers["admin"].(*controllers.AdminController)

	// Admin-specific routes
	if adminRouter != nil {
		adminRouter.Route("/admin", func(r chi.Router) {
			r.Get("/dashboard", ac.DashboardHandler)
			r.Get("/settings", ac.SettingsHandler)
			r.Post("/settings", ac.SettingsHandler)

			r.Get("/logout", ac.LogoutHandler)
			r.Post("/login", ac.LoginHandler)
			r.Get("/login", ac.LoginHandler)
			r.Get("/register", ac.RegisterHandler)
			r.Post("/register", ac.RegisterHandler)
			r.Get("/forgot-password", ac.ForgotPasswordHandler)
			r.Post("/forgot-password", ac.ForgotPasswordHandler)
			r.Get("/reset-password", ac.ResetPasswordHandler)
			r.Post("/reset-password", ac.ResetPasswordHandler)
			r.Get("/profile", ac.ProfileHandler)
			r.Post("/profile", ac.ProfileHandler)
			r.Get("/change-password", ac.ChangePasswordHandler)
			r.Post("/change-password", ac.ChangePasswordHandler)

			// Additional admin routes can be added here
		})
	}

	return nil
}
