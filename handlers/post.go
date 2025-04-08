package handlers

import (
	"forumynov/database"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(filepath.Join("./Templates", "create_post.html"))
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Erreur ParseFiles >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Erreur ParseForm >>>", err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

		// Temporairement fixé à l'utilisateur ID 1
		// Correspondra à l'UUID de l'utilisateur
		// Il faudra associer l'UUID au username pour le récupérer et l'afficher dans le post
		userID := "1"

		err = database.CreatePost(userID, categoryID, title, content)
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Erreur CreatePost >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

// =========================

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Erreur QueryID >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Erreur GetPostByID >>>", err)
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	comments, err := database.GetCommentsByPostID(postID)
	if err != nil {
		log.Println("[handlers/post.go][ViewPostHandler] Erreur chargement commentaires :", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	data := struct {
		Post     database.Posts
		Comments []database.Comments
	}{
		Post:     post,
		Comments: comments,
	}

	tmpl, err := template.ParseFiles(filepath.Join("./Templates", "view_post.html"))
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Erreur ParseFiles >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

// =========================

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("[handlers/post.go] [DeletePostHandler] Méthode non autorisée >>>", r.Method)
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] ID invalide >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	err = database.DeletePostByID(postID)
	if err != nil {
		log.Println("[handlers/post.go] [DeletePostHandler] Erreur suppression post >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

// =========================

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetCompletePostList()
	if err != nil {
		log.Println("[handlers/post.go] [PostsHandler] Erreur GetCompletePostList >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("./Templates", "post.html"))
	if err != nil {
		log.Println("[handlers/post.go] [PostsHandler] Erreur ParseFiles >>>", err)
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
		log.Println("[handlers/post.go] [EditPostHandler] ID invalide >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		post, err := database.GetPostByID(postID)
		if err != nil {
			log.Println("[handlers/post.go] [EditPostHandler] Post introuvable >>>", err)
			ErrorHandler(w, http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles(filepath.Join("./Templates", "edit_post.html"))
		if err != nil {
			log.Println("[handlers/post.go] [EditPostHandler] Erreur ParseFile >>>", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, post)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("[handlers/post.go][EditPostHandler] Erreur ParseForm >>>", err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		err = database.UpdatePost(postID, title, content)
		if err != nil {
			log.Println("[handlers/post.go][EditPostHandler] Erreur update BDD :", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts/view?id="+strconv.Itoa(postID), http.StatusSeeOther)
	}
}

// =========================

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("./Templates/", "post-detail.html"))
	if err != nil {
		log.Println(err)
		return
	}
	tmpl.Execute(w, nil)
}

func PostListHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("./Templates/", "post-list.html"))
	if err != nil {
		log.Println(err)
		return
	}
	tmpl.Execute(w, nil)
}
