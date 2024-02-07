package data

import "errors"

type MemoryDataService struct {
}

type Item struct {
	Value string `json:"value"`
}

func (svc *MemoryDataService) FindItems(term string, limit int) ([]Item, error) {
	// Typically unhelpful error messages
	if term == "" {
		return []Item{}, errors.New("unexpected empty value")
	}
	if limit == 0 {
		return []Item{}, errors.New("unexpected empty value")
	}

	return []Item{
		{"A"},
		{"B"},
		{"C"},
	}, nil
}
