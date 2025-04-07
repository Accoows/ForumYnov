package database

import "database/sql"

type Categories struct {
	ID        int
	Name      string
	Parent_id int
}

type Users struct {
	ID            string
	Email         string
	Username      string
	Password_hash string
	Created_at    string
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
}

type LikesDislikes struct {
	ID         int
	User_id    string
	Post_id    sql.NullInt64
	Comment_id sql.NullInt64
	TypeValue  int
}
