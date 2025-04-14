package handlers

import (
	"forumynov/database"
<<<<<<< HEAD
=======
	"forumynov/models"
>>>>>>> main
	"html/template"
	"net/http"
)

// ProfilePage handles the profile page for users.
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_cookie") // get the cookie from the request using the cookie name
	if err != nil {
<<<<<<< HEAD
		http.Error(w, "Unauthorized : cookie not found", http.StatusUnauthorized)
=======
		models.SetNotification(w, "You are logged out", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
>>>>>>> main
		return
	}

	var userID string
	// get the user ID from the session using the session ID from the cookie
	err = database.SQL.QueryRow("SELECT user_id FROM Sessions WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var user database.Users
	// get the user data from the database using the user ID
	err = database.SQL.QueryRow("SELECT username, email FROM Users WHERE id = ?", userID).Scan(&user.Username, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

<<<<<<< HEAD
	// Get posts created by the user
	createdPosts, err := database.GetPostsByUser(userID)
	if err != nil {
		http.Error(w, "Error fetching created posts", http.StatusInternalServerError)
		return
	}

	// Get posts liked by the user
	likedPosts, err := database.GetLikedPostsByUser(userID)
	if err != nil {
		http.Error(w, "Error fetching liked posts", http.StatusInternalServerError)
		return
	}

	// data to be passed to the template
	data := struct {
		Username     string
		Email        string
		CreatedPosts []database.Posts
		LikedPosts   []database.Posts
	}{
		Username:     user.Username,
		Email:        user.Email,
		CreatedPosts: createdPosts,
		LikedPosts:   likedPosts,
=======
	data := struct { // data to be passed to the template
		Username string
	}{
		Username: user.Username,
>>>>>>> main
	}

	tmpl, err := template.ParseFiles("templates/edit-profile.html") // parse the HTML template file
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data) // execute the template with the data and write it to the response
}
