package database

import (
	"database/sql"
)

// Checks if a like or dislike already exists by a user for a specific post or comment
func GetExistingLikeDislike(userID string, postID, commentID int) (*LikesDislikes, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, type
		FROM Likes_Dislikes
		WHERE user_id = ?
		  AND ((? IS NOT NULL AND post_id = ? AND comment_id IS NULL)
			OR (? IS NOT NULL AND comment_id = ? AND post_id IS NULL))
	`
	// We check if the user has liked either a post (if commentID is null) or a comment (if postID is null)
	row := SQL.QueryRow(query,
		userID,
		sqlNullableInt(postID), postID,
		sqlNullableInt(commentID), commentID,
	)

	var like LikesDislikes
	err := row.Scan(&like.ID, &like.User_id, &like.Post_id, &like.Comment_id, &like.TypeValue)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Return the existing like/dislike if found, or nil if not found
	return &like, nil
}

// Converts an integer to a SQL-compatible nullable value (nil if 0)
func sqlNullableInt(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

// Counts the number of likes and dislikes for a given post.
func CountLikesForPost(db *sql.DB, postID int) (likeCount, dislikeCount int, err error) {
	query := `
		SELECT
			SUM(CASE WHEN type = 1 THEN 1 ELSE 0 END),
			SUM(CASE WHEN type = -1 THEN 1 ELSE 0 END)
		FROM Likes_Dislikes
		WHERE post_id = ? AND comment_id IS NULL
	`
	// We sum all the "likes" (type = 1) and "dislikes" (type = -1) associated with the post

	row := db.QueryRow(query, postID)
	err = row.Scan(&likeCount, &dislikeCount)
	return
}

// Counts the number of likes and dislikes for a given comment
func CountLikesForComment(db *sql.DB, commentID int) (likeCount, dislikeCount int, err error) {
	query := `
		SELECT
			SUM(CASE WHEN type = 1 THEN 1 ELSE 0 END),
			SUM(CASE WHEN type = -1 THEN 1 ELSE 0 END)
		FROM Likes_Dislikes
		WHERE comment_id = ? AND post_id IS NULL
	`
	// Same logic as for posts, but here we count the likes/dislikes associated with a comment

	row := db.QueryRow(query, commentID)
	err = row.Scan(&likeCount, &dislikeCount)
	return
}
