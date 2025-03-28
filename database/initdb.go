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
var err error

func InitDatabase() {
	SQL, err = sql.Open("sqlite3", "database/forum.db")

	ErrorTest(err)
}

func CloseDatabase() {
	if SQL != nil {
		SQL.Close()
	}
}

func InsertUsersData(users *Users) error {

	/*row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Users")
	err = row.Scan(&users.ID)*/

	if err != nil {
		return err
	}

	insertUsersInSql := `INSERT OR IGNORE INTO Users(id, email, username, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertUsersInSql, users.ID, users.Email, users.Username, users.Password_hash, users.Created_at)

	ErrorTest(err)

	return err
}

func InsertPostsData(posts *Posts) error {

	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Posts")
	err = row.Scan(&posts.ID)

	if err != nil {
		return err
	}

	insertPostsInSql := `INSERT OR IGNORE INTO Posts(id, user_id, category_id, title, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertPostsInSql, posts.ID, posts.User_id, posts.Category_id, posts.Title, posts.Content, posts.Created_at)

	ErrorTest(err)

	return err
}

func InsertCommentsData(comments *Comments) error {

	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Comments")
	err = row.Scan(&comments.ID)

	if err != nil {
		return err
	}

	insertCommentsInSql := `INSERT OR IGNORE INTO Comments(id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertCommentsInSql, comments.ID, comments.Post_id, comments.User_id, comments.Content, comments.Created_at)

	ErrorTest(err)

	return err
}

func InsertLikesDislikesData(likesDislikes *LikesDislikes) error {

	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Likes_Dislikes")
	err = row.Scan(&likesDislikes.ID)

	if err != nil {
		return err
	}

	insertLikesDislikesInSql := `INSERT OR IGNORE INTO Likes_Dislikes(id, user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?, ?)`

	_, err = SQL.Exec(insertLikesDislikesInSql, likesDislikes.ID, likesDislikes.User_id, likesDislikes.Post_id, likesDislikes.Comment_id, likesDislikes.TypeValue)

	ErrorTest(err)

	return err
}

func GetCategoriesData() ([]Categories, error) {
	rows, err := SQL.Query("SELECT id, name, parent_id FROM Categories")
	if err != nil {
		return nil, err
	}

	var categories []Categories
	for rows.Next() {
		var categorie Categories
		err := rows.Scan(&categorie.ID, &categorie.Name, &categorie.Parent_id)
		if err != nil {
			return nil, err
		}
		categories = append(categories, categorie)
	}

	return categories, nil
}

func GetUsersData() ([]Users, error) {
	rows, err := SQL.Query("SELECT id, email, username, password_hash, created_at FROM Users")
	if err != nil {
		return nil, err
	}

	var users []Users
	for rows.Next() {
		var user Users
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Password_hash, &user.Created_at)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetPostsData() ([]Posts, error) {
	rows, err := SQL.Query("SELECT id, user_id, category_id, title, content, created_at FROM Posts")
	if err != nil {
		return nil, err
	}

	var posts []Posts
	for rows.Next() {
		var post Posts
		err := rows.Scan(&post.ID, &post.User_id, &post.Category_id, &post.Title, &post.Content, &post.Created_at)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func GetCommentsData() ([]Comments, error) {
	rows, err := SQL.Query("SELECT id, post_id, user_id, content, created_at FROM Comments")
	if err != nil {
		return nil, err
	}

	var comments []Comments
	for rows.Next() {
		var comment Comments
		err := rows.Scan(&comment.ID, &comment.Post_id, &comment.User_id, &comment.Content, &comment.Created_at)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func GetLikesDislikesData() ([]LikesDislikes, error) {
	rows, err := SQL.Query("SELECT id, user_id, post_id, comment_id, type FROM Likes_Dislikes")
	if err != nil {
		return nil, err
	}

	var likesDislikes []LikesDislikes
	for rows.Next() {
		var likeDislike LikesDislikes
		err := rows.Scan(&likeDislike.ID, &likeDislike.User_id, &likeDislike.Post_id, &likeDislike.Comment_id, &likeDislike.TypeValue)
		if err != nil {
			return nil, err
		}
		likesDislikes = append(likesDislikes, likeDislike)
	}

	return likesDislikes, nil
}
