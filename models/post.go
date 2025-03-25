package models

import (
	"database/sql"
	"time"
)

// Post struct + CRUD

type Post struct {
	ID         int
	UserID     int
	Title      string
	Content    string
	CreatedAt  string
	CategoryID int
	Author     string
	Category   string
}

func CreatePost(db *sql.DB, userID string, CategoryID int, title, content string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := createdAt

	_, err := db.Exec(`
		INSERT INTO posts (user_id, category_id, title, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, userID, title, content, createdAt, updatedAt)

	return err
}

// GetAllPosts retourne tous les posts depuis la base
func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, p.category_id, p.title, p.content, p.created_at,
		       u.username, c.name
		FROM Posts p
		INNER JOIN Users u ON p.user_id = u.id
		INNER JOIN Categories c ON p.category_id = c.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.CategoryID, &p.Title, &p.Content, &p.CreatedAt, &p.Author, &p.Category)
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
		SELECT p.id, p.user_id, p.category_id, p.title, p.content, p.created_at,
		       u.username, c.name
		FROM Posts p
		INNER JOIN Users u ON p.user_id = u.id
		INNER JOIN Categories c ON p.category_id = c.id
		WHERE p.id = ?
	`, id).Scan(&p.ID, &p.UserID, &p.CategoryID, &p.Title, &p.Content, &p.CreatedAt, &p.Author, &p.Category)

	return p, err
}

// DeletePost supprime un post par son ID
func DeletePost(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}
