package models

import (
	"database/sql"
	"forumynov/database"
	"net/http"
)

func VerifyCookieValidity(r *http.Request, userID string) (bool, error) {
	var cookieName string
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
