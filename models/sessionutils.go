package models

import (
	"database/sql"
	"forumynov/database"
	"net/http"
)

// GetUserIDFromRequest retrieves the user ID from the cookie if valid
func GetUserIDFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_cookie") // retrieve the cookie from the request using the cookie name
	if err != nil {
		return "", err
	}

	var userID string
	// query the database to get the user ID associated with the session ID from the cookie
	err = database.SQL.QueryRow("SELECT user_id FROM Sessions WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows { // no session found for this cookie
			return "", nil
		}
		return "", err
	}

	return userID, nil
}
