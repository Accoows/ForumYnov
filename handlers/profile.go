package handlers

import (
	"forumynov/database"
	"html/template"
	"net/http"
)

// ProfilePage handles the profile page for users.
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_cookie") // get the cookie from the request using the cookie name
	if err != nil {
		http.Error(w, "Unauthorized : cookie not found", http.StatusUnauthorized)
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

	data := struct { // data to be passed to the template
		Username string
	}{
		Username: user.Username,
	}

	tmpl, err := template.ParseFiles("templates/edit-profile.html") // parse the HTML template file
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data) // execute the template with the data and write it to the response
}
