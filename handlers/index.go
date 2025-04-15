package handlers

import (
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Page d'accueil, général

// Gestionnaire pour servir la page index
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	userID, err := models.GetUserIDFromRequest(r)
	isLoggedIn := err == nil && userID != ""

	data := struct {
		IsLoggedIn bool
	}{
		IsLoggedIn: isLoggedIn,
	}

	tmpl, err := template.ParseFiles(filepath.Join("./templates/", "index.html"))
	if err != nil {
		log.Println("[handlers/index.go] [IndexHandler] Erreur de chargement du template :", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("[handlers/index.go] [IndexHandler] Erreur d'exécution du template :", err)
		ErrorHandler(w, http.StatusInternalServerError)
	}
}
