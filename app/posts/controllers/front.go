package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"base/blog/app/posts/models"
	"base/blog/config"
	"base/blog/utils"
)

type FrontController struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewFrontController(db *gorm.DB, cfg *config.Config) *FrontController {
	return &FrontController{DB: db, Config: cfg}
}

func (fc *FrontController) ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	result := fc.DB.Find(&posts)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "All Posts",
		"Posts": posts,
	}

	err := utils.RenderTemplate(w, "home.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FrontController) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var post models.Post
	result := fc.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	data := map[string]interface{}{
		"Title": post.Title,
		"Post":  post,
	}

	err := utils.RenderTemplate(w, "post.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
