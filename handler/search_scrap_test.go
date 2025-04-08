package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ddg-search/config"
	"ddg-search/service"
)

// MockSearchServiceWithScraping is a mock implementation of SearchService that supports scraping.
type MockSearchServiceWithScraping struct {
	results []service.SearchResult
	err     error
}

func (m *MockSearchServiceWithScraping) Search(query string, limit int) ([]service.SearchResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func TestSearchHandler_Scraping(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name            string
		scrapParam      string
		mockResults     []service.SearchResult
		mockError       error
		expectedContent []string
	}{
		{
			name:       "Scraping enabled, scrap query param present",
			scrapParam: "true",
			mockResults: []service.SearchResult{
				{Title: "Result 1", URL: "https://example.com/1", Snippet: "Snippet 1"},
				{Title: "Result 2", URL: "https://example.com/2", Snippet: "Snippet 2"},
			},
			expectedContent: []string{"Example Domain", "Example Domain"}, // Expecting markdown content
		},
		{
			name:       "Scraping enabled, scrap query param missing",
			scrapParam: "",
			mockResults: []service.SearchResult{
				{Title: "Result 1", URL: "https://example.com/1", Snippet: "Snippet 1"},
				{Title: "Result 2", URL: "https://example.com/2", Snippet: "Snippet 2"},
			},
			expectedContent: []string{"", ""}, // Expecting empty content
		},
		{
			name:       "Scraping fails for a URL",
			scrapParam: "true",
			mockResults: []service.SearchResult{
				{Title: "Result 1", URL: "https://example.com/1", Snippet: "Snippet 1"},
				{Title: "Result 2", URL: "invalid-url", Snippet: "Snippet 2"}, // Invalid URL
			},
			expectedContent: []string{"Example Domain", ""}, // Expecting empty content for invalid URL
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock config
			cfg := &config.Config{
				LocalMode: true,
			}

			// Create a mock search service
			mockSvc := &MockSearchServiceWithScraping{
				results: tc.mockResults,
				err:     tc.mockError,
			}

			// Create a search handler with the mock config and service
			handler := NewSearchHandler(cfg, mockSvc)

			// Create a request with the scrap query parameter
			req, err := http.NewRequest("GET", "/search?q=test&scrap="+tc.scrapParam, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a recorder to capture the response
			recorder := httptest.NewRecorder()

			// Serve the request
			handler.Handle(recorder, req)

			// Check the response status code
			if recorder.Code != http.StatusOK {
				t.Fatalf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
			}

			// Decode the response body
			var response []SearchResultResponse
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			// Check the number of results
			if len(response) != len(tc.mockResults) {
				t.Fatalf("Expected %d results, got %d", len(tc.mockResults), len(response))
			}

			// Check the content of each result
			for i, result := range response {
				if !strings.Contains(result.Content, tc.expectedContent[i]) {
					t.Errorf("Expected content to contain %q, got %q", tc.expectedContent[i], result.Content)
				}
			}
		})
	}
}

// Mock HTTP server for testing scraping
func init() {
	http.HandleFunc("https://example.com/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "<html><body><h1>Example Domain</h1></body></html>")
	})
	http.HandleFunc("https://example.com/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "<html><body><h1>Example Domain</h1></body></html>")
	})
	go http.ListenAndServe(":8081", nil)
}
