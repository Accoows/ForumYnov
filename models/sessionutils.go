package models

import (
	"database/sql"
	"forumynov/database"
	"net/http"
)

// GetUserIDFromRequest retrieves the user ID from the cookie if valid
func GetUserIDFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_cookie")
	if err != nil {
		return "", err
	}

	var userID string
	err = database.SQL.QueryRow("SELECT user_id FROM Sessions WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return userID, nil
}
