package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func main() {
	// Parse the HTML template file
	tmpl, err := template.ParseFiles(filepath.Join(".", "Test.html"))
	if err != nil {
		panic(err)
	}

	// Handler function to serve the template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Serve static files (CSS, images, etc.) from the current directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}
