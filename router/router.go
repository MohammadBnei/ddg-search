package router

import (
	"net/http"

	"ddg-search/config"
	"ddg-search/handler"
	"ddg-search/service"

	httpSwagger "github.com/swaggo/http-swagger" // Import http-swagger
)

// Router handles HTTP routing.
type Router struct {
	Mux *http.ServeMux
}

// New creates a new router.
func New(cfg *config.Config) *Router {
	mux := http.NewServeMux()

	// Create service with retry configuration
	searchService := service.NewDuckDuckGoService().WithRetryConfig(cfg.MaxRetries, cfg.RetryBackoff)

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

	// Serve Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return &Router{
		Mux: mux,
	}
}

// Handler returns the HTTP handler for the router.
func (r *Router) Handler() http.Handler {
	return r.Mux
}
