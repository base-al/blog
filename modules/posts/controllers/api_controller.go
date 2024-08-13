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

type APIController struct {
	DB     *sql.DB
	Config *config.Config
}

func NewAPIController(db *sql.DB, cfg *config.Config) *APIController {
	return &APIController{DB: db, Config: cfg}
}

func (api *APIController) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts(api.DB)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, posts)
}

func (api *APIController) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	post, err := models.GetPostByID(api.DB, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

func (api *APIController) CreatePost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := &models.Post{Title: title, Content: content}
	err := post.Save(api.DB)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, post)
}

func (api *APIController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := &models.Post{ID: id, Title: title, Content: content}
	err := post.Save(api.DB)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

func (api *APIController) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := models.DeletePost(api.DB, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
