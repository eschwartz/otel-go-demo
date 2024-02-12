package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eschwartz/otel-go-demo/internal/pkg/data"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
	"strconv"
)

var dataService = &data.MemoryDataService{}
var tracer = otel.Tracer("example.app")

/*
Configure OpenTelemetry to export traces to honeycomb
See https://docs.honeycomb.io/getting-data-in/opentelemetry/go-distro/#using-opentelemetry-without-the-honeycomb-distribution

export OTEL_SERVICE_NAME="your-service-name"
export OTEL_EXPORTER_OTLP_ENDPOINT="https://api.honeycomb.io:443" # US instance
export OTEL_EXPORTER_OTLP_HEADERS="x-honeycomb-team=your-api-key"
*/

func main() {

	err, cleanupOtel := setupOtel()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanupOtel()

	// Setup HTTP server
	http.HandleFunc("/items", HandleGetItems)

	log.Println("Listening on http://localhost:8000")
	err = http.ListenAndServe(":8000", nil)
	log.Fatal(err)
}

// Adapted from https://docs.honeycomb.io/getting-data-in/opentelemetry/go-distro/#configure
func setupOtel() (error, func()) {
	ctx := context.Background()

	// Configure a new OTLP exporter using environment variables for sending data to Honeycomb over gRPC
	otelClient := otlptracegrpc.NewClient()
	otelExporter, err := otlptrace.New(ctx, otelClient)
	if err != nil {
		return fmt.Errorf("failed to initialize exporter: %w", err), func() {}
	}

	// Create a new tracer provider with a batch span processor and the otlp exporter
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(otelExporter),
	)

	// Handle shutdown to ensure all sub processes are closed correctly and telemetry is exported
	cleanupFn := func() {
		_ = otelExporter.Shutdown(ctx)
		_ = tracerProvider.Shutdown(ctx)
	}

	// Register the global Tracer provider
	otel.SetTracerProvider(tracerProvider)

	return nil, cleanupFn
}
func HandleGetItems(w http.ResponseWriter, r *http.Request) {
	// Start a new trace, creating a "parent span"
	// This span will describe the entire GET /items request
	_, span := tracer.Start(context.Background(), "GET /items")
	defer span.End()

	// Add attributes to the span (similar to structured log values)
	span.SetAttributes(
		// We'll use the `app.` prefix to distinguish from other instrumented attributes
		attribute.String("app.url", "/items"),
		attribute.String("app.method", "GET"),
	)

	// Parse query params
	searchTerm := r.URL.Query().Get("q")
	limit := r.URL.Query().Get("limit")

	// As we continue processing the request,
	//we'll keep adding attributes to the span
	limitInt, _ := strconv.Atoi(limit)
	span.SetAttributes(
		attribute.String("app.searchTerm", searchTerm),
		attribute.Int("app.limit", limitInt),
	)

	// Execute DB query
	items, err := dataService.FindItems(searchTerm, limitInt)
	if err != nil {
		// Errors are just another attribute of the span!
		span.SetAttributes(
			attribute.String("app.error", fmt.Sprintf("Database query failed! %s", err)),
			attribute.Int("app.response.status", 500),
		)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	span.SetAttributes(
		attribute.Int("app.resultCount", len(items)),
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
