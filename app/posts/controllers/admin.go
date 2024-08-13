package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"base/blog/app/posts/models"
	"base/blog/config"
	"base/blog/utils"
)

type AdminController struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewAdminController(db *gorm.DB, cfg *config.Config) *AdminController {
	return &AdminController{DB: db, Config: cfg}
}

// Index lists all posts
func (ac *AdminController) Index(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	result := ac.DB.Find(&posts)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.RenderTemplate(w, "admin/posts/index.html", posts)
}

// Show displays a specific post
func (ac *AdminController) Show(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var post models.Post
	result := ac.DB.First(&post, id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.RenderTemplate(w, "admin/posts/show.html", post)
}

// New displays a form for a new post
func (ac *AdminController) New(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "admin/posts/new.html", nil)
}

// Create adds a new post
func (ac *AdminController) Create(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	result := ac.DB.Create(&post)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, post)
}

// Edit displays a form to edit an existing post
func (ac *AdminController) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var post models.Post
	result := ac.DB.First(&post, id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.RenderTemplate(w, "admin/posts/edit.html", post)
}

// Update modifies an existing post
func (ac *AdminController) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var post models.Post
	result := ac.DB.First(&post, id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	result = ac.DB.Save(&post)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

// Destroy removes a post
func (ac *AdminController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	result := ac.DB.Delete(&models.Post{}, id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
