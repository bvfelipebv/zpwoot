# 📡 API Endpoints - zpwoot (zpmeow)

## ✅ Status: 14/14 Endpoints Funcionando (100%)

---

## 📨 Mensagens de Texto

### 1. Enviar Texto
```bash
POST /sessions/{id}/message/text
```
```json
{
  "phone": "559981769536",
  "message": "Olá, mundo!"
}
```

---

## 🖼️ Mensagens de Mídia

### 2. Enviar Imagem
```bash
POST /sessions/{id}/message/image
```
```json
{
  "phone": "559981769536",
  "image": "https://example.com/image.jpg",
  "caption": "Legenda opcional"
}
```

### 3. Enviar Áudio
```bash
POST /sessions/{id}/message/audio
```
```json
{
  "phone": "559981769536",
  "audio": "https://example.com/audio.mp3"
}
```

### 4. Enviar Vídeo
```bash
POST /sessions/{id}/message/video
```
```json
{
  "phone": "559981769536",
  "video": "https://example.com/video.mp4",
  "caption": "Legenda opcional"
}
```

### 5. Enviar Documento
```bash
POST /sessions/{id}/message/document
```
```json
{
  "phone": "559981769536",
  "document": "https://example.com/doc.pdf",
  "fileName": "documento.pdf",
  "caption": "Legenda opcional"
}
```

### 6. Enviar Sticker ⭐ NOVO
```bash
POST /sessions/{id}/message/sticker
```
**URL:**
```json
{
  "phone": "559981769536",
  "sticker": "https://example.com/sticker.webp"
}
```
**Base64:**
```json
{
  "phone": "559981769536",
  "stickerBase64": "data:image/webp;base64,UklGRiQAAABXRUJQ..."
}
```

---

## 👤 Contatos e Localização

### 7. Enviar Contato(s) ⭐ MELHORADO
```bash
POST /sessions/{id}/message/contact
```
**Contato Único:**
```json
{
  "phone": "559981769536",
  "contacts": [
    {
      "name": "João Silva",
      "phone": "559981769536"
    }
  ]
}
```
**Lista de Contatos:**
```json
{
  "phone": "559981769536",
  "contacts": [
    {"name": "João Silva", "phone": "559981769536"},
    {"name": "Maria Santos", "phone": "5511999999999"},
    {"name": "Pedro Costa", "phone": "5521888888888"}
  ]
}
```

### 8. Enviar Localização
```bash
POST /sessions/{id}/message/location
```
```json
{
  "phone": "559981769536",
  "latitude": -23.5505,
  "longitude": -46.6333,
  "name": "São Paulo"
}
```

---

## 🗳️ Enquetes e Reações

### 9. Enviar Enquete
```bash
POST /sessions/{id}/message/poll
```
```json
{
  "phone": "559981769536",
  "question": "Qual sua cor favorita?",
  "options": ["Vermelho", "Azul", "Verde"],
  "selectableCount": 1
}
```

### 10. Enviar Reação
```bash
POST /sessions/{id}/message/reaction
```
```json
{
  "phone": "559981769536",
  "messageId": "3EB02C7ADA457D1D68D8AC",
  "emoji": "👍"
}
```

---

## 👁️ Status e Presença

### 11. Enviar Presença (Digitando/Gravando)
```bash
POST /sessions/{id}/message/presence
```
```json
{
  "phone": "559981769536",
  "presence": "composing"
}
```
**Opções:** `composing` (digitando), `recording` (gravando áudio), `paused` (parou)

---

## ✏️ Ações em Mensagens

### 12. Marcar como Lida
```bash
POST /sessions/{id}/message/read
```
```json
{
  "phone": "559981769536",
  "messageIds": ["3EB02C7ADA457D1D68D8AC"]
}
```

### 13. Editar Mensagem
```bash
PUT /sessions/{id}/message/edit
```
```json
{
  "phone": "559981769536",
  "messageId": "3EB02C7ADA457D1D68D8AC",
  "newMessage": "Mensagem editada"
}
```

### 14. Revogar/Apagar Mensagem
```bash
DELETE /sessions/{id}/message/revoke
```
```json
{
  "phone": "559981769536",
  "messageId": "3EB02C7ADA457D1D68D8AC"
}
```

---

## 🔑 Autenticação

Todas as requisições requerem header:
```
apikey: your-secret-api-key-here
```

---

## ✨ Funcionalidades Especiais

### Contatos
- ✅ Geração automática de vCard com `waid` (WhatsApp ID)
- ✅ Formatação automática de telefone brasileiro
- ✅ Botão "Conversar" funcionando
- ✅ Suporte a contato único ou lista de contatos

### Sticker
- ✅ Suporte a URL
- ✅ Suporte a Base64
- ✅ Aceita WebP, PNG e JPEG

### Mídia
- ✅ Suporte a URL
- ✅ Suporte a Base64 (data:image/jpeg;base64,...)
- ✅ Download automático
- ✅ Upload automático para WhatsApp

