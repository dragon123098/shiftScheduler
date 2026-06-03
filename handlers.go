package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

)

// Home handler to render the home page
func student(w http.ResponseWriter, r *http.Request) {
	_, ok := verifySession(w, r, 1)
	if !ok {
		log.Println("Unauthorized access attempt to student page")
		return
	}
	
	ts, err := template.ParseFiles("./templates/student.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)

		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)

		return
	}
}

// Save handler to process form input
func save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	inputValue := data["data"]

	// Write input value to schedule.txt
	file, err := os.OpenFile("schedule.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(inputValue + "\n")
	if err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		log.Println("Error writing to file:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}
	ts.Execute(w, nil)

}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifySession(w, r, 2)
	if !ok {
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", username)))
}


func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	
	log.Println("Login attempt received")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	
	// Read form values instead of JSON
    username := r.FormValue("username")
    password := r.FormValue("password")

	user, exists := users[username]
	if !exists || user.Password != password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//This addeds a session to our map. 
	//You can access the session by using the Value in the cookie, which is the session token.
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[sessionToken] = session{username: username, expiry: expiresAt, role: user.Role}

	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: sessionToken,
		Expires: expiresAt,
	})
	
	// If the user is authenticated, redirect to the appropriate page based on their role
	switch user.Role {
	case "student":
		http.Redirect(w, r, "/student", http.StatusSeeOther)
	case "admin":
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	default:
		http.Error(w, "Unknown role", http.StatusInternalServerError)
	}


}


func verifySession(w http.ResponseWriter, r *http.Request, role int) (string, bool) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("No session cookie found")
			return "", false
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return "", false
	}

	
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	
	var access int
	switch userSession.role {
	case "student":
		access = 1
	case "admin":
		access = 2
	default:
		access = 0
	}

	//log.Printf("Session token: %s, Username: %s, Role: %d, Access level: %d\n", sessionToken, userSession.username, role, access)

	if !exists || userSession.isExpired() {
		log.Println("Invalid or expired session token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return "", false
	}

	fmt.Printf("Role: %d type: %T\n", role, role)
	fmt.Printf("Access: %d type: %T\n", access, access)

	if access < role {
		log.Printf("Insufficient access level for user %s. Required: %d, User's access: %d\n", userSession.username, role, access)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return "", false
	}
	return userSession.username, true
}