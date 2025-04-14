package handlers

import (
	"database/sql"
	"fmt"
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Page d'accueil, général

// Gestionnaire pour servir la page de login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(filepath.Join("./templates/", "login.html"))
		if err != nil {
			log.Println(err)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	if r.Method == http.MethodPost {
		// Traiter les données du formulaire
		LoginUsers(w, r)
		return
	}
}

// LoginUsers handles the login process for users. It verifies the user's credentials and creates a session if they are valid.
func LoginUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // check if the request method is POST
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed) // http.Error sends an HTTP response with the specified status code and message.
		return
	}

	err := r.ParseForm() // ParseForm() parses the form data from the request body and populates r.Form with the parsed values.
	if err != nil {
		http.Error(w, "Error during form processing", http.StatusBadRequest)
		return
	}

	newUser := database.Users{ // create a new user struct to hold the login data without modifying the database.
		Email:         r.FormValue("email"),    // get the email from the form data
		Password_hash: r.FormValue("password"), // get the password from the form data
	}

	if newUser.Email == "" || newUser.Password_hash == "" { // check if the email and password fields are empty
		http.Error(w, "All fields are mandatory", http.StatusBadRequest)
		return
	}

	emailRight := VerifyEmailConformity(&newUser) // check if the email is valid
	if !emailRight {
		http.Error(w, "Email does not comply", http.StatusBadRequest)
		return
	}
	userUUID, hashedPassword, err := GetHashedPasswordAndUUID(newUser.Email) // get the hashed password and user UUID from the database using the email
	if err != nil {
		http.Error(w, "Incorrect email", http.StatusNotFound)
		return
	}

	if !models.CheckPasswordHash(newUser.Password_hash, hashedPassword) { // check if the provided password matches the hashed password from the database
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	newID := uuid.NewV4() // generate a new UUID for the session

	newSessions := database.Sessions{ // create a new session struct to hold the session data
		ID:         newID.String(),                 // set the session ID to the generated UUID
		User_id:    userUUID,                       // set the user ID to the UUID of the logged-in user
		Expires_at: time.Now().Add(24 * time.Hour), // set the expiration time of the session to 24 hours from now
	}

	err = database.InsertSessionsData(&newSessions) // insert the session data into the database
	if err != nil {
		http.Error(w, "Error during Session database integration", http.StatusInternalServerError)
		return
	}

	userIdCookie := &http.Cookie{ // create a new cookie to store the session ID. The cookie is used to identify the user session on the server side.
		Name:       "user_cookie",  // set the name of the cookie
		Value:      newSessions.ID, // set the value of the cookie to the session ID
		Path:       "/",            // set the path for which the cookie is valid
		Domain:     "",
		Expires:    newSessions.Expires_at, // set the expiration time of the cookie to the session expiration time
		RawExpires: "",
		MaxAge:     86400, // set the maximum age of the cookie to 24 hours in seconds
		Secure:     false,
		HttpOnly:   true,                    // set to true if the cookie should not be accessible via JavaScript for security reasons
		SameSite:   http.SameSiteStrictMode, // set the SameSite attribute to prevent CSRF attacks
		Raw:        "",
		Unparsed:   []string{},
	}

	http.SetCookie(w, userIdCookie) // set the cookie in the response header to send it to the client

	http.Redirect(w, r, "/profile", http.StatusFound) // redirect the user to the profile page after successful login
}

// GetHashedPasswordAndUUID retrieves the hashed password and UUID of a user from the database using their email address.
func GetHashedPasswordAndUUID(email string) (string, string, error) {
	var userUUID string
	var hashedPassword string
	query := "SELECT id, password_hash FROM Users WHERE email = ?"
	err := database.SQL.QueryRow(query, email).Scan(&userUUID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user found")
		}
		return "", "", err
	}

	return userUUID, hashedPassword, nil // return the user UUID and hashed password
}
