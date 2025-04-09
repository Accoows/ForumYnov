package database

import "time"

/*
*** Data structures for the database tables ***
 */

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

type Sessions struct {
	ID         string
	User_id    string
	Expires_at time.Time
}

type Posts struct {
	ID          int
	User_id     string
	Category_id int
	Title       string
	Content     string
	Created_at  string
}

type Comments struct {
	ID         int
	Post_id    int
	User_id    string
	Content    string
	Created_at string
}

type LikesDislikes struct {
	ID         int
	User_id    string
	Post_id    int
	Comment_id int
	TypeValue  int
}
