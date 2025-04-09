package handlers

import (
	"forumynov/database"
	"forumynov/models"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

/*
RegisterUsers handles the registration process for new users. It verifies the provided data, hashes the password, and stores the user in the database.
It also checks for the uniqueness of the email and username before inserting the new user into the database.
*/
func RegisterUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error during form processing", http.StatusBadRequest)
		return
	}

	newUser := database.Users{ // create a new user struct to hold the registration data without modifying the database.
		Email:      r.FormValue("email"),
		Username:   r.FormValue("username"),
		Created_at: time.Now().Format("2006-01-02 15:04:05"), // format the current time to a string
	}
	password := r.FormValue("password")

	if newUser.Email == "" || newUser.Username == "" || password == "" { // check if the email, username, and password fields are empty
		http.Error(w, "All fields are mandatory", http.StatusBadRequest)
		return
	}

	emailRight := VerifyEmailConformity(&newUser) // check if the email is valid
	if emailRight == false {
		http.Error(w, "Email does not comply", http.StatusBadRequest)
		return
	}

	exists, err := VerifyEmailAndUsernameUnicity(newUser.Email, newUser.Username) // check if the email or username already exists in the database
	if err != nil {
		http.Error(w, "Error during user verification", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email or username already used", http.StatusConflict)
		return
	}

	password_hash, err := models.HashPassword(password) // hash the password using the HashPassword function from the models package
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password_hash = password_hash // set the hashed password in the user struct

	newID := uuid.NewV4()       // generate a new UUID for the user
	newUser.ID = newID.String() // set the user ID to the generated UUID

	err = database.InsertUsersData(&newUser) // insert the new user into the database using the InsertUsersData function from the database package
	if err != nil {
		http.Error(w, "Error during user registration", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound) // redirect the user to the home page after registration
}

// VerifyEmailAndUsernameUnicity checks if the provided email or username already exists in the database.
func VerifyEmailAndUsernameUnicity(email string, username string) (bool, error) {
	var exists bool
	// SQL query to check if the email or username already exists in the Users table
	query := "SELECT EXISTS (SELECT 1 FROM Users WHERE email = ? OR username = ?)"
	err := database.SQL.QueryRow(query, email, username).Scan(&exists) // Execute the query and scan the result into the exists variable
	if err != nil {
		log.Println("Erreur SQL dans VerifyEmailAndUsernameUnicity:", err)
		return false, err
	}
	return exists, nil // return the result of the query
}

func VerifyEmailConformity(users *database.Users) bool {
	email := []rune(users.Email)       // convert the email string to a slice of runes for better handling of Unicode characters
	emailLenght := len(email)          // get the length of the email string
	emailRight := false                // initialize a boolean variable to check if the email is valid
	for i := 0; i < emailLenght; i++ { // iterate through each character of the email string
		if email[i] == '@' { // check if the character is '@'
			for j := i; j < emailLenght; j++ { // iterate through the rest of the email string
				if email[j] == '.' { // check if the character is '.'
					emailRight = true // set the emailRight variable to true if both '@' and '.' are found in the correct order
				}
			}
		}
	}
	return emailRight // return the result of the email validation
}
