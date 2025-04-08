package main

import (
	"fmt"
	"forumynov/database"
	"forumynov/handlers"
	"net/http"
)

func main() {
	database.InitDatabase()
	defer database.CloseDatabase()

	// Gère les requêtes vers le dossier "Scripts", de manière similaire au dossier "Styles".
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))

	// Serve static files (CSS, images, etc.) from the current directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/Templates/", http.StripPrefix("/Templates/", http.FileServer(http.Dir("./Templates"))))

	// ========================

	// Routes

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	//http.HandleFunc("/logout", handlers.LogoutHandler)
	//http.HandleFunc("/reset-password", handlers.ResetPasswordHandler)
	//http.HandleFunc("/forgot-username", handlers.ForgotUsernameHandler)

	// CRUD pour les posts
	http.HandleFunc("/posts", handlers.PostsHandler)
	http.HandleFunc("/posts/create", handlers.CreatePostHandler)
	http.HandleFunc("/posts/view", handlers.ViewPostHandler)
	http.HandleFunc("/posts/edit", handlers.EditPostHandler)
	http.HandleFunc("/posts/delete", handlers.DeletePostHandler)
	//http.HandleFunc("/post-detail", handlers.PostDetailHandler)
	http.HandleFunc("/category", handlers.CategoryPostsHandler)
	http.HandleFunc("/post-list", handlers.PostListHandler)

	// CRUD pour les commentaires
	// L'affichage des commentaires est géré dans la page de post (ViewPostHandler)
	http.HandleFunc("/comments/create", handlers.CreateCommentHandler)
	http.HandleFunc("/comments/delete", handlers.DeleteCommentHandler)

	http.HandleFunc("/like", handlers.LikeHandler)

	//http.HandleFunc("/filter", handlers.FilterHandler)

	fmt.Println("Starting server at port 8080")
	fmt.Println(">>>> http://localhost:8080 <<<<")

	// Démarre le serveur HTTP sur le port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
