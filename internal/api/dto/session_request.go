package dto

// CreateSessionRequest - Criar nova sessão
type CreateSessionRequest struct {
	Name          string                 `json:"name" binding:"required,min=3,max=100"`
	WebhookURL    string                 `json:"webhook_url" binding:"omitempty,url"`
	WebhookEvents []string               `json:"webhook_events" binding:"omitempty"`
	Metadata      map[string]interface{} `json:"metadata" binding:"omitempty"`
}

// PairPhoneRequest - Parear com número de telefone
type PairPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"` // Formato: +5511999999999
}

// UpdateWebhookRequest - Atualizar webhook
type UpdateWebhookRequest struct {
	WebhookURL    string   `json:"webhook_url" binding:"required,url"`
	WebhookEvents []string `json:"webhook_events" binding:"required,min=1"`
	WebhookSecret string   `json:"webhook_secret" binding:"omitempty,min=16"`
}

// ConnectSessionRequest - Conectar sessão (opcional, pode ser vazio)
type ConnectSessionRequest struct {
	AutoReconnect bool `json:"auto_reconnect" binding:"omitempty"`
}

