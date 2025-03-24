package models

// Post struct + CRUD

type Post struct {
	ID         int
	UserID     int
	Title      string
	Content    string
	CreatedAt  string
	UpdatedAt  string
	Categories []Category
}
