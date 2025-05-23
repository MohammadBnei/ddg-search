package router

import (
	"log"
	"net/http"
	"net/http/pprof" // Import pprof package

	httpSwagger "github.com/swaggo/http-swagger" // Import http-swagger

	"ddg-search/config"
	"ddg-search/handler"
	"ddg-search/middleware" // Import the middleware package
	"ddg-search/service"
)

// Router handles HTTP routing.
type Router struct {
	Mux *http.ServeMux
	cfg *config.Config
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
	mux.HandleFunc("/search", searchHandler.Handle)
	mux.HandleFunc("/health", healthHandler.Handle)

	// Serve Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Enable pprof endpoints if debug mode is enabled
	if cfg.DebugMode {
		log.Println("Debug mode enabled: Registering pprof endpoints")
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		// Alternatively, for more fine-grained control:
		// mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	}

	return &Router{
		Mux: mux,
		cfg: cfg,
	}
}

// Handler returns the HTTP handler for the router.
func (r *Router) Handler() http.Handler {
	// Wrap the ServeMux with the logging middleware
	return middleware.LoggingMiddleware(r.cfg)(r.Mux)
}
