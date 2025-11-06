# 🔄 Refatoração do Sistema de Logging - zpwoot

## 📋 Resumo

Refatoração completa do sistema de logging para usar **zerolog** de forma profissional e consistente em todo o projeto.

## ✅ O que foi feito

### 1. **pkg/logger/logger.go** - Sistema de Logging Aprimorado

#### Novos Recursos:
- ✅ **Configuração avançada** via `Config` struct
- ✅ **Múltiplos formatos**: Console (desenvolvimento) e JSON (produção)
- ✅ **Log sampling** para alta performance em produção
- ✅ **Caller information** opcional (arquivo:linha)
- ✅ **Níveis dinâmicos** - altere em runtime com `SetLevel()`
- ✅ **Loggers contextuais**:
  - `WithContext(ctx)` - extrai request_id e session_id do contexto
  - `WithComponent(name)` - logger por componente
  - `WithSession(sessionID)` - logger por sessão
  - `WithWorker(workerID)` - logger por worker
  - `WithFields(map)` - logger com campos customizados

#### Constantes de Campos:
```go
const (
    FieldSessionID   = "session_id"
    FieldWorkerID    = "worker_id"
    FieldEvent       = "event"
    FieldURL         = "url"
    FieldPhone       = "phone"
    FieldMessageID   = "message_id"
    FieldAttempt     = "attempt"
    FieldStatus      = "status"
    FieldDuration    = "duration"
    FieldError       = "error"
    FieldRequestID   = "request_id"
    FieldMethod      = "method"
    FieldPath        = "path"
    FieldIP          = "ip"
    FieldUserAgent   = "user_agent"
    FieldSubject     = "subject"
    FieldQueue       = "queue"
    FieldComponent   = "component"
    FieldEnvironment = "environment"
    FieldLogLevel    = "log_level"
)
```

### 2. **pkg/logger/fields.go** - Helpers de Logging

Funções helper para padrões comuns:
- `SessionFields()` - logs de sessão
- `WebhookFields()` - logs de webhook
- `MessageFields()` - logs de mensagem
- `HTTPFields()` - logs HTTP
- `WorkerFields()` - logs de worker
- `ErrorFields()` - logs de erro com contexto
- `NATSFields()` - logs NATS
- `PerformanceFields()` - métricas de performance

### 3. **internal/api/middleware/auth.go** - Middleware HTTP Melhorado

#### RequestLogger():
- ✅ Gera **request_id único** (UUID) para cada requisição
- ✅ Mede **latência** da requisição
- ✅ Log com **nível apropriado** baseado no status:
  - 2xx → Info
  - 4xx → Warn
  - 5xx → Error
- ✅ Campos estruturados: method, path, ip, user_agent, status, duration
- ✅ Emojis para facilitar visualização: `→ Incoming` / `← Completed`

#### RequestLoggerWithSkip():
- ✅ Permite pular logging de certos paths (ex: `/health`)

#### AuthenticateGlobal():
- ✅ Usa constantes de campos para consistência

### 4. **internal/service/webhook_worker.go** - Worker com Logger Contextual

#### Melhorias:
- ✅ Cada worker tem seu próprio **logger contextual** com `worker_id`
- ✅ Cada mensagem processada cria **logger de sessão** com contexto completo
- ✅ Usa **constantes de campos** em todos os logs
- ✅ Emojis para status visual:
  - `✅` - Sucesso
  - `⚠️` - Retry
  - `❌` - Falha permanente
- ✅ Logs estruturados com campos relevantes

#### Antes:
```go
logger.Log.Info().
    Int("worker_id", w.id).
    Str("session_id", sessionID).
    Msg("Processing webhook")
```

#### Depois:
```go
sessionLog := w.log.With().
    Str(logger.FieldSessionID, sessionID).
    Str(logger.FieldEvent, event).
    Logger()

sessionLog.Info().Msg("Processing webhook")
```

### 5. **pkg/logger/README.md** - Documentação Completa

Documentação abrangente com:
- 📖 Guia de uso básico
- 🎨 Exemplos de loggers contextuais
- 📊 Helpers de campos
- 🎯 Boas práticas (DO's e DON'Ts)
- 🔍 Exemplos práticos
- 📈 Dicas de performance

## 🎯 Benefícios

### Performance
- ⚡ **Zerolog** é um dos loggers mais rápidos do Go
- ⚡ **Sampling** reduz volume em alta carga
- ⚡ **Logs estruturados** são mais eficientes que concatenação de strings

### Observabilidade
- 🔍 **Campos consistentes** facilitam busca e análise
- 🔍 **Request ID** permite rastrear requisições end-to-end
- 🔍 **Contexto automático** (session_id, worker_id) facilita debugging
- 🔍 **Formato JSON** em produção permite integração com ferramentas (ELK, Datadog, etc)

### Manutenibilidade
- 📝 **Constantes de campos** evitam typos
- 📝 **Loggers contextuais** reduzem repetição de código
- 📝 **Helpers** padronizam logging em todo o projeto
- 📝 **Documentação** facilita onboarding de novos desenvolvedores

## 🚀 Próximos Passos Sugeridos

### 1. Atualizar Componentes Restantes
Aplicar o mesmo padrão em:
- [ ] `internal/service/webhook_processor.go`
- [ ] `internal/service/webhook_delivery.go`
- [ ] `internal/service/event_handler.go`
- [ ] `internal/service/message_service.go`
- [ ] `internal/nats/client.go`
- [ ] `internal/api/handlers/*.go`

### 2. Configuração por Ambiente
```go
// .env
LOG_LEVEL=info
LOG_FORMAT=json        # console em dev, json em prod
LOG_CALLER=false       # true em dev, false em prod
LOG_SAMPLE_RATE=0      # 0 em dev, 10 em prod (1 a cada 10)
```

### 3. Integração com Ferramentas
- **ELK Stack**: Logs JSON → Elasticsearch → Kibana
- **Datadog**: Enviar logs estruturados
- **Grafana Loki**: Agregação e visualização

### 4. Métricas e Alertas
- Alertar em logs de ERROR
- Dashboard de latência (usando `duration`)
- Monitorar taxa de retry de webhooks

## 📊 Exemplo de Uso Completo

```go
// Inicialização (main.go)
cfg := logger.Config{
    Level:       config.AppConfig.LogLevel,
    Format:      config.AppConfig.LogFormat,
    AddCaller:   config.AppConfig.Environment == "development",
    SampleRate:  config.AppConfig.LogSampleRate,
    Environment: config.AppConfig.Environment,
    Service:     "zpwoot",
}
logger.InitWithConfig(cfg)

// Worker (webhook_worker.go)
type WebhookWorker struct {
    log zerolog.Logger
}

func NewWebhookWorker(id int) *WebhookWorker {
    return &WebhookWorker{
        log: logger.WithWorker(id),
    }
}

func (w *WebhookWorker) Process(msg WebhookMessage) {
    sessionLog := w.log.With().
        Str(logger.FieldSessionID, msg.SessionID).
        Str(logger.FieldEvent, msg.Event).
        Logger()

    sessionLog.Info().Msg("Processing webhook")
    
    result := w.deliver(msg)
    
    if result.Success {
        sessionLog.Info().
            Dur(logger.FieldDuration, result.Duration).
            Msg("✅ Webhook delivered")
    } else {
        sessionLog.Error().
            Err(result.Error).
            Msg("❌ Webhook failed")
    }
}
```

## ✅ Build Status

✅ **Build bem-sucedido** - Código compila sem erros
✅ **Todas as tarefas concluídas**
✅ **Documentação completa**
✅ **Pronto para uso**

