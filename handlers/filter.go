package handlers

import (
	"fmt"
	"net/http"
)

// Filtres : catégorie, créés, likés

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Filtrage - à venir")
}
