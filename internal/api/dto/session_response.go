package dto

import "time"

type SessionResponse struct {
	ID               string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name             string     `json:"name" example:"Minha Sessão WhatsApp"`
	JID              string     `json:"jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	Status           string     `json:"status" example:"connected"`
	PushName         string     `json:"push_name,omitempty" example:"João Silva"`
	Platform         string     `json:"platform,omitempty" example:"android"`
	BusinessName     string     `json:"business_name,omitempty" example:"Minha Empresa LTDA"`
	WebhookURL       string     `json:"webhook_url,omitempty" example:"https://seu-webhook.com/whatsapp"`
	WebhookEvents    []string   `json:"webhook_events,omitempty" example:"message,qr,connected"`
	LastConnected    *time.Time `json:"last_connected,omitempty" example:"2025-11-05T18:30:00Z"`
	LastDisconnected *time.Time `json:"last_disconnected,omitempty" example:"2025-11-05T17:00:00Z"`
	CreatedAt        time.Time  `json:"created_at" example:"2025-11-05T10:00:00Z"`
	UpdatedAt        time.Time  `json:"updated_at" example:"2025-11-05T18:30:00Z"`
}

type SessionListResponse struct {
	Sessions []SessionResponse `json:"sessions"`
	Total    int               `json:"total" example:"3"`
}

type SessionStatusResponse struct {
	ID             string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status         string     `json:"status" example:"connected"`
	IsConnected    bool       `json:"is_connected" example:"true"`
	IsLoggedIn     bool       `json:"is_logged_in" example:"true"`
	JID            string     `json:"jid,omitempty" example:"5511999999999@s.whatsapp.net"`
	PushName       string     `json:"push_name,omitempty" example:"João Silva"`
	Platform       string     `json:"platform,omitempty" example:"android"`
	LastConnected  *time.Time `json:"last_connected,omitempty" example:"2025-11-05T18:30:00Z"`
	ConnectionTime string     `json:"connection_time,omitempty" example:"2h 30m 15s"` // Duração formatada
	NeedsPairing   bool       `json:"needs_pairing" example:"false"`
	CanConnect     bool       `json:"can_connect" example:"true"`
}

type PairQRResponse struct {
	SessionID string    `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	QRCode    string    `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."` // Base64 data URL
	ExpiresAt time.Time `json:"expires_at" example:"2025-11-05T18:35:00Z"`
	Message   string    `json:"message" example:"Scan the QR code with your WhatsApp"`
}

type PairPhoneResponse struct {
	SessionID   string `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PhoneNumber string `json:"phone_number" example:"+5511999999999"`
	PairingCode string `json:"pairing_code" example:"ABCD-1234"`
	Message     string `json:"message" example:"Enter the pairing code on your phone"`
}

type ErrorResponse struct {
	Error   string                 `json:"error" example:"invalid_request"`
	Message string                 `json:"message" example:"Nome da sessão é obrigatório"`
	Details map[string]interface{} `json:"details,omitempty" swaggertype:"object"`
}

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operação realizada com sucesso"`
	Data    interface{} `json:"data,omitempty" swaggertype:"object"`
}

// WebhookConfigResponse - Resposta com configuração de webhook
type WebhookConfigResponse struct {
	SessionID string   `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Enabled   bool     `json:"enabled" example:"true"`
	URL       string   `json:"url" example:"https://hooks.exemplo.com/whatsapp"`
	Events    []string `json:"events" example:"message,status,qr,connected,disconnected"`
	Token     string   `json:"token,omitempty" example:"Bearer secret-token-123"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-11-06T10:30:00Z"`
}

