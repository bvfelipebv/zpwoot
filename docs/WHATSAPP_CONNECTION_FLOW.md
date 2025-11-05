# 📱 Fluxo de Conexão WhatsApp - ZPWoot

## 🎯 Processo Implementado

### ✅ **CORREÇÃO APLICADA**

O processo de conexão foi corrigido para seguir a documentação oficial do whatsmeow:

**Antes (❌ ERRADO):**
```go
// Verificava se estava pareado ANTES de conectar
if session.DeviceJID == "" {
    return fmt.Errorf("session not paired yet")
}
client.Connect()
```

**Agora (✅ CORRETO):**
```go
// Conecta PRIMEIRO, depois gera QR Code ou pareia
if session.DeviceJID == "" {
    // Obter canal de QR ANTES de conectar
    qrChan, _ := client.GetQRChannel(ctx)
    client.Connect()
    // Processar QR codes
    go handleQRCodes(ctx, sessionID, qrChan)
} else {
    // Já pareado, apenas conectar
    client.Connect()
}
```

---

## 📋 Fluxo Completo

### 1. **Criar Sessão**
```bash
POST /sessions/create
{
  "name": "minha-sessao"
}
```

**Resposta:**
```json
{
  "id": "60bab65d-d00c-46a7-ba99-d580017f690a",
  "name": "minha-sessao",
  "status": "disconnected",
  "connected": false
}
```

---

### 2. **Conectar Sessão (Gera QR Code)**
```bash
POST /sessions/60bab65d-d00c-46a7-ba99-d580017f690a/connect
```

**O que acontece:**
1. ✅ Cria cliente whatsmeow
2. ✅ Obtém canal de QR Code (`GetQRChannel`)
3. ✅ Conecta ao WhatsApp (`Connect`)
4. ✅ Processa QR codes em background
5. ✅ Salva QR code no banco
6. ✅ Atualiza status para `qr_code`

**Resposta:**
```json
{
  "success": true,
  "message": "Session connecting"
}
```

---

### 3. **Obter QR Code**
```bash
GET /sessions/60bab65d-d00c-46a7-ba99-d580017f690a/info
```

**Resposta:**
```json
{
  "id": "60bab65d-d00c-46a7-ba99-d580017f690a",
  "name": "minha-sessao",
  "status": "qr_code",
  "qr_code": "2@abc123...",
  "connected": false
}
```

---

### 4. **Escanear QR Code**

O usuário escaneia o QR code com o WhatsApp no celular.

**Eventos automáticos:**
1. ✅ WhatsApp envia confirmação
2. ✅ Event handler salva `device_jid`
3. ✅ Atualiza status para `connected`
4. ✅ Marca `connected = true`

---

### 5. **Verificar Status**
```bash
GET /sessions/60bab65d-d00c-46a7-ba99-d580017f690a/status
```

**Resposta (Conectado):**
```json
{
  "session_id": "60bab65d-d00c-46a7-ba99-d580017f690a",
  "name": "minha-sessao",
  "status": "connected",
  "connected": true,
  "device_jid": "5511999999999:1@s.whatsapp.net",
  "qr_code": ""
}
```

---

## 🔄 Método Alternativo: Pair Phone

### 1. **Conectar Sessão**
```bash
POST /sessions/{id}/connect
```

### 2. **Parear com Telefone**
```bash
POST /sessions/{id}/pair
{
  "phone_number": "+5511999999999"
}
```

**Resposta:**
```json
{
  "code": "ABCD-1234"
}
```

### 3. **Inserir Código no WhatsApp**

O usuário insere o código no WhatsApp:
- WhatsApp > Aparelhos conectados > Conectar aparelho
- Inserir código manualmente

---

## 🎯 Estados da Sessão

| Status | Descrição |
|--------|-----------|
| `disconnected` | Sessão criada, não conectada |
| `connecting` | Conectando ao WhatsApp |
| `qr_code` | QR Code gerado, aguardando scan |
| `pairing` | Aguardando código de pareamento |
| `connected` | Conectado e autenticado |
| `failed` | Falha na conexão |
| `logged_out` | Deslogado pelo usuário |

---

## 📊 Diagrama de Fluxo

```
[Criar Sessão] 
    ↓
[Conectar] → GetQRChannel() → Connect()
    ↓
[QR Code Gerado]
    ↓
[Usuário Escaneia]
    ↓
[Event: Paired]
    ↓
[Salvar device_jid]
    ↓
[Status: connected]
```

---

## 🔧 Implementação Técnica

### Código Principal (session_manager.go)

```go
func (m *SessionManager) ConnectSession(ctx context.Context, sessionID string) error {
    // ... validações ...
    
    client := m.whatsappSvc.NewClient(device)
    m.eventHandler.RegisterHandlers(client, sessionID)
    
    if session.DeviceJID == "" {
        // NÃO PAREADO: Gerar QR Code
        qrChan, err := client.GetQRChannel(ctx)
        if err != nil {
            return fmt.Errorf("failed to get QR channel: %w", err)
        }
        
        if err := client.Connect(); err != nil {
            return fmt.Errorf("failed to connect: %w", err)
        }
        
        go m.handleQRCodes(ctx, sessionID, qrChan)
    } else {
        // JÁ PAREADO: Apenas conectar
        if err := client.Connect(); err != nil {
            return fmt.Errorf("failed to connect: %w", err)
        }
    }
    
    // Armazenar cliente ativo
    m.clients[sessionID] = client
    
    return nil
}
```

---

## ✅ Testes

### Teste 1: Criar e Conectar
```bash
# 1. Criar sessão
curl -X POST http://localhost:8080/sessions/create \
  -H "apikey: sua-chave" \
  -H "Content-Type: application/json" \
  -d '{"name": "teste-qr"}'

# 2. Conectar (gera QR)
curl -X POST http://localhost:8080/sessions/{id}/connect \
  -H "apikey: sua-chave"

# 3. Ver QR Code
curl http://localhost:8080/sessions/{id}/info \
  -H "apikey: sua-chave"
```

---

## 🎉 Resultado

✅ **Conexão WhatsApp funcionando corretamente!**

- ✅ QR Code gerado automaticamente
- ✅ Processo assíncrono
- ✅ Status atualizado em tempo real
- ✅ Suporte a pair phone
- ✅ Reconexão automática

---

**Documentação baseada em:** https://pkg.go.dev/go.mau.fi/whatsmeow

