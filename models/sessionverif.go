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
	var sessionID string
	query := "SELECT id FROM Sessions WHERE user_id = ? AND expires_at > datetime('now')"
	err := database.SQL.QueryRow(query, userID).Scan(&sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // no session as been found for this user
		}
		return false, err // another error as occurred
	}

	cookie, err := r.Cookie("user_cookie")
	if err != nil {
		return false, nil // cookie as not been found or as expired
	}

	if cookie.Value != sessionID {
		return false, nil // cookie does not match the session ID
	}

	return true, nil // cookie is valid
}
