package database

import (
	"database/sql"
)

func GetExistingLikeDislike(userID string, postID, commentID int) (*LikesDislikes, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, type
		FROM Likes_Dislikes
		WHERE user_id = ?
		  AND ((? IS NOT NULL AND post_id = ? AND comment_id IS NULL)
		    OR (? IS NOT NULL AND comment_id = ? AND post_id IS NULL))
	`

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

	return &like, nil
}

func sqlNullableInt(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

func CountLikesForPost(db *sql.DB, postID int) (likeCount, dislikeCount int, err error) {
	query := `
		SELECT
			SUM(CASE WHEN type = 1 THEN 1 ELSE 0 END),
			SUM(CASE WHEN type = -1 THEN 1 ELSE 0 END)
		FROM Likes_Dislikes
		WHERE post_id = ? AND comment_id IS NULL
	`
	row := db.QueryRow(query, postID)
	err = row.Scan(&likeCount, &dislikeCount)
	return
}

func CountLikesForComment(db *sql.DB, commentID int) (likeCount, dislikeCount int, err error) {
	query := `
		SELECT
			SUM(CASE WHEN type = 1 THEN 1 ELSE 0 END),
			SUM(CASE WHEN type = -1 THEN 1 ELSE 0 END)
		FROM Likes_Dislikes
		WHERE comment_id = ? AND post_id IS NULL
	`
	row := db.QueryRow(query, commentID)
	err = row.Scan(&likeCount, &dislikeCount)
	return
}
