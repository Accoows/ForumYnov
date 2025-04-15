package database

// GetPostsByUser retrieves all posts created by a specific user.
func GetPostsByUser(userID string) ([]Posts, error) {
	var posts []Posts
	query := `SELECT id, title, content, created_at FROM Posts WHERE user_id = ?`
	rows, err := SQL.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Posts
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_at)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetPostsByCategory retrieves all posts associated with a specific category and its subcategories.
func GetLikedPostsByUser(userID string) ([]Posts, error) {
	var posts []Posts
	query := `SELECT Posts.id, Posts.title, Posts.content, Posts.created_at FROM Posts JOIN Likes_Dislikes ON Posts.id = Likes_Dislikes.post_id WHERE Likes_Dislikes.user_id = ? AND Likes_Dislikes.type = 1`
	rows, err := SQL.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Posts
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_at)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
