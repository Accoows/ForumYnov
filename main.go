package main

import (
	"fmt"
	"forumynov/database"
	"forumynov/handlers"
	"log"
	"net/http"
)

func main() {
	database.InitTempDB()
	defer database.DB.Close()

	// Handler pour les fichiers statiques (ex: CSS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/reset-password", handlers.ResetPasswordHandler)
	http.HandleFunc("/forgot-username", handlers.ForgotUsernameHandler)

	http.HandleFunc("/posts", handlers.PostsHandler)
	http.HandleFunc("/posts/create", handlers.CreatePostHandler)
	http.HandleFunc("/posts/delete", handlers.DeletePostHandler)
	http.HandleFunc("/posts/edit", handlers.EditPostHandler)
	http.HandleFunc("/posts/view", handlers.ViewPostHandler)

	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/dislike", handlers.DislikeHandler)

	http.HandleFunc("/filter", handlers.FilterHandler)

	port := ":8080"
	fmt.Println("Serveur lancé sur http://localhost" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}
