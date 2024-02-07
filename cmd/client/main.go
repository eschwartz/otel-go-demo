package main

import (
	"encoding/json"
	"fmt"
	"github.com/eschwartz/otel-go-demo/internal/pkg/data"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var dataService = &data.MemoryDataService{}

func init() {
	// We can use logrus to support structured logs: https://github.com/sirupsen/logrus
	// Configure logrus to output logs as JSON
	log.SetFormatter(&log.JSONFormatter{})
}
func HandleGetItems(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"url":    "/items",
		"method": "GET",
	}).Info("API Request")

	// Parse query params
	searchTerm := r.URL.Query().Get("")
	limit := r.URL.Query().Get("limit")
	log.WithFields(log.Fields{
		"searchTerm": searchTerm,
		"limit":      limit,
	}).Info("Query params")

	// Execute DB query
	limitInt, _ := strconv.Atoi(limit)
	items, err := dataService.FindItems(searchTerm, limitInt)
	if err != nil {
		log.WithFields(log.Fields{
			"message": fmt.Sprintf("Database query failed! %s", err),
		}).Error("error")
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}
	log.WithFields(log.Fields{
		"resultCount": len(items),
	}).Info("DB query")

	// Write JSON response
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		log.WithFields(log.Fields{
			"message": fmt.Sprintf("Failed to encode json: %s", err),
		}).Error("error")
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	w.WriteHeader(200)
	log.WithFields(log.Fields{
		"method":     "GET",
		"url":        "/items",
		"statusCode": 200,
	}).Info("API Response")
}
