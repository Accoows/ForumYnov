package models

import (
	"database/sql"
	"time"
)

// Post struct + CRUD

type Post struct {
	ID         int
	UserID     string
	Title      string
	Content    string
	CreatedAt  string
	UpdatedAt  string
	Categories []Category
	Author     string
}

func CreatePost(db *sql.DB, userID string, title, content string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := createdAt

	_, err := db.Exec(`
		INSERT INTO posts (user_id, title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, userID, title, content, createdAt, updatedAt)

	return err
}

// GetAllPosts retourne tous les posts depuis la base
func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
		SELECT id, user_id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// GetPostByID retourne un post spécifique par son ID
func GetPostByID(db *sql.DB, id int) (Post, error) {
	var p Post
	err := db.QueryRow(`
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts WHERE id = ?
	`, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)

	return p, err
}

// DeletePost supprime un post par son ID
func DeletePost(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}
