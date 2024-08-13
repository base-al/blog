package utils

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

var templates map[string]*template.Template

// InitTemplates initializes the templates with the given theme

func InitTemplates(theme string) error {
	templates = make(map[string]*template.Template)

	// Load main layout
	mainLayout := filepath.Join("themes", theme, "layouts", "main.html")

	// Load partials
	partials, err := filepath.Glob(filepath.Join("themes", theme, "partials", "*.html"))
	if err != nil {
		log.Printf("Failed to load partials: %v", err)
		return err
	}

	// Load content templates
	contentTemplates, err := filepath.Glob(filepath.Join("modules", "*", "views", "*", "*.html"))
	if err != nil {
		log.Printf("Failed to load content templates: %v", err)
		return err
	}

	// Parse each content template with the main layout and partials
	for _, contentTemplate := range contentTemplates {
		files := append([]string{mainLayout}, partials...)
		files = append(files, contentTemplate)

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", contentTemplate, err)
		}

		name := filepath.Base(contentTemplate)
		templates[name] = tmpl
	}

	return nil
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

// RenderTemplate renders a template with the given name and data
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		log.Printf("Template %s not found", name)
		return fmt.Errorf("template %s not found", name)
	}

	w.Header().Set("Content-Type", "text/html")
	return tmpl.ExecuteTemplate(w, "main", data)
}
func RenderAdminTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("core/admin/views", tmpl)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
