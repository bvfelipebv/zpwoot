# 📋 Plano de Implementação - Serviço de Entrega de Webhooks

## 🎯 Objetivo

Implementar um sistema completo de entrega de webhooks que:
1. **Recebe** eventos do WhatsApp via whatsmeow
2. **Processa** e formata os eventos
3. **Entrega** via HTTP POST para URLs configuradas pelos clientes
4. **Garante** confiabilidade com retry e fila

---

## 📊 Análise da Estrutura Atual

### **Estrutura de Pastas Existente**
```
internal/
├── api/              # Handlers HTTP e rotas
│   ├── dto/          # Data Transfer Objects
│   ├── handlers/     # HTTP handlers
│   └── middleware/   # Middlewares
├── config/           # Configurações
├── constants/        # Constantes (eventos de webhook)
├── db/               # Database e migrations
├── model/            # Modelos de dados
├── repository/       # Camada de acesso a dados
└── service/          # Lógica de negócio
    ├── event_handler.go      # ✅ JÁ EXISTE - Handler de eventos
    ├── session_manager.go    # ✅ JÁ EXISTE - Gerenciador de sessões
    ├── whatsapp_service.go   # ✅ JÁ EXISTE - Serviço WhatsApp
    └── message_service.go    # ✅ JÁ EXISTE - Serviço de mensagens
```

### **Pontos de Integração Identificados**

1. **`event_handler.go`** - Já recebe eventos do whatsmeow
   - ✅ Tem TODOs para enviar webhooks
   - ✅ Já processa eventos: Connected, Disconnected, LoggedOut, Message, Receipt, etc.

2. **`model.Session`** - Já tem `WebhookConfig`
   - ✅ Enabled, URL, Events, Token

3. **`constants/webhook_events.go`** - Já tem todos os eventos mapeados
   - ✅ 60+ eventos do whatsmeow
   - ✅ Funções de validação

---

## 🏗️ Arquitetura Proposta

### **Camadas do Sistema**

```
┌─────────────────────────────────────────────────────────────┐
│                    WhatsApp (whatsmeow)                      │
└────────────────────────┬────────────────────────────────────┘
                         │ eventos
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              EventHandler (event_handler.go)                 │
│  - Recebe eventos do whatsmeow                              │
│  - Filtra por sessão                                        │
└────────────────────────┬────────────────────────────────────┘
                         │ eventos filtrados
                         ▼
┌─────────────────────────────────────────────────────────────┐
│          WebhookProcessor (webhook_processor.go)             │
│  - Verifica se sessão tem webhook configurado               │
│  - Valida se evento está subscrito                          │
│  - Formata payload do webhook                               │
└────────────────────────┬────────────────────────────────────┘
                         │ webhook payload
                         ▼
┌─────────────────────────────────────────────────────────────┐
│          WebhookQueue (webhook_queue.go)                     │
│  - Enfileira webhooks para entrega                          │
│  - Gerencia prioridades                                     │
│  - Persiste em memória/Redis (futuro)                       │
└────────────────────────┬────────────────────────────────────┘
                         │ webhooks enfileirados
                         ▼
┌─────────────────────────────────────────────────────────────┐
│          WebhookDelivery (webhook_delivery.go)               │
│  - Envia HTTP POST para URL configurada                     │
│  - Implementa retry com backoff exponencial                 │
│  - Registra logs de entrega                                 │
│  - Trata timeouts e erros                                   │
└────────────────────────┬────────────────────────────────────┘
                         │ resultado
                         ▼
┌─────────────────────────────────────────────────────────────┐
│          WebhookLog (webhook_log.go)                         │
│  - Registra tentativas de entrega                           │
│  - Armazena status (sucesso/falha)                          │
│  - Permite consulta de histórico                            │
└─────────────────────────────────────────────────────────────┘
```

---

## 📁 Estrutura de Arquivos a Criar

### **1. Services** (`internal/service/`)

```
internal/service/
├── webhook_processor.go    # Processa e formata eventos para webhook
├── webhook_queue.go         # Fila de webhooks (em memória)
├── webhook_delivery.go      # Entrega HTTP dos webhooks
└── webhook_formatter.go     # Formata payloads por tipo de evento
```

