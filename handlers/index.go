package handlers

import (
	"forumynov/database"
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

	topCategories, err := database.GetMostsPostsCategoriesOfTheWeek()
	if err != nil {
		http.Error(w, "Error fetching top categories", http.StatusInternalServerError)
		return
	}

	topCategoriesAllTime, err := database.GetMostsPostsCategories()
	if err != nil {
		http.Error(w, "Error fetching all time top categories", http.StatusInternalServerError)
		return
	}

	latestPosts, err := database.GetLatestPosts() // Récupérer les 3 derniers posts
	if err != nil {
		http.Error(w, "Error fetching latest posts", http.StatusInternalServerError)
		return
	}

	data := struct { // Define data passed to the template
		IsLoggedIn           bool
		TopCategories        []database.Categories
		TopCategoriesAllTime []database.Categories
		LatestPosts          []database.Posts
	}{
		IsLoggedIn:           isLoggedIn,
		TopCategories:        topCategories,
		TopCategoriesAllTime: topCategoriesAllTime,
		LatestPosts:          latestPosts,
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
