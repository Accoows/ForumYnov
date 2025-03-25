package models

import "time"

// Commentaire struct + CRUD

type Comment struct {
	ID        int
	UserID    int
	PostID    int
	Content   string
	CreatedAt time.Time
	Username  string
}
