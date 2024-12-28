package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "auth-session")

	// Generate OAuth URL
	url := Config.AuthCodeURL("state")

	// Store the intended URL
	session.Values["intended"] = r.URL.Path
	session.Save(r, w)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Callback received") // Debug

	code := r.URL.Query().Get("code")
	fmt.Printf("Auth code: %s\n", code) // Debug

	token, err := Config.Exchange(r.Context(), code)
	if err != nil {
		fmt.Printf("Token exchange error: %v\n", err) // Debug
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info
	client := Config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Printf("User info error: %v\n", err) // Debug
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		fmt.Printf("JSON parse error: %v\n", err) // Debug
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	fmt.Printf("User email: %s\n", userInfo.Email) // Debug

	// Check if email is allowed
	if !AllowedEmails[userInfo.Email] {
		fmt.Printf("Unauthorized email: %s\n", userInfo.Email) // Debug
		http.Error(w, "Unauthorized email", http.StatusUnauthorized)
		return
	}

	// Set session
	session, _ := Store.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["email"] = userInfo.Email
	if err := session.Save(r, w); err != nil {
		fmt.Printf("Session save error: %v\n", err) // Debug
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	fmt.Println("Auth successful, redirecting") // Debug
	http.Redirect(w, r, "http://localhost:1313/", http.StatusTemporaryRedirect)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "auth-session")
	session.Values["authenticated"] = false
	session.Values["email"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)

		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