### **2. Models** (`internal/model/`)

```
internal/model/
├── webhook_payload.go       # Estrutura do payload de webhook
└── webhook_log.go           # Modelo de log de webhook
```

### **3. Repository** (`internal/repository/`)

```
internal/repository/
└── webhook_log_repo.go      # Persistência de logs de webhook
```

### **4. Database** (`internal/db/migrations/`)

```
internal/db/migrations/
└── 003_create_webhook_logs.sql  # Tabela de logs de webhook
```

### **5. DTOs** (`internal/api/dto/`)

```
internal/api/dto/
└── webhook_dto.go           # DTOs para consulta de logs
```

### **6. Handlers** (`internal/api/handlers/`)

```
internal/api/handlers/
└── webhook_log_handler.go   # Endpoints para consultar logs
```

---

## 🗄️ Schema de Banco de Dados

### **Tabela: webhook_logs**

```sql
CREATE TABLE webhook_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    webhook_url TEXT NOT NULL,
    
    -- Payload e resposta
    payload JSONB NOT NULL,
    response_status INT,
    response_body TEXT,
    
    -- Controle de tentativas
    attempt INT DEFAULT 1,
    max_attempts INT DEFAULT 3,
    
    -- Status
    status VARCHAR(50) NOT NULL, -- pending, success, failed, retrying
    error_message TEXT,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP,
    completed_at TIMESTAMP,
    next_retry_at TIMESTAMP,
    
    -- Índices
    CONSTRAINT fk_session FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE INDEX idx_webhook_logs_session ON webhook_logs(session_id);
CREATE INDEX idx_webhook_logs_status ON webhook_logs(status);
CREATE INDEX idx_webhook_logs_created ON webhook_logs(created_at DESC);
CREATE INDEX idx_webhook_logs_event_type ON webhook_logs(event_type);
```

---

## 📦 Modelos de Dados

### **WebhookPayload** (Estrutura Padrão)

```go
type WebhookPayload struct {
    Event     string                 `json:"event"`      // Tipo do evento
    SessionID string                 `json:"session_id"` // ID da sessão
    Timestamp time.Time              `json:"timestamp"`  // Quando ocorreu
    Data      map[string]interface{} `json:"data"`       // Dados do evento
}
```

### **WebhookLog**

```go
type WebhookLog struct {
    ID            string
    SessionID     string
    EventType     string
    WebhookURL    string
    Payload       WebhookPayload
    ResponseStatus int
    ResponseBody  string
    Attempt       int
    MaxAttempts   int
    Status        string // pending, success, failed, retrying
    ErrorMessage  string
    CreatedAt     time.Time
    SentAt        *time.Time
    CompletedAt   *time.Time
    NextRetryAt   *time.Time
}
```

---

## 🔄 Fluxo de Processamento

### **1. Recebimento de Evento**

```go
// event_handler.go
func (h *EventHandler) handleMessage(sessionID string, evt *events.Message) {
    // 1. Log do evento
    logger.Log.Debug().Msg("Message received")
    
    // 2. Enviar para processador de webhook
    h.webhookProcessor.ProcessEvent(sessionID, constants.EventMessage, evt)
}
```

### **2. Processamento**

```go
// webhook_processor.go
func (p *WebhookProcessor) ProcessEvent(sessionID string, eventType constants.WebhookEventType, data interface{}) {
    // 1. Buscar configuração de webhook da sessão
    session := p.getSession(sessionID)
    if !session.WebhookConfig.Enabled {
        return // Webhook desabilitado
    }
    
    // 2. Verificar se evento está subscrito
    if !p.isEventSubscribed(session.WebhookConfig.Events, eventType) {
        return // Evento não subscrito
    }
    
    // 3. Formatar payload
    payload := p.formatter.Format(eventType, sessionID, data)
    
    // 4. Enfileirar para entrega
    p.queue.Enqueue(session.WebhookConfig.URL, payload, session.WebhookConfig.Token)
}
```

