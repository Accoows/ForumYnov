package handlers

import (
	"forumynov/database"
	"net/http"
	"time"
)

func LogoutUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error during form processing", http.StatusBadRequest)
		return
	}

	logout := r.FormValue("logout")
	if logout != "true" {
		return
	}

	cookie, err := r.Cookie("user_cookie")
	if err != nil {
		http.Error(w, "Unauthorized : cookie not found", http.StatusUnauthorized)
		return
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, cookie)

	err = database.DeleteSession(cookie.Value)

	if err != nil {
		http.Error(w, "Error during Session database erasure", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
