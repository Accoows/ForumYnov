package handlers

import (
	"fmt"
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Afficher le formulaire
		tmpl, err := template.ParseFiles("Templates/register.html")
		if err != nil {
			log.Println("[handlers/register.go] Erreur chargement template :", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Traiter les données du formulaire
		RegisterUsers(w, r)
		return
	}

	ErrorHandler(w, http.StatusMethodNotAllowed)
}

func RegisterUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	newUser := database.Users{
		Email:      r.FormValue("email"),
		Username:   r.FormValue("username"),
		Created_at: time.Now().Format("2006-01-02 15:04:05"),
	}
	password := r.FormValue("password")

	if newUser.Email == "" || newUser.Username == "" || password == "" {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	emailRight := VerifyEmailConformity(&newUser)
	if !emailRight {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	exists, err := VerifyEmailAndUsernameUnicity(newUser.Email, newUser.Username)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	if exists {
		ErrorHandler(w, http.StatusConflict)
		return
	}

	password_hash, err := models.HashPassword(password)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	newUser.Password_hash = password_hash

	newID := uuid.NewV4()
	newUser.ID = newID.String()

	err = database.InsertUsersData(&newUser)
	if err != nil {
		http.Error(w, "Error during user registration", http.StatusInternalServerError)
		//ErrorHandler(w, http.StatusInternalServerError) SUREMENT UNE POPUP D'ERREUR
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
