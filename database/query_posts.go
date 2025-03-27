package database

import "time"

func CreatePost(userID, categoryID int, title, content string) error {
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
		err := rows.Scan(&post.ID, &post.User_id, &post.Category_id, &post.Title, &post.Content, &post.Created_at, &post.AuthorUsername, &post.CategoryName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
