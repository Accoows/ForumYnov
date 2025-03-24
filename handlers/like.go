package handlers

import (
	"fmt"
	"net/http"
)

// Like/dislike posts et commentaires

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Like - à venir")
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Dislike - à venir")
}
