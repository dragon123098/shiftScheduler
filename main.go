package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Defines a user struct to hold the user's name, password and schedule
// type User struct {
// 	Name     string
// 	Password string
// 	Schedule Schedule
// }

// Home handler to render the home page
func home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./templates/home.html")
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

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/save", save)

	fmt.Println("Server is running on http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
