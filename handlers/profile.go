package handlers

import (
	"database/sql"
	"fmt"
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ProfilePage handles the profile page for users.
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_cookie") // get the cookie from the request using the cookie name
	if err != nil {
		models.SetNotification(w, "You are logged out", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var userID string
	// get the user ID from the session using the session ID from the cookie
	err = database.SQL.QueryRow("SELECT user_id FROM Sessions WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		models.SetNotification(w, "Invalid session", "error")
		return
	}

	var user database.Users
	var profilePicture sql.NullString
	// get the user data from the database using the user ID
	err = database.SQL.QueryRow("SELECT username, email, profilepicture FROM Users WHERE id = ?", userID).Scan(&user.Username, &user.Email, &profilePicture)
	if err != nil {
		models.SetNotification(w, "User not found", "error")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	} else {
		user.ProfilePicture = ""
	}

	// Get posts created by the user
	createdPosts, err := database.GetPostsByUser(userID)
	if err != nil {
		http.Error(w, "Error fetching created posts", http.StatusInternalServerError)
		return
	}

	// Get posts liked by the user
	likedPosts, err := database.GetLikedPostsByUser(userID)
	if err != nil {
		http.Error(w, "Error fetching liked posts", http.StatusInternalServerError)
		return
	}

	// data to be passed to the template
	data := struct {
		Username       string
		Email          string
		ProfilePicture string
		CreatedPosts   []database.Posts
		LikedPosts     []database.Posts
	}{
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		CreatedPosts:   createdPosts,
		LikedPosts:     likedPosts,
	}

	tmpl, err := template.ParseFiles("templates/edit-profile.html") // parse the HTML template file
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data) // execute the template with the data and write it to the response
}

func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := models.GetUserIDFromRequest(r)
	if err != nil || userID == "" {
		models.SetNotification(w, "Unauthorized", "error")
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		models.SetNotification(w, "File too large (10MB)", "error")
		return
	}

	file, handler, err := r.FormFile("profile_picture")
	if err != nil {
		models.SetNotification(w, "Invalid file", "error")

		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(handler.Filename))
	path := filepath.Join("static/uploads", filename)

	dst, err := os.Create(path)
	if err != nil {
		models.SetNotification(w, "Unable to save file", "error")
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		models.SetNotification(w, "Error saving file", "error")
		return
	}

	err = database.UpdateProfilePicture(userID, filename)
	if err != nil {
		models.SetNotification(w, "Error updating user", "error")
		return
	}

	models.SetNotification(w, "Profile picture updated successfully", "info")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func DeleteProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := models.GetUserIDFromRequest(r)
	if err != nil || userID == "" {
		models.SetNotification(w, "Unauthorized", "error")
		return
	}

	var filename string
	err = database.SQL.QueryRow("SELECT profilepicture FROM Users WHERE id = ?", userID).Scan(&filename)
	if err != nil {
		models.SetNotification(w, "Error retrieving image", "error")
		return
	}

	// Supprimer le fichier s’il existe
	if filename != "" {
		path := filepath.Join("static/uploads", filename)
		log.Println("Attempting to delete file:", path)

		if err := os.Remove(path); err != nil {
			if os.IsNotExist(err) {
				log.Println("File does not exist:", path)
			} else {
				log.Println("Failed to delete file:", err)
				models.SetNotification(w, "Error deleting file", "error")
				return
			}
		} else {
			log.Println("File deleted successfully:", path)
		}
	} else {
		log.Println("No profile picture to delete")
	}

	// Réinitialiser la valeur en DB
	err = database.UpdateProfilePicture(userID, "")
	if err != nil {
		models.SetNotification(w, "Error updating database", "error")
		return
	}

	models.SetNotification(w, "Profile picture deleted successfully", "info")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
