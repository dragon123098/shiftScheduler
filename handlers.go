package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Home handler to render the home page
func student(w http.ResponseWriter, r *http.Request) {
	username, ok := verifySession(w, r, 1)
	if !ok {
		log.Println("Unauthorized access attempt to student page")
		return
	}
	//Check if the user has a schedule, if not create an empty one for them
	if _, exists := shifts[username]; !exists {
		shifts[username] = Schedule{}
	}

	s, err := template.ParseFiles("./templates/student.html", "./templates/weekly_schedule.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)

		return
	}

	schedule := shifts[username]
	err = s.Execute(w, WeeklyScheduleComponentData{Schedule: &schedule, CanWrite: true})
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

	var schedule Schedule
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	username, ok := verifySession(w, r, 1)
	if !ok {
		log.Println("Unauthorized access attempt to save schedule")
		return
	}
	shifts[username] = schedule
	schedule.PrintSchedule()

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
		Name:    "session_token",
		Value:   sessionToken,
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
