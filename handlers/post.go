package handlers

import (
	"fmt"
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Création, affichage des posts

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	// Données mockées (fictives)
	posts, err := models.GetAllPosts(database.SQL)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
		return
	}

	// Chargement du template
	tmpl, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		log.Println("Erreur template :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	// Injection des données dans le template
	err = tmpl.Execute(w, struct {
		Posts []models.Post
	}{Posts: posts})
	if err != nil {
		log.Println("Erreur rendering :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
}
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Affiche le formulaire de création
		tmpl, err := template.ParseFiles("templates/create_post.html")
		if err != nil {
			log.Println("Erreur template création :", err)
			http.Error(w, "Erreur interne", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)

	case http.MethodPost:
		// Traitement du formulaire
		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author") // temporaire

		if title == "" || content == "" || author == "" {
			http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
			return
		}

		// TEMP : récupérer user_id fictif à partir du nom
		// Plus tard ce sera via la session utilisateur
		userID := author // on suppose que c’est un UUID pour test

		err := models.CreatePost(database.SQL, userID, title, content)
		if err != nil {
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	err = models.DeletePost(database.SQL, id)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("id")
		fmt.Fprintf(w, "Formulaire d’édition pour le post ID : %s", id)
		return
	}
	if r.Method == http.MethodPost {
		id := r.URL.Query().Get("id")
		fmt.Fprintf(w, "Post ID %s édité (simulation)", id)
		return
	}
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	post, err := models.GetPostByID(database.SQL, id)
	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Println("Erreur template post :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, post)
}
