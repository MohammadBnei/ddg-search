package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
func (h *SearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Skip authentication in local mode
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
	if err = json.NewEncoder(w).Encode(response); err != nil {
		writeError(w, fmt.Errorf("failed to encode response: %w", err), http.StatusInternalServerError)
		return
	}
}

func writeError(w http.ResponseWriter, err error, code int) {
	slog.Error("Error handling request", "error", err, "status_code", code)
	w.WriteHeader(code)
	if encodeErr := json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}); encodeErr != nil {
		slog.Error("Failed to encode error response", "error", encodeErr)
	}
}
