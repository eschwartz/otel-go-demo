package data

type MemoryDataService struct {
}

type Item struct {
	Value string `json:"value"`
}

func (svc *MemoryDataService) FindItems(term string, limit int) ([]Item, error) {
	return []Item{
		{"A"},
		{"B"},
		{"C"},
	}, nil
}
