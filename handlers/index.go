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
	tmpl, err := template.ParseFiles(filepath.Join("./Templates/", "index.html"))
	if err != nil {
		log.Println(err)
		return
	}
	tmpl.Execute(w, nil)
}
