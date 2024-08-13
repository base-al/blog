package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(db *sql.DB, id int) (User, error) {
	var u User
	err := db.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?", id).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (u *User) Save(db *sql.DB) error {
	if u.ID == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		result, err := db.Exec("INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", u.Username, u.Email, string(hashedPassword), time.Now(), time.Now())
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		u.ID = int(id)
		return err
	}
	_, err := db.Exec("UPDATE users SET username = ?, email = ?, updated_at = ? WHERE id = ?", u.Username, u.Email, time.Now(), u.ID)
	return err
}

func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func AuthenticateUser(db *sql.DB, email, password string) (User, error) {
	var u User
	err := db.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return User{}, err
	}

	return u, nil
}
