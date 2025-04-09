package handlers

import (
	"forumynov/database"
	"html/template"
	"net/http"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_cookie")
	if err != nil {
		http.Error(w, "Unauthorized : cookie not found", http.StatusUnauthorized)
		return
	}

	var userID string
	err = database.SQL.QueryRow("SELECT user_id FROM Sessions WHERE cookie_name = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var user database.Users
	err = database.SQL.QueryRow("SELECT username, email FROM Users WHERE id = ?", userID).Scan(&user.Username, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data := struct {
		Username string
	}{
		Username: user.Username,
	}

	tmpl, err := template.ParseFiles("templates/edit-profile.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
