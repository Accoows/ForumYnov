package database

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

func GetLikedPostsByUser(userID string) ([]Posts, error) {
	var posts []Posts
	query := `SELECT Posts.id, Posts.title, Posts.content, Posts.created_at FROM Posts
			  JOIN Likes ON Posts.id = Likes.post_id WHERE Likes.user_id = ?`
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

func GetDislikedPostsByUser(userID string) ([]Posts, error) {
	var posts []Posts
	query := `SELECT Posts.id, Posts.title, Posts.content, Posts.created_at FROM Posts
			  JOIN Dislikes ON Posts.id = Dislikes.post_id WHERE Dislikes.user_id = ?`
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
