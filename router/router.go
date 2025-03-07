package router

import (
	"net/http"

	"ddg-search/config"
	"ddg-search/handler"
	"ddg-search/service"
)

// Router handles HTTP routing.
type Router struct {
	mux *http.ServeMux
}

// New creates a new router.
func New(cfg *config.Config) *Router {
	mux := http.NewServeMux()

	// Create service
	searchService := service.NewDuckDuckGoService()

	// Create handlers
	searchHandler := handler.NewSearchHandler(cfg, searchService)
	healthHandler := handler.NewHealthHandler()

	// Register routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte("Hello from Go server!"))
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/health", healthHandler.Handle)
	mux.HandleFunc("/search", searchHandler.Handle)

	return &Router{
		mux: mux,
	}
}

// Handler returns the HTTP handler for the router.
func (r *Router) Handler() http.Handler {
	return r.mux
}
