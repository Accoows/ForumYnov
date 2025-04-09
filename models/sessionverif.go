package models

import (
	"database/sql"
	"forumynov/database"
	"net/http"
)

/*
VerifyCookieValidity checks if the cookie is valid for the given user ID.
It returns true if the cookie is valid, false if it is not found or expired, and an error if there was an issue querying the database.
*/
func VerifyCookieValidity(r *http.Request, userID string) (bool, error) {
	var cookieName string
	// get the cookie name from the database using the user ID
	err := database.SQL.QueryRow("SELECT cookie_name FROM Sessions WHERE user_id = ?", userID).Scan(&cookieName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // no session as been found for this user
		}
		return false, err // another error as occurred
	}

	if _, err := r.Cookie(cookieName); err != nil {
		return false, nil // cookie as not been found or as expired
	}

	return true, nil // cookie is valid
}
