package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

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

	// NATS Configuration
	NATSURL           string
	NATSMaxReconnect  int
	NATSReconnectWait time.Duration

	// Webhook Configuration
	WebhookWorkers        int
	WebhookTimeout        time.Duration
	WebhookMaxRetries     int
	WebhookRetryBaseDelay time.Duration
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

		// NATS
		NATSURL:           getEnv("NATS_URL", "nats://localhost:4222"),
		NATSMaxReconnect:  getEnvInt("NATS_MAX_RECONNECT", 10),
		NATSReconnectWait: getEnvDuration("NATS_RECONNECT_WAIT", 2*time.Second),

		// Webhooks
		WebhookWorkers:        getEnvInt("WEBHOOK_WORKERS", 10),
		WebhookTimeout:        getEnvDuration("WEBHOOK_TIMEOUT", 30*time.Second),
		WebhookMaxRetries:     getEnvInt("WEBHOOK_MAX_RETRIES", 3),
		WebhookRetryBaseDelay: getEnvDuration("WEBHOOK_RETRY_BASE_DELAY", 5*time.Second),
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

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultVal
	}
	return d
}
