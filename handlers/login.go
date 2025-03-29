package handlers

import (
	"database/sql"
	"fmt"
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Page d'accueil, général

// Gestionnaire pour servir la page de login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(filepath.Join("./Templates/", "login.html"))
		if err != nil {
			log.Println(err)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	if r.Method == http.MethodPost {
		// Traiter les données du formulaire
		LoginUsers(w, r)
		return
	}
}

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
	if !emailRight {
		http.Error(w, "Email does not comply", http.StatusBadRequest)
		return
	}
	hashedPassword, err := GetHashedPassword(newUser.Email)
	if err != nil {
		http.Error(w, "Incorrect email", http.StatusNotFound)
		return
	}

	if models.CheckPasswordHash(newUser.Password_hash, hashedPassword) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successful connection !"))
	} else {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

}

func GetHashedPassword(email string) (string, error) {
	var hashedPassword string
	query := "SELECT password_hash FROM Users WHERE email = ?"
	err := database.SQL.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user found")
		}
		return "", err
	}

	return hashedPassword, nil
}
