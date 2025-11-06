package logger

import (
	"time"

	"github.com/rs/zerolog"
)

// LogFields provides helper functions for common logging patterns

// SessionFields creates a logger with session-related fields
func SessionFields(sessionID, status string, connected bool) *zerolog.Event {
	return Log.Info().
		Str(FieldSessionID, sessionID).
		Str(FieldStatus, status).
		Bool("connected", connected)
}

// WebhookFields creates a logger with webhook-related fields
func WebhookFields(sessionID, event, url string, attempt int) *zerolog.Event {
	return Log.Info().
		Str(FieldSessionID, sessionID).
		Str(FieldEvent, event).
		Str(FieldURL, url).
		Int(FieldAttempt, attempt)
}

// MessageFields creates a logger with message-related fields
func MessageFields(sessionID, phone, messageID string) *zerolog.Event {
	return Log.Info().
		Str(FieldSessionID, sessionID).
		Str(FieldPhone, phone).
		Str(FieldMessageID, messageID)
}

// HTTPFields creates a logger with HTTP request fields
func HTTPFields(method, path, ip, userAgent string, status int, duration time.Duration) *zerolog.Event {
	return Log.Info().
		Str(FieldMethod, method).
		Str(FieldPath, path).
		Str(FieldIP, ip).
		Str(FieldUserAgent, userAgent).
		Int(FieldStatus, status).
		Dur(FieldDuration, duration)
}

// WorkerFields creates a logger with worker-related fields
func WorkerFields(workerID int, subject, queue string) *zerolog.Event {
	return Log.Info().
		Int(FieldWorkerID, workerID).
		Str(FieldSubject, subject).
		Str(FieldQueue, queue)
}

// ErrorFields creates a logger with error context
func ErrorFields(err error, component string, fields map[string]interface{}) *zerolog.Event {
	event := Log.Error().
		Err(err).
		Str(FieldComponent, component)

	for key, value := range fields {
		event = event.Interface(key, value)
	}

	return event
}

// NATSFields creates a logger with NATS-related fields
func NATSFields(subject, queue string) *zerolog.Event {
	return Log.Info().
		Str(FieldSubject, subject).
		Str(FieldQueue, queue)
}

// PerformanceFields creates a logger with performance metrics
func PerformanceFields(operation string, duration time.Duration, success bool) *zerolog.Event {
	return Log.Info().
		Str("operation", operation).
		Dur(FieldDuration, duration).
		Bool("success", success)
}

