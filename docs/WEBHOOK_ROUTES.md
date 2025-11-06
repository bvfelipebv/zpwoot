# 🔗 Rotas de Webhook - zpmeow

## 📋 Visão Geral

As rotas de webhook permitem configurar, consultar e gerenciar webhooks para receber eventos do WhatsApp em tempo real.

---

## 🆕 Novas Rotas (Recomendadas)

### 1. **Configurar Webhook** - `/sessions/:id/webhook/set`

**Método:** `POST`

**Descrição:** Configura ou atualiza o webhook de uma sessão específica com seus eventos subscritos.

**Request Body:**
```json
{
  "enabled": true,
  "url": "https://hooks.exemplo.com/whatsapp",
  "events": ["message", "status", "qr", "connected", "disconnected"],
  "token": "Bearer secret-token-123"
}
```

**Campos:**
- `enabled` (boolean, obrigatório): Habilita ou desabilita o webhook
- `url` (string, obrigatório se enabled=true): URL do webhook
- `events` (array, opcional): Lista de eventos a serem enviados. Se vazio, usa eventos padrão
- `token` (string, opcional): Token de autenticação para o webhook

**Eventos Disponíveis:**
- `message` - Mensagens recebidas
- `status` - Atualizações de status de mensagens
- `qr` - QR Code gerado
- `connected` - Sessão conectada
- `disconnected` - Sessão desconectada
- `presence` - Atualizações de presença
- `receipt` - Recibos de leitura
- `history_sync` - Sincronização de histórico

**Response (200 OK):**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "enabled": true,
  "url": "https://hooks.exemplo.com/whatsapp",
  "events": ["message", "status", "qr", "connected", "disconnected"],
  "token": "Bearer secret-token-123",
  "updated_at": "2025-11-06T10:30:00Z"
}
```

**Exemplo cURL:**
```bash
curl -X POST "http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/webhook/set" \
  -H "apikey: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "url": "https://hooks.exemplo.com/whatsapp",
    "events": ["message", "connected"],
    "token": "Bearer secret-token-123"
  }'
```

---

### 2. **Obter Configuração de Webhook** - `/sessions/:id/webhook/find`

**Método:** `GET`

**Descrição:** Retorna a configuração atual de webhook e seus eventos subscritos.

**Response (200 OK):**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "enabled": true,
  "url": "https://hooks.exemplo.com/whatsapp",
  "events": ["message", "status", "qr", "connected", "disconnected"],
  "token": "Bearer secret-token-123",
  "updated_at": "2025-11-06T10:30:00Z"
}
```

**Exemplo cURL:**
```bash
curl -X GET "http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/webhook/find" \
  -H "apikey: your-api-key"
```

---

### 3. **Limpar Webhook** - `/sessions/:id/webhook/clear`

**Método:** `DELETE`

**Descrição:** Remove/desabilita a configuração de webhook de uma sessão.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Webhook configuration cleared successfully"
}
```

**Exemplo cURL:**
```bash
curl -X DELETE "http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/webhook/clear" \
  -H "apikey: your-api-key"
```

---

## 🔄 Rota Antiga (Deprecated)

### **Atualizar Webhook** - `/sessions/:id/webhook` (PUT)

**Status:** ⚠️ DEPRECATED - Use `/webhook/set` ao invés

**Método:** `PUT`

**Descrição:** Atualiza webhook (mantida para compatibilidade)

---

## 📊 Comparação: Antes vs Depois

| Aspecto | Rota Antiga | Novas Rotas |
|---------|-------------|-------------|
| **Configurar** | `PUT /sessions/:id/webhook` | `POST /sessions/:id/webhook/set` |
| **Consultar** | ❌ Não existe | `GET /sessions/:id/webhook/find` |
| **Limpar** | ❌ Não existe | `DELETE /sessions/:id/webhook/clear` |
| **Semântica** | ⚠️ Confusa | ✅ Clara e explícita |
| **RESTful** | ⚠️ Parcial | ✅ Completo |

---

## 🎯 Casos de Uso

### **Caso 1: Configurar webhook pela primeira vez**
```bash
POST /sessions/{id}/webhook/set
{
  "enabled": true,
  "url": "https://meu-servidor.com/webhook",
  "events": ["message", "connected"]
}
```

### **Caso 2: Atualizar eventos subscritos**
```bash
POST /sessions/{id}/webhook/set
{
  "enabled": true,
  "url": "https://meu-servidor.com/webhook",
  "events": ["message", "status", "receipt"]
}
```

### **Caso 3: Desabilitar temporariamente**
```bash
POST /sessions/{id}/webhook/set
{
  "enabled": false,
  "url": "https://meu-servidor.com/webhook"
}
```

### **Caso 4: Verificar configuração atual**
```bash
GET /sessions/{id}/webhook/find
```

### **Caso 5: Remover webhook completamente**
```bash
DELETE /sessions/{id}/webhook/clear
```

---

## ✅ Validações Implementadas

1. ✅ **URL obrigatória** quando `enabled=true`
2. ✅ **Eventos válidos** - Apenas eventos da lista permitida
3. ✅ **Sessão existe** - Verifica se a sessão existe antes de configurar
4. ✅ **Eventos padrão** - Se não fornecidos, usa lista padrão
5. ✅ **Token opcional** - Pode ser usado para autenticação no webhook

---

## 🔐 Autenticação

Todas as rotas requerem header de autenticação:
```
apikey: your-secret-api-key-here
```

---

## 📝 Notas Importantes

- A rota antiga `PUT /sessions/:id/webhook` ainda funciona mas está marcada como deprecated
- Use as novas rotas para novos desenvolvimentos
- O token do webhook é opcional e pode ser usado para validar requisições no seu servidor
- Se eventos não forem especificados, o sistema usa: `["message", "status", "qr", "connected", "disconnected"]`

