package handlers

import (
	"database/sql"
	"forumynov/database"
	"log"
	"net/http"
	"strconv"
)

// LikeHandler gère les likes/dislikes d'un post ou d'un commentaire
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	userID := "1" // Remplacer par la logique de récupération de l'utilisateur connecté avec les cookies/sessions
	if userID == "" {
		http.Error(w, "Vous devez être connecté pour liker", http.StatusUnauthorized) // à vérifier
		return
	}

	// Récupérer les valeurs du formulaire
	postIDStr := r.FormValue("post_id")
	commentIDStr := r.FormValue("comment_id")
	action := r.FormValue("action") // "like" ou "dislike" avec en SQL 1 ou -1

	postID, _ := strconv.Atoi(postIDStr)
	commentID, _ := strconv.Atoi(commentIDStr)

	var typeValue int
	switch action {
	case "like":
		typeValue = 1
	case "dislike":
		typeValue = -1
	default:
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	// Structure existante dans datastruct.go de LikesDislikes
	like := &database.LikesDislikes{
		User_id:    userID,
		Post_id:    sql.NullInt64{Int64: int64(postID), Valid: postID != 0},
		Comment_id: sql.NullInt64{Int64: int64(commentID), Valid: commentID != 0},
		TypeValue:  typeValue,
	}

	err := InsertOrUpdateLikeDislike(like) // Mise à jour des likes/dislikes
	if err != nil {
		log.Println("[HandleLike] Erreur InsertOrUpdateLikeDislike:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	// Redirection après le like
	if postID != 0 {
		http.Redirect(w, r, "/posts/view?id="+postIDStr, http.StatusSeeOther)
	} else {
		// Dans le cas d’un commentaire, il faut que le parent_post_id soit fourni dans le formulaire A VERIFIER
		parentPostID := r.FormValue("parent_post_id")
		http.Redirect(w, r, "/posts/view?id="+parentPostID, http.StatusSeeOther)
	}
}

// InsertOrUpdateLikeDislike insère ou met à jour un like/dislike dans la base de données
func InsertOrUpdateLikeDislike(like *database.LikesDislikes) error {
	postID := 0
	commentID := 0

	if like.Post_id.Valid {
		postID = int(like.Post_id.Int64)
	}
	if like.Comment_id.Valid {
		commentID = int(like.Comment_id.Int64)
	}

	existing, err := database.GetExistingLikeDislike(like.User_id, postID, commentID)
	if err != nil {
		log.Println("[InsertOrUpdateLikeDislike] Erreur GetExistingLikeDislike:", err)
		return err
	}

	if existing == nil {
		row := database.SQL.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM Likes_Dislikes")
		err = row.Scan(&like.ID)
		if err != nil {
			log.Println("[InsertOrUpdateLikeDislike] Erreur Scan ID:", err)
			return err
		}

		_, err = database.SQL.Exec(`
			INSERT INTO Likes_Dislikes(id, user_id, post_id, comment_id, type)
			VALUES (?, ?, ?, ?, ?)`,
			like.ID, like.User_id, like.Post_id, like.Comment_id, like.TypeValue)
		if err != nil {
			log.Println("[InsertOrUpdateLikeDislike] Erreur INSERT:", err)
		}
		return err
	}

	if existing.TypeValue == like.TypeValue {
		_, err = database.SQL.Exec(`DELETE FROM Likes_Dislikes WHERE id = ?`, existing.ID)
		if err != nil {
			log.Println("[InsertOrUpdateLikeDislike] Erreur DELETE toggle:", err)
		}
		return err
	}

	_, err = database.SQL.Exec(`UPDATE Likes_Dislikes SET type = ? WHERE id = ?`, like.TypeValue, existing.ID)
	if err != nil {
		log.Println("[InsertOrUpdateLikeDislike] Erreur UPDATE:", err)
	}
	return err
}
