package handlers

import (
	"forumynov/database"
	"log"
	"net/http"
	"strconv"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		log.Println("[handlers/comment.go] [CreateCommentHandler] Méthode non autorisée >>>", r.Method)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur de parsing", http.StatusBadRequest)
		log.Println("[handlers/comment.go] [CreateCommentHandler] ParseForm échoué >>>", err)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	userID, err2 := strconv.Atoi(r.FormValue("user_id"))
	content := r.FormValue("content")

	if err != nil || err2 != nil || content == "" {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		log.Println("[handlers/comment.go] [CreateCommentHandler] Données invalides >>>", err, err2)
		return
	}

	err = database.CreateComment(userID, postID, content)
	if err != nil {
		http.Error(w, "Erreur en base de données", http.StatusInternalServerError)
		log.Println("[handlers/comment.go] [CreateCommentHandler] Erreur CreateComment >>", err)
		return
	}

	http.Redirect(w, r, "/posts/view?id="+strconv.Itoa(postID), http.StatusSeeOther)
}
