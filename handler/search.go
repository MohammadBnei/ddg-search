package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"ddh-search/config"
	"ddh-search/service"
)

// SearchHandler handles search requests
type SearchHandler struct {
	config  *config.Config
	service service.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(cfg *config.Config, svc service.SearchService) *SearchHandler {
	return &SearchHandler{
		config:  cfg,
		service: svc,
	}
}

// Handle processes search requests
func (h *SearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Skip authentication in local mode
	if !h.config.LocalMode {
		// Authenticate request
		user, pass, ok := r.BasicAuth()
		if !ok {
			writeError(w, fmt.Errorf("missing credentials"), http.StatusUnauthorized)
			return
		}

		if user != h.config.AuthUsername {
			writeError(w, fmt.Errorf("invalid username"), http.StatusUnauthorized)
			return
		}

		if pass != h.config.AuthPassword {
			writeError(w, fmt.Errorf("invalid password"), http.StatusUnauthorized)
			return
		}
	}

	// Get search query
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		writeError(w, fmt.Errorf("missing query"), http.StatusBadRequest)
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

	// Format response
	response := make([]struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Snippet string `json:"snippet"`
	}, len(results))

	for i, r := range results {
		response[i] = struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Snippet string `json:"snippet"`
		}{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Snippet,
		}
	}

	// Send response
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, fmt.Errorf("failed to encode response: %w", err), http.StatusInternalServerError)
		return
	}
}
