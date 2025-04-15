package database

// GetPostsByCategory retrieves all posts associated with a specific category and its subcategories.
func GetPostsByCategory(categoryID int) ([]Posts, error) {
	var posts []Posts

	//RECURSIVE query is used to get all subcategories of the given category
	//UNION ALL is used to combine the results of the main category and its subcategories
	//INNER JOIN is used to join the Posts table with the Categories table to get the posts associated with the categories
	query := `WITH RECURSIVE Subcategories AS (SELECT id FROM Categories WHERE id = ? UNION ALL SELECT c.id FROM Categories c INNER JOIN Subcategories s ON c.parent_id = s.id)
    SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Posts.user_id, Categories.name AS category_name
    FROM Posts JOIN Categories ON Posts.category_id = Categories.id WHERE Posts.category_id IN (SELECT id FROM Subcategories)`
	rows, err := SQL.Query(query, categoryID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Posts
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_at, &post.User_id, &post.CategoryName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetCategoryIDByName retrieves the ID of a category based on its name.
func GetCategoryIDByName(categoryName string) (int, error) {
	var categoryID int
	query := `SELECT id FROM Categories WHERE name = ?`
	err := SQL.QueryRow(query, categoryName).Scan(&categoryID)
	if err != nil {
		return 0, err
	}
	return categoryID, nil
}
