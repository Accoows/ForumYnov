package handlers

import (
	"database/sql"
	"forumynov/database"
	"forumynov/models"
	"log"
	"net/http"
	"strconv"
)

// LikeHandler handles likes/dislikes for a post or comment
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	userID, err := getConnectedUserID(r)
	if err != nil {
		models.SetNotification(w, "You must be logged in to like or dislike", "error")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// Retrieve form values
	postIDStr := r.FormValue("post_id")
	commentIDStr := r.FormValue("comment_id")
	action := r.FormValue("action") // "like" or "dislike" with SQL values 1 or -1

	postID, _ := strconv.Atoi(postIDStr)
	commentID, _ := strconv.Atoi(commentIDStr)

	var typeValue int
	switch action {
	case "like":
		typeValue = 1
	case "dislike":
		typeValue = -1
	default:
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	// Existing structure in datastruct.go for LikesDislikes
	like := &database.LikesDislikes{
		User_id:    userID,
		Post_id:    sql.NullInt64{Int64: int64(postID), Valid: postID != 0},
		Comment_id: sql.NullInt64{Int64: int64(commentID), Valid: commentID != 0},
		TypeValue:  typeValue,
	}

	err = InsertOrUpdateLikeDislike(like) // Update likes/dislikes
	if err != nil {
		log.Println("[HandleLike] Error InsertOrUpdateLikeDislike:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	// Redirect after the like
	if postID != 0 {
		http.Redirect(w, r, "/posts/view?id="+postIDStr, http.StatusSeeOther)
	} else {
		// In the case of a comment, the parent_post_id must be provided in the form TO VERIFY
		parentPostID := r.FormValue("parent_post_id")
		http.Redirect(w, r, "/posts/view?id="+parentPostID, http.StatusSeeOther)
	}
}

// InsertOrUpdateLikeDislike inserts or updates a like/dislike in the database
func InsertOrUpdateLikeDislike(like *database.LikesDislikes) error {
	postID := 0
	commentID := 0

	if like.Post_id.Valid {
		postID = int(like.Post_id.Int64)
	}
	if like.Comment_id.Valid {
		commentID = int(like.Comment_id.Int64)
	}

	existing, err := database.GetExistingLikeDislike(like.User_id, postID, commentID)
	if err != nil {
		log.Println("[InsertOrUpdateLikeDislike] Error GetExistingLikeDislike:", err)
		return err
	}

	if existing == nil {
		row := database.SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Likes_Dislikes")
		err = row.Scan(&like.ID)
		if err != nil {
			log.Println("[InsertOrUpdateLikeDislike] Error Scan ID:", err)
			return err
		}

		_, err = database.SQL.Exec(`
			INSERT INTO Likes_Dislikes(id, user_id, post_id, comment_id, type)
			VALUES (?, ?, ?, ?, ?)`,
			like.ID, like.User_id, like.Post_id, like.Comment_id, like.TypeValue)
		if err != nil {
			log.Println("[InsertOrUpdateLikeDislike] Error INSERT:", err)
		}
		return err
	}

	if existing.TypeValue == like.TypeValue {
		_, err = database.SQL.Exec(`DELETE FROM Likes_Dislikes WHERE id = ?`, existing.ID)
		if err != nil {
			log.Println("[InsertOrUpdateLikeDislike] Error DELETE toggle:", err)
		}
		return err
	}

	_, err = database.SQL.Exec(`UPDATE Likes_Dislikes SET type = ? WHERE id = ?`, like.TypeValue, existing.ID)
	if err != nil {
		log.Println("[InsertOrUpdateLikeDislike] Error UPDATE:", err)
	}
	return err
}
