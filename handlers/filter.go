package handlers

import (
	"forumynov/database"
	"html/template"
	"net/http"
	"path/filepath"
)

// FilterPostsByCategories handles the filtering of posts by categories.
func FilterPostsByCategories(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Query().Get("category") // Get the category name from the URL query parameters
	if categoryName == "" {                       // If no category is provided, redirect to the posts page
		ErrorHandler(w, http.StatusBadRequest)
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	categoryID, err := database.GetCategoryIDByName(categoryName) // Get the category ID from the database based on the category name
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de la catégorie", http.StatusInternalServerError)
		return
	}

	posts, err := database.GetPostsByCategory(categoryID) // Get the posts associated with the category ID
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("./templates", "post.html")) // Parse the HTML template for displaying posts
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, struct { //Execute the template with the posts and category name
		Posts    []database.Posts
		Category string
	}{
		Posts:    posts,
		Category: categoryName,
	})
}
