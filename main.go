package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Gestionnaire pour servir la page index
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {

	// Handler pour les fichiers statiques (ex: CSS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", indexHandler)

	port := ":8080"
	fmt.Println("Serveur lancé sur http://localhost" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}
