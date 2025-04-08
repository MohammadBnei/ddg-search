package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"

	"ddg-search/config"
	"ddg-search/service"
)

// SearchHandler handles search requests.
type SearchHandler struct {
	config  *config.Config
	service service.SearchService
}

// NewSearchHandler creates a new search handler.
func NewSearchHandler(cfg *config.Config, svc service.SearchService) *SearchHandler {
	return &SearchHandler{
		config:  cfg,
		service: svc,
	}
}

// Handle processes search requests.
//
//	@Summary		Search DuckDuckGo
//	@Description	Search DuckDuckGo with optional limit
//	@Tags			search
//	@Security		BasicAuth
//	@Param			q		query	string	true	"Search query"
//	@Param			limit	query	int		false	"Maximum number of results to return"
//	@Produce		json
//	@Success		200	{array}		SearchResultResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/search [get]
func (h *SearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Skip authentication in local mode
	enableScrap := false
	if h.config.EnableScraping {
		enableScrap = r.URL.Query().Get("scrap") == "true"
	}

	if !h.config.LocalMode {
		// Authenticate request
		user, pass, ok := r.BasicAuth()
		if !ok {
			writeError(w, errors.New("missing credentials"), http.StatusUnauthorized)
			return
		}

		if user != h.config.AuthUsername {
			writeError(w, errors.New("invalid username"), http.StatusUnauthorized)
			return
		}

		if pass != h.config.AuthPassword {
			writeError(w, errors.New("invalid password"), http.StatusUnauthorized)
			return
		}
	}

	// Get search query
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		writeError(w, errors.New("missing query"), http.StatusBadRequest)
		return
	}

	// Get limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			writeError(w, fmt.Errorf("invalid limit %s: %w", limitStr, err), http.StatusBadRequest)
			return
		}
		limit = l
	}

	// Perform search
	results, err := h.service.Search(searchQuery, limit)
	if err != nil {
		writeError(w, fmt.Errorf("failed to search: %w", err), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	if enableScrap {
		wg.Add(len(results))
		for i := range results {
			go func(i int) {
				defer wg.Done()
				markdown, err := h.scrapURL(results[i].URL)
				if err != nil {
					slog.Error("Error scraping URL", "url", results[i].URL, "error", err)
					return // Don't return, continue with other results
				}

				mu.Lock()
				results[i].Snippet = markdown
				mu.Unlock()
			}(i)
		}
		wg.Wait()
	}

	// Format response
	response := make([]SearchResultResponse, len(results))

	for i, r := range results {
		response[i] = SearchResultResponse{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Snippet,
		}
	}

	// Send response
	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		writeError(w, fmt.Errorf("failed to encode response: %w", err), http.StatusInternalServerError)
		return
	}
}

func (h *SearchHandler) scrapURL(URL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Fetch HTML content from the URL
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL, status code: %d", resp.StatusCode)
	}

	// Convert HTML to Markdown
	converter := htmltomarkdown.NewConverter("", true, nil)
	markdown, err := converter.ConvertReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to Markdown: %w", err)
	}

	return markdown, nil
}

// SearchResultResponse is the response format for search results.
type SearchResultResponse struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// ErrorResponse is the response format for errors.
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error, code int) {
	slog.Error("Error handling request", "error", err, "status_code", code)
	w.WriteHeader(code)
	if encodeErr := json.NewEncoder(w).Encode(ErrorResponse{
		Error: err.Error(),
	}); encodeErr != nil {
		slog.Error("Failed to encode error response", "error", encodeErr)
	}
}
