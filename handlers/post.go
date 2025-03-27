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
			http.Error(w, "Erreur chargement formulaire", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Erreur ParseForm >>>", err)
			http.Error(w, "Erreur parsing du formulaire", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

		// Temporairement fixé à l'utilisateur ID 1
		// Correspondra à l'UUID de l'utilisateur
		// Il faudra associer l'UUID au username pour le récupérer et l'afficher dans le post
		userID := 1

		err = database.CreatePost(userID, categoryID, title, content)
		if err != nil {
			log.Println("[handlers/post.go] [CreatePostHandler] Erreur CreatePost >>>", err)
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
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
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Erreur GetPostByID >>>", err)
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("./Templates", "view_post.html"))
	if err != nil {
		log.Println("[handlers/post.go] [ViewPostHandler] Erreur ParseFiles >>>", err)
		http.Error(w, "Erreur chargement template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, post)
}

// =========================

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		log.Println("[handlers/post.go] [DeletePostHandler] Méthode non autorisée >>>", r.Method)
		return
	}

	idStr := r.FormValue("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		log.Println("[handlers/post.go] [DeletePostHandler] ID invalide >>>", err)
		return
	}

	err = database.DeletePostByID(postID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		log.Println("[handlers/post.go] [DeletePostHandler] Erreur suppression post >>>", err)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

// =========================

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetCompletePostList()
	if err != nil {
		log.Println("[handlers/post.go] [PostsHandler] Erreur GetCompletePostList >>>", err)
		http.Error(w, "Erreur lors du chargement des posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("./Templates", "post.html"))
	if err != nil {
		log.Println("[handlers/post.go] [PostsHandler] Erreur ParseFiles >>>", err)
		http.Error(w, "Erreur chargement template", http.StatusInternalServerError)
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
		http.Error(w, "ID invalide", http.StatusBadRequest)
		log.Println("[handlers/post.go] [EditPostHandler] ID invalide >>>", err)
		return
	}

	if r.Method == http.MethodGet {
		post, err := database.GetPostByID(postID)
		if err != nil {
			http.Error(w, "Post introuvable", http.StatusNotFound)
			log.Println("[handlers/post.go] [EditPostHandler] Post introuvable >>>", err)
			return
		}

		tmpl, err := template.ParseFiles(filepath.Join("./Templates", "edit_post.html"))
		if err != nil {
			http.Error(w, "Erreur template", http.StatusInternalServerError)
			log.Println("[handlers/post.go] [EditPostHandler] Erreur ParseFile >>>", err)
			return
		}

		tmpl.Execute(w, post)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur parsing", http.StatusBadRequest)
			log.Println("[handlers/post.go][EditPostHandler] Erreur ParseForm >>>", err)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		err = database.UpdatePost(postID, title, content)
		if err != nil {
			http.Error(w, "Erreur update", http.StatusInternalServerError)
			log.Println("[handlers/post.go][EditPostHandler] Erreur update BDD :", err)
			return
		}

		http.Redirect(w, r, "/posts/view?id="+strconv.Itoa(postID), http.StatusSeeOther)
	}
}
