package main

import (
	"ddg-search/config"
	"os"
	"testing"
)

func TestConfigLoading(t *testing.T) {
	// Save original env var value to restore later
	originalValue := os.Getenv("LOCAL_MODE")
	defer os.Setenv("LOCAL_MODE", originalValue)

	// Test local mode setting
	os.Setenv("LOCAL_MODE", "true")
	cfg := config.New()

	if !cfg.LocalMode {
		t.Error("LocalMode should be true when LOCAL_MODE env var is set to 'true'")
	}
}
