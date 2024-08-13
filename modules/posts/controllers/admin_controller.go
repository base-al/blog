package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"base/blog/config"
	"base/blog/modules/posts/models"
	"base/blog/utils"
)

type AdminController struct {
	DB     *sql.DB
	Config *config.Config
}

func NewAdminController(db *sql.DB, cfg *config.Config) *AdminController {
	return &AdminController{DB: db, Config: cfg}
}

func (ac *AdminController) ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts(ac.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, "admin/dashboard.html", posts)
}

func (ac *AdminController) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	post, err := models.GetPostByID(ac.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, "admin/edit_post.html", post)
}

func (ac *AdminController) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := &models.Post{ID: id, Title: title, Content: content}
	err := post.Save(ac.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (ac *AdminController) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := models.DeletePost(ac.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (ac *AdminController) NewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")

		post := &models.Post{Title: title, Content: content}
		err := post.Save(ac.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, "admin/new_post.html", nil)
}
