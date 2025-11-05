package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Ensure no env overrides
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("DATABASE_DRIVER")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("WHATSAPP_DATA_DIR")

	if err := Load(); err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if AppConfig.Port != "8080" {
		t.Fatalf("expected default port 8080, got %s", AppConfig.Port)
	}
	if AppConfig.DatabaseDriver != "sqlite" {
		t.Fatalf("expected default driver sqlite, got %s", AppConfig.DatabaseDriver)
	}
}
