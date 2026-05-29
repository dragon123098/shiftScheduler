package main

import (
	"net/http"
	"log"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World ALIENS TEST"))
}


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}