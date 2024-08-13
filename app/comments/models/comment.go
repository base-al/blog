package models

import (
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	CreatedAt time.Time
	Username  string // We'll join this from the users table
}
