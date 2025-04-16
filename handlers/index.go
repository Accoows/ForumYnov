package handlers

import (
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Home page, general

// Handler to serve the index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	userID, err := models.GetUserIDFromRequest(r)
	isLoggedIn := err == nil && userID != ""

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

	data := struct {
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
