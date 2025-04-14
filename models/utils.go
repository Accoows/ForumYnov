package models

import (
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword crée un hash à partir du mot de passe brut
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compare un mot de passe brut avec son hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SetNotification définit un cookie temporaire contenant un message de notification
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
