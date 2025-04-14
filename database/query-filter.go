package database

func GetPostsByCategory(category string) ([]Posts, error) {
	var posts []Posts
	query := `SELECT id, title, content, created_at FROM Posts WHERE category = ?`
	rows, err := SQL.Query(query, category)
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
