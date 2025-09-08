package util

import (
	"testing"
)

func TestGetRandomHeaders(t *testing.T) {
	headers := GetRandomHeaders()

	// Check that all expected keys are present
	expectedKeys := []string{
		"User-Agent", "Accept", "Accept-Language", "Referer",
		"Cache-Control", "Pragma", "Sec-Fetch-Dest", "Sec-Fetch-Mode",
		"Sec-Fetch-Site", "Sec-Fetch-User", "Upgrade-Insecure-Requests",
	}
	for _, key := range expectedKeys {
		if _, exists := headers[key]; !exists {
			t.Errorf("Expected header key %s not found", key)
		}
	}

	// Check that values are non-empty
	for key, value := range headers {
		if value == "" {
			t.Errorf("Header %s has empty value", key)
		}
	}

	// Check randomization by generating multiple times and ensuring variety
	userAgents := make(map[string]bool)
	for i := 0; i < 10; i++ {
		h := GetRandomHeaders()
		userAgents[h["User-Agent"]] = true
	}
	if len(userAgents) < 2 {
		t.Error("User-Agent randomization not working; expected variety")
	}
}

func TestGetRandomHeadersFallback(t *testing.T) {
	// This test is hard to trigger without mocking rand.Reader failure,
	// but we can at least ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetRandomHeaders panicked: %v", r)
		}
	}()
	headers := GetRandomHeaders()
	if len(headers) == 0 {
		t.Error("Headers map is empty")
	}
}
