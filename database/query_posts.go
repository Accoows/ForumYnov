package database

import "time"

func CreatePost(userID string, categoryID int, title, content string) error {
	post := &Posts{
		User_id:     userID,
		Category_id: categoryID,
		Title:       title,
		Content:     content,
		Created_at:  time.Now().Format("2006-01-02 15:04:05"),
	}

	return InsertPostsData(post)
}

func GetCompletePostList() ([]Posts, error) {
	rows, err := SQL.Query(`
		SELECT Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at,
		       Users.username, Categories.name
		FROM Posts
		JOIN Users ON Posts.user_id = Users.id
		JOIN Categories ON Posts.category_id = Categories.id
		ORDER BY Posts.created_at DESC
	`)
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

func GetPostByID(id int) (Posts, error) {
	var post Posts
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

	post.LikeCount, post.DislikeCount, _ = CountLikesForPost(SQL, post.ID)

	return post, err
}

func UpdatePost(id int, title, content string) error {
	_, err := SQL.Exec(`
		UPDATE Posts SET title = ?, content = ? WHERE id = ?
	`, title, content, id)
	return err
}

func DeletePostWithDependencies(postID int) error {
	// Supprimer les likes associés au post
	_, err := SQL.Exec(`DELETE FROM Likes_Dislikes WHERE post_id = ?`, postID)
	if err != nil {
		return err
	}

	// Supprimer les likes associés aux commentaires de ce post
	_, err = SQL.Exec(`
		DELETE FROM Likes_Dislikes 
		WHERE comment_id IN (SELECT id FROM Comments WHERE post_id = ?)`, postID)
	if err != nil {
		return err
	}

	// Supprimer les commentaires associés au post
	_, err = SQL.Exec(`DELETE FROM Comments WHERE post_id = ?`, postID)
	if err != nil {
		return err
	}

	// Supprimer le post lui-même
	_, err = SQL.Exec(`DELETE FROM Posts WHERE id = ?`, postID)
	return err
}

func GetLatestPosts() ([]Posts, error) {
	query := `SELECT Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.name, Categories.category_photos, COALESCE(SUM(CASE WHEN Likes_Dislikes.type = 1 THEN 1 ELSE 0 END), 0) AS like_count,
    COALESCE(SUM(CASE WHEN Likes_Dislikes.type = -1 THEN 1 ELSE 0 END), 0) AS dislike_count FROM Posts JOIN Users ON Posts.user_id = Users.id JOIN Categories ON Posts.category_id = Categories.id LEFT JOIN Likes_Dislikes ON Posts.id = Likes_Dislikes.post_id 
	GROUP BY Posts.id, Posts.user_id, Posts.category_id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.name, Categories.category_photos ORDER BY Posts.created_at DESC LIMIT 3;`

	rows, err := SQL.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		err := rows.Scan(&post.ID, &post.User_id, &post.Category_id, &post.Title, &post.Content, &post.Created_at, &post.AuthorUsername, &post.CategoryName, &post.CategoryPhotos, &post.LikeCount, &post.DislikeCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
