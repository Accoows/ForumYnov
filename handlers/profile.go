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
		log.Println("[handlers/profile.go] Invalid session ID in cookie:", cookie.Value)
		models.SetNotification(w, "Invalid session", "error")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var user database.Users
	var profilePicture sql.NullString
	// get the user data from the database using the user ID
	err = database.SQL.QueryRow("SELECT username, email, profilepicture FROM Users WHERE id = ?", userID).Scan(&user.Username, &user.Email, &profilePicture)
	if err != nil {
		log.Println("[handlers/profile.go] User not found for session ID:", userID)
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
		log.Println("[handlers/profile.go] Error fetching created posts:", err)
		http.Error(w, "Error fetching created posts", http.StatusInternalServerError)
		return
	}

	// Get posts liked by the user
	likedPosts, err := database.GetLikedPostsByUser(userID)
	if err != nil {
		log.Println("[handlers/profile.go] Error fetching liked posts:", err)
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
		log.Println("[handlers/profile.go] Error loading template:", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("[handlers/profile.go] Error executing template:", err)
	}
	//tmpl.Execute(w, data) // execute the template with the data and write it to the response
}

// UpdateProfileHandler handles the profile picture update for users.
func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := models.GetUserIDFromRequest(r) // get the user ID from the request using the cookie
	if err != nil || userID == "" {
		models.SetNotification(w, "Unauthorized", "error")
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB // set the maximum memory for the form data
	if err != nil {
		models.SetNotification(w, "File too large (10MB)", "error")
		return
	}

	file, handler, err := r.FormFile("profile_picture") // get the file from the form data
	if err != nil {
		models.SetNotification(w, "Invalid file", "error")

		return
	}
	defer file.Close()

	// create a unique filename using the current timestamp and the original filename
	// and save it in the static/uploads directory
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(handler.Filename))
	path := filepath.Join("static/uploads", filename)

	dst, err := os.Create(path) // create the file in the static/uploads directory
	if err != nil {
		models.SetNotification(w, "Unable to save file", "error")
		return
	}
	defer dst.Close() // close the file after writing

	_, err = io.Copy(dst, file) // copy the file content to the new file
	if err != nil {
		models.SetNotification(w, "Error saving file", "error")
		return
	}

	err = database.UpdateProfilePicture(userID, filename) // update the profile picture in the database
	if err != nil {
		models.SetNotification(w, "Error updating user", "error")
		return
	}

	models.SetNotification(w, "Profile picture updated successfully", "info")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Deletes a post only if the connected user is its author
func DeleteProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := models.GetUserIDFromRequest(r)
	if err != nil || userID == "" {
		models.SetNotification(w, "Unauthorized", "error")
		return
	}

	var filename string
	// Get the current profile picture filename from the database
	err = database.SQL.QueryRow("SELECT profilepicture FROM Users WHERE id = ?", userID).Scan(&filename)
	if err != nil {
		models.SetNotification(w, "Error retrieving image", "error")
		return
	}

	// Delete the file if it exists
	if filename != "" {
		path := filepath.Join("static/uploads", filename) // construct the file path
		log.Println("Attempting to delete file:", path)

		if err := os.Remove(path); err != nil { // delete the file
			if os.IsNotExist(err) { // check if the file does not exist
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

	// Reset the value in the DB
	err = database.UpdateProfilePicture(userID, "") // update the profile picture in the database to an empty string
	if err != nil {
		models.SetNotification(w, "Error updating database", "error")
		return
	}

	// Notify user of success and redirect to the profile page
	models.SetNotification(w, "Profile picture deleted successfully", "info")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func UpdateProfileInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	userID, err := models.GetUserIDFromRequest(r)
	if err != nil || userID == "" {
		models.SetNotification(w, "Unauthorized", "error")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	newUsername := r.FormValue("username")
	newEmail := r.FormValue("email")

	if newUsername == "" || newEmail == "" {
		models.SetNotification(w, "Fields cannot be empty", "error")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	_, err = database.SQL.Exec("UPDATE Users SET username = ?, email = ? WHERE id = ?", newUsername, newEmail, userID)
	if err != nil {
		log.Println("[UpdateProfileInfoHandler] Error updating user info:", err)
		models.SetNotification(w, "Error updating info", "error")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	models.SetNotification(w, "Profile updated successfully", "success")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
