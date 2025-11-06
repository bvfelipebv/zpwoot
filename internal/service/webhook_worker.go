package service

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/nats-io/nats.go"
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
}

// NewWebhookWorker creates a new webhook worker
func NewWebhookWorker(
	id int,
	natsClient *natsclient.Client,
	delivery *WebhookDelivery,
	maxRetries int,
	retryBaseDelay time.Duration,
) *WebhookWorker {
	return &WebhookWorker{
		id:             id,
		natsClient:     natsClient,
		delivery:       delivery,
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
	}
}

// Start starts the worker
func (w *WebhookWorker) Start() error {
	logger.Log.Info().
		Int("worker_id", w.id).
		Msg("Starting webhook worker")

	// Subscribe to webhooks.* with queue group for load balancing
	sub, err := w.natsClient.QueueSubscribe("webhooks.*", "webhook-workers", w.handleMessage)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	w.subscription = sub

	logger.Log.Info().
		Int("worker_id", w.id).
		Msg("Webhook worker started and listening")

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
		logger.Log.Error().
			Err(err).
			Int("worker_id", w.id).
			Msg("Failed to unmarshal webhook message")
		msg.Ack() // ACK to remove from queue
		return
	}

	logger.Log.Debug().
		Int("worker_id", w.id).
		Str("session_id", webhookMsg.SessionID).
		Str("url", webhookMsg.WebhookURL).
		Int("attempt", webhookMsg.Attempt).
		Msg("Processing webhook message")

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(webhookMsg.Payload)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Int("worker_id", w.id).
			Str("session_id", webhookMsg.SessionID).
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
		logger.Log.Info().
			Int("worker_id", w.id).
			Str("session_id", webhookMsg.SessionID).
			Str("event", webhookMsg.Payload.Event).
			Int("status", result.StatusCode).
			Dur("duration", result.Duration).
			Msg("Webhook delivered successfully")
	} else {
		// Failed - check if should retry
		if webhookMsg.Attempt < w.maxRetries && IsRetryableError(result) {
			// Retry with exponential backoff
			w.retryWebhook(msg, &webhookMsg)
		} else {
			// Max retries reached or non-retryable error - move to DLQ
			w.moveToDLQ(&webhookMsg, result)
			msg.Ack() // ACK original message
		}
	}
}

// retryWebhook retries a failed webhook with exponential backoff
func (w *WebhookWorker) retryWebhook(msg *nats.Msg, webhookMsg *WebhookMessage) {
	// Calculate delay: 5s, 25s, 125s (exponential)
	delay := w.calculateRetryDelay(webhookMsg.Attempt)

	logger.Log.Warn().
		Int("worker_id", w.id).
		Str("session_id", webhookMsg.SessionID).
		Int("attempt", webhookMsg.Attempt).
		Int("max_retries", w.maxRetries).
		Dur("retry_delay", delay).
		Msg("Webhook delivery failed, scheduling retry")

	// Increment attempt
	webhookMsg.Attempt++

	// Wait for delay
	time.Sleep(delay)

	// Re-publish to NATS for retry
	data, err := json.Marshal(webhookMsg)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", webhookMsg.SessionID).
			Msg("Failed to marshal webhook for retry")
		msg.Ack()
		return
	}

	subject := fmt.Sprintf("webhooks.%s", webhookMsg.SessionID)
	err = w.natsClient.Publish(subject, data)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", webhookMsg.SessionID).
			Msg("Failed to republish webhook for retry")
	}

	// ACK original message
	msg.Ack()
}

// moveToDLQ moves a failed webhook to the dead letter queue
func (w *WebhookWorker) moveToDLQ(webhookMsg *WebhookMessage, result *DeliveryResult) {
	logger.Log.Error().
		Int("worker_id", w.id).
		Str("session_id", webhookMsg.SessionID).
		Str("event", webhookMsg.Payload.Event).
		Int("attempts", webhookMsg.Attempt).
		Str("url", webhookMsg.WebhookURL).
		Str("error", fmt.Sprintf("%v", result.Error)).
		Int("status", result.StatusCode).
		Msg("Webhook delivery failed permanently, moving to DLQ")

	// Publish to DLQ
	data, err := json.Marshal(webhookMsg)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Failed to marshal webhook for DLQ")
		return
	}

	err = w.natsClient.Publish("webhooks.dlq", data)
	if err != nil {
		logger.Log.Error().
			Err(err).
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

