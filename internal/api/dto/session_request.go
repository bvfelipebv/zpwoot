package dto

// ProxyConfig - Configuração de proxy
type ProxyConfig struct {
	Enabled  bool   `json:"enabled" example:"true"`
	Protocol string `json:"protocol" binding:"omitempty,oneof=http https socks5" example:"http"`
	Host     string `json:"host" binding:"required_if=Enabled true" example:"10.0.0.1"`
	Port     int    `json:"port" binding:"required_if=Enabled true,min=1,max=65535" example:"3128"`
	Username string `json:"username,omitempty" example:"proxyuser"`
	Password string `json:"password,omitempty" example:"proxypass"`
}

// WebhookConfig - Configuração de webhook
type WebhookConfig struct {
	Enabled bool     `json:"enabled" example:"true"`
	URL     string   `json:"url" binding:"required_if=Enabled true,omitempty,url" example:"https://hooks.exemplo.com/wuz"`
	Events  []string `json:"events" binding:"omitempty" example:"message,status,qr"`
	Token   string   `json:"token,omitempty" example:"secreto-opcional"`
}

// CreateSessionRequest - Criar nova sessão
type CreateSessionRequest struct {
	Name    string         `json:"name" binding:"required,min=3,max=100" example:"sessao-atendimento-1"`
	APIKey  *string        `json:"apikey" example:"null"`
	Proxy   *ProxyConfig   `json:"proxy,omitempty"`
	Webhook *WebhookConfig `json:"webhook,omitempty"`
}

// PairPhoneRequest - Parear com número de telefone
type PairPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"+5511999999999"` // Formato: +5511999999999
}

// UpdateWebhookRequest - Atualizar webhook
type UpdateWebhookRequest struct {
	Webhook WebhookConfig `json:"webhook" binding:"required"`
}

// ConnectSessionRequest - Conectar sessão (opcional, pode ser vazio)
type ConnectSessionRequest struct {
	AutoReconnect bool `json:"auto_reconnect" binding:"omitempty" example:"true"`
}

