# TaskList - Implementação de Envio de Mensagens

Baseado no estudo do código wuzapi (wmiau.go e handlers.go)

## 📋 Análise do Código Wuzapi

### Principais Aprendizados:

1. **Download/Decode de Mídia**:
   - Usa `dataurl.DecodeString()` para base64
   - Formato: `data:image/jpeg;base64,XXXXX`
   - Detecta MIME type automaticamente

2. **Upload de Mídia**:
   - `client.Upload(ctx, data, whatsmeow.MediaImage/Audio/Video/Document)`
   - Retorna `UploadResponse` com URL, DirectPath, MediaKey, FileEncSHA256, FileSHA256

3. **Estrutura de Mensagens**:
   - Usa protobuf `waE2E.Message`
   - Cada tipo tem seu próprio campo (ImageMessage, AudioMessage, etc)
   - Campos importantes: URL, DirectPath, MediaKey, Mimetype, FileLength, Caption

4. **ContextInfo** (Reply/Quote):
   - StanzaID: ID da mensagem original
   - Participant: JID do participante (em grupos)
   - QuotedMessage: Mensagem citada
   - MentionedJID: Array de JIDs mencionados

---

## ✅ Tasks

### 1. [ ] Implementar Helper de Download/Decode de Mídia
**Arquivo**: `internal/service/media_helper.go`

**Funções**:
```go
// downloadOrDecodeMedia baixa de URL ou decodifica base64
func downloadOrDecodeMedia(mediaURL string) ([]byte, string, error)

// detectMimeType detecta o tipo MIME dos dados
func detectMimeType(data []byte) string

// downloadFromURL faz download de URL HTTP/HTTPS
func downloadFromURL(url string) ([]byte, string, error)

// decodeBase64Media decodifica data URL base64
func decodeBase64Media(dataURL string) ([]byte, string, error)
```

**Referência wuzapi**:
- `dataurl.DecodeString()` para base64
- `http.DetectContentType()` para MIME type
- Validação de prefixos: `data:image/`, `data:audio/`, etc

---

### 2. [ ] Implementar Envio de Imagem Completo
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendImageFromURL()`

**Campos ImageMessage**:
- URL (do upload)
- DirectPath
- MediaKey
- Mimetype
- FileEncSHA256
- FileSHA256
- FileLength
- Caption (opcional)
- JPEGThumbnail (opcional)
- ContextInfo (para reply)

**Referência wuzapi**: handlers.go linha ~1500

---

### 3. [ ] Implementar Envio de Áudio Completo
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendAudioFromURL()`

**Campos AudioMessage**:
- URL
- DirectPath
- MediaKey
- Mimetype
- FileEncSHA256
- FileSHA256
- FileLength
- PTT (bool - push to talk, default true)
- Seconds (duração)
- Waveform ([]byte - forma de onda)
- ContextInfo

**Referência wuzapi**: handlers.go linha ~1700

---

### 4. [ ] Implementar Envio de Vídeo Completo
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendVideoFromURL()`

**Campos VideoMessage**:
- URL
- DirectPath
- MediaKey
- Mimetype
- FileEncSHA256
- FileSHA256
- FileLength
- Caption
- JPEGThumbnail
- Seconds
- ContextInfo

**Referência wuzapi**: handlers.go linha ~1900

---

### 5. [ ] Implementar Envio de Documento Completo
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendDocumentFromURL()`

**Campos DocumentMessage**:
- URL
- DirectPath
- MediaKey
- Mimetype
- FileEncSHA256
- FileSHA256
- FileLength
- FileName
- Caption
- ContextInfo

**Referência wuzapi**: handlers.go linha ~2100

---

### 6. [ ] Implementar Envio de Sticker
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendSticker()`

**Campos StickerMessage**:
- URL
- DirectPath
- MediaKey
- Mimetype (image/webp)
- FileEncSHA256
- FileSHA256
- FileLength
- IsAnimated (bool)

**Referência wuzapi**: wmiau.go linha ~800

---

### 7. [ ] Implementar Envio de Localização
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendLocation()`

**Campos LocationMessage**:
- DegreesLatitude (float64)
- DegreesLongitude (float64)
- Name (string)
- Address (string)

**Referência**: Protobuf waE2E.LocationMessage

---

### 8. [ ] Implementar Envio de Contato
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendContact()`

**Campos ContactMessage**:
- DisplayName (string)
- Vcard (string - formato vCard 3.0)

**Exemplo vCard**:
```
BEGIN:VCARD
VERSION:3.0
FN:Nome Completo
TEL;type=CELL:+5511999999999
END:VCARD
```

---

### 9. [ ] Implementar Envio de Enquete (Poll)
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendPoll()`

**Campos PollCreationMessage**:
- Name (pergunta)
- Options ([]PollOption)
- SelectableOptionsCount (uint32)

---

### 10. [ ] Implementar Envio de Reação
**Arquivo**: `internal/service/message_service.go`

**Função**: `SendReaction()`

**Campos ReactionMessage**:
- Key (MessageKey com ID da mensagem)
- Text (emoji)
- SenderTimestampMS

---

### 11. [ ] Implementar Marcar Como Lida
**Arquivo**: `internal/service/message_service.go`

**Função**: `MarkAsRead()`

**Método**: `client.MarkRead(messageIDs, timestamp, chat, sender)`

---

### 12. [ ] Implementar Revogar Mensagem
**Arquivo**: `internal/service/message_service.go`

**Função**: `RevokeMessage()`

**Campos ProtocolMessage**:
- Type: REVOKE (0)
- Key: MessageKey com ID da mensagem

---

### 13. [ ] Implementar Editar Mensagem
**Arquivo**: `internal/service/message_service.go`

**Função**: `EditMessage()`

**Campos EditedMessage**:
- Message: Nova mensagem
- Key: MessageKey com ID original
- TimestampMS

---

### 14. [ ] Atualizar DTOs com ContextInfo
**Arquivo**: `internal/api/dto/message_dto.go`

**Adicionar em todos os DTOs**:
```go
type ContextInfo struct {
    StanzaID      *string  `json:"stanzaId,omitempty"`
    Participant   *string  `json:"participant,omitempty"`
    MentionedJID  []string `json:"mentionedJid,omitempty"`
    IsForwarded   *bool    `json:"isForwarded,omitempty"`
}
```

---

### 15. [ ] Testar Todos os Endpoints
**Número de teste**: 559981769536

**Checklist**:
- [ ] Texto
- [ ] Imagem (URL e base64)
- [ ] Áudio (URL e base64)
- [ ] Vídeo (URL e base64)
- [ ] Documento (URL e base64)
- [ ] Sticker
- [ ] Localização
- [ ] Contato
- [ ] Enquete
- [ ] Reação
- [ ] Presença
- [ ] Marcar como lida
- [ ] Revogar
- [ ] Editar
- [ ] Reply (com ContextInfo)
- [ ] Menções

---

## 📦 Dependências Necessárias

```bash
go get github.com/vincent-petithory/dataurl
```

---

## 🎯 Ordem de Implementação Sugerida

1. ✅ Helper de download/decode (base para tudo)
2. ✅ Imagem (mais comum)
3. ✅ Áudio
4. ✅ Vídeo
5. ✅ Documento
6. ✅ Localização (simples, sem upload)
7. ✅ Contato (simples, sem upload)
8. ✅ Reação (simples)
9. ✅ Marcar como lida (simples)
10. ✅ Sticker
11. ✅ Enquete
12. ✅ Revogar
13. ✅ Editar
14. ✅ ContextInfo (reply/menções)
15. ✅ Testes completos

