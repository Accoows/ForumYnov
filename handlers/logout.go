package handlers

import (
	"forumynov/database"
	"net/http"
	"time"
)

// LogoutUsers handles the logout process for users.
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

	logout := r.FormValue("logout") // get the logout value from the form data
	if logout != "true" {
		return
	}

	cookie, err := r.Cookie("user_cookie") // get the cookie from the request using the cookie name
	if err != nil {
		http.Error(w, "Unauthorized : cookie not found", http.StatusUnauthorized)
		return
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour) // set the expiration time of the cookie to a past time to invalidate it
	http.SetCookie(w, cookie)                       // set the cookie in the response to invalidate it

	err = database.DeleteSession(cookie.Value) // delete the session from the database using the session ID from the cookie

	if err != nil {
		http.Error(w, "Error during Session database erasure", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound) // redirect the user to the home page after logout
}
