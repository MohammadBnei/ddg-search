package config

import (
	"log"
	"os"
)

// Config holds all application configuration
type Config struct {
	Port         string
	AuthUsername string
	AuthPassword string
	DebugMode    bool
	LocalMode    bool // When true, bypasses authentication for local testing
}

// New creates a new Config with values from environment variables
func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	password := os.Getenv("AUTH_PASSWORD")
	if password == "" {
		log.Println("Warning: AUTH_PASSWORD environment variable not set")
	}

	// Debug mode defaults to false
	debugMode := os.Getenv("DEBUG_MODE") == "true"
	
	// Local mode defaults to false
	localMode := os.Getenv("LOCAL_MODE") == "true"
	
	if localMode {
		log.Println("Running in LOCAL_MODE - authentication is disabled")
	}

	username := os.Getenv("AUTH_USERNAME")
	if username == "" {
		username = "duckduckgo-api" // Default value if not set
		log.Println("Warning: AUTH_USERNAME environment variable not set, using default")
	}

	return &Config{
		Port:         port,
		AuthUsername: username,
		AuthPassword: password,
		DebugMode:    debugMode,
		LocalMode:    localMode,
	}
}
