package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Basic logging example
func HandleGetItems(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to GET /items")

	// Parse query params
	searchTerm := r.URL.Query().Get("")
	limit := r.URL.Query()["limit"]
	log.Printf("search term: %s, limit: %s", searchTerm, limit)

	// Execute DB query
	log.Println("Querying database....")
	items, err := dataService.FindItems(searchTerm, limit)
	if err != nil {
		log.Println("Database query failed! %s", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}
	log.Println("DB query complete! Found %d items", len(items))

	// Write JSON response
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		log.Printf("Failed to encode json: %s", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	w.WriteHeader(200)
	log.Println("Request to GET /items successful!")
}
