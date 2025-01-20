package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:8080,http://localhost:6060")

	conf := New()

	if len(conf.AllowedOrigins) != 2 {
		t.Fatalf("Expected 2 allowed origins, got %d", len(conf.AllowedOrigins))
	}

	if !conf.AllowedOrigins["http://localhost:8080"] {
		t.Fatal("http://localhost:8080 should be allowed to connect.")
	}
}
