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

func Database() (*sql.DB, error) {
	var err error
	SQL, err = sql.Open("sqlite3", "database/forum.db")

	if err != nil {
		return nil, err
	}

	return SQL, nil
}

func UsersData(usersID int, email string, username string, password_hash string, usersCreated_at string) {

	SQL, err := Database()

	ErrorTest(err)

	defer SQL.Close()

	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Users")
	err = row.Scan(&usersID)

	ErrorTest(err)

	insertUsersInSql := `INSERT OR IGNORE INTO Users(id, email, username, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertUsersInSql, usersID, email, username, password_hash, usersCreated_at)

	ErrorTest(err)
}

func PostsData(postsID int, postsUser_id int, category_id int, title string, postsContent string, postsCreated_at string) {

	SQL, err := Database()

	ErrorTest(err)

	defer SQL.Close()

	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Posts")
	err = row.Scan(&postsID)

	ErrorTest(err)

	insertPostsInSql := `INSERT OR IGNORE INTO Posts(id, user_id, category_id, title, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertPostsInSql, postsID, postsUser_id, category_id, title, postsContent, postsCreated_at)

	ErrorTest(err)
}

func CommentsData(commentsID int, commentsPost_id int, commentsUser_id int, commentsContent string, commentsCreated_at string) {

	SQL, err := Database()

	ErrorTest(err)

	defer SQL.Close()

	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Comments")
	err = row.Scan(&commentsID)

	ErrorTest(err)

	insertCommentsInSql := `INSERT OR IGNORE INTO Comments(id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertCommentsInSql, commentsID, commentsPost_id, commentsUser_id, commentsContent, commentsCreated_at)

	ErrorTest(err)
}

func LikesDislikesData(likesDislikesID int, likesDislikesUser_id int, likesDislikesPost_id int, comment_id int, typeValue int) {

	SQL, err := Database()

	ErrorTest(err)

	defer SQL.Close()

	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Likes_Dislikes")
	err = row.Scan(&likesDislikesID)

	ErrorTest(err)

	insertLikesDislikesInSql := `INSERT OR IGNORE INTO Likes_Dislikes(id, user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertLikesDislikesInSql, likesDislikesID, likesDislikesUser_id, likesDislikesPost_id, comment_id, typeValue)

	ErrorTest(err)
}
