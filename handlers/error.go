package handlers

import (
	"forumynov/database"
	"html/template"
	"log"
	"net/http"
)

// Handles displaying a custom error page with the appropriate HTTP status code.
func ErrorHandler(w http.ResponseWriter, statusCode int) {
	// Get the standard HTTP status message (e.g., "Not Found", "Internal Server Error", etc.)
	msg := http.StatusText(statusCode)
	if msg == "" { // Fallback if the code is unknown
		statusCode = http.StatusInternalServerError
		msg = "Unknown error"
	}

	w.WriteHeader(statusCode) // Set the HTTP status code in the response

	tmpl, err := template.ParseFiles("./templates/error.html")
	if err != nil {
		log.Println("[handlers/error.go] [ErrorHandler] Error ParseFiles >>>", err)
		http.Error(w, msg, statusCode)
		return
	}

	data := database.ErrorPageData{ // Prepare the data to send to the template
		Code:    statusCode,
		Message: msg,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("[handlers/error.go] [ErrorHandler] Error Execute >>>", err)
		http.Error(w, msg, statusCode)
		return
	}
}
