package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type SessionStatus string

const (
	SessionStatusDisconnected SessionStatus = "disconnected"
	SessionStatusConnecting   SessionStatus = "connecting"
	SessionStatusConnected    SessionStatus = "connected"
	SessionStatusPairing      SessionStatus = "pairing"
	SessionStatusFailed       SessionStatus = "failed"
	SessionStatusLoggedOut    SessionStatus = "logged_out"
)

type PairingMethod string

const (
	PairingMethodQR    PairingMethod = "qr"
	PairingMethodPhone PairingMethod = "phone"
)

type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Protocol string `json:"protocol"` // http, https, socks5
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type WebhookConfig struct {
	Enabled bool     `json:"enabled"`
	URL     string   `json:"url"`
	Events  []string `json:"events"`
	Token   string   `json:"token,omitempty"`
}

type Session struct {
	ID        string // UUID gerado automaticamente
	Name      string
	DeviceJID string // JID do WhatsApp (após pareamento)
	Status    string // disconnected, connecting, connected, pairing, failed, logged_out
	Connected bool   // Flag rápida de conexão

	// WhatsApp data
	QRCode string // Base64 QR code para pareamento

	// Configuration (JSON)
	ProxyConfig   *ProxyConfig   // Configuração de proxy
	WebhookConfig *WebhookConfig // Configuração de webhook

	// Authentication
	APIKey *string // API key para autenticação da sessão (opcional)

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Session) IsConnected() bool {
	return s.Connected && s.Status == "connected"
}

func (s *Session) CanConnect() bool {
	return s.DeviceJID != "" && (s.Status == "disconnected" || s.Status == "failed")
}

func (s *Session) NeedsPairing() bool {
	return s.DeviceJID == ""
}

func (p *ProxyConfig) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func (p *ProxyConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, p)
}

func (w *WebhookConfig) Value() (driver.Value, error) {
	if w == nil {
		return nil, nil
	}
	return json.Marshal(w)
}

func (w *WebhookConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, w)
}

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		*s = []string{}
		return nil
	}

	return json.Unmarshal(bytes, s)
}

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		*j = make(map[string]interface{})
		return nil
	}

	return json.Unmarshal(bytes, j)
}
