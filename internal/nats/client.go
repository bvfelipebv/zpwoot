package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"zpwoot/pkg/logger"
)

// Client wraps NATS connection
type Client struct {
	conn          *nats.Conn
	url           string
	maxReconnect  int
	reconnectWait time.Duration
}

// Config holds NATS client configuration
type Config struct {
	URL           string
	MaxReconnect  int
	ReconnectWait time.Duration
}

// NewClient creates a new NATS client
func NewClient(cfg Config) *Client {
	return &Client{
		url:           cfg.URL,
		maxReconnect:  cfg.MaxReconnect,
		reconnectWait: cfg.ReconnectWait,
	}
}

// Connect establishes connection to NATS server
func (c *Client) Connect() error {
	opts := []nats.Option{
		nats.MaxReconnects(c.maxReconnect),
		nats.ReconnectWait(c.reconnectWait),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				logger.Log.Error().Err(err).Msg("NATS disconnected")
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Log.Info().
				Str("url", nc.ConnectedUrl()).
				Msg("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Log.Warn().Msg("NATS connection closed")
		}),
	}

	conn, err := nats.Connect(c.url, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	c.conn = conn

	logger.Log.Info().
		Str("url", c.url).
		Str("server_id", conn.ConnectedServerId()).
		Msg("Connected to NATS")

	return nil
}

// Publish publishes a message to a subject
func (c *Client) Publish(subject string, data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("NATS connection not established")
	}

	err := c.conn.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish to %s: %w", subject, err)
	}

	return nil
}

// Subscribe subscribes to a subject with a handler
func (c *Client) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("NATS connection not established")
	}

	sub, err := c.conn.Subscribe(subject, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	logger.Log.Info().
		Str("subject", subject).
		Msg("Subscribed to NATS subject")

	return sub, nil
}

// QueueSubscribe subscribes to a subject with queue group
func (c *Client) QueueSubscribe(subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("NATS connection not established")
	}

	sub, err := c.conn.QueueSubscribe(subject, queue, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to queue subscribe to %s: %w", subject, err)
	}

	logger.Log.Info().
		Str("subject", subject).
		Str("queue", queue).
		Msg("Queue subscribed to NATS subject")

	return sub, nil
}

// IsConnected checks if connection is active
func (c *Client) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

// Stats returns connection statistics
func (c *Client) Stats() nats.Statistics {
	if c.conn == nil {
		return nats.Statistics{}
	}
	return c.conn.Stats()
}

// Close closes the NATS connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		logger.Log.Info().Msg("NATS connection closed")
	}
}

// Drain drains the connection (graceful shutdown)
func (c *Client) Drain() error {
	if c.conn != nil {
		return c.conn.Drain()
	}
	return nil
}

