package handlers

import (
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Home page handler — serves the main index page of the site
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // Check if the requested path is exactly "/"
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	// Try to get the user ID from the request (via session/cookie)
	userID, err := models.GetUserIDFromRequest(r)
	isLoggedIn := err == nil && userID != "" // true if a valid user ID is found

	data := struct { // Define data passed to the template
		IsLoggedIn bool
	}{
		IsLoggedIn: isLoggedIn,
	}

	// Load the index.html template
	tmpl, err := template.ParseFiles(filepath.Join("./templates/", "index.html"))
	if err != nil {
		log.Println("[handlers/index.go] [IndexHandler] Error loading template:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("[handlers/index.go] [IndexHandler] Error executing template:", err)
		ErrorHandler(w, http.StatusInternalServerError)
	}
}
