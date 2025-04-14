package handlers

import (
	"forumynov/database"
<<<<<<< HEAD
=======
	"forumynov/models"
>>>>>>> main
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
<<<<<<< HEAD
		http.Error(w, "Unauthorized : cookie not found", http.StatusUnauthorized)
=======
		models.SetNotification(w, "You have been logged out", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther) // redirect the user to the home page if the cookie is not found
>>>>>>> main
		return
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour) // set the expiration time of the cookie to a past time to invalidate it
	http.SetCookie(w, cookie)                       // set the cookie in the response to invalidate it

	err = database.DeleteSession(cookie.Value) // delete the session from the database using the session ID from the cookie

	if err != nil {
		http.Error(w, "Error during Session database erasure", http.StatusInternalServerError)
		return
	}

<<<<<<< HEAD
=======
	models.SetNotification(w, "You have been logged out", "error") // set a notification message to inform the user about the logout status

>>>>>>> main
	http.Redirect(w, r, "/", http.StatusFound) // redirect the user to the home page after logout
}
