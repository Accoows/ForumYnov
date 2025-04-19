package handlers

import (
	"forumynov/database"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// Handles the request to view posts by category
func CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("id")       // Get the category ID from the query parameter
	categoryID, err := strconv.Atoi(categoryIDStr) // Convert the ID from string to int
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		log.Println("[handlers/post.go][CategoryPostsHandler] Invalid ID >>>", err)
		return
	}

	posts, err := database.GetPostsByCategoryID(categoryID) // Fetch posts belonging to this category
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("[handlers/post.go][CategoryPostsHandler] Error GetPostsByCategoryID >>>", err)
		return
	}

	category, err := database.GetCategoryByID(categoryID) // Fetch the category details (name, etc.)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("[handlers/post.go][CategoryPostsHandler] Error GetCategoryByID >>>", err)
		return
	}

	data := database.CategoryPageData{
		CategoryName: category.Name,
		CategoryID:   categoryID,
		Posts:        posts,
	}

	// Create and configure the template, adding a formatter to convert line breaks
	tmpl := template.New("post-detail.html").Funcs(template.FuncMap{
		"format": func(s string) template.HTML {
			escaped := template.HTMLEscapeString(s)
			withBreaks := strings.ReplaceAll(escaped, "\n", "<br>")
			return template.HTML(withBreaks)
		},
	})
	tmpl, err = tmpl.ParseFiles(filepath.Join("./templates/", "post-detail.html"))
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("[handlers/post.go][CategoryPostsHandler] Error ParseFiles >>>", err)
		return
	}

	tmpl.Execute(w, data)
}
