package handlers

import (
	"forumynov/database"
	"forumynov/models"
	"log"
	"net/http"
	"strconv"
)

// Handles the creation of a new comment (only via POST request)
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Check request method
		log.Println("[handlers/comment.go] [CreateCommentHandler] Method not allowed >>>", r.Method)
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm() // Parse form values
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] ParseForm failed >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID, err := getConnectedUserID(r) // Retrieve logged-in user ID from session
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] User not connected >>>", err)
		models.SetNotification(w, "You must be logged in to comment", "error")
		return
	}

	// Extract post ID and comment content from form
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	content := r.FormValue("content")

	if err != nil || content == "" {
		log.Println("[handlers/comment.go] [CreateCommentHandler] Invalid data >>>", err)
		models.SetNotification(w, "Missing or invalid comment data", "error")
		return
	}

	err = database.CreateComment(userID, postID, content) // Create comment in the database
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] CreateComment error >>>", err)
		models.SetNotification(w, "Could not post your comment", "error")
		return
	}

	// Notify user of success and redirect back to the post
	models.SetNotification(w, "Comment successfully posted", "success")
	http.Redirect(w, r, "/posts/view?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

// Handles deletion of a comment, only if the user is the author
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("[handlers/post.go] [DeleteCommentHandler] Method not allowed >>>", r.Method)
		return
	}

	userID, err := getConnectedUserID(r) // Check if the user is logged in
	if err != nil {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] User not connected >>>", err)
		http.Error(w, "Login required", http.StatusUnauthorized)
		return
	}

	idStr := r.FormValue("comment_id")    // Get comment ID from form
	commentID, err := strconv.Atoi(idStr) // Convert the comment ID string to an integer
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		log.Println("[handlers/post.go] [DeleteCommentHandler] Invalid ID >>>", err)
		return
	}

	comment, err := database.GetCommentByID(commentID) // Fetch comment from DB
	if err != nil {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] Error retrieving comment >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if comment.User_id != userID { // Check if the logged-in user is the author of the comment
		log.Println("[handlers/comment.go] [DeleteCommentHandler] Deletion denied for user:", userID)
		models.SetNotification(w, "You are not allowed to delete this comment", "error")
		return
	}

	postID := r.FormValue("post_id") // Retrieve post ID to redirect correctly after deletion

	err = database.DeleteCommentByID(commentID) // Delete the comment from DB
	if err != nil {
		log.Println("[handlers/post.go] [DeleteCommentHandler] Error deleting comment >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	// Notify user and redirect
	models.SetNotification(w, "Comment successfully deleted", "success")
	http.Redirect(w, r, "/posts/view?id="+postID, http.StatusSeeOther)
}
