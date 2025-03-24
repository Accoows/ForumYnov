package handlers

import (
	"fmt"
	"net/http"
)

// Création, affichage des posts

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des posts - à venir")
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Formulaire de création de post")
		return
	}
	if r.Method == http.MethodPost {
		fmt.Fprintln(w, "Post créé (simulation)")
		return
	}
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.URL.Query().Get("id")
		fmt.Fprintf(w, "Post ID %s supprimé (simulation)", id)
		return
	}
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
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
	fmt.Fprintln(w, "Visualisation d'un post - à venir")
}
