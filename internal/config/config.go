package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	DatabaseDriver  string
	APIKey          string
	Environment     string
	LogLevel        string
	WhatsAppDataDir string

	// WhatsApp Session Configuration
	MaxSessions         int
	ConnectionTimeout   int
	PairingTimeout      int
	AutoRestoreSessions bool
}

var AppConfig *Config

func Load() error {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		DatabaseDriver:      getEnv("DATABASE_DRIVER", "postgres"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://zpwoot:zpwoot_password_change_in_production@localhost:5432/zpwoot?sslmode=disable"),
		Environment:         getEnv("ENVIRONMENT", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		WhatsAppDataDir:     getEnv("WHATSAPP_DATA_DIR", "./data"),
		APIKey:              os.Getenv("API_KEY"),
		MaxSessions:         getEnvInt("MAX_SESSIONS", 10),
		ConnectionTimeout:   getEnvInt("CONNECTION_TIMEOUT", 30),
		PairingTimeout:      getEnvInt("PAIRING_TIMEOUT", 120),
		AutoRestoreSessions: getEnvBool("AUTO_RESTORE_SESSIONS", true),
	}

	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL cannot be empty")
	}

	AppConfig = cfg
	return nil
}

func GetDatabaseDSN() string {
	if AppConfig == nil {
		return ""
	}
	return AppConfig.DatabaseURL
}

func getEnv(key, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}

func getEnvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return i
}

func getEnvBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}
