package service

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	natsclient "zpwoot/internal/nats"
	"zpwoot/pkg/logger"
)

// WebhookWorker consumes webhook messages from NATS and delivers them
type WebhookWorker struct {
	id              int
	natsClient      *natsclient.Client
	delivery        *WebhookDelivery
	maxRetries      int
	retryBaseDelay  time.Duration
	subscription    *nats.Subscription
	log             zerolog.Logger // Worker-specific logger with context
}

// NewWebhookWorker creates a new webhook worker
func NewWebhookWorker(
	id int,
	natsClient *natsclient.Client,
	delivery *WebhookDelivery,
	maxRetries int,
	retryBaseDelay time.Duration,
) *WebhookWorker {
	// Create worker-specific logger with context
	workerLog := logger.WithWorker(id)

	return &WebhookWorker{
		id:             id,
		natsClient:     natsClient,
		delivery:       delivery,
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		log:            workerLog,
	}
}

// Start starts the worker
func (w *WebhookWorker) Start() error {
	w.log.Info().
		Str(logger.FieldSubject, "webhooks.*").
		Str(logger.FieldQueue, "webhook-workers").
		Msg("Starting webhook worker")

	// Subscribe to webhooks.* with queue group for load balancing
	sub, err := w.natsClient.QueueSubscribe("webhooks.*", "webhook-workers", w.handleMessage)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	w.subscription = sub

	w.log.Info().Msg("✅ Webhook worker started and listening")

	return nil
}

// Stop stops the worker
func (w *WebhookWorker) Stop() error {
	if w.subscription != nil {
		return w.subscription.Unsubscribe()
	}
	return nil
}

// handleMessage processes a webhook message from NATS
func (w *WebhookWorker) handleMessage(msg *nats.Msg) {
	// Parse webhook message
	var webhookMsg WebhookMessage
	err := json.Unmarshal(msg.Data, &webhookMsg)
	if err != nil {
		w.log.Error().
			Err(err).
			Msg("Failed to unmarshal webhook message")
		msg.Ack() // ACK to remove from queue
		return
	}

	// Create session-specific logger
	sessionLog := w.log.With().
		Str(logger.FieldSessionID, webhookMsg.SessionID).
		Str(logger.FieldURL, webhookMsg.WebhookURL).
		Str(logger.FieldEvent, webhookMsg.Payload.Event).
		Int(logger.FieldAttempt, webhookMsg.Attempt).
		Logger()

	sessionLog.Debug().Msg("Processing webhook message")

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(webhookMsg.Payload)
	if err != nil {
		sessionLog.Error().
			Err(err).
			Msg("Failed to marshal webhook payload")
		msg.Ack() // ACK to remove from queue
		return
	}

	// Deliver webhook
	result := w.delivery.Send(webhookMsg.WebhookURL, payloadBytes, webhookMsg.WebhookToken)

	// Handle result
	if result.Success {
		// Success - ACK message
		msg.Ack()
		sessionLog.Info().
			Int(logger.FieldStatus, result.StatusCode).
			Dur(logger.FieldDuration, result.Duration).
			Msg("✅ Webhook delivered successfully")
	} else {
		// Failed - check if should retry
		if webhookMsg.Attempt < w.maxRetries && IsRetryableError(result) {
			// Retry with exponential backoff
			w.retryWebhook(msg, &webhookMsg, sessionLog)
		} else {
			// Max retries reached or non-retryable error - move to DLQ
			w.moveToDLQ(&webhookMsg, result, sessionLog)
			msg.Ack() // ACK original message
		}
	}
}

// retryWebhook retries a failed webhook with exponential backoff
func (w *WebhookWorker) retryWebhook(msg *nats.Msg, webhookMsg *WebhookMessage, sessionLog zerolog.Logger) {
	// Calculate delay: 5s, 25s, 125s (exponential)
	delay := w.calculateRetryDelay(webhookMsg.Attempt)

	sessionLog.Warn().
		Int("max_retries", w.maxRetries).
		Dur("retry_delay", delay).
		Msg("⚠️ Webhook delivery failed, scheduling retry")

	// Increment attempt
	webhookMsg.Attempt++

	// Wait for delay
	time.Sleep(delay)

	// Re-publish to NATS for retry
	data, err := json.Marshal(webhookMsg)
	if err != nil {
		sessionLog.Error().
			Err(err).
			Msg("Failed to marshal webhook for retry")
		msg.Ack()
		return
	}

	subject := fmt.Sprintf("webhooks.%s", webhookMsg.SessionID)
	err = w.natsClient.Publish(subject, data)
	if err != nil {
		sessionLog.Error().
			Err(err).
			Str(logger.FieldSubject, subject).
			Msg("Failed to republish webhook for retry")
	}

	// ACK original message
	msg.Ack()
}

// moveToDLQ moves a failed webhook to the dead letter queue
func (w *WebhookWorker) moveToDLQ(webhookMsg *WebhookMessage, result *DeliveryResult, sessionLog zerolog.Logger) {
	sessionLog.Error().
		Int("attempts", webhookMsg.Attempt).
		Str("error", fmt.Sprintf("%v", result.Error)).
		Int(logger.FieldStatus, result.StatusCode).
		Msg("❌ Webhook delivery failed permanently, moving to DLQ")

	// Publish to DLQ
	data, err := json.Marshal(webhookMsg)
	if err != nil {
		sessionLog.Error().
			Err(err).
			Msg("Failed to marshal webhook for DLQ")
		return
	}

	err = w.natsClient.Publish("webhooks.dlq", data)
	if err != nil {
		sessionLog.Error().
			Err(err).
			Str(logger.FieldSubject, "webhooks.dlq").
			Msg("Failed to publish to DLQ")
	}
}

// calculateRetryDelay calculates exponential backoff delay
func (w *WebhookWorker) calculateRetryDelay(attempt int) time.Duration {
	// Attempt 1: 5s
	// Attempt 2: 25s (5s * 5)
	// Attempt 3: 125s (25s * 5)
	multiplier := math.Pow(5, float64(attempt-1))
	return time.Duration(float64(w.retryBaseDelay) * multiplier)
}

