package models

import (
	"base/blog/app/comments/models"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	Comments []models.Comment
}
