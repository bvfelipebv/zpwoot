package dto

import "time"

// SessionResponse - Resposta padrão de sessão
type SessionResponse struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	JID              string     `json:"jid,omitempty"`
	Status           string     `json:"status"`
	PushName         string     `json:"push_name,omitempty"`
	Platform         string     `json:"platform,omitempty"`
	BusinessName     string     `json:"business_name,omitempty"`
	WebhookURL       string     `json:"webhook_url,omitempty"`
	WebhookEvents    []string   `json:"webhook_events,omitempty"`
	LastConnected    *time.Time `json:"last_connected,omitempty"`
	LastDisconnected *time.Time `json:"last_disconnected,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// SessionListResponse - Lista de sessões
type SessionListResponse struct {
	Sessions []SessionResponse `json:"sessions"`
	Total    int               `json:"total"`
}

// SessionStatusResponse - Status detalhado da sessão
type SessionStatusResponse struct {
	ID             string     `json:"id"`
	Status         string     `json:"status"`
	IsConnected    bool       `json:"is_connected"`
	IsLoggedIn     bool       `json:"is_logged_in"`
	JID            string     `json:"jid,omitempty"`
	PushName       string     `json:"push_name,omitempty"`
	Platform       string     `json:"platform,omitempty"`
	LastConnected  *time.Time `json:"last_connected,omitempty"`
	ConnectionTime string     `json:"connection_time,omitempty"` // Duração formatada
	NeedsPairing   bool       `json:"needs_pairing"`
	CanConnect     bool       `json:"can_connect"`
}

// PairQRResponse - Resposta com QR code
type PairQRResponse struct {
	SessionID string    `json:"session_id"`
	QRCode    string    `json:"qr_code"` // Base64 data URL
	ExpiresAt time.Time `json:"expires_at"`
	Message   string    `json:"message"`
}

// PairPhoneResponse - Resposta de pareamento por telefone
type PairPhoneResponse struct {
	SessionID   string `json:"session_id"`
	PhoneNumber string `json:"phone_number"`
	PairingCode string `json:"pairing_code"`
	Message     string `json:"message"`
}

// ErrorResponse - Resposta de erro padrão
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse - Resposta de sucesso genérica
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

