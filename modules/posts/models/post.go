package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(db *sql.DB, id int) (Post, error) {
	var p Post
	err := db.QueryRow("SELECT id, title, content, created_at, updated_at FROM posts WHERE id = ?", id).Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (p *Post) Save(db *sql.DB) error {
	if p.ID == 0 {
		result, err := db.Exec("INSERT INTO posts (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)", p.Title, p.Content, time.Now(), time.Now())
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		p.ID = int(id)
		return err
	}
	_, err := db.Exec("UPDATE posts SET title = ?, content = ?, updated_at = ? WHERE id = ?", p.Title, p.Content, time.Now(), p.ID)
	return err
}

func DeletePost(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}
