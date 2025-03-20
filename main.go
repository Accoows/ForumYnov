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

	// Routes
	http.HandleFunc("/", indexHandler)

	port := ":8080"
	fmt.Println("Serveur lancé sur http://localhost" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}
