package database

import (
	"database/sql"
	"time"
)

/*
*** Data structures for the database tables ***
 */

type Categories struct {
	ID       int
	Name     string
	ParentID sql.NullInt64
}

type CategoryPageData struct {
	CategoryName string
	CategoryID   int
	Posts        []Posts
}

type Users struct {
	ID            string
	Email         string
	Username      string
	Password_hash string
	Created_at    string
}

type Sessions struct {
	ID         string
	User_id    string
	Expires_at time.Time
}

type Posts struct {
	ID             int
	User_id        string
	Category_id    int
	CategoryName   string
	Title          string
	Content        string
	Created_at     string
	AuthorUsername string
	LikeCount      int
	DislikeCount   int
}

type Comments struct {
	ID             int
	Post_id        int
	User_id        string
	Content        string
	Created_at     string
	AuthorUsername string
	LikeCount      int
	DislikeCount   int
}

type LikesDislikes struct {
	ID         int
	User_id    string
	Post_id    sql.NullInt64
	Comment_id sql.NullInt64
	TypeValue  int
}

type CreatePostPageData struct {
	CategoryID    int
	CategoryName  string
	AllCategories []Categories
}

type ErrorPageData struct {
	Code    int
	Message string
}
