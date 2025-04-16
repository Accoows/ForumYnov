package models

import (
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a hash from the raw password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a raw password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SetNotification sets a temporary cookie containing a notification message
func SetNotification(w http.ResponseWriter, message string, notifType string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "notif_msg",
		Value:  url.QueryEscape(message),
		Path:   "/",
		MaxAge: 5,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "notif_type",
		Value:  notifType, // "success", "info", "error"
		Path:   "/",
		MaxAge: 5,
	})
}
