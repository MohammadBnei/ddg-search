package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"ddh-search/duckduckgogo"
)

var password string

const User = "telegram-wait-qa"

func init() {
	password = os.Getenv("auth-password")

	log.Println("ddg-search handler initialized")
}

func Handle(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		writeError(w, fmt.Errorf("missing credentials"), http.StatusUnauthorized)
		return
	}

	if user != User {
		writeError(w, fmt.Errorf("invalid username"), http.StatusUnauthorized)
		return
	}

	if pass != password {
		writeError(w, fmt.Errorf("invalid password"), http.StatusUnauthorized)
		return
	}

	ddg := duckduckgogo.NewDuckDuckGoSearchClient()

	searchQuery := r.URL.Query().Get("q")

	if searchQuery == "" {
		writeError(w, fmt.Errorf("missing query"), http.StatusBadRequest)
		return
	}

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

	res, err := ddg.SearchLimited(searchQuery, limit)
	if err != nil {
		writeError(w, fmt.Errorf("failed to search: %w", err), http.StatusInternalServerError)
		return
	}

	response := make([]struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Snippet string `json:"snippet"`
	}, len(res))

	for i, r := range res {
		response[i] = struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Snippet string `json:"snippet"`
		}{
			Title:   r.Title,
			URL:     r.FormattedUrl,
			Snippet: r.Snippet,
		}
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, fmt.Errorf("failed to encode response: %w", err), http.StatusInternalServerError)
		return
	}
}

func writeError(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
