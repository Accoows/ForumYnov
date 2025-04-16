package handlers

import (
	"forumynov/database"
	"forumynov/models"
	"log"
	"net/http"
	"strconv"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("[handlers/comment.go] [CreateCommentHandler] Method not allowed >>>", r.Method)
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] ParseForm failed >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID, err := getConnectedUserID(r)
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] User not connected >>>", err)
		models.SetNotification(w, "You must be logged in to comment", "error")
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	content := r.FormValue("content")

	if err != nil || content == "" {
		log.Println("[handlers/comment.go] [CreateCommentHandler] Invalid data >>>", err)
		models.SetNotification(w, "Missing or invalid comment data", "error")
		return
	}

	err = database.CreateComment(userID, postID, content)
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] CreateComment error >>>", err)
		models.SetNotification(w, "Could not post your comment", "error")
		return
	}

	models.SetNotification(w, "Comment successfully posted", "success")
	http.Redirect(w, r, "/posts/view?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("[handlers/post.go] [DeleteCommentHandler] Method not allowed >>>", r.Method)
		return
	}

	userID, err := getConnectedUserID(r)
	if err != nil {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] User not connected >>>", err)
		http.Error(w, "Login required", http.StatusUnauthorized)
		return
	}

	idStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		log.Println("[handlers/post.go] [DeleteCommentHandler] Invalid ID >>>", err)
		return
	}

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] Error retrieving comment >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if comment.User_id != userID {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] Deletion denied for user:", userID)
		models.SetNotification(w, "You are not allowed to delete this comment", "error")
		return
	}

	postID := r.FormValue("post_id")

	err = database.DeleteCommentByID(commentID)
	if err != nil {
		log.Println("[handlers/post.go] [DeleteCommentHandler] Error deleting comment >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	models.SetNotification(w, "Comment successfully deleted", "success")
	http.Redirect(w, r, "/posts/view?id="+postID, http.StatusSeeOther)
}
