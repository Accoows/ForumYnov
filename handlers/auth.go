package handlers

import (
	"fmt"
	"net/http"
)

// Login, register, logout

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Page de connexion")
		return
	}
	if r.Method == http.MethodPost {
		remember := r.FormValue("remember")
		if remember == "on" {
			fmt.Fprintln(w, "Connexion avec les cookies")
		} else {
			fmt.Fprintln(w, "Connexion normale (session courte)")
		}
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Page d'inscription")
		return
	}
	if r.Method == http.MethodPost {
		fmt.Fprintln(w, "Inscription effectué")
		return
	}
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Déconnexion de l'utilisateur")
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Page de mot de passe oublié")
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		fmt.Fprintf(w, "Mail de récupération envoyé à %s (simulation)", email)
		return
	}
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func ForgotUsernameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Page de nom d'utilisateur oublié")
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		fmt.Fprintf(w, "Nom d'utilisateur envoyé à %s (simulation)", email)
		return
	}
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}
