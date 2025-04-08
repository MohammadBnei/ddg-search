package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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
			assert.NoError(t, err, "Failed to create request")

			// Create a recorder to capture the response
			recorder := httptest.NewRecorder()

			// Serve the request
			handler.Handle(recorder, req)

			// Check the response status code
			assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code %d, got %d", http.StatusOK, recorder.Code)

			// Decode the response body
			var response []SearchResultResponse
			err = json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err, "Failed to decode response body")

			// Check the number of results
			assert.Equal(t, len(tc.mockResults), len(response), "Expected %d results, got %d", len(tc.mockResults), len(response))

			// Check the content of each result
			for i, result := range response {
				assert.Contains(t, result.Content, tc.expectedContent[i], "Expected content to contain %q, got %q", tc.expectedContent[i], result.Content)
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
