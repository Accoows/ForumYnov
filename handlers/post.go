package handlers

import (
	"forumynov/models"
	"html/template"
	"net/http"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("Templates/create_post.html")
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.FormValue("title")
		content := r.FormValue("content")

		// On met une catégorie par défaut (ex : ID 1) ou NULL si géré dans SQL
		categoryID := 1

		// TEMPORAIRE : ID utilisateur en dur
		userID := 1

		err := models.CreatePost(userID, categoryID, title, content)
		if err != nil {
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}
