package config

import (
	"errors"
	"log/slog"
	"os"
)

// Config holds all application configuration.
type Config struct {
	Port         string
	AuthUsername string
	AuthPassword string
	DebugMode    bool
	LocalMode    bool // When true, bypasses authentication for local testing
}

// New creates a new Config with values from environment variables.
// Returns an error if required authentication credentials are missing.
func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Debug mode defaults to false
	debugMode := os.Getenv("DEBUG_MODE") == "true"

	// Local mode defaults to false
	localMode := os.Getenv("LOCAL_MODE") == "true"

	if localMode {
		slog.Warn("Running in LOCAL_MODE - authentication is disabled")
	}

	// Only check auth credentials if not in local mode
	username := os.Getenv("AUTH_USERNAME")
	password := os.Getenv("AUTH_PASSWORD")

	if !localMode {
		if username == "" {
			return nil, errors.New("AUTH_USERNAME environment variable not set")
		}

		if password == "" {
			return nil, errors.New("AUTH_PASSWORD environment variable not set")
		}
	}

	return &Config{
		Port:         port,
		AuthUsername: username,
		AuthPassword: password,
		DebugMode:    debugMode,
		LocalMode:    localMode,
	}, nil
}
