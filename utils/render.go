package utils

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates map[string]*template.Template

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

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		log.Printf("Template %s not found", name)
		return fmt.Errorf("template %s not found", name)
	}

	w.Header().Set("Content-Type", "text/html")
	return tmpl.ExecuteTemplate(w, "main", data) // Change "content" to "main"
}
