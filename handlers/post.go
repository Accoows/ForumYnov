package handlers

import (
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// Helper to retrieve user ID from session or cookie
func getConnectedUserID(r *http.Request) (string, error) {
	return models.GetUserIDFromRequest(r)
}

// Handles both GET (form display) and POST (form submission) for post creation
// CRUD for Posts
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := models.GetUserIDFromRequest(r)
	if err != nil || userID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login if not authenticated
		return
	}

	if r.Method == http.MethodGet {
		categoryIDStr := r.URL.Query().Get("category_id") // Handle optional preselected category
		categoryID := 0
		var categoryName string

		if categoryIDStr != "" {
			id, err := strconv.Atoi(categoryIDStr)
			if err == nil {
				categoryID = id
				category, err := database.GetCategoryByID(id) // Fetch the category from the database using its ID
				if err == nil {
					categoryName = category.Name // If successful, store the category name
				}
			}
		}

		// Retrieve the full list of categories for the dropdown in the form
		allCategories, err := database.GetAllCategories()
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Error GetAllCategories >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		data := database.CreatePostPageData{ // Pass data to template
			CategoryID:    categoryID,
			CategoryName:  categoryName,
			AllCategories: allCategories,
		}

		tmpl, err := template.ParseFiles(filepath.Join("./templates", "create_post.html"))
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Error ParseFiles >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm() // Handle form submission
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Error ParseForm >>>", err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

		err = database.CreatePost(userID, categoryID, title, content) // Insert post into DB
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Error CreatePost >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		// Notify success with a popup and redirect to the category page
		models.SetNotification(w, "Post successfully created", "success")
		http.Redirect(w, r, "/category?id="+strconv.Itoa(categoryID), http.StatusSeeOther)
	}
}

// =========================

// Displays a post and its associated comments
func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Error QueryID >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID) // Fetch the post from the database using its ID
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Error GetPostByID >>>", err)
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	comments, err := database.GetCommentsByPostID(postID) // Fetch comments associated with the post
	if err != nil {
		log.Println("[handlers/post.go][ViewPostHandler] Error loading comments:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	userID, _ := getConnectedUserID(r) // Retrieve the logged-in user ID from the request
	isAuthor := userID == post.User_id // Check if the logged-in user is the author of the post

	data := database.ViewPostPageData{ // Prepare data for the template
		Post:          post,
		Comments:      comments,
		IsAuthor:      isAuthor,
		ConnectedUser: userID,
	}

	// Create and configure the template, adding a formatter to convert line breaks
	tmpl := template.New("view_post.html").Funcs(template.FuncMap{
		"format": func(s string) template.HTML {
			escaped := template.HTMLEscapeString(s)
			withBreaks := strings.ReplaceAll(escaped, "\n", "<br>")
			return template.HTML(withBreaks)
		},
	})
	tmpl, err = tmpl.ParseFiles(filepath.Join("./templates/", "view_post.html"))
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Error ParseFiles >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

// =========================

// Deletes a post only if the connected user is its author
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("[handlers/post.go] [DeletePostHandler] Unauthorized method >>>", r.Method)
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id") // Get the post ID from the form
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] Invalid ID >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID) // Fetch the post from the database using its ID
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] Post not found >>>", err)
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	userID, err := getConnectedUserID(r) // Retrieve the logged-in user ID from the request
	if err != nil || userID != post.User_id {
		log.Println("[handlers/post.go] [DeletePostHandler] Forbidden access for user:", userID)
		models.SetNotification(w, "Unauthorized deletion", "error")
		return
	}

	err = database.DeletePostWithDependencies(postID) // Delete post and related data (comments + likes)
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] Error deleting post >>>", err)
		models.SetNotification(w, "Error deleting post", "error")
		return
	}

	// Notify user of success and redirect to the category page
	models.SetNotification(w, "Post successfully deleted", "success")
	http.Redirect(w, r, "/category?id="+strconv.Itoa(post.Category_id), http.StatusSeeOther)
}

// =========================

// Displays the full list of posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetCompletePostList() // Fetch all posts from the database
	if err != nil {
		log.Println("[handlers/post.go] [PostsHandler] Error GetCompletePostList >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("./templates", "post.html"))
	if err != nil {
		log.Println("[handlers/post.go] [PostsHandler] Error ParseFiles >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, struct {
		Posts []database.Posts // List of posts to be displayed
	}{Posts: posts})
}

// =========================

// Handles post edition (GET = display form, POST = apply changes)
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // Get the post ID from the query parameter
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [EditPostHandler] Invalid ID >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID) // Fetch the post from the database using its ID
	if err != nil {
		log.Println("[handlers/post.go] [EditPostHandler] Post not found >>>", err)
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	userID, err := getConnectedUserID(r)      // Retrieve the logged-in user ID from the request
	if err != nil || userID != post.User_id { // Check if the logged-in user is the author of the post
		log.Println("[handlers/post.go] [EditPostHandler] Forbidden access for user:", userID)
		models.SetNotification(w, "Unauthorized modification", "error")
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(filepath.Join("./templates", "edit_post.html"))
		if err != nil {
			log.Println("[handlers/post.go] [EditPostHandler] Error ParseFile >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, post)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm() // Handle form submission
		if err != nil {
			log.Println("[handlers/post.go][EditPostHandler] Error ParseForm >>>", err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		err = database.UpdatePost(postID, title, content) // Update the post in the database
		if err != nil {
			log.Println("[handlers/post.go][EditPostHandler] Error updating DB:", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		// Notify user of success and redirect to the post view page
		models.SetNotification(w, "Post successfully updated", "info")
		http.Redirect(w, r, "/posts/view?id="+strconv.Itoa(postID), http.StatusSeeOther)
	}
}

// =========================

func PostListHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := models.GetUserIDFromRequest(r)
	isLoggedIn := err == nil && userID != ""

	topCategories, err := database.GetMostsPostsCategoriesOfTheWeek()
	if err != nil {
		http.Error(w, "Error fetching top categories", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsLoggedIn    bool
		TopCategories []database.Categories
	}{
		IsLoggedIn:    isLoggedIn,
		TopCategories: topCategories,
	}

	tmpl, err := template.ParseFiles(filepath.Join("./templates/", "post-list.html"))
	if err != nil {
		log.Println(err)
		return
	}
	tmpl.Execute(w, data)
}
