package config

import (
	"os"
	"testing"
)

func TestLocalMode(t *testing.T) {
	// Save original env vars to restore later
	originalLocalMode := os.Getenv("LOCAL_MODE")
	originalUsername := os.Getenv("AUTH_USERNAME")
	originalPassword := os.Getenv("AUTH_PASSWORD")
	
	defer func() {
		os.Setenv("LOCAL_MODE", originalLocalMode)
		os.Setenv("AUTH_USERNAME", originalUsername)
		os.Setenv("AUTH_PASSWORD", originalPassword)
	}()

	// Set auth credentials for non-local mode tests
	os.Setenv("AUTH_USERNAME", "testuser")
	os.Setenv("AUTH_PASSWORD", "testpass")

	// Test when LOCAL_MODE is not set
	os.Unsetenv("LOCAL_MODE")
	cfg, err := New()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.LocalMode {
		t.Error("LocalMode should be false when LOCAL_MODE env var is not set")
	}

	// Test when LOCAL_MODE is set to true
	os.Setenv("LOCAL_MODE", "true")
	cfg, err = New()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !cfg.LocalMode {
		t.Error("LocalMode should be true when LOCAL_MODE env var is set to 'true'")
	}

	// Test when LOCAL_MODE is set to something else
	os.Setenv("LOCAL_MODE", "yes")
	cfg, err = New()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.LocalMode {
		t.Error("LocalMode should be false when LOCAL_MODE env var is not 'true'")
	}
}

func TestAuthCredentials(t *testing.T) {
	// Save original env vars to restore later
	originalLocalMode := os.Getenv("LOCAL_MODE")
	originalUsername := os.Getenv("AUTH_USERNAME")
	originalPassword := os.Getenv("AUTH_PASSWORD")
	
	defer func() {
		os.Setenv("LOCAL_MODE", originalLocalMode)
		os.Setenv("AUTH_USERNAME", originalUsername)
		os.Setenv("AUTH_PASSWORD", originalPassword)
	}()

	// Test missing username in non-local mode
	os.Unsetenv("LOCAL_MODE")
	os.Unsetenv("AUTH_USERNAME")
	os.Setenv("AUTH_PASSWORD", "testpass")
	
	_, err := New()
	if err == nil {
		t.Error("Expected error for missing username in non-local mode")
	}

	// Test missing password in non-local mode
	os.Setenv("AUTH_USERNAME", "testuser")
	os.Unsetenv("AUTH_PASSWORD")
	
	_, err = New()
	if err == nil {
		t.Error("Expected error for missing password in non-local mode")
	}

	// Test local mode with missing credentials (should not error)
	os.Setenv("LOCAL_MODE", "true")
	os.Unsetenv("AUTH_USERNAME")
	os.Unsetenv("AUTH_PASSWORD")
	
	cfg, err := New()
	if err != nil {
		t.Errorf("Unexpected error in local mode: %v", err)
	}
	if !cfg.LocalMode {
		t.Error("LocalMode should be true")
	}
}
