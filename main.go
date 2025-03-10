package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ddg-search/config"
	"ddg-search/router"
)

func main() {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Setup router
	r := router.New(cfg)

	// Configure server
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r.Handler(),
		ReadHeaderTimeout: 10 * time.Second, // Fix for G112: Potential Slowloris Attack
	}

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		log.Printf("Server starting on port %s...\n", cfg.Port)
		if cfg.LocalMode {
			log.Println("LOCAL_MODE enabled: Authentication is bypassed")
		}
		log.Printf("Search retry configuration: max_retries=%d, retry_backoff=%dms\n", 
			cfg.MaxRetries, cfg.RetryBackoff)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	// Wait for shutdown signal
	<-stop

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		// Don't use Fatalf after defer to ensure defer runs
		log.Printf("Server shutdown failed: %v", err)
		defer os.Exit(1)
	}
	log.Println("Server stopped")
}
