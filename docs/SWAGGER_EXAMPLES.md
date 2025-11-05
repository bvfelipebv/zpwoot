# 📝 Exemplos Completos do Swagger - ZPWoot API

## ✅ DTOs com Exemplos Completos

Todos os modelos de dados (DTOs) agora possuem exemplos completos que aparecem automaticamente no Swagger UI.

---

## 📥 REQUEST MODELS

### 1. CreateSessionRequest
**Endpoint:** `POST /sessions/create`

```json
{
  "name": "Minha Sessão WhatsApp",
  "webhook_url": "https://seu-webhook.com/whatsapp",
  "webhook_events": [
    "message",
    "qr",
    "connected",
    "disconnected"
  ],
  "metadata": {
    "cliente": "Empresa XYZ",
    "ambiente": "producao",
    "responsavel": "João Silva"
  }
}
```

**Campos:**
- `name` (obrigatório): Nome da sessão (3-100 caracteres)
- `webhook_url` (opcional): URL para receber eventos
- `webhook_events` (opcional): Lista de eventos a receber
- `metadata` (opcional): Dados customizados em formato chave-valor

---

### 2. PairPhoneRequest
**Endpoint:** `POST /sessions/{id}/pair`

```json
{
  "phone_number": "+5511999999999"
}
```

**Campos:**
- `phone_number` (obrigatório): Número com código do país (formato E.164)

**Exemplos de formatos válidos:**
- Brasil: `+5511999999999`
- EUA: `+15551234567`
- Portugal: `+351912345678`

---

### 3. UpdateWebhookRequest
**Endpoint:** `PUT /sessions/{id}/webhook`

```json
{
  "webhook_url": "https://novo-webhook.com/whatsapp",
  "webhook_events": [
    "message",
    "qr",
    "connected"
  ],
  "webhook_secret": "meu-secret-super-seguro-123"
}
```

**Campos:**
- `webhook_url` (obrigatório): Nova URL do webhook
- `webhook_events` (obrigatório): Lista de eventos (mínimo 1)
- `webhook_secret` (opcional): Secret para validação (mínimo 16 caracteres)

**Eventos disponíveis:**
- `message` - Mensagens recebidas
- `qr` - QR Code gerado
- `connected` - Sessão conectada
- `disconnected` - Sessão desconectada
- `status` - Mudanças de status

---

### 4. ConnectSessionRequest
**Endpoint:** `POST /sessions/{id}/connect`

```json
{
  "auto_reconnect": true
}
```

**Campos:**
- `auto_reconnect` (opcional): Reconectar automaticamente se desconectar

---

## 📤 RESPONSE MODELS

### 1. SessionResponse
**Retornado em:** `POST /sessions/create`, `GET /sessions/{id}/info`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Minha Sessão WhatsApp",
  "jid": "5511999999999@s.whatsapp.net",
  "status": "connected",
  "push_name": "João Silva",
  "platform": "android",
  "business_name": "Minha Empresa LTDA",
  "webhook_url": "https://seu-webhook.com/whatsapp",
  "webhook_events": [
    "message",
    "qr",
    "connected"
  ],
  "last_connected": "2025-11-05T18:30:00Z",
  "last_disconnected": "2025-11-05T17:00:00Z",
  "created_at": "2025-11-05T10:00:00Z",
  "updated_at": "2025-11-05T18:30:00Z"
}
```

**Campos:**
- `id`: UUID da sessão
- `name`: Nome da sessão
- `jid`: ID do WhatsApp (quando conectado)
- `status`: Status atual (`connected`, `disconnected`, `connecting`)
- `push_name`: Nome do usuário no WhatsApp
- `platform`: Plataforma (`android`, `ios`, `web`)
- `business_name`: Nome comercial (se for WhatsApp Business)
- `webhook_url`: URL do webhook configurada
- `webhook_events`: Eventos configurados
- `last_connected`: Última conexão
- `last_disconnected`: Última desconexão
- `created_at`: Data de criação
- `updated_at`: Última atualização

---

### 2. SessionListResponse
**Retornado em:** `GET /sessions/list`

```json
{
  "sessions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Minha Sessão WhatsApp",
      "status": "connected",
      "created_at": "2025-11-05T10:00:00Z"
    }
  ],
  "total": 3
}
```

---

### 3. SessionStatusResponse
**Retornado em:** `GET /sessions/{id}/status`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "connected",
  "is_connected": true,
  "is_logged_in": true,
  "jid": "5511999999999@s.whatsapp.net",
  "push_name": "João Silva",
  "platform": "android",
  "last_connected": "2025-11-05T18:30:00Z",
  "connection_time": "2h 30m 15s",
  "needs_pairing": false,
  "can_connect": true
}
```

**Campos:**
- `is_connected`: Se está conectado ao WhatsApp
- `is_logged_in`: Se está autenticado
- `connection_time`: Tempo de conexão formatado
- `needs_pairing`: Se precisa parear novamente
- `can_connect`: Se pode conectar agora

---

### 4. PairPhoneResponse
**Retornado em:** `POST /sessions/{id}/pair`

```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "phone_number": "+5511999999999",
  "pairing_code": "ABCD-1234",
  "message": "Enter the pairing code on your phone"
}
```

**Como usar:**
1. Faça a requisição com seu número
2. Copie o `pairing_code`
3. No WhatsApp: Configurações > Aparelhos conectados > Conectar aparelho
4. Digite o código

---

### 5. SuccessResponse
**Retornado em:** Operações de sucesso genéricas

```json
{
  "success": true,
  "message": "Operação realizada com sucesso",
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

---

### 6. ErrorResponse
**Retornado em:** Erros (400, 404, 500)

```json
{
  "error": "invalid_request",
  "message": "Nome da sessão é obrigatório",
  "details": {
    "field": "name",
    "error": "required"
  }
}
```

**Tipos de erro comuns:**
- `invalid_request` - Requisição inválida
- `not_found` - Recurso não encontrado
- `unauthorized` - Não autenticado
- `internal_error` - Erro interno do servidor

---

## 🎯 Como Ver no Swagger UI

1. **Acesse:** http://localhost:8080/swagger/index.html

2. **Clique em qualquer endpoint** (ex: POST /sessions/create)

3. **Clique em "Try it out"**

4. **Veja o exemplo pré-preenchido** no campo de request

5. **Modifique conforme necessário**

6. **Execute** e veja a resposta com exemplos

---

## 📋 Benefícios dos Exemplos

✅ **Documentação Visual** - Veja exatamente o formato esperado
✅ **Testes Rápidos** - Exemplos prontos para testar
✅ **Menos Erros** - Formato correto já mostrado
✅ **Aprendizado Fácil** - Entenda a API rapidamente
✅ **Integração Simples** - Copie e cole os exemplos

---

## 🔄 Regenerar Documentação

Se você modificar os DTOs:

```bash
make swagger
make rebuild
```

---

## 📚 Documentação Relacionada

- `docs/SWAGGER_GUIDE.md` - Guia de uso do Swagger
- `docs/QUICK_START.md` - Início rápido
- `README.md` - Documentação principal

