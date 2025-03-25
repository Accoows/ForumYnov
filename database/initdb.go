package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ErrorTest(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var SQL *sql.DB

func Database() {
	SQL, err := sql.Open("sqlite3", "database/forum.db")

	ErrorTest(err)

	defer SQL.Close()

	var email, username, password_hash, usersCreated_at string
	var usersID int
	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Users")
	err = row.Scan(&usersID)

	ErrorTest(err)

	var title, postsContent, postsCreated_at string
	var postsID, postsUser_id, category_id int
	row = SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Posts")
	err = row.Scan(&postsID)

	ErrorTest(err)

	var commentsContent, commentsCreated_at string
	var commentsID, commentsPost_id, commentsUser_id int
	row = SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Comments")
	err = row.Scan(&commentsID)

	ErrorTest(err)

	var likesDislikesID, likesDislikesUser_id, likesDislikesPost_id, comment_id, typeValue int
	row = SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Likes_Dislikes")
	err = row.Scan(&likesDislikesID)

	ErrorTest(err)

	insertUsersInSQL := `INSERT OR IGNORE INTO Users(id, email, username, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertUsersInSQL, usersID, email, username, password_hash, usersCreated_at)

	ErrorTest(err)

	insertPostsInSQL := `INSERT OR IGNORE INTO Posts(id, user_id, category_id, title, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertPostsInSQL, postsID, postsUser_id, category_id, title, postsContent, postsCreated_at)

	ErrorTest(err)

	insertCommentsInSQL := `INSERT OR IGNORE INTO Comments(id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertCommentsInSQL, commentsID, commentsPost_id, commentsUser_id, commentsContent, commentsCreated_at)

	ErrorTest(err)

	insertLikesDislikesInSQL := `INSERT OR IGNORE INTO Likes_Dislikes(id, user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertLikesDislikesInSQL, likesDislikesID, likesDislikesUser_id, likesDislikesPost_id, comment_id, typeValue)

	ErrorTest(err)
}