### **3. Enfileiramento**

```go
// webhook_queue.go
func (q *WebhookQueue) Enqueue(url string, payload WebhookPayload, token string) {
    item := &QueueItem{
        URL:     url,
        Payload: payload,
        Token:   token,
        Attempt: 0,
    }
    q.queue <- item
}
```

### **4. Entrega**

```go
// webhook_delivery.go
func (d *WebhookDelivery) Send(item *QueueItem) error {
    // 1. Preparar request
    body, _ := json.Marshal(item.Payload)
    req, _ := http.NewRequest("POST", item.URL, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    if item.Token != "" {
        req.Header.Set("Authorization", item.Token)
    }
    
    // 2. Enviar com timeout
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    
    // 3. Registrar log
    d.logDelivery(item, resp, err)
    
    // 4. Retry se necessário
    if err != nil && item.Attempt < 3 {
        d.scheduleRetry(item)
    }
    
    return err
}
```

---

## ⚙️ Configurações

### **Adicionar ao `config.go`**

```go
type WebhookConfig struct {
    MaxRetries      int           `env:"WEBHOOK_MAX_RETRIES" envDefault:"3"`
    RetryDelay      time.Duration `env:"WEBHOOK_RETRY_DELAY" envDefault:"5s"`
    Timeout         time.Duration `env:"WEBHOOK_TIMEOUT" envDefault:"30s"`
    QueueSize       int           `env:"WEBHOOK_QUEUE_SIZE" envDefault:"1000"`
    Workers         int           `env:"WEBHOOK_WORKERS" envDefault:"10"`
    EnableLogs      bool          `env:"WEBHOOK_ENABLE_LOGS" envDefault:"true"`
    LogRetention    time.Duration `env:"WEBHOOK_LOG_RETENTION" envDefault:"168h"` // 7 dias
}
```

---

## 🎨 Formatação de Payloads

### **Payload Padrão**

```json
{
  "event": "message",
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2025-11-06T10:30:00Z",
  "data": {
    // Dados específicos do evento
  }
}
```

### **Exemplo: Evento Message**

```json
{
  "event": "message",
  "session_id": "abc-123",
  "timestamp": "2025-11-06T10:30:00Z",
  "data": {
    "message_id": "3EB0123456789ABCDEF",
    "from": "5511999999999@s.whatsapp.net",
    "from_me": false,
    "chat": "5511999999999@s.whatsapp.net",
    "timestamp": "2025-11-06T10:30:00Z",
    "type": "conversation",
    "body": "Olá, tudo bem?",
    "media_type": null
  }
}
```

### **Exemplo: Evento Connected**

```json
{
  "event": "connected",
  "session_id": "abc-123",
  "timestamp": "2025-11-06T10:30:00Z",
  "data": {
    "status": "connected"
  }
}
```

---

## 🔐 Segurança

### **1. Assinatura de Webhooks (HMAC)**

```go
func SignPayload(payload []byte, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write(payload)
    return hex.EncodeToString(h.Sum(nil))
}

// Header: X-Webhook-Signature: sha256=abc123...
```

### **2. Validação de URL**

```go
func ValidateWebhookURL(url string) error {
    // Apenas HTTPS em produção
    // Não permitir localhost/127.0.0.1
    // Validar formato de URL
}
```

---

## 📊 Monitoramento e Métricas

### **Métricas a Coletar**

```go
// Prometheus metrics
var (
    webhooksTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "webhooks_total",
            Help: "Total de webhooks enviados",
        },
        []string{"event_type", "status"},
    )
    
    webhookDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "webhook_duration_seconds",
            Help: "Duração do envio de webhook",
        },
        []string{"event_type"},
    )
)
```

---

## 🧪 Testes

### **Testes Unitários**

```
internal/service/
├── webhook_processor_test.go
├── webhook_queue_test.go
├── webhook_delivery_test.go
└── webhook_formatter_test.go
```

### **Testes de Integração**

```
tests/integration/
└── webhook_delivery_test.go
```

---

## 📝 Próximos Passos

Continua no próximo arquivo...

