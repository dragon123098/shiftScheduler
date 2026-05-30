package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
)

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


func main() {
	mux := http.NewServeMux()
	
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	mux.HandleFunc("/", home)

	fmt.Println("Server is running on http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
	

}