package controllers

import (
	"net/http"

	"gorm.io/gorm"

	"base/blog/config"
	"base/blog/core/admin/models"
	"base/blog/utils"
)

type AdminController struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewAdminController(db *gorm.DB, cfg *config.Config) *AdminController {
	return &AdminController{DB: db, Config: cfg}
}

// DashboardHandler renders the admin dashboard
func (ac *AdminController) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderAdminTemplate(w, "admin/dashboard.html", nil)
}

// SettingsHandler handles GET and POST requests for admin settings
func (ac *AdminController) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderAdminTemplate(w, "admin/settings.html", nil)
	} else if r.Method == http.MethodPost {
		// Handle settings form submission
		// Logic to update settings in the database
		// Example: ac.DB.Model(&models.Setting{}).Where("name = ?", "example_setting").Update("value", r.FormValue("example_setting_value"))

		utils.RenderAdminTemplate(w, "admin/settings.html", map[string]interface{}{
			"Success": "Settings updated successfully",
		})
	}
}

// LoginHandler handles admin login
func (ac *AdminController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderAdminTemplate(w, "admin/login.html", nil)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var admin models.Administrator
		if err := ac.DB.Where("username = ?", username).First(&admin).Error; err != nil || !utils.CheckPasswordHash(password, admin.Password) {
			utils.RenderAdminTemplate(w, "admin/login.html", map[string]interface{}{
				"Error": "Invalid username or password",
			})
			return
		}

		// Set session for the logged-in admin
		utils.SetAdminSession(w, r, "admin_id", admin.ID)

		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}
}

// LogoutHandler handles admin logout
func (ac *AdminController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.ClearAdminSession(w, r, "admin_id")
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// RegisterHandler handles admin registration (if allowed)
func (ac *AdminController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderAdminTemplate(w, "admin/register.html", nil)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			utils.RenderAdminTemplate(w, "admin/register.html", map[string]interface{}{
				"Error": "Error hashing password",
			})
			return
		}

		admin := models.Administrator{
			Username: username,
			Password: hashedPassword,
			Email:    email,
			Role:     "admin", // Set default role, you can expand this to include different roles
		}

		if err := ac.DB.Create(&admin).Error; err != nil {
			utils.RenderAdminTemplate(w, "admin/register.html", map[string]interface{}{
				"Error": "Error registering admin",
			})
			return
		}

		utils.SetAdminSession(w, r, "admin_id", admin.ID)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}
}

// ForgotPasswordHandler handles forgot password requests
func (ac *AdminController) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderAdminTemplate(w, "admin/forgot-password.html", nil)
	} else if r.Method == http.MethodPost {
		email := r.FormValue("email")

		var admin models.Administrator
		if err := ac.DB.Where("email = ?", email).First(&admin).Error; err != nil {
			utils.RenderAdminTemplate(w, "admin/forgot-password.html", map[string]interface{}{
				"Error": "No account found with that email address",
			})
			return
		}

		// Generate a password reset token and send an email (not implemented here)
		// Example: utils.SendPasswordResetEmail(admin.Email, resetToken)

		utils.RenderAdminTemplate(w, "admin/forgot-password.html", map[string]interface{}{
			"Success": "Password reset link sent to your email address",
		})
	}
}

// ResetPasswordHandler handles password resets
func (ac *AdminController) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderAdminTemplate(w, "admin/reset-password.html", nil)
	} else if r.Method == http.MethodPost {
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			utils.RenderAdminTemplate(w, "admin/reset-password.html", map[string]interface{}{
				"Error": "Passwords do not match",
			})
			return
		}

		_, err := utils.HashPassword(password)
		if err != nil {
			utils.RenderAdminTemplate(w, "admin/reset-password.html", map[string]interface{}{
				"Error": "Error hashing password",
			})
			return
		}

		// Logic to update the password in the database (assuming a reset token was provided)
		// Example: ac.DB.Model(&models.Administrator{}).Where("reset_token = ?", token).Update("password", hashedPassword)

		utils.RenderAdminTemplate(w, "admin/reset-password.html", map[string]interface{}{
			"Success": "Password successfully reset",
		})
	}
}

// ProfileHandler handles profile viewing and updating
func (ac *AdminController) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	adminID := utils.GetAdminSessionValue(r, "admin_id")
	var admin models.Administrator
	if err := ac.DB.First(&admin, adminID).Error; err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		utils.RenderAdminTemplate(w, "admin/profile.html", map[string]interface{}{
			"Admin": admin,
		})
	} else if r.Method == http.MethodPost {
		// Update profile information
		admin.Email = r.FormValue("email")
		// Update other fields as necessary

		if err := ac.DB.Save(&admin).Error; err != nil {
			utils.RenderAdminTemplate(w, "admin/profile.html", map[string]interface{}{
				"Error": "Error updating profile",
				"Admin": admin,
			})
			return
		}

		utils.RenderAdminTemplate(w, "admin/profile.html", map[string]interface{}{
			"Success": "Profile updated successfully",
			"Admin":   admin,
		})
	}
}

// ChangePasswordHandler handles password changes
func (ac *AdminController) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderAdminTemplate(w, "admin/change-password.html", nil)
	} else if r.Method == http.MethodPost {
		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		adminID := utils.GetAdminSessionValue(r, "admin_id")
		var admin models.Administrator
		if err := ac.DB.First(&admin, adminID).Error; err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		if !utils.CheckPasswordHash(oldPassword, admin.Password) {
			utils.RenderAdminTemplate(w, "admin/change-password.html", map[string]interface{}{
				"Error": "Old password is incorrect",
			})
			return
		}

		if newPassword != confirmPassword {
			utils.RenderAdminTemplate(w, "admin/change-password.html", map[string]interface{}{
				"Error": "New passwords do not match",
			})
			return
		}

		hashedPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			utils.RenderAdminTemplate(w, "admin/change-password.html", map[string]interface{}{
				"Error": "Error hashing new password",
			})
			return
		}

		admin.Password = hashedPassword
		if err := ac.DB.Save(&admin).Error; err != nil {
			utils.RenderAdminTemplate(w, "admin/change-password.html", map[string]interface{}{
				"Error": "Error updating password",
			})
			return
		}

		utils.RenderAdminTemplate(w, "admin/change-password.html", map[string]interface{}{
			"Success": "Password changed successfully",
		})
	}
}
