package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eschwartz/otel-go-demo/internal/pkg/data"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"net/http"
	"strconv"
)

var dataService = &data.MemoryDataService{}
var tracer = otel.Tracer("example.app")

func main() {
	// Setup HTTP server
	http.HandleFunc("/items", HandleGetItems)

	log.Println("Listening on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	log.Fatal(err)
}
func HandleGetItems(w http.ResponseWriter, r *http.Request) {
	// Start a new trace, creating a "parent span"
	// This span will describe the entire GET /items request
	_, span := tracer.Start(context.Background(), "GET /items")

	// Add attributes to the span (similar to structured log values)
	span.SetAttributes(
		attribute.String("url", "/items"),
		attribute.String("method", "GET"),
	)

	// Parse query params
	searchTerm := r.URL.Query().Get("")
	limit := r.URL.Query().Get("limit")

	// As we continue processing the request,
	//we'll keep adding attributes to the span
	limitInt, _ := strconv.Atoi(limit)
	span.SetAttributes(
		attribute.String("searchTerm", searchTerm),
		attribute.Int("limit", limitInt),
	)

	// Execute DB query
	items, err := dataService.FindItems(searchTerm, limitInt)
	if err != nil {
		// Errors are just another attribute of the span!
		span.SetAttributes(
			attribute.String("error", fmt.Sprintf("Database query failed! %s", err)),
			attribute.Int("response.status", 500),
		)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	span.SetAttributes(
		attribute.Int("resultCount", len(items)),
	)

	// Write JSON response
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		span.SetAttributes(
			attribute.String("error", fmt.Sprintf("Failed to encode json: %s", err)),
			attribute.Int("response.status", 500),
		)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	span.SetAttributes(
		attribute.Int("response.status", 200),
	)
}
