package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"base/blog/config"
	"base/blog/modules/users/models"
	"base/blog/utils"
)

type AdminController struct {
	DB     *sql.DB
	Config *config.Config
}

func NewAdminController(db *sql.DB, cfg *config.Config) *AdminController {
	return &AdminController{DB: db, Config: cfg}
}

func (ac *AdminController) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers(ac.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "User List",
		"Users": users,
	}

	err = utils.RenderTemplate(w, "admin/user_list.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ac *AdminController) EditUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	user, err := models.GetUserByID(ac.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Edit User",
		"User":  user,
	}

	err = utils.RenderTemplate(w, "admin/edit_user.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ac *AdminController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	username := r.FormValue("username")
	email := r.FormValue("email")

	user := &models.User{ID: id, Username: username, Email: email}
	err := user.Save(ac.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (ac *AdminController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := models.DeleteUser(ac.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
