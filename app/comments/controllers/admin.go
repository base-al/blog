package controllers

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"base/blog/app/comments/models"
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

func (ac *AdminController) ManageCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("post_id"))

	var comments []models.Comment
	result := ac.DB.Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":    "Manage Comments",
		"Comments": comments,
	}

	utils.RenderTemplate(w, "admin/manage_comments.html", data)
}

func (ac *AdminController) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	result := ac.DB.Delete(&models.Comment{}, commentID)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
