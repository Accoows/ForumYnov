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

func CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		log.Println("[handlers/post.go][CategoryPostsHandler] ID invalide >>>", err)
		return
	}

	posts, err := database.GetPostsByCategoryID(categoryID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("[handlers/post.go][CategoryPostsHandler] Erreur GetPostsByCategoryID >>>", err)
		return
	}

	category, err := database.GetCategoryByID(categoryID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("[handlers/post.go][CategoryPostsHandler] Erreur GetCategoryByID >>>", err)
		return
	}

	data := database.CategoryPageData{
		CategoryName: category.Name,
		CategoryID:   categoryID,
		Posts:        posts,
	}

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
		log.Println("[handlers/post.go][CategoryPostsHandler] Erreur ParseFiles >>>", err)
		return
	}

	tmpl.Execute(w, data)
}
