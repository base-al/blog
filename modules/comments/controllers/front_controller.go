package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"base/blog/config"
	"base/blog/modules/comments/models"
	"base/blog/utils"
)

type FrontController struct {
	DB     *sql.DB
	Config *config.Config
}

func NewFrontController(db *sql.DB, cfg *config.Config) *FrontController {
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

	err := comment.Save(fc.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := models.GetCommentsByPostID(fc.DB, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, "comments/comment_list.html", comments)
}

func (fc *FrontController) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("post_id"))

	comments, err := models.GetCommentsByPostID(fc.DB, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, "comments/comment_list.html", comments)
}
