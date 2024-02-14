package data

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"time"
)

var tracer = otel.Tracer("example.app")

type MemoryDataService struct {
}

type Item struct {
	Value string `json:"value"`
}

// This is here to repesent a backend data system, like SQL database or data API.
// The data and behavior is stubbed out, for simplicity
func (svc *MemoryDataService) FindItems(term string, limit int, ctx context.Context) ([]Item, error) {
	_, span := tracer.Start(ctx, "DataService.FindItems")
	defer span.End()

	// Typically unhelpful error messages
	if term == "" {
		return []Item{}, errors.New("unexpected empty value")
	}
	if limit == 0 {
		return []Item{}, errors.New("unexpected empty value")
	}

	// Simulate some slow db operations
	if limit > 70 {
		time.Sleep(time.Millisecond * 2000)
	} else if limit > 50 {
		time.Sleep(time.Millisecond * 750)
	} else if limit > 30 {
		time.Sleep(time.Millisecond * 100)
	} else {
		time.Sleep(time.Millisecond * 5)
	}

	return []Item{
		{"A"},
		{"B"},
		{"C"},
	}, nil
}
