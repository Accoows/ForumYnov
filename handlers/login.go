package handlers

import (
	"database/sql"
	"fmt"
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

var tpl *template.Template

func LoginUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error during form processing", http.StatusBadRequest)
		return
	}

	newUser := database.Users{
		Email:         r.FormValue("email"),
		Password_hash: r.FormValue("password"),
	}

	if newUser.Email == "" || newUser.Password_hash == "" {
		http.Error(w, "All fields are mandatory", http.StatusBadRequest)
		return
	}

	emailRight := VerifyEmailConformity(&newUser)
	if emailRight == false {
		http.Error(w, "Email does not comply", http.StatusBadRequest)
		return
	}
	userUUID, hashedPassword, err := GetHashedPasswordAndUUID(newUser.Email)
	if err != nil {
		http.Error(w, "Incorrect email", http.StatusNotFound)
		return
	}

	if !models.CheckPasswordHash(newUser.Password_hash, hashedPassword) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	newID := uuid.NewV4()

	newSessions := database.Sessions{
		ID:         newID.String(),
		User_id:    userUUID,
		Expires_at: time.Now().Add(24 * time.Hour),
	}

	err = database.InsertSessionsData(&newSessions)
	if err != nil {
		http.Error(w, "Error during Session database integration", http.StatusInternalServerError)
		return
	}

	userIdCookie := &http.Cookie{
		Name:       "user_cookie",
		Value:      newSessions.ID,
		Path:       "/",
		Domain:     "",
		Expires:    newSessions.Expires_at,
		RawExpires: "",
		MaxAge:     86400,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   []string{},
	}

	http.SetCookie(w, userIdCookie)

	http.Redirect(w, r, "/profile", http.StatusFound)
}

func GetHashedPasswordAndUUID(email string) (string, string, error) {
	var userUUID string
	var hashedPassword string
	query := "SELECT id, password_hash FROM Users WHERE email = ?"
	err := database.SQL.QueryRow(query, email).Scan(&userUUID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user found")
		}
		return "", "", err
	}

	return userUUID, hashedPassword, nil
}
