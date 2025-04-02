package database

import "time"

func CreateComment(userID string, postID int, content string) error {
	comment := &Comments{
		User_id:    userID,
		Post_id:    postID,
		Content:    content,
		Created_at: time.Now().Format("2006-01-02 15:04:05"),
	}
	return InsertCommentsData(comment)
}

func DeleteCommentByID(id int) error {
	_, err := SQL.Exec("DELETE FROM Comments WHERE id = ?", id)
	return err
}

func GetCommentsByPostID(postID int) ([]Comments, error) {
	rows, err := SQL.Query(`
		SELECT Comments.id, Comments.post_id, Comments.user_id, Comments.content, Comments.created_at,
		       Users.username
		FROM Comments
		JOIN Users ON Comments.user_id = Users.id
		WHERE Comments.post_id = ?
		ORDER BY Comments.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		err := rows.Scan(
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
		comments = append(comments, comment)
	}
	return comments, nil
}
