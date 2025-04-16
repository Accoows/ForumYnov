package main

import (
	"fmt"
	"forumynov/database"
	"forumynov/handlers"
	"net/http"
	"time"
)

func main() {
	database.InitDatabase()
	defer database.CloseDatabase()

	go func() { // Periodically delete expired sessions
		for {
			database.DeleteExpiredSessions() // delete expired sessions from the database
			time.Sleep(5 * time.Minute)      // wait for 5 minutes before checking again
		}
	}()

	// Handles requests to the "Scripts" folder, similar to the "Styles" folder
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))

	// Serve static files (CSS, images, etc.) from the current directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	http.Handle("/static/uploads/", http.StripPrefix("/static/uploads/", http.FileServer(http.Dir("static/uploads"))))

	// ========================

	// Routes

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutUsers)
	http.HandleFunc("/profile", handlers.ProfilePage)
	http.HandleFunc("/upload-profile-picture", handlers.UploadProfilePictureHandler)
	http.HandleFunc("/delete-profile-picture", handlers.DeleteProfilePictureHandler)

	// CRUD for posts
	http.HandleFunc("/posts", handlers.PostsHandler)
	http.HandleFunc("/posts/create", handlers.CreatePostHandler)
	http.HandleFunc("/posts/view", handlers.ViewPostHandler)
	http.HandleFunc("/posts/edit", handlers.EditPostHandler)
	http.HandleFunc("/posts/delete", handlers.DeletePostHandler)
	http.HandleFunc("/posts/filter", handlers.FilterPostsByCategories)
	http.HandleFunc("/category", handlers.CategoryPostsHandler)
	http.HandleFunc("/post-list", handlers.PostListHandler)

	// CRUD for comments
	// Comment display is handled in the post page (ViewPostHandler)
	http.HandleFunc("/comments/create", handlers.CreateCommentHandler)
	http.HandleFunc("/comments/delete", handlers.DeleteCommentHandler)

	http.HandleFunc("/like", handlers.LikeHandler)

	fmt.Println("Starting server at port 8080")
	fmt.Println(">>>> http://localhost:8080 <<<<")

	// Starts the HTTP server on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
