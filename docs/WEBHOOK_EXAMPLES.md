# 🔔 Exemplos de Uso - Webhook Routes

## 📚 Exemplos Práticos

### **Exemplo 1: Configurar webhook para receber mensagens**

```bash
curl -X POST "http://localhost:8080/sessions/minha-sessao-123/webhook/set" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "url": "https://meu-servidor.com/webhook/whatsapp",
    "events": ["message"],
    "token": "Bearer meu-token-secreto"
  }'
```

**Resposta:**
```json
{
  "session_id": "minha-sessao-123",
  "enabled": true,
  "url": "https://meu-servidor.com/webhook/whatsapp",
  "events": ["message"],
  "token": "Bearer meu-token-secreto",
  "updated_at": "2025-11-06T10:30:00Z"
}
```

---

### **Exemplo 2: Configurar webhook com múltiplos eventos**

```bash
curl -X POST "http://localhost:8080/sessions/minha-sessao-123/webhook/set" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "url": "https://meu-servidor.com/webhook/whatsapp",
    "events": ["message", "status", "connected", "disconnected", "qr"]
  }'
```

---

### **Exemplo 3: Consultar configuração atual**

```bash
curl -X GET "http://localhost:8080/sessions/minha-sessao-123/webhook/find" \
  -H "apikey: minha-chave-secreta"
```

**Resposta:**
```json
{
  "session_id": "minha-sessao-123",
  "enabled": true,
  "url": "https://meu-servidor.com/webhook/whatsapp",
  "events": ["message", "status", "connected", "disconnected", "qr"],
  "token": "",
  "updated_at": "2025-11-06T10:30:00Z"
}
```

---

### **Exemplo 4: Atualizar apenas a URL do webhook**

```bash
curl -X POST "http://localhost:8080/sessions/minha-sessao-123/webhook/set" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "url": "https://novo-servidor.com/webhook",
    "events": ["message", "status"]
  }'
```

---

### **Exemplo 5: Desabilitar webhook temporariamente**

```bash
curl -X POST "http://localhost:8080/sessions/minha-sessao-123/webhook/set" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": false,
    "url": "https://meu-servidor.com/webhook/whatsapp"
  }'
```

---

### **Exemplo 6: Remover webhook completamente**

```bash
curl -X DELETE "http://localhost:8080/sessions/minha-sessao-123/webhook/clear" \
  -H "apikey: minha-chave-secreta"
```

**Resposta:**
```json
{
  "success": true,
  "message": "Webhook configuration cleared successfully"
}
```

---

## 🔄 Fluxo Completo de Uso

### **1. Criar sessão**
```bash
curl -X POST "http://localhost:8080/sessions/create" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Minha Sessão WhatsApp"
  }'
```

### **2. Configurar webhook**
```bash
curl -X POST "http://localhost:8080/sessions/{session_id}/webhook/set" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "url": "https://meu-servidor.com/webhook",
    "events": ["message", "connected"]
  }'
```

### **3. Conectar sessão**
```bash
curl -X POST "http://localhost:8080/sessions/{session_id}/connect" \
  -H "apikey: minha-chave-secreta"
```

### **4. Verificar webhook configurado**
```bash
curl -X GET "http://localhost:8080/sessions/{session_id}/webhook/find" \
  -H "apikey: minha-chave-secreta"
```

---

## 🎯 Casos de Uso Avançados

### **Webhook com autenticação**
```json
{
  "enabled": true,
  "url": "https://meu-servidor.com/webhook",
  "events": ["message"],
  "token": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### **Webhook apenas para eventos de conexão**
```json
{
  "enabled": true,
  "url": "https://meu-servidor.com/webhook/status",
  "events": ["connected", "disconnected", "qr"]
}
```

### **Webhook para todos os eventos**
```json
{
  "enabled": true,
  "url": "https://meu-servidor.com/webhook/all",
  "events": ["message", "status", "qr", "connected", "disconnected", "presence", "receipt", "history_sync"]
}
```

---

## ⚠️ Tratamento de Erros

### **Erro: URL obrigatória quando enabled=true**
```bash
curl -X POST "http://localhost:8080/sessions/minha-sessao-123/webhook/set" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "url": ""
  }'
```

**Resposta (400 Bad Request):**
```json
{
  "error": "invalid_request",
  "message": "URL is required when webhook is enabled"
}
```

### **Erro: Evento inválido**
```bash
curl -X POST "http://localhost:8080/sessions/minha-sessao-123/webhook/set" \
  -H "apikey: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "url": "https://meu-servidor.com/webhook",
    "events": ["invalid_event"]
  }'
```

**Resposta (400 Bad Request):**
```json
{
  "error": "invalid_event",
  "message": "Invalid event type: invalid_event"
}
```

