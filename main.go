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
	cfg := config.New()

	// Setup router
	r := router.New(cfg)

	// Configure server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r.Handler(),
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
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}
