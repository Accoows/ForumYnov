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

func getConnectedUserID(r *http.Request) (string, error) {
	return models.GetUserIDFromRequest(r)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := models.GetUserIDFromRequest(r)
	if err != nil || userID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		// Read the "category_id" parameter if it exists
		categoryIDStr := r.URL.Query().Get("category_id")
		categoryID := 0
		var categoryName string

		if categoryIDStr != "" {
			id, err := strconv.Atoi(categoryIDStr)
			if err == nil {
				categoryID = id
				category, err := database.GetCategoryByID(id)
				if err == nil {
					categoryName = category.Name
				}
			}
		}

		allCategories, err := database.GetAllCategories()
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Error GetAllCategories >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		data := database.CreatePostPageData{
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
		err := r.ParseForm()
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Error ParseForm >>>", err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

		err = database.CreatePost(userID, categoryID, title, content)
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Error CreatePost >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		models.SetNotification(w, "Post successfully created", "success")
		http.Redirect(w, r, "/category?id="+strconv.Itoa(categoryID), http.StatusSeeOther)
	}
}

// =========================

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Error QueryID >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Error GetPostByID >>>", err)
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	comments, err := database.GetCommentsByPostID(postID)
	if err != nil {
		log.Println("[handlers/post.go][ViewPostHandler] Error loading comments:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	userID, _ := getConnectedUserID(r)
	isAuthor := userID == post.User_id

	data := database.ViewPostPageData{
		Post:          post,
		Comments:      comments,
		IsAuthor:      isAuthor,
		ConnectedUser: userID,
	}

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

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("[handlers/post.go] [DeletePostHandler] Unauthorized method >>>", r.Method)
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] Invalid ID >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] Post not found >>>", err)
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	userID, err := getConnectedUserID(r)
	if err != nil || userID != post.User_id {
		log.Println("[handlers/post.go] [DeletePostHandler] Forbidden access for user:", userID)
		http.Error(w, "Unauthorized deletion", http.StatusUnauthorized)
		return
	}

	err = database.DeletePostWithDependencies(postID)
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] Error deleting post >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	models.SetNotification(w, "Post successfully deleted", "success")
	http.Redirect(w, r, "/category?id="+strconv.Itoa(post.Category_id), http.StatusSeeOther)
}

// =========================

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetCompletePostList()
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
		Posts []database.Posts
	}{Posts: posts})
}

// =========================

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [EditPostHandler] Invalid ID >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		log.Println("[handlers/post.go] [EditPostHandler] Post not found >>>", err)
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	userID, err := getConnectedUserID(r)
	if err != nil || userID != post.User_id {
		log.Println("[handlers/post.go] [EditPostHandler] Forbidden access for user:", userID)
		http.Error(w, "Unauthorized modification", http.StatusUnauthorized)
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
		err := r.ParseForm()
		if err != nil {
			log.Println("[handlers/post.go][EditPostHandler] Error ParseForm >>>", err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		err = database.UpdatePost(postID, title, content)
		if err != nil {
			log.Println("[handlers/post.go][EditPostHandler] Error updating DB:", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

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
