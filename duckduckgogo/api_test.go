package duckduckgogo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockSearchClient implements the SearchClient interface for testing.
type MockSearchClient struct {
	results []Result
	err     error
}

func (m *MockSearchClient) Search(ctx context.Context, query string) ([]Result, error) {
	return m.results, m.err
}

func (m *MockSearchClient) SearchLimited(ctx context.Context, query string, limit int) ([]Result, error) {
	if limit <= 0 || limit > len(m.results) {
		return m.results, m.err
	}
	return m.results[:limit], m.err
}

func TestDuckDuckGoSearchClient_SearchLimited(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		if !strings.Contains(r.URL.String(), "?q=test") {
			t.Errorf("Expected query parameter 'q=test', got %s", r.URL.String())
		}

		// Check user agent is set
		if r.Header.Get("User-Agent") == "" {
			t.Error("User-Agent header not set")
		}

		// Return a simple HTML response with search results
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<div class="results">
				<div class="web-result">
					<a class="result__a">Test Title</a>
					<div class="result__snippet">Test Snippet</div>
					<a class="result__url">https://example.com</a>
					<img class="result__icon__img" src="icon.png" width="16" height="16" />
				</div>
			</div>
		`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &DuckDuckGoSearchClient{
		baseUrl: server.URL + "/",
	}

	// Test search
	results, err := client.SearchLimited(t.Context(), "test", 1)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Verify results
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// Check result fields
	result := results[0]
	if result.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", result.Title)
	}
	if result.Snippet != "Test Snippet" {
		t.Errorf("Expected snippet 'Test Snippet', got '%s'", result.Snippet)
	}
	if result.FormattedURL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got '%s'", result.FormattedURL)
	}
}

func TestCleanFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  test  ", "test"},
		{"test\ntest", "testtest"},
		{"\n test \n", "test"},
		{"", ""},
	}

	for _, test := range tests {
		result := clean(test.input)
		if result != test.expected {
			t.Errorf("clean(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestToIntFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"-10", -10},
		{"abc", 0}, // Invalid number should return 0
		{"", 0},    // Empty string should return 0
	}

	for _, test := range tests {
		result := toInt(test.input)
		if result != test.expected {
			t.Errorf("toInt(%q) = %d, expected %d", test.input, result, test.expected)
		}
	}
}
