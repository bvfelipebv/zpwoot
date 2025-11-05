package dto

// CreateSessionRequest - Criar nova sessão
type CreateSessionRequest struct {
	Name          string                 `json:"name" binding:"required,min=3,max=100" example:"Minha Sessão WhatsApp"`
	WebhookURL    string                 `json:"webhook_url" binding:"omitempty,url" example:"https://seu-webhook.com/whatsapp"`
	WebhookEvents []string               `json:"webhook_events" binding:"omitempty" example:"message,qr,connected,disconnected"`
	Metadata      map[string]interface{} `json:"metadata" binding:"omitempty" swaggertype:"object" example:"cliente:Empresa XYZ,ambiente:producao"`
}

// PairPhoneRequest - Parear com número de telefone
type PairPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"+5511999999999"` // Formato: +5511999999999
}

// UpdateWebhookRequest - Atualizar webhook
type UpdateWebhookRequest struct {
	WebhookURL    string   `json:"webhook_url" binding:"required,url" example:"https://novo-webhook.com/whatsapp"`
	WebhookEvents []string `json:"webhook_events" binding:"required,min=1" example:"message,qr,connected"`
	WebhookSecret string   `json:"webhook_secret" binding:"omitempty,min=16" example:"meu-secret-super-seguro-123"`
}

// ConnectSessionRequest - Conectar sessão (opcional, pode ser vazio)
type ConnectSessionRequest struct {
	AutoReconnect bool `json:"auto_reconnect" binding:"omitempty" example:"true"`
}

