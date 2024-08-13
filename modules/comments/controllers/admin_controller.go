package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"base/blog/config"
	"base/blog/modules/comments/models"
	"base/blog/utils"
)

type AdminController struct {
	DB     *sql.DB
	Config *config.Config
}

func NewAdminController(db *sql.DB, cfg *config.Config) *AdminController {
	return &AdminController{DB: db, Config: cfg}
}

func (ac *AdminController) ManageCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("post_id"))

	comments, err := models.GetCommentsByPostID(ac.DB, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, "admin/manage_comments.html", comments)
}

func (ac *AdminController) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	err := models.DeleteComment(ac.DB, commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
