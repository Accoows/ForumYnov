package main

import (
	"fmt"
	"forumynov/database"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

func openBrowser(url string) {
	var err error
	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

func main() {
	database.Database()

	// Parse the HTML template file
	tmpl, err := template.ParseFiles(filepath.Join("./Templates/", "mobile.html"))
	if err != nil {
		panic(err)
	}
	tmplLogin, err := template.ParseFiles(filepath.Join("./Templates/", "login.html"))
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
	// Handler function to serve the login template
	http.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		err := tmplLogin.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Gère les requêtes vers le dossier "Scripts", de manière similaire au dossier "Styles".
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))

	// Serve static files (CSS, images, etc.) from the current directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	// Ouvre automatiquement le navigateur web à l'adresse "http://localhost:8080"
	fmt.Println("Starting server at port 8080")

	const url = "http://localhost:8080"
	// Démarre le serveur HTTP sur le port 8080
	go openBrowser("http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
