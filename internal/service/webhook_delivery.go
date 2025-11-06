package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"zpwoot/pkg/logger"
)

type WebhookDelivery struct {
	timeout time.Duration
	client  *http.Client
}

func NewWebhookDelivery(timeout time.Duration) *WebhookDelivery {
	return &WebhookDelivery{
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type DeliveryResult struct {
	Success      bool
	StatusCode   int
	ResponseBody string
	Error        error
	Duration     time.Duration
}

func (d *WebhookDelivery) Send(url string, payload []byte, token string) *DeliveryResult {
	startTime := time.Now()

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to create webhook request")

		return &DeliveryResult{
			Success:  false,
			Error:    fmt.Errorf("failed to create request: %w", err),
			Duration: time.Since(startTime),
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "zpwoot-webhook/1.0")

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	// Send request
	resp, err := d.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("url", url).
			Dur("duration", duration).
			Msg("Failed to send webhook")

		return &DeliveryResult{
			Success:  false,
			Error:    fmt.Errorf("failed to send request: %w", err),
			Duration: duration,
		}
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Warn().
			Err(err).
			Str("url", url).
			Int("status", resp.StatusCode).
			Msg("Failed to read webhook response body")
		bodyBytes = []byte{}
	}

	responseBody := string(bodyBytes)

	// Check if successful (2xx status codes)
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	if success {
		logger.Log.Info().
			Str("url", url).
			Int("status", resp.StatusCode).
			Dur("duration", duration).
			Msg("Webhook delivered successfully")
	} else {
		logger.Log.Warn().
			Str("url", url).
			Int("status", resp.StatusCode).
			Str("response", responseBody).
			Dur("duration", duration).
			Msg("Webhook delivery failed with non-2xx status")
	}

	return &DeliveryResult{
		Success:      success,
		StatusCode:   resp.StatusCode,
		ResponseBody: responseBody,
		Duration:     duration,
	}
}

func IsRetryableError(result *DeliveryResult) bool {
	// Retry on network errors
	if result.Error != nil {
		return true
	}

	// Retry on 5xx server errors
	if result.StatusCode >= 500 && result.StatusCode < 600 {
		return true
	}

	// Retry on 429 Too Many Requests
	if result.StatusCode == 429 {
		return true
	}

	// Retry on 408 Request Timeout
	if result.StatusCode == 408 {
		return true
	}

	// Don't retry on 4xx client errors (except 408 and 429)
	return false
}
