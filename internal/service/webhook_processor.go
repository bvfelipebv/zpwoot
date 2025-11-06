package service

import (
	"context"
	"encoding/json"
	"fmt"

	"zpwoot/internal/constants"
	"zpwoot/internal/model"
	natsclient "zpwoot/internal/nats"
	"zpwoot/internal/repository"
	"zpwoot/pkg/logger"
)

type WebhookProcessor struct {
	natsClient  *natsclient.Client
	formatter   *WebhookFormatter
	sessionRepo *repository.SessionRepository
}

func NewWebhookProcessor(
	natsClient *natsclient.Client,
	formatter *WebhookFormatter,
	sessionRepo *repository.SessionRepository,
) *WebhookProcessor {
	return &WebhookProcessor{
		natsClient:  natsClient,
		formatter:   formatter,
		sessionRepo: sessionRepo,
	}
}

type WebhookMessage struct {
	SessionID    string          `json:"session_id"`
	WebhookURL   string          `json:"webhook_url"`
	WebhookToken string          `json:"webhook_token,omitempty"`
	Attempt      int             `json:"attempt"`
	Payload      *WebhookPayload `json:"payload"`
}

func (p *WebhookProcessor) ProcessEvent(sessionID string, eventType constants.WebhookEventType, payload *WebhookPayload) error {
	// 1. Get session from database
	ctx := context.Background()
	session, err := p.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get session for webhook")
		return fmt.Errorf("failed to get session: %w", err)
	}

	// 2. Check if webhook is configured and enabled
	if session.WebhookConfig == nil || !session.WebhookConfig.Enabled {
		logger.Log.Debug().
			Str("session_id", sessionID).
			Str("event", string(eventType)).
			Msg("Webhook not enabled for session")
		return nil
	}

	// 3. Validate webhook URL
	if session.WebhookConfig.URL == "" {
		logger.Log.Warn().
			Str("session_id", sessionID).
			Msg("Webhook enabled but URL is empty")
		return nil
	}

	// 4. Check if event is subscribed
	if !p.isEventSubscribed(session.WebhookConfig.Events, eventType) {
		logger.Log.Debug().
			Str("session_id", sessionID).
			Str("event", string(eventType)).
			Msg("Event not subscribed")
		return nil
	}

	// 5. Create webhook message
	webhookMsg := &WebhookMessage{
		SessionID:    sessionID,
		WebhookURL:   session.WebhookConfig.URL,
		WebhookToken: session.WebhookConfig.Token,
		Attempt:      1,
		Payload:      payload,
	}

	// 6. Marshal to JSON
	data, err := json.Marshal(webhookMsg)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to marshal webhook message")
		return fmt.Errorf("failed to marshal webhook message: %w", err)
	}

	// 7. Publish to NATS
	subject := fmt.Sprintf("webhooks.%s", sessionID)
	err = p.natsClient.Publish(subject, data)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("subject", subject).
			Msg("Failed to publish webhook to NATS")
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Str("event", string(eventType)).
		Str("subject", subject).
		Str("url", session.WebhookConfig.URL).
		Msg("Webhook published to NATS")

	return nil
}

func (p *WebhookProcessor) isEventSubscribed(subscribedEvents []string, eventType constants.WebhookEventType) bool {
	// If no events specified, subscribe to all
	if len(subscribedEvents) == 0 {
		return true
	}

	eventStr := string(eventType)
	for _, evt := range subscribedEvents {
		if evt == eventStr {
			return true
		}
		// Support wildcard "*" to subscribe to all events
		if evt == "*" {
			return true
		}
	}

	return false
}

func ValidateWebhookConfig(config *model.WebhookConfig) error {
	if config == nil {
		return fmt.Errorf("webhook config is nil")
	}

	if config.Enabled && config.URL == "" {
		return fmt.Errorf("webhook URL cannot be empty when enabled")
	}

	// Validate events
	for _, evt := range config.Events {
		if evt != "*" && !constants.IsValidEventType(evt) {
			return fmt.Errorf("invalid event type: %s", evt)
		}
	}

	return nil
}

