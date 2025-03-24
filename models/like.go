package models

import (
	"database/sql"
	"errors"
	"time"
)

// Système de like/dislike

type LikeDislike struct {
	ID        int
	UserID    int
	PostID    *int
	CommentID *int
	Like      bool // true si c'est un like
	Dislike   bool // true si c'est un dislike, mais les 2 ne sont pas possibles
	CreatedAt string
}

// AddLikeDislike enregistre un vote (like ou dislike) pour un post ou un commentaire
func AddLikeDislike(db *sql.DB, userID int, postID *int, commentID *int, like bool, dislike bool) error {
	// Validation logique
	if like && dislike {
		return errors.New("impossible d'aimer et de ne pas aimer en même temps")
	}
	if !like && !dislike {
		return errors.New("vote vide : like ou dislike requis")
	}
	if (postID == nil && commentID == nil) || (postID != nil && commentID != nil) {
		return errors.New("le vote doit concerner un post OU un commentaire")
	}

	// Insertion du vote dans la table
	// EN ATTENTE DE LA MISE EN PLACE DE LA TABLE SQL
	_, err := db.Exec(`
		INSERT INTO likes_dislikes (user_id, post_id, comment_id, like, dislike, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, postID, commentID, like, dislike, time.Now().Format(time.RFC3339))

	return err
}

// GetLikeCount retourne le nombre de likes pour un post ou un commentaire
// EN ATTENTE DE BASE SQL
func GetLikeCount(db *sql.DB, postID *int, commentID *int) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM likes_dislikes
		WHERE like = 1 AND post_id IS ?
		AND comment_id IS ?
	`, postID, commentID).Scan(&count)
	return count, err
}

// GetDislikeCount retourne le nombre de dislikes pour un post ou un commentaire
// EN ATTENTE DE BASE SQL
func GetDislikeCount(db *sql.DB, postID *int, commentID *int) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM likes_dislikes
		WHERE dislike = 1 AND post_id IS ?
		AND comment_id IS ?
	`, postID, commentID).Scan(&count)
	return count, err
}
