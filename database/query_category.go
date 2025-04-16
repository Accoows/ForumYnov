package database

func GetCategoryByID(id int) (Categories, error) {
	var c Categories
	query := `SELECT id, name, parent_id FROM Categories WHERE id = ?`
	err := SQL.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.ParentID)
	return c, err
}

func GetPostsByCategoryID(categoryID int) ([]Posts, error) {
	rows, err := SQL.Query(`
		SELECT Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at,
		       Users.username, Categories.name
		FROM Posts
		JOIN Users ON Posts.user_id = Users.id
		JOIN Categories ON Posts.category_id = Categories.id
		WHERE Posts.category_id = ?
		ORDER BY Posts.created_at DESC
	`, categoryID)
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

func GetAllCategories() ([]Categories, error) {
	rows, err := SQL.Query("SELECT id, name, parent_id FROM Categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Categories
	for rows.Next() {
		var cat Categories
		err := rows.Scan(&cat.ID, &cat.Name, &cat.ParentID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func GetMostsPostsCategoriesOfTheWeek() ([]Categories, error) {
	query := `SELECT Categories.id, Categories.name, Categories.category_photos FROM Categories JOIN Posts ON Categories.id = Posts.category_id
    WHERE Posts.created_at >= datetime('now', '-7 days') GROUP BY Categories.id, Categories.name, Categories.category_photos ORDER BY COUNT(Posts.id) DESC LIMIT 3;`

	rows, err := SQL.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Categories
	if !rows.Next() {
		return categories, nil
	}

	for {
		var category Categories
		err := rows.Scan(&category.ID, &category.Name, &category.CategoryPhotos)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)

		if !rows.Next() {
			break
		}
	}

	return categories, nil
}

func GetMostsPostsCategories() ([]Categories, error) {
	query := `SELECT Categories.id, Categories.name, Categories.category_photos FROM Categories JOIN Posts ON Categories.id = Posts.category_id
    GROUP BY Categories.id, Categories.name, Categories.category_photos ORDER BY COUNT(Posts.id) DESC LIMIT 3;`

	rows, err := SQL.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Categories
	if !rows.Next() {
		return categories, nil // Retourne une liste vide
	}

	for {
		var category Categories
		err := rows.Scan(&category.ID, &category.Name, &category.CategoryPhotos)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)

		if !rows.Next() {
			break
		}
	}

	return categories, nil
}
