package logger

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// Log is the global logger instance
var Log zerolog.Logger

// Field name constants for consistency across the application
const (
	FieldSessionID   = "session_id"
	FieldWorkerID    = "worker_id"
	FieldEvent       = "event"
	FieldURL         = "url"
	FieldPhone       = "phone"
	FieldMessageID   = "message_id"
	FieldAttempt     = "attempt"
	FieldStatus      = "status"
	FieldDuration    = "duration"
	FieldError       = "error"
	FieldRequestID   = "request_id"
	FieldMethod      = "method"
	FieldPath        = "path"
	FieldIP          = "ip"
	FieldUserAgent   = "user_agent"
	FieldSubject     = "subject"
	FieldQueue       = "queue"
	FieldComponent   = "component"
	FieldEnvironment = "environment"
	FieldLogLevel    = "log_level"
)

// Config holds logger configuration
type Config struct {
	Level       string
	Format      string // "console" or "json"
	Output      io.Writer
	AddCaller   bool
	SampleRate  int // 0 = no sampling, N = log 1 out of N messages
	Environment string
	Service     string
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:       "info",
		Format:      "console",
		Output:      os.Stdout,
		AddCaller:   false,
		SampleRate:  0,
		Environment: "development",
		Service:     "zpwoot",
	}
}

// Init initializes the global logger with default settings
func Init(level string) {
	cfg := DefaultConfig()
	cfg.Level = level
	InitWithConfig(cfg)
}

// InitWithConfig initializes the global logger with custom configuration
func InitWithConfig(cfg Config) {
	// Parse log level
	lvl := parseLevel(cfg.Level)

	// Create writer based on format
	var writer io.Writer = cfg.Output
	if cfg.Format == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        cfg.Output,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
	}

	// Create base logger
	baseLogger := zerolog.New(writer).
		With().
		Timestamp().
		Str("service", cfg.Service).
		Str(FieldEnvironment, cfg.Environment)

	// Add caller information if enabled
	if cfg.AddCaller {
		baseLogger = baseLogger.Caller()
	}

	// Build logger
	Log = baseLogger.Logger().Level(lvl)

	// Apply sampling if configured
	if cfg.SampleRate > 0 {
		sampler := &zerolog.BurstSampler{
			Burst:       uint32(cfg.SampleRate),
			Period:      1 * time.Second,
			NextSampler: &zerolog.BasicSampler{N: uint32(cfg.SampleRate)},
		}
		Log = Log.Sample(sampler)
	}

	// Configure global logger
	zlog.Logger = Log
	zerolog.SetGlobalLevel(lvl)

	// Set time field format
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Log initialization
	Log.Info().
		Str(FieldLogLevel, cfg.Level).
		Str("format", cfg.Format).
		Bool("caller", cfg.AddCaller).
		Int("sample_rate", cfg.SampleRate).
		Msg("Logger initialized")
}

// parseLevel converts string level to zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// WithContext creates a logger with context fields
func WithContext(ctx context.Context) zerolog.Logger {
	logger := Log

	// Extract common context values
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		logger = logger.With().Str(FieldRequestID, requestID).Logger()
	}

	if sessionID, ok := ctx.Value("session_id").(string); ok && sessionID != "" {
		logger = logger.With().Str(FieldSessionID, sessionID).Logger()
	}

	return logger
}

// WithComponent creates a logger for a specific component
func WithComponent(component string) zerolog.Logger {
	return Log.With().Str(FieldComponent, component).Logger()
}

// WithSession creates a logger with session context
func WithSession(sessionID string) zerolog.Logger {
	return Log.With().Str(FieldSessionID, sessionID).Logger()
}

// WithWorker creates a logger with worker context
func WithWorker(workerID int) zerolog.Logger {
	return Log.With().Int(FieldWorkerID, workerID).Logger()
}

// WithFields creates a logger with custom fields
func WithFields(fields map[string]interface{}) zerolog.Logger {
	logger := Log.With()
	for key, value := range fields {
		logger = logger.Interface(key, value)
	}
	return logger.Logger()
}

// SetLevel dynamically changes the log level
func SetLevel(level string) {
	lvl := parseLevel(level)
	zerolog.SetGlobalLevel(lvl)
	Log = Log.Level(lvl)
	Log.Info().Str(FieldLogLevel, level).Msg("Log level changed")
}

// GetLevel returns the current log level as string
func GetLevel() string {
	return zerolog.GlobalLevel().String()
}
