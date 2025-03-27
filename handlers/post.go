package handlers

import (
	"forumynov/database"
	"html/template"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("Templates/create_post.html")
		if err != nil {
			http.Error(w, "Erreur chargement formulaire", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur parsing du formulaire", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

		// Temporairement fixé à l'utilisateur ID 1
		userID := 1

		err = database.CreatePost(userID, categoryID, title, content)
		if err != nil {
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetCompletePostList()
	if err != nil {
		http.Error(w, "Erreur lors du chargement des posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("Templates/posts.html")
	if err != nil {
		http.Error(w, "Erreur chargement template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, struct {
		Posts []database.Posts
	}{Posts: posts})
}
