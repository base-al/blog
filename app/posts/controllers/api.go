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

type APIController struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewAPIController(db *gorm.DB, cfg *config.Config) *APIController {
	return &APIController{DB: db, Config: cfg}
}

// Index lists all posts
func (api *APIController) Index(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	result := api.DB.Find(&posts)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, posts)
}

// Show displays a specific post
func (api *APIController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	result := api.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

// Create adds a new post
func (api *APIController) Create(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	result := api.DB.Create(&post)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, post)
}

// Update modifies an existing post
func (api *APIController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	result := api.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	result = api.DB.Save(&post)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

// Destroy removes a post
func (api *APIController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	result := api.DB.Delete(&models.Post{}, id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
