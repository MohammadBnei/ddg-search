package service

import (
	"context"
	"ddg-search/duckduckgogo"
)

// SearchService defines the interface for search operations.
type SearchService interface {
	Search(query string, limit int) ([]SearchResult, error)
}

// SearchResult represents a search result.
type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// DuckDuckGoService implements SearchService using DuckDuckGo.
type DuckDuckGoService struct {
	client duckduckgogo.SearchClient
}

// NewDuckDuckGoService creates a new DuckDuckGo search service.
func NewDuckDuckGoService() *DuckDuckGoService {
	return &DuckDuckGoService{
		client: duckduckgogo.NewDuckDuckGoSearchClient(),
	}
}

// Search performs a search with the given query and limit.
func (s *DuckDuckGoService) Search(query string, limit int) ([]SearchResult, error) {
	results, err := s.client.SearchLimited(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}

	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		searchResults[i] = SearchResult{
			Title:   r.Title,
			URL:     r.FormattedURL,
			Snippet: r.Snippet,
		}
	}

	return searchResults, nil
}
