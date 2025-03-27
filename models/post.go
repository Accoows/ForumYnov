package models

import (
	"forumynov/database"
	"time"
)

func CreatePost(userID, categoryID int, title, content string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := database.SQL.Exec(`
		INSERT INTO Posts (user_id, category_id, title, content, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		userID, categoryID, title, content, createdAt)
	return err
}
