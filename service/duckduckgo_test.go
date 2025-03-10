package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"ddg-search/duckduckgogo"
)

// MockDDGClient implements the SearchClient interface for testing.
type MockDDGClient struct {
	results []duckduckgogo.Result
	err     error
	retryConfigCalled bool
	maxRetries int
	retryBackoff int
}

func (m *MockDDGClient) Search(ctx context.Context, query string) ([]duckduckgogo.Result, error) {
	return m.results, m.err
}

func (m *MockDDGClient) SearchLimited(ctx context.Context, query string, limit int) ([]duckduckgogo.Result, error) {
	if limit <= 0 || limit > len(m.results) {
		return m.results, m.err
	}
	return m.results[:limit], m.err
}

// Mock implementation of WithRetryConfig for testing
func (m *MockDDGClient) WithRetryConfig(maxRetries, retryBackoff int) duckduckgogo.SearchClient {
	m.retryConfigCalled = true
	m.maxRetries = maxRetries
	m.retryBackoff = retryBackoff
	return m
}

func TestDuckDuckGoService_Search(t *testing.T) {
	// Test data
	mockResults := []duckduckgogo.Result{
		{
			Title:        "Test Result 1",
			FormattedURL: "https://example.com/1",
			Snippet:      "This is test result 1",
		},
		{
			Title:        "Test Result 2",
			FormattedURL: "https://example.com/2",
			Snippet:      "This is test result 2",
		},
	}

	expectedResults := []SearchResult{
		{
			Title:   "Test Result 1",
			URL:     "https://example.com/1",
			Snippet: "This is test result 1",
		},
		{
			Title:   "Test Result 2",
			URL:     "https://example.com/2",
			Snippet: "This is test result 2",
		},
	}

	tests := []struct {
		name     string
		query    string
		limit    int
		mockData []duckduckgogo.Result
		mockErr  error
		want     []SearchResult
		wantErr  bool
	}{
		{
			name:     "successful search with all results",
			query:    "test query",
			limit:    0,
			mockData: mockResults,
			mockErr:  nil,
			want:     expectedResults,
			wantErr:  false,
		},
		{
			name:     "successful search with limited results",
			query:    "test query",
			limit:    1,
			mockData: mockResults,
			mockErr:  nil,
			want:     expectedResults[:1],
			wantErr:  false,
		},
		{
			name:     "search with error",
			query:    "test query",
			limit:    0,
			mockData: nil,
			mockErr:  errors.New("search failed"),
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockDDGClient{
				results: tt.mockData,
				err:     tt.mockErr,
			}

			// Create service with mock client
			service := &DuckDuckGoService{
				client: mockClient,
			}

			// Call the method being tested
			got, err := service.Search(tt.query, tt.limit)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("DuckDuckGoService.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check results
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DuckDuckGoService.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestWithRetryConfig(t *testing.T) {
	// Create a mock client
	mockClient := &MockDDGClient{}
	
	// Create service with mock client
	service := &DuckDuckGoService{client: mockClient}
	
	// Apply retry configuration
	service.WithRetryConfig(5, 100)
	
	// Verify that retry config was called on the mock
	if !mockClient.retryConfigCalled {
		t.Error("WithRetryConfig was not called on the client")
	}
	
	// Verify the values were passed correctly
	if mockClient.maxRetries != 5 {
		t.Errorf("Expected maxRetries to be 5, got %d", mockClient.maxRetries)
	}
	
	if mockClient.retryBackoff != 100 {
		t.Errorf("Expected retryBackoff to be 100, got %d", mockClient.retryBackoff)
	}
}
