package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// SessionStatus representa os possíveis estados de uma sessão
type SessionStatus string

const (
	SessionStatusDisconnected SessionStatus = "disconnected"
	SessionStatusConnecting   SessionStatus = "connecting"
	SessionStatusConnected    SessionStatus = "connected"
	SessionStatusPairing      SessionStatus = "pairing"
	SessionStatusFailed       SessionStatus = "failed"
	SessionStatusLoggedOut    SessionStatus = "logged_out"
)

// PairingMethod representa o método de pareamento
type PairingMethod string

const (
	PairingMethodQR    PairingMethod = "qr"
	PairingMethodPhone PairingMethod = "phone"
)

// ProxyConfig representa a configuração de proxy
type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Protocol string `json:"protocol"` // http, https, socks5
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// WebhookConfig representa a configuração de webhook
type WebhookConfig struct {
	Enabled bool     `json:"enabled"`
	URL     string   `json:"url"`
	Events  []string `json:"events"`
	Token   string   `json:"token,omitempty"`
}

// Session representa uma sessão WhatsApp
// Usamos structs simples sem tags GORM, pois usamos database/sql nativo
type Session struct {
	ID            string // UUID gerado automaticamente
	Name          string
	DeviceJID     string // JID do WhatsApp (após pareamento)
	Status        string // disconnected, connecting, connected, pairing, failed, logged_out
	Connected     bool   // Flag rápida de conexão

	// WhatsApp data
	QRCode        string // Base64 QR code para pareamento

	// Configuration (JSON)
	ProxyConfig   *ProxyConfig   // Configuração de proxy
	WebhookConfig *WebhookConfig // Configuração de webhook

	// Authentication
	APIKey        *string // API key para autenticação da sessão (opcional)

	// Timestamps
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// IsConnected verifica se a sessão está conectada
func (s *Session) IsConnected() bool {
	return s.Connected && s.Status == "connected"
}

// CanConnect verifica se a sessão pode conectar
func (s *Session) CanConnect() bool {
	return s.DeviceJID != "" && (s.Status == "disconnected" || s.Status == "failed")
}

// NeedsPairing verifica se precisa de pareamento
func (s *Session) NeedsPairing() bool {
	return s.DeviceJID == ""
}

// Value implementa driver.Valuer para ProxyConfig
func (p *ProxyConfig) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan implementa sql.Scanner para ProxyConfig
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

// Value implementa driver.Valuer para WebhookConfig
func (w *WebhookConfig) Value() (driver.Value, error) {
	if w == nil {
		return nil, nil
	}
	return json.Marshal(w)
}

// Scan implementa sql.Scanner para WebhookConfig
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

// StringArray é um helper para converter []string para JSON
type StringArray []string

// Value implementa driver.Valuer para salvar no banco
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan implementa sql.Scanner para ler do banco
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

// JSONMap é um helper para converter map[string]interface{} para JSON
type JSONMap map[string]interface{}

// Value implementa driver.Valuer para salvar no banco
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}
	return json.Marshal(j)
}

// Scan implementa sql.Scanner para ler do banco
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
