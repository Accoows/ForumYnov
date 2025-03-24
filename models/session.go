package models

import "time"

// Création/validation des sessions

type Session struct {
	ID          int // UUID
	UserID      int
	ExpiresUser time.Time
}
