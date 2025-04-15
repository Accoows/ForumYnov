package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/*
ErrorTest is a helper function to check for errors and log them.
It's used to handle errors that may occur during database operations
*/
func ErrorTest(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var SQL *sql.DB // SQL is a global variable that holds the database connection
var err error   // err is a global variable that holds the error returned by database operations

// InitDatabase initializes the database connection
func InitDatabase() {
	SQL, err = sql.Open("sqlite3", "database/forum.db") // Open() is used to open the database file

	ErrorTest(err)
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if SQL != nil {
		SQL.Close() // Close() is used to close the database connection
	}
}

/* Inserts functions allows to inserts data into the sql tables */

func InsertUsersData(users *Users) error {
	// INSERT OR IGNORE is used to insert the data in the tables avoiding inserting duplicate entries
	insertUsersInSql := `INSERT OR IGNORE INTO Users(id, email, username, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`

	// sql.Exec() is used to execute the SQL statement
	_, err = SQL.Exec(insertUsersInSql, users.ID, users.Email, users.Username, users.Password_hash, users.Created_at)

	ErrorTest(err)

	return err // err is returned to the caller to check if the operation was successful
}

func InsertSessionsData(sessions *Sessions) error {
	insertSessionsSql := `INSERT OR IGNORE INTO Sessions(id, user_id, expires_at) VALUES (?, ?, ?)`

	_, err = SQL.Exec(insertSessionsSql, sessions.ID, sessions.User_id, sessions.Expires_at)

	ErrorTest(err)

	return err
}

func InsertPostsData(posts *Posts) error {
	/*
		sql.QueryRow() is used to execute a query that returns a single row.
		COALESCE() is used to return the first non-null value in the list, so if there are no rows in the table, it will return 0 + 1 = 1
	*/
	row := SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Posts")

	err = row.Scan(&posts.ID) // Scan() is used to read the result of the query into the variable given as an argument

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

/* Get functions allows to get data from the sql tables */

func GetCategoriesData() ([]Categories, error) {
	rows, err := SQL.Query("SELECT id, name, parent_id FROM Categories") // SELECT is used to select data from the table
	if err != nil {
		return nil, err
	}

	var categories []Categories // Categories is a slice of Categories struct
	for rows.Next() {           // rows.Next() is used to iterate over the rows returned by the query
		var categorie Categories // categorie is a variable of type Categories struct
		err := rows.Scan(&categorie.ID, &categorie.Name, &categorie.ParentID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, categorie) // append() is used to add the categorie to the slice
	}

	return categories, nil // the data is returned to the caller
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

func GetSessionsData() ([]Sessions, error) {
	rows, err := SQL.Query("SELECT id, user_id, expires_at FROM Sessions")
	if err != nil {
		return nil, err
	}

	var sessions []Sessions
	for rows.Next() {
		var session Sessions
		err := rows.Scan(&session.ID, &session.User_id, &session.Expires_at)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
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

// Delete functions allows to delete data from the sql tables
func DeleteSession(cookieName string) error {
	_, err := SQL.Exec("DELETE FROM Sessions WHERE id = ?", cookieName) // DELETE is used to delete data from the table
	if err != nil {
		return err
	}
	return nil
}

// DeleteExpiredSessions deletes expired sessions from the Sessions table
func DeleteExpiredSessions() {
	query := "DELETE FROM Sessions WHERE expires_at < ?"
	_, err := SQL.Exec(query, time.Now()) // Deletes sessions where the expiration date is less than the current time
	if err != nil {
		log.Printf("Error deleting expired sessions: %v\n", err)
	}
}
