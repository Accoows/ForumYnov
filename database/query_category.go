package database

// Retrieves a single category from the database using its ID.
func GetCategoryByID(id int) (Categories, error) {
	var c Categories
	query := `SELECT id, name, parent_id FROM Categories WHERE id = ?`
	// Execute the query and store the result in the 'c' struct
	err := SQL.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.ParentID)
	return c, err
}

// Retrieves all posts associated with a specific category ID.
// Joins with Users and Categories tables to include author's username and category name.
func GetPostsByCategoryID(categoryID int) ([]Posts, error) {
	// Execute the SQL query with JOINs to get complete post info
	rows, err := SQL.Query(`
		SELECT Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at,
		       Users.username, Categories.name
		FROM Posts
		JOIN Users ON Posts.user_id = Users.id					-- Link each post to its author's username
		JOIN Categories ON Posts.category_id = Categories.id	-- Link each post to its category name
		WHERE Posts.category_id = ?
		ORDER BY Posts.created_at DESC							-- Order posts from newest to oldest
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		err := rows.Scan(
			// Scan each row into the 'post' struct
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

// Fetches all categories from the database, ordered by their ID.
func GetAllCategories() ([]Categories, error) {
	// Query to fetch all categories ordered by ID
	rows, err := SQL.Query("SELECT id, name, parent_id FROM Categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Categories
	for rows.Next() {
		var cat Categories
		// Extract column values into the category struct
		err := rows.Scan(&cat.ID, &cat.Name, &cat.ParentID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}
