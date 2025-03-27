package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Gestionnaire pour servir la page de login
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("./Templates/", "register.html"))
	if err != nil {
		log.Println(err)
		return
	}
	tmpl.Execute(w, nil)
}
