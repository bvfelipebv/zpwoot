# 📦 Package: constants

## 📋 Descrição

Este pacote contém todas as constantes de eventos de webhook suportados pelo zpmeow, baseadas na biblioteca oficial [whatsmeow](https://pkg.go.dev/go.mau.fi/whatsmeow/types/events).

## 🎯 Uso

### Importar o pacote

```go
import "zpwoot/internal/constants"
```

### Validar um evento

```go
if constants.IsValidEventType("message") {
    fmt.Println("Evento válido!")
}
```

### Obter eventos por categoria

```go
messageEvents := constants.GetEventsByCategory("messages")
for _, event := range messageEvents {
    fmt.Println(event)
}
```

### Validar lista de eventos

```go
events := []string{"message", "invalid_event", "connected"}
valid, invalid := constants.ValidateEventList(events)

fmt.Printf("Válidos: %v\n", valid)     // [message, connected]
fmt.Printf("Inválidos: %v\n", invalid) // [invalid_event]
```

### Verificar se evento é crítico

```go
if constants.IsCriticalEvent("logged_out") {
    fmt.Println("Evento crítico! Tomar ação imediata.")
}
```

### Obter descrição de um evento

```go
desc := constants.GetEventDescription("message")
fmt.Println(desc) // "Mensagem recebida (texto, mídia, documentos, etc)"
```

### Obter categoria de um evento

```go
category := constants.GetEventCategory("call_offer")
fmt.Println(category) // "calls"
```

## 📊 Constantes Disponíveis

### Listas de Eventos

- **`SupportedEventTypes`** - Lista plana de todos os eventos suportados
- **`DefaultWebhookEvents`** - Eventos padrão quando nenhum é especificado
- **`CriticalEvents`** - Eventos críticos de conexão
- **`RecommendedEvents`** - Eventos recomendados para maioria dos casos
- **`MessageEvents`** - Apenas eventos de mensagens
- **`ConnectionEvents`** - Apenas eventos de conexão

### Mapas

- **`AllWebhookEvents`** - Eventos organizados por categoria
- **`EventTypeMap`** - Mapa para validação rápida

## 🧪 Testes

Execute os testes com:

```bash
go test -v ./internal/constants/
```

Cobertura de testes:

```bash
go test -cover ./internal/constants/
```

## 📚 Categorias de Eventos

1. **messages** - Mensagens e comunicação (5 eventos)
2. **groups_contacts** - Grupos e contatos (8 eventos)
3. **connection** - Conexão e sessão (15 eventos)
4. **privacy** - Privacidade e configurações (4 eventos)
5. **sync** - Sincronização e estado (16 eventos)
6. **calls** - Chamadas de voz/vídeo (9 eventos)
7. **presence** - Presença e atividade (2 eventos)
8. **identity** - Identidade e segurança (2 eventos)
9. **newsletter** - Canais do WhatsApp (4 eventos)
10. **facebook** - Facebook/Instagram bridge (1 evento)
11. **special** - Eventos especiais (1 evento)

**Total:** 60+ eventos

## 🔍 Funções Auxiliares

### `IsValidEventType(eventType string) bool`
Verifica se um tipo de evento é válido.

### `IsCriticalEvent(eventType string) bool`
Verifica se um evento é crítico para a conexão.

### `IsMessageEvent(eventType string) bool`
Verifica se um evento é relacionado a mensagens.

### `IsConnectionEvent(eventType string) bool`
Verifica se um evento é relacionado a conexão.

### `GetEventsByCategory(category string) []WebhookEventType`
Retorna eventos de uma categoria específica.

### `GetAllCategories() []string`
Retorna todas as categorias disponíveis.

### `GetEventDescription(eventType string) string`
Retorna descrição amigável de um evento.

### `ValidateEventList(events []string) (valid []string, invalid []string)`
Valida uma lista de eventos e separa válidos de inválidos.

### `GetEventCategory(eventType string) string`
Retorna a categoria de um evento.

## 📖 Exemplos

### Exemplo 1: Validar configuração de webhook

```go
func ValidateWebhookConfig(events []string) error {
    valid, invalid := constants.ValidateEventList(events)
    
    if len(invalid) > 0 {
        return fmt.Errorf("eventos inválidos: %v", invalid)
    }
    
    // Verificar se tem pelo menos um evento crítico
    hasCritical := false
    for _, event := range valid {
        if constants.IsCriticalEvent(event) {
            hasCritical = true
            break
        }
    }
    
    if !hasCritical {
        log.Warn("Nenhum evento crítico configurado")
    }
    
    return nil
}
```

### Exemplo 2: Listar eventos por categoria

```go
func ListEventsByCategory() {
    categories := constants.GetAllCategories()
    
    for _, category := range categories {
        events := constants.GetEventsByCategory(category)
        fmt.Printf("\n%s (%d eventos):\n", category, len(events))
        
        for _, event := range events {
            desc := constants.GetEventDescription(string(event))
            fmt.Printf("  - %s: %s\n", event, desc)
        }
    }
}
```

### Exemplo 3: Usar eventos padrão

```go
func SetupDefaultWebhook(sessionID string, url string) error {
    config := &WebhookConfig{
        Enabled: true,
        URL:     url,
        Events:  constants.DefaultWebhookEvents,
    }
    
    return SaveWebhookConfig(sessionID, config)
}
```

## 🔗 Referências

- [Documentação whatsmeow](https://pkg.go.dev/go.mau.fi/whatsmeow/types/events)
- [Documentação completa de eventos](../../../docs/WEBHOOK_EVENTS.md)
- [Exemplos de webhook](../../../docs/WEBHOOK_EXAMPLES.md)

## 📝 Notas

- Todos os eventos são baseados na versão mais recente do whatsmeow
- Eventos críticos devem sempre ser monitorados
- Use `EventAll` com cuidado - pode gerar muito tráfego
- Alguns eventos requerem subscrição ou configuração adicional (ex: `presence`)

