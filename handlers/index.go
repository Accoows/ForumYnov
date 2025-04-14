package handlers

import (
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

	tmpl, err := template.ParseFiles(filepath.Join("./templates/", "index.html"))
	if err != nil {
		log.Println("[handlers/index.go] [IndexHandler] Erreur de chargement du template :", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("[handlers/index.go] [IndexHandler] Erreur d'exécution du template :", err)
		ErrorHandler(w, http.StatusInternalServerError)
	}
}
