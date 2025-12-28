package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up!")
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", healthcheckHandler).Methods("GET")

	fmt.Println("Server starting on porn :3838")
	log.Fatal(http.ListenAndServe(":3838", router))
}
