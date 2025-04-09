package main

import (
	"fmt"
	"forumynov/database"
	"forumynov/handlers"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
	"time"
)

func openBrowser(url string) {
	var err error
	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

func main() {
	database.InitDatabase()        // Initialize the database connection
	defer database.CloseDatabase() // Close the database connection when the program exits

	go func() { // Periodically delete expired sessions
		for {
			database.DeleteExpiredSessions() // delete expired sessions from the database
			time.Sleep(5 * time.Minute)      // wait for 5 minutes before checking again
		}
	}()

	// Parse templates
	tmpl := template.Must(template.ParseFiles(filepath.Join("./templates/", "index.html")))
	tmplLogin := template.Must(template.ParseFiles(filepath.Join("./templates/", "login.html")))
	tmplregister := template.Must(template.ParseFiles(filepath.Join("./templates/", "register.html")))
	tmplpost := template.Must(template.ParseFiles(filepath.Join("./templates/", "post.html")))
	tmplpostlist := template.Must(template.ParseFiles(filepath.Join("./templates/", "post-list.html")))
	tmplpostdetail := template.Must(template.ParseFiles(filepath.Join("./templates/", "post-detail.html")))

	http.HandleFunc("/login", handlers.LoginUsers)
	http.HandleFunc("/profile", handlers.ProfilePage)
	http.HandleFunc("/logout", handlers.LogoutUsers)
	http.HandleFunc("/register", handlers.RegisterUsers)

	// Template handlers
	http.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		_ = tmplLogin.Execute(w, nil)
	})
	http.HandleFunc("/profile.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/profile", http.StatusFound)
	})
	http.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) {
		_ = tmplregister.Execute(w, nil)
	})
	http.HandleFunc("/post.html", func(w http.ResponseWriter, r *http.Request) {
		_ = tmplpost.Execute(w, nil)
	})
	http.HandleFunc("/post-list.html", func(w http.ResponseWriter, r *http.Request) {
		_ = tmplpostlist.Execute(w, nil)
	})
	http.HandleFunc("/post-detail.html", func(w http.ResponseWriter, r *http.Request) {
		_ = tmplpostdetail.Execute(w, nil)
	})

	// Static files
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/Templates/", http.StripPrefix("/Templates/", http.FileServer(http.Dir("./Templates"))))

	// Root handler - must be last to avoid conflicts
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an API request
		if r.URL.Path == "/" {
			_ = tmpl.Execute(w, nil)
			return
		}
		http.NotFound(w, r)
	})

	fmt.Println("Starting server at 127.0.0.1:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
