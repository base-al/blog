package core

import (
	"base/blog/config"
	"log"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// ModuleMigrator interface for modules that support database migration
type ModuleMigrator interface {
	Migrate(db *gorm.DB) error
}

// Controller interface
type Controller interface{}

// Module structure
type Module struct {
	DB          *gorm.DB
	Cfg         *config.Config
	Controllers map[string]Controller
}

// ModuleSetupRoutes interface for setting up module routes
type ModuleSetupRoutes interface {
	SetupRoutes(r chi.Router, apiRouter, adminRouter chi.Router) error
}

// Registry structure to hold modules
type Registry struct {
	Modules map[string]ModuleSetupRoutes
}

// NewRegistry creates a new module registry
func NewRegistry() *Registry {
	return &Registry{
		Modules: make(map[string]ModuleSetupRoutes),
	}
}

// RegisterModule registers a new module in the registry
func (r *Registry) RegisterModule(name string, module ModuleSetupRoutes) {
	r.Modules[name] = module
}

// GetModule retrieves a module by name
func (r *Registry) GetModule(name string) ModuleSetupRoutes {
	return r.Modules[name]
}

// MigrateModules migrates all registered modules that implement the ModuleMigrator interface
func (r *Registry) MigrateModules(db *gorm.DB) {
	for name, module := range r.Modules {
		if migrator, ok := module.(ModuleMigrator); ok {
			if err := migrator.Migrate(db); err != nil {
				log.Printf("Error migrating module %s: %v", name, err)
			} else {
				log.Printf("Successfully migrated module %s", name)
			}
		}
	}
}
