# 📝 Formato de Requisição da API - ZPWoot

## ✅ Novo Formato Simplificado

A API foi atualizada para usar um formato mais limpo e estruturado com objetos JSON para proxy e webhook.

---

## 📥 POST /sessions/create

### Formato Completo

```json
{
  "name": "sessao-atendimento-1",
  "apikey": null,
  "proxy": {
    "enabled": true,
    "protocol": "http",
    "host": "10.0.0.1",
    "port": 3128,
    "username": "proxyuser",
    "password": "proxypass"
  },
  "webhook": {
    "enabled": true,
    "url": "https://hooks.exemplo.com/wuz",
    "events": ["message", "status", "qr"],
    "token": "secreto-opcional"
  }
}
```

### Campos

#### `name` (obrigatório)
- **Tipo**: string
- **Tamanho**: 3-100 caracteres
- **Descrição**: Nome identificador da sessão
- **Exemplo**: `"sessao-atendimento-1"`

#### `apikey` (opcional)
- **Tipo**: string ou null
- **Descrição**: API Key específica para esta sessão (se null, usa a global)
- **Exemplo**: `null` ou `"minha-chave-especifica"`

#### `proxy` (opcional)
- **Tipo**: objeto
- **Descrição**: Configuração de proxy para conexão WhatsApp

**Campos do proxy:**
- `enabled` (boolean): Se o proxy está ativado
- `protocol` (string): Protocolo (`http`, `https`, `socks5`)
- `host` (string): Endereço do servidor proxy
- `port` (integer): Porta do proxy (1-65535)
- `username` (string, opcional): Usuário para autenticação
- `password` (string, opcional): Senha para autenticação

#### `webhook` (opcional)
- **Tipo**: objeto
- **Descrição**: Configuração de webhook para receber eventos

**Campos do webhook:**
- `enabled` (boolean): Se o webhook está ativado
- `url` (string): URL que receberá os eventos
- `events` (array): Lista de eventos a receber
- `token` (string, opcional): Token para validação de requisições

**Eventos disponíveis:**
- `message` - Mensagens recebidas
- `status` - Mudanças de status
- `qr` - QR Code gerado
- `connected` - Sessão conectada
- `disconnected` - Sessão desconectada

---

## 📋 Exemplos de Uso

### Exemplo 1: Sessão Simples (sem proxy/webhook)

```json
{
  "name": "sessao-basica",
  "apikey": null
}
```

### Exemplo 2: Sessão com Webhook

```json
{
  "name": "sessao-com-webhook",
  "apikey": null,
  "webhook": {
    "enabled": true,
    "url": "https://meu-servidor.com/webhook",
    "events": ["message", "status"],
    "token": "meu-token-secreto"
  }
}
```

### Exemplo 3: Sessão com Proxy

```json
{
  "name": "sessao-com-proxy",
  "apikey": null,
  "proxy": {
    "enabled": true,
    "protocol": "http",
    "host": "proxy.empresa.com",
    "port": 8080,
    "username": "usuario",
    "password": "senha123"
  }
}
```

### Exemplo 4: Sessão Completa

```json
{
  "name": "sessao-completa",
  "apikey": "chave-especifica-desta-sessao",
  "proxy": {
    "enabled": true,
    "protocol": "socks5",
    "host": "10.0.0.1",
    "port": 1080,
    "username": "proxyuser",
    "password": "proxypass"
  },
  "webhook": {
    "enabled": true,
    "url": "https://hooks.exemplo.com/whatsapp",
    "events": ["message", "qr", "connected", "disconnected", "status"],
    "token": "webhook-secret-token-123"
  }
}
```

---

## 🔄 PUT /sessions/{id}/webhook

### Formato

```json
{
  "webhook": {
    "enabled": true,
    "url": "https://novo-webhook.com/eventos",
    "events": ["message", "status"],
    "token": "novo-token"
  }
}
```

---

## 📤 Resposta da API

### Sucesso (201 Created)

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "sessao-atendimento-1",
  "jid": null,
  "status": "disconnected",
  "webhook_url": "https://hooks.exemplo.com/wuz",
  "webhook_events": ["message", "status", "qr"],
  "created_at": "2025-11-05T19:00:00Z",
  "updated_at": "2025-11-05T19:00:00Z"
}
```

### Erro (400 Bad Request)

```json
{
  "error": "invalid_request",
  "message": "Key: 'CreateSessionRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

---

## 🧪 Testando com cURL

### Criar Sessão Simples

```bash
curl -X POST http://localhost:8080/sessions/create \
  -H "apikey: sua-chave-aqui" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "minha-sessao"
  }'
```

### Criar Sessão com Webhook

```bash
curl -X POST http://localhost:8080/sessions/create \
  -H "apikey: sua-chave-aqui" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sessao-webhook",
    "webhook": {
      "enabled": true,
      "url": "https://meu-webhook.com/eventos",
      "events": ["message", "qr"]
    }
  }'
```

### Criar Sessão com Proxy

```bash
curl -X POST http://localhost:8080/sessions/create \
  -H "apikey: sua-chave-aqui" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sessao-proxy",
    "proxy": {
      "enabled": true,
      "protocol": "http",
      "host": "10.0.0.1",
      "port": 3128
    }
  }'
```

---

## ✅ Validações

### Nome da Sessão
- ✅ Obrigatório
- ✅ Mínimo 3 caracteres
- ✅ Máximo 100 caracteres

### Proxy
- ✅ Protocol deve ser: `http`, `https` ou `socks5`
- ✅ Port deve estar entre 1 e 65535
- ✅ Host é obrigatório se enabled=true

### Webhook
- ✅ URL deve ser válida se enabled=true
- ✅ Events é opcional (padrão: todos os eventos)

---

## 📚 Documentação Relacionada

- Swagger UI: http://localhost:8080/swagger/index.html
- `docs/QUICK_START.md` - Início rápido
- `docs/SWAGGER_EXAMPLES.md` - Exemplos completos

