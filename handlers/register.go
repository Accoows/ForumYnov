package handlers

import (
	"fmt"
	"forumynov/database"
	"forumynov/models"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func RegisterUsers(w http.ResponseWriter, r *http.Request) {
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
		Email:      r.FormValue("email"),
		Username:   r.FormValue("username"),
		Created_at: time.Now().Format("2006-01-02 15:04:05"),
	}
	password := r.FormValue("password")

	if newUser.Email == "" || newUser.Username == "" || password == "" {
		http.Error(w, "All fields are mandatory", http.StatusBadRequest)
		return
	}

	emailRight := VerifyEmailConformity(&newUser)
	if emailRight == false {
		http.Error(w, "Email does not comply", http.StatusBadRequest)
		return
	}

	exists, err := VerifyEmailAndUsernameUnicity(newUser.Email, newUser.Username)
	if err != nil {
		http.Error(w, "Error during user verification", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email or username already used", http.StatusConflict)
		return
	}

	password_hash, err := models.HashPassword(password)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password_hash = password_hash

	newID := uuid.NewV4()
	newUser.ID = newID.String()

	err = database.InsertUsersData(&newUser)
	if err != nil {
		http.Error(w, "Error during user registration", http.StatusInternalServerError)
		return
	}

	// Afficher les valeurs dans la console
	fmt.Println("Username:", newUser.Username)
	fmt.Println("Email:", newUser.Email)
	fmt.Println("Password:", newUser.Password_hash)

	// Répondre au client
	w.Write([]byte("Succesfull inscription !"))
}

func VerifyEmailAndUsernameUnicity(email string, username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM Users WHERE email = ? OR username = ?)"
	err := database.SQL.QueryRow(query, email, username).Scan(&exists)
	if err != nil {
		log.Println("Erreur SQL dans VerifyEmailAndUsernameUnicity:", err)
		return false, err
	}
	return exists, nil
}

func VerifyEmailConformity(users *database.Users) bool {
	email := []rune(users.Email)
	emailLenght := len(email)
	emailRight := false
	for i := 0; i < emailLenght; i++ {
		if email[i] == '@' {
			for j := i; j < emailLenght; j++ {
				if email[j] == '.' {
					emailRight = true
				}
			}
		}
	}
	return emailRight
}
