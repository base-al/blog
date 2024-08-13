package models

import (
	"database/sql"
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

func GetCommentsByPostID(db *sql.DB, postID int) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (c *Comment) Save(db *sql.DB) error {
	if c.ID == 0 {
		result, err := db.Exec("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
			c.PostID, c.UserID, c.Content, time.Now())
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		c.ID = int(id)
	} else {
		_, err := db.Exec("UPDATE comments SET content = ? WHERE id = ?", c.Content, c.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteComment(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}
