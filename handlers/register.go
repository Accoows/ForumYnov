package handlers

import (
	"forumynov/database"
	"html/template"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Afficher le formulaire
		tmpl, err := template.ParseFiles("Templates/register.html")
		if err != nil {
			log.Println("[handlers/register.go] Erreur chargement template :", err)
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Traiter les données du formulaire
		database.RegisterUsers(w, r)
		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}
