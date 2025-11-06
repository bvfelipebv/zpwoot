# 📝 Sistema de Logging - zpwoot

Sistema de logging profissional baseado em **zerolog** com suporte a logs estruturados, contextuais e de alta performance.

## 🎯 Características

- ✅ **Logs estruturados** - Campos consistentes e pesquisáveis
- ✅ **Contexto automático** - Session ID, Worker ID, Request ID
- ✅ **Alta performance** - Zerolog é um dos loggers mais rápidos do Go
- ✅ **Múltiplos formatos** - Console (desenvolvimento) e JSON (produção)
- ✅ **Sampling** - Reduz volume de logs em alta carga
- ✅ **Níveis dinâmicos** - Altere o nível de log em runtime

## 🚀 Uso Básico

### Inicialização

```go
import "zpwoot/pkg/logger"

// Inicialização simples
logger.Init("info")

// Inicialização com configuração customizada
cfg := logger.Config{
    Level:       "debug",
    Format:      "json",        // "console" ou "json"
    AddCaller:   true,          // Adiciona arquivo:linha nos logs
    SampleRate:  10,            // Log 1 a cada 10 mensagens (0 = sem sampling)
    Environment: "production",
    Service:     "zpwoot",
}
logger.InitWithConfig(cfg)
```

### Logs Simples

```go
// Níveis de log
logger.Log.Trace().Msg("Mensagem de trace")
logger.Log.Debug().Msg("Mensagem de debug")
logger.Log.Info().Msg("Mensagem informativa")
logger.Log.Warn().Msg("Aviso")
logger.Log.Error().Err(err).Msg("Erro")
logger.Log.Fatal().Msg("Erro fatal - encerra aplicação")
```

### Logs com Campos Estruturados

```go
// Use as constantes de campos para consistência
logger.Log.Info().
    Str(logger.FieldSessionID, sessionID).
    Str(logger.FieldPhone, phone).
    Str(logger.FieldMessageID, messageID).
    Msg("Mensagem enviada")

// Campos comuns
logger.Log.Info().
    Str(logger.FieldURL, url).
    Int(logger.FieldStatus, 200).
    Dur(logger.FieldDuration, duration).
    Bool("success", true).
    Msg("Requisição completada")
```

## 🎨 Loggers Contextuais

### Logger por Componente

```go
// Cria logger com contexto de componente
webhookLog := logger.WithComponent("webhook-processor")
webhookLog.Info().Msg("Processando webhook")
```

### Logger por Sessão

```go
// Cria logger com contexto de sessão
sessionLog := logger.WithSession(sessionID)
sessionLog.Info().Msg("Sessão conectada")
```

### Logger por Worker

```go
// Cria logger com contexto de worker
workerLog := logger.WithWorker(workerID)
workerLog.Info().Msg("Worker iniciado")
```

### Logger com Contexto HTTP

```go
// Extrai contexto da requisição
ctx := c.Request.Context()
reqLog := logger.WithContext(ctx)
reqLog.Info().Msg("Processando requisição")
```

## 📊 Helpers de Campos

Use os helpers para padrões comuns de logging:

```go
// Logs de webhook
logger.WebhookFields(sessionID, "message", url, attempt).
    Msg("Webhook enviado")

// Logs de mensagem
logger.MessageFields(sessionID, phone, messageID).
    Msg("Mensagem processada")

// Logs HTTP
logger.HTTPFields(method, path, ip, userAgent, status, duration).
    Msg("Requisição HTTP")

// Logs de performance
logger.PerformanceFields("send_message", duration, true).
    Msg("Operação completada")
```

## 🔧 Constantes de Campos

Use sempre as constantes para nomes de campos:

```go
logger.FieldSessionID   // "session_id"
logger.FieldWorkerID    // "worker_id"
logger.FieldEvent       // "event"
logger.FieldURL         // "url"
logger.FieldPhone       // "phone"
logger.FieldMessageID   // "message_id"
logger.FieldAttempt     // "attempt"
logger.FieldStatus      // "status"
logger.FieldDuration    // "duration"
logger.FieldRequestID   // "request_id"
logger.FieldMethod      // "method"
logger.FieldPath        // "path"
logger.FieldIP          // "ip"
logger.FieldUserAgent   // "user_agent"
logger.FieldSubject     // "subject"
logger.FieldQueue       // "queue"
logger.FieldComponent   // "component"
```

## 🎯 Boas Práticas

### ✅ FAÇA

```go
// Use campos estruturados
logger.Log.Info().
    Str(logger.FieldSessionID, sessionID).
    Str(logger.FieldPhone, phone).
    Msg("Mensagem enviada")

// Use loggers contextuais
sessionLog := logger.WithSession(sessionID)
sessionLog.Info().Msg("Evento processado")

// Use emojis para facilitar visualização
logger.Log.Info().Msg("✅ Operação bem-sucedida")
logger.Log.Warn().Msg("⚠️ Tentando novamente")
logger.Log.Error().Msg("❌ Falha na operação")
```

### ❌ NÃO FAÇA

```go
// NÃO use fmt.Sprintf na mensagem
logger.Log.Info().Msgf("Sessão %s conectada", sessionID) // ❌

// FAÇA assim:
logger.Log.Info().
    Str(logger.FieldSessionID, sessionID).
    Msg("Sessão conectada") // ✅

// NÃO use nomes de campos inconsistentes
logger.Log.Info().Str("session", sessionID) // ❌
logger.Log.Info().Str("sid", sessionID)     // ❌

// FAÇA assim:
logger.Log.Info().Str(logger.FieldSessionID, sessionID) // ✅
```

## 🔍 Exemplos Práticos

### Webhook Worker

```go
type WebhookWorker struct {
    id  int
    log zerolog.Logger
}

func NewWebhookWorker(id int) *WebhookWorker {
    return &WebhookWorker{
        id:  id,
        log: logger.WithWorker(id),
    }
}

func (w *WebhookWorker) Process(msg WebhookMessage) {
    // Logger com contexto de sessão
    sessionLog := w.log.With().
        Str(logger.FieldSessionID, msg.SessionID).
        Str(logger.FieldEvent, msg.Event).
        Logger()

    sessionLog.Info().Msg("Processando webhook")
    
    // ... processamento ...
    
    sessionLog.Info().
        Dur(logger.FieldDuration, duration).
        Msg("✅ Webhook entregue")
}
```

## 📈 Performance

- **Sampling**: Use para reduzir volume em alta carga
- **Níveis**: Use DEBUG apenas em desenvolvimento
- **Campos**: Adicione apenas campos relevantes
- **Contexto**: Reutilize loggers contextuais

## 🔄 Mudança Dinâmica de Nível

```go
// Alterar nível em runtime
logger.SetLevel("debug")

// Obter nível atual
currentLevel := logger.GetLevel()
```

