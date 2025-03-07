package config

import (
	"os"
	"testing"
)

func TestLocalMode(t *testing.T) {
	// Save original env var value to restore later
	originalValue := os.Getenv("LOCAL_MODE")
	defer os.Setenv("LOCAL_MODE", originalValue)

	// Test when LOCAL_MODE is not set
	os.Unsetenv("LOCAL_MODE")
	cfg := New()
	if cfg.LocalMode {
		t.Error("LocalMode should be false when LOCAL_MODE env var is not set")
	}

	// Test when LOCAL_MODE is set to true
	os.Setenv("LOCAL_MODE", "true")
	cfg = New()
	if !cfg.LocalMode {
		t.Error("LocalMode should be true when LOCAL_MODE env var is set to 'true'")
	}

	// Test when LOCAL_MODE is set to something else
	os.Setenv("LOCAL_MODE", "yes")
	cfg = New()
	if cfg.LocalMode {
		t.Error("LocalMode should be false when LOCAL_MODE env var is not 'true'")
	}
}
