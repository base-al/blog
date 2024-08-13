package controllers

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"base/blog/app/comments/models"
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

func (fc *FrontController) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	userID, _ := strconv.Atoi(r.FormValue("user_id")) // In reality, get this from the session
	content := r.FormValue("content")

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	result := fc.DB.Create(comment)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	var comments []models.Comment
	result = fc.DB.Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, "comments/comment_list.html", comments)
}

func (fc *FrontController) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("post_id"))

	var comments []models.Comment
	result := fc.DB.Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, "comments/comment_list.html", comments)
}
