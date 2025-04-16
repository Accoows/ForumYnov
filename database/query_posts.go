package database

import "time"

// Creates a new post and inserts it into the database.
func CreatePost(userID string, categoryID int, title, content string) error {
	post := &Posts{
		User_id:     userID,
		Category_id: categoryID,
		Title:       title,
		Content:     content,
		Created_at:  time.Now().Format("2006-01-02 15:04:05"), // SQL format
	}

	return InsertPostsData(post) // Insert post into the database
}

// Retrieves the complete list of posts, with their authors and categories.
func GetCompletePostList() ([]Posts, error) {
	rows, err := SQL.Query(`
		SELECT Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at,
		       Users.username, Categories.name
		FROM Posts
		JOIN Users ON Posts.user_id = Users.id					-- Link each post to its author's username
		JOIN Categories ON Posts.category_id = Categories.id	-- Link post to category name
		ORDER BY Posts.created_at DESC							-- Order from newest to oldest
	`) // Execute SQL query to fetch posts with user and category data
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		err := rows.Scan(
			&post.ID,
			&post.User_id,
			&post.Category_id,
			&post.Title,
			&post.Content,
			&post.Created_at,
			&post.AuthorUsername,
			&post.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Retrieves a single post by its ID, including author, category, and like/dislike counts.
func GetPostByID(id int) (Posts, error) {
	var post Posts
	// Fetch the post with joined user and category info
	err := SQL.QueryRow(`
		SELECT Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at,
		       Users.username, Categories.name
		FROM Posts
		JOIN Users ON Posts.user_id = Users.id
		JOIN Categories ON Posts.category_id = Categories.id
		WHERE Posts.id = ?
	`, id).Scan(
		&post.ID,
		&post.User_id,
		&post.Category_id,
		&post.Title,
		&post.Content,
		&post.Created_at,
		&post.AuthorUsername,
		&post.CategoryName,
	)

	// Adds like/dislike counts to the post, which will be directly updated on the HTML page
	post.LikeCount, post.DislikeCount, _ = CountLikesForPost(SQL, post.ID)

	return post, err
}

// Updates a post's title and content.
func UpdatePost(id int, title, content string) error {
	_, err := SQL.Exec(`
		UPDATE Posts SET title = ?, content = ? WHERE id = ?
	`, title, content, id)
	return err
}

// Deletes a post and all related data: likes, comment likes, comments.
func DeletePostWithDependencies(postID int) error {
	// Delete likes linked to the post
	_, err := SQL.Exec(`DELETE FROM Likes_Dislikes WHERE post_id = ?`, postID)
	if err != nil {
		return err
	}

	// Delete likes linked to comments of this post
	_, err = SQL.Exec(`
		DELETE FROM Likes_Dislikes 
		WHERE comment_id IN (SELECT id FROM Comments WHERE post_id = ?)`, postID)
	if err != nil {
		return err
	}

	// Delete comments linked to the post
	_, err = SQL.Exec(`DELETE FROM Comments WHERE post_id = ?`, postID)
	if err != nil {
		return err
	}

	// Delete the post itself
	_, err = SQL.Exec(`DELETE FROM Posts WHERE id = ?`, postID)
	return err
}

func GetLatestPosts() ([]Posts, error) {
	query := `SELECT Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.name, Categories.category_photos 
	FROM Posts JOIN Users ON Posts.user_id = Users.id JOIN Categories ON Posts.category_id = Categories.id ORDER BY Posts.created_at DESC LIMIT 3;`

	rows, err := SQL.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		err := rows.Scan(&post.ID, &post.User_id, &post.Category_id, &post.Title, &post.Content, &post.Created_at, &post.AuthorUsername, &post.CategoryName, &post.CategoryPhotos)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
