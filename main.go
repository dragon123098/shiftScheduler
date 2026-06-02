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

//Defines a user struct to hold the user's name, password and schedule
type User struct {
	Password string
	Role    string
	//Schedule Schedule
}

type session struct {
	username string
	expiry time.Time
	role string
}

var sessions = map[string]session{} // Map to store session tokens and associated usernames



func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}


//Map that maps the username to the user struct, this will be used for authentication and to store the schedule for each user
//Student 1 has password "pass123" and role "student"
//Admin 1 has password "adminpass" and role "admin"
var users = map[string]User{
	"student1": {Password: "pass123", Role: "student"},
	"admin1":   {Password: "adminpass", Role: "admin"},
}

// Home handler to render the home page
func student(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists || userSession.isExpired() {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	
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





func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", loginHandler)
	mux.HandleFunc("/student", student)
	mux.HandleFunc("/save", save)
	//mux.HandleFunc("/admin", adminHandler)
	mux.HandleFunc("POST /login", loginHandlerPost)

	//I need to start a cookie session, and have each handler check if the user is authenticated
	//If not, use an if statment to simply redirect to the login page.

	fmt.Println("Server is running on http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
