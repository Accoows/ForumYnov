package database

import "time"

// Creates a new comment and inserts it into the database
func CreateComment(userID string, postID int, content string) error {
	comment := &Comments{
		User_id:    userID,
		Post_id:    postID,
		Content:    content,
		Created_at: time.Now().Format("2006-01-02 15:04:05"), // SQL format
	}
	return InsertCommentsData(comment)
}

// Deletes all likes/dislikes associated with a specific comment
func DeleteLikesByCommentID(commentID int) error {
	// Remove all entries from Likes_Dislikes where the comment_id matches
	_, err := SQL.Exec("DELETE FROM Likes_Dislikes WHERE comment_id = ?", commentID)
	return err
}

// Deletes a comment by its ID, including its associated likes/dislikes
func DeleteCommentByID(id int) error {
	// First, delete all likes/dislikes tied to the comment
	err := DeleteLikesByCommentID(id)
	if err != nil {
		return err
	}
	// Then delete the comment itself from the Comments table
	_, err = SQL.Exec("DELETE FROM Comments WHERE id = ?", id)
	return err
}

// Retrieves a single comment by its ID
func GetCommentByID(id int) (Comments, error) {
	var comment Comments
	// Execute the query to get the comment data
	err := SQL.QueryRow(`
		SELECT id, post_id, user_id, content, created_at
		FROM Comments
		WHERE id = ?
	`, id).Scan(
		&comment.ID,
		&comment.Post_id,
		&comment.User_id,
		&comment.Content,
		&comment.Created_at,
	)
	return comment, err
}

// Retrieves all comments for a given post ID, with their author's username and like/dislike counts
func GetCommentsByPostID(postID int) ([]Comments, error) {
	// Query to get all comments for the given post, joined with user data
	rows, err := SQL.Query(`
		SELECT Comments.id, Comments.post_id, Comments.user_id, Comments.content, Comments.created_at,
		       Users.username
		FROM Comments
		JOIN Users ON Comments.user_id = Users.id			-- Get the username of the comment's author
		WHERE Comments.post_id = ?
		ORDER BY Comments.created_at ASC					-- Sort comments from oldest to newest
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		err := rows.Scan( // Map each column to the corresponding field in the struct
			&comment.ID,
			&comment.Post_id,
			&comment.User_id,
			&comment.Content,
			&comment.Created_at,
			&comment.AuthorUsername,
		)
		if err != nil {
			return nil, err
		}

		// Count the likes and dislikes of the comment directly using the CountLikesForComment function in query_likes.go
		comment.LikeCount, comment.DislikeCount, _ = CountLikesForComment(SQL, comment.ID)
		comments = append(comments, comment)
	}
	return comments, nil
}
