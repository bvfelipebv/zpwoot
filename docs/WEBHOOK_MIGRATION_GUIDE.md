# 🔄 Guia de Migração - Rotas de Webhook

## 📋 Resumo das Mudanças

As rotas de webhook foram refatoradas para serem mais semânticas e RESTful.

---

## 🆚 Comparação: Antes vs Depois

### **ANTES (Rota Antiga)**

```bash
# Atualizar webhook
PUT /sessions/:id/webhook

# Consultar webhook
❌ NÃO EXISTIA

# Limpar webhook
❌ NÃO EXISTIA
```

### **DEPOIS (Novas Rotas)**

```bash
# Configurar/Atualizar webhook
POST /sessions/:id/webhook/set

# Consultar webhook
GET /sessions/:id/webhook/find

# Limpar webhook
DELETE /sessions/:id/webhook/clear
```

---

## 📊 Tabela Comparativa

| Ação | Rota Antiga | Nova Rota | Status |
|------|-------------|-----------|--------|
| **Configurar** | `PUT /sessions/:id/webhook` | `POST /sessions/:id/webhook/set` | ✅ Recomendado |
| **Consultar** | ❌ Não existe | `GET /sessions/:id/webhook/find` | ✅ Novo |
| **Limpar** | ❌ Não existe | `DELETE /sessions/:id/webhook/clear` | ✅ Novo |

---

## 🔧 Como Migrar

### **Cenário 1: Você usa PUT /sessions/:id/webhook**

**Código Antigo:**
```javascript
// JavaScript/Node.js
const response = await fetch(`${API_URL}/sessions/${sessionId}/webhook`, {
  method: 'PUT',
  headers: {
    'apikey': API_KEY,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    webhook: {
      enabled: true,
      url: 'https://meu-webhook.com',
      events: ['message']
    }
  })
});
```

**Código Novo (Recomendado):**
```javascript
// JavaScript/Node.js
const response = await fetch(`${API_URL}/sessions/${sessionId}/webhook/set`, {
  method: 'POST',
  headers: {
    'apikey': API_KEY,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    enabled: true,
    url: 'https://meu-webhook.com',
    events: ['message'],
    token: 'Bearer meu-token' // Opcional
  })
});
```

**Mudanças:**
1. ✅ Método: `PUT` → `POST`
2. ✅ Endpoint: `/webhook` → `/webhook/set`
3. ✅ Body: Remover wrapper `webhook`, enviar campos diretamente
4. ✅ Adicionar campo `token` (opcional)

---

### **Cenário 2: Você precisa consultar a configuração**

**Antes:**
```javascript
// ❌ Não era possível consultar
```

**Agora:**
```javascript
// ✅ Consultar configuração atual
const response = await fetch(`${API_URL}/sessions/${sessionId}/webhook/find`, {
  method: 'GET',
  headers: {
    'apikey': API_KEY
  }
});

const config = await response.json();
console.log(config);
// {
//   "session_id": "...",
//   "enabled": true,
//   "url": "https://meu-webhook.com",
//   "events": ["message"],
//   "updated_at": "2025-11-06T10:30:00Z"
// }
```

---

### **Cenário 3: Você precisa remover o webhook**

**Antes:**
```javascript
// ❌ Tinha que desabilitar manualmente
await fetch(`${API_URL}/sessions/${sessionId}/webhook`, {
  method: 'PUT',
  body: JSON.stringify({
    webhook: {
      enabled: false,
      url: '',
      events: []
    }
  })
});
```

**Agora:**
```javascript
// ✅ Rota específica para limpar
await fetch(`${API_URL}/sessions/${sessionId}/webhook/clear`, {
  method: 'DELETE',
  headers: {
    'apikey': API_KEY
  }
});
```

---

## 🎯 Exemplos de Migração por Linguagem

### **Python**

**Antes:**
```python
import requests

response = requests.put(
    f"{API_URL}/sessions/{session_id}/webhook",
    headers={"apikey": API_KEY},
    json={
        "webhook": {
            "enabled": True,
            "url": "https://meu-webhook.com",
            "events": ["message"]
        }
    }
)
```

**Depois:**
```python
import requests

# Configurar
response = requests.post(
    f"{API_URL}/sessions/{session_id}/webhook/set",
    headers={"apikey": API_KEY},
    json={
        "enabled": True,
        "url": "https://meu-webhook.com",
        "events": ["message"],
        "token": "Bearer meu-token"
    }
)

# Consultar
config = requests.get(
    f"{API_URL}/sessions/{session_id}/webhook/find",
    headers={"apikey": API_KEY}
).json()

# Limpar
requests.delete(
    f"{API_URL}/sessions/{session_id}/webhook/clear",
    headers={"apikey": API_KEY}
)
```

---

### **PHP**

**Antes:**
```php
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "$apiUrl/sessions/$sessionId/webhook");
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "PUT");
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode([
    'webhook' => [
        'enabled' => true,
        'url' => 'https://meu-webhook.com',
        'events' => ['message']
    ]
]));
```

**Depois:**
```php
// Configurar
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "$apiUrl/sessions/$sessionId/webhook/set");
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode([
    'enabled' => true,
    'url' => 'https://meu-webhook.com',
    'events' => ['message'],
    'token' => 'Bearer meu-token'
]));

// Consultar
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "$apiUrl/sessions/$sessionId/webhook/find");
curl_setopt($ch, CURLOPT_HTTPGET, true);

// Limpar
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "$apiUrl/sessions/$sessionId/webhook/clear");
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "DELETE");
```

---

## ⚠️ Compatibilidade

- ✅ **Rota antiga ainda funciona** - Marcada como `@Deprecated`
- ✅ **Sem breaking changes** - Código antigo continua funcionando
- ⚠️ **Recomendação** - Migre para as novas rotas quando possível
- 📅 **Deprecação futura** - A rota antiga pode ser removida em versões futuras

---

## 🎁 Benefícios das Novas Rotas

1. ✅ **Semântica clara** - `/set`, `/find`, `/clear` são auto-explicativos
2. ✅ **RESTful** - Usa métodos HTTP corretos (POST, GET, DELETE)
3. ✅ **Separação de responsabilidades** - Cada rota faz uma coisa
4. ✅ **Validações melhoradas** - Eventos válidos, URL obrigatória
5. ✅ **Token de autenticação** - Novo campo para segurança
6. ✅ **Resposta detalhada** - Retorna configuração completa

---

## 📞 Suporte

Se tiver dúvidas sobre a migração:
1. Consulte a documentação completa em `docs/WEBHOOK_ROUTES.md`
2. Veja exemplos práticos em `docs/WEBHOOK_EXAMPLES.md`
3. Execute o script de teste: `./test_webhook_routes.sh`

