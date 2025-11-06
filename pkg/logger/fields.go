package logger

import (
	"time"

	"github.com/rs/zerolog"
)

func SessionFields(sessionID, status string, connected bool) *zerolog.Event {
	return Log.Info().
		Str(FieldSessionID, sessionID).
		Str(FieldStatus, status).
		Bool("connected", connected)
}

func WebhookFields(sessionID, event, url string, attempt int) *zerolog.Event {
	return Log.Info().
		Str(FieldSessionID, sessionID).
		Str(FieldEvent, event).
		Str(FieldURL, url).
		Int(FieldAttempt, attempt)
}

func MessageFields(sessionID, phone, messageID string) *zerolog.Event {
	return Log.Info().
		Str(FieldSessionID, sessionID).
		Str(FieldPhone, phone).
		Str(FieldMessageID, messageID)
}

func HTTPFields(method, path, ip, userAgent string, status int, duration time.Duration) *zerolog.Event {
	return Log.Info().
		Str(FieldMethod, method).
		Str(FieldPath, path).
		Str(FieldIP, ip).
		Str(FieldUserAgent, userAgent).
		Int(FieldStatus, status).
		Dur(FieldDuration, duration)
}

func WorkerFields(workerID int, subject, queue string) *zerolog.Event {
	return Log.Info().
		Int(FieldWorkerID, workerID).
		Str(FieldSubject, subject).
		Str(FieldQueue, queue)
}

func ErrorFields(err error, component string, fields map[string]interface{}) *zerolog.Event {
	event := Log.Error().
		Err(err).
		Str(FieldComponent, component)

	for key, value := range fields {
		event = event.Interface(key, value)
	}

	return event
}

func NATSFields(subject, queue string) *zerolog.Event {
	return Log.Info().
		Str(FieldSubject, subject).
		Str(FieldQueue, queue)
}

func PerformanceFields(operation string, duration time.Duration, success bool) *zerolog.Event {
	return Log.Info().
		Str("operation", operation).
		Dur(FieldDuration, duration).
		Bool("success", success)
}
