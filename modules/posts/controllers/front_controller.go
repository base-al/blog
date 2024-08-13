package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"base/blog/config"
	"base/blog/modules/posts/models"
	"base/blog/utils"
)

type FrontController struct {
	DB     *sql.DB
	Config *config.Config
}

func NewFrontController(db *sql.DB, cfg *config.Config) *FrontController {
	return &FrontController{DB: db, Config: cfg}
}

func (fc *FrontController) ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts(fc.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "All Posts",
		"Posts": posts,
	}

	err = utils.RenderTemplate(w, "home.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FrontController) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	post, err := models.GetPostByID(fc.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": post.Title,
		"Post":  post,
	}

	err = utils.RenderTemplate(w, "post.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
