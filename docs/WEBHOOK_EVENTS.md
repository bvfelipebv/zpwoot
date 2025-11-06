# 📡 Eventos de Webhook - zpmeow

## 📚 Documentação Completa dos Eventos

Este documento lista todos os eventos de webhook suportados pelo zpmeow, baseados na biblioteca oficial [whatsmeow](https://pkg.go.dev/go.mau.fi/whatsmeow/types/events).

---

## 🎯 Eventos Recomendados

Para a maioria dos casos de uso, recomendamos subscrever aos seguintes eventos:

```json
{
  "events": [
    "message",
    "receipt",
    "connected",
    "disconnected",
    "logged_out",
    "qr",
    "group_info",
    "joined_group",
    "presence",
    "chat_presence"
  ]
}
```

---

## 📋 Categorias de Eventos

### 1️⃣ **Mensagens** (`messages`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `message` | Mensagem recebida (texto, mídia, documentos, etc) | `*events.Message` |
| `undecryptable_message` | Mensagem que não pôde ser descriptografada | `*events.UndecryptableMessage` |
| `receipt` | Confirmação de entrega/leitura de mensagem | `*events.Receipt` |
| `media_retry` | Resposta a solicitação de reenvio de mídia | `*events.MediaRetry` |
| `delete_for_me` | Mensagem deletada apenas para o usuário | `*events.DeleteForMe` |

**Exemplo de uso:**
```json
{
  "enabled": true,
  "url": "https://seu-servidor.com/webhook",
  "events": ["message", "receipt"]
}
```

---

### 2️⃣ **Grupos e Contatos** (`groups_contacts`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `group_info` | Metadados de grupo alterados (nome, descrição, participantes) | `*events.GroupInfo` |
| `joined_group` | Entrou ou foi adicionado a um grupo | `*events.JoinedGroup` |
| `picture` | Foto de perfil de usuário ou grupo alterada | `*events.Picture` |
| `blocklist_change` | Mudança individual na lista de bloqueados | `*events.BlocklistChange` |
| `blocklist` | Lista completa de bloqueados atualizada | `*events.Blocklist` |
| `contact` | Entrada na lista de contatos modificada | `*events.Contact` |
| `push_name` | Nome de exibição de contato mudou | `*events.PushName` |
| `business_name` | Nome comercial verificado mudou | `*events.BusinessName` |

---

### 3️⃣ **Conexão e Sessão** (`connection`) ⚠️ CRÍTICO

| Evento | Descrição | Tipo whatsmeow | Crítico |
|--------|-----------|----------------|---------|
| `connected` | Conectado e autenticado com sucesso | `*events.Connected` | ✅ |
| `disconnected` | WebSocket fechado pelo servidor | `*events.Disconnected` | ✅ |
| `connect_failure` | Servidor rejeitou a conexão | `*events.ConnectFailure` | ⚠️ |
| `keepalive_restored` | Pings keepalive restaurados após timeout | `*events.KeepAliveRestored` | ⚠️ |
| `keepalive_timeout` | Ping keepalive expirou | `*events.KeepAliveTimeout` | ✅ |
| `logged_out` | Desconectado do telefone | `*events.LoggedOut` | ✅ |
| `client_outdated` | Cliente desatualizado | `*events.ClientOutdated` | ✅ |
| `temporary_ban` | Conta temporariamente banida | `*events.TemporaryBan` | ✅ |
| `stream_error` | Erro de stream desconhecido | `*events.StreamError` | ⚠️ |
| `stream_replaced` | Login em outro dispositivo | `*events.StreamReplaced` | ✅ |
| `pair_success` | QR code escaneado com sucesso | `*events.PairSuccess` | ℹ️ |
| `pair_error` | Erro no pareamento | `*events.PairError` | ⚠️ |
| `qr` | QR code gerado para pareamento | `*events.QR` | ℹ️ |
| `qr_scanned_without_multidevice` | QR escaneado mas telefone sem multidevice | `*events.QRScannedWithoutMultidevice` | ⚠️ |
| `manual_login_reconnect` | Reconexão manual necessária após login | `*events.ManualLoginReconnect` | ℹ️ |

**⚠️ IMPORTANTE:** Sempre monitore os eventos críticos de conexão!

---

### 4️⃣ **Privacidade e Configurações** (`privacy`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `privacy_settings` | Configurações de privacidade alteradas | `*events.PrivacySettings` |
| `push_name_setting` | Push name alterado de outro dispositivo | `*events.PushNameSetting` |
| `user_about` | Status/sobre de usuário alterado | `*events.UserAbout` |
| `user_status_mute` | Atualizações de status silenciadas/dessilenciadas | `*events.UserStatusMute` |

---

### 5️⃣ **Sincronização e Estado** (`sync`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `app_state` | Novos dados de sincronização de app state | `*events.AppState` |
| `app_state_sync_complete` | App state foi ressincronizado | `*events.AppStateSyncComplete` |
| `history_sync` | Telefone enviou blob de mensagens históricas | `*events.HistorySync` |
| `offline_sync_completed` | Servidor terminou de enviar eventos perdidos | `*events.OfflineSyncCompleted` |
| `offline_sync_preview` | Preview de eventos que serão sincronizados | `*events.OfflineSyncPreview` |
| `archive` | Chat arquivado/desarquivado de outro dispositivo | `*events.Archive` |
| `pin` | Chat fixado/desfixado de outro dispositivo | `*events.Pin` |
| `mute` | Chat silenciado/dessilenciado de outro dispositivo | `*events.Mute` |
| `mark_chat_as_read` | Chat marcado como lido/não lido | `*events.MarkChatAsRead` |
| `delete_chat` | Chat deletado de outro dispositivo | `*events.DeleteChat` |
| `clear_chat` | Chat limpo de outro dispositivo | `*events.ClearChat` |
| `star` | Mensagem favoritada/desfavoritada | `*events.Star` |
| `unarchive_chats_setting` | Configuração "Manter chats arquivados" alterada | `*events.UnarchiveChatsSetting` |
| `label_edit` | Label editada | `*events.LabelEdit` |
| `label_association_chat` | Chat etiquetado/desetiquetado | `*events.LabelAssociationChat` |
| `label_association_message` | Mensagem etiquetada/desetiquetada | `*events.LabelAssociationMessage` |

---

### 6️⃣ **Chamadas** (`calls`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `call_offer` | Chamada recebida (1:1) | `*events.CallOffer` |
| `call_accept` | Chamada aceita | `*events.CallAccept` |
| `call_terminate` | Chamada terminada | `*events.CallTerminate` |
| `call_offer_notice` | Notificação de oferta de chamada (grupos) | `*events.CallOfferNotice` |
| `call_relay_latency` | Latência do relay de chamada | `*events.CallRelayLatency` |
| `call_pre_accept` | Chamada pré-aceita | `*events.CallPreAccept` |
| `call_reject` | Chamada rejeitada | `*events.CallReject` |
| `call_transport` | Transporte de chamada | `*events.CallTransport` |
| `unknown_call_event` | Evento de chamada desconhecido | `*events.UnknownCallEvent` |

---

### 7️⃣ **Presença e Atividade** (`presence`)

| Evento | Descrição | Tipo whatsmeow | Requisitos |
|--------|-----------|----------------|------------|
| `presence` | Atualização de presença (online/offline/last seen) | `*events.Presence` | Requer subscrição |
| `chat_presence` | Estado de chat (digitando, gravando áudio) | `*events.ChatPresence` | Requer estar online |

**Nota:** Para receber eventos de presença, você precisa:
- `presence`: Subscrever com `client.SubscribePresence(userJID)`
- `chat_presence`: Estar online com `client.SendPresence(types.PresenceAvailable)`

---

### 8️⃣ **Identidade e Segurança** (`identity`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `identity_change` | Outro usuário mudou seu dispositivo principal | `*events.IdentityChange` |
| `cat_refresh_error` | Erro ao atualizar CAT (Client Access Token) | `*events.CATRefreshError` |

---

### 9️⃣ **Newsletter/Canais** (`newsletter`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `newsletter_join` | Entrou em um canal | `*events.NewsletterJoin` |
| `newsletter_leave` | Saiu de um canal | `*events.NewsletterLeave` |
| `newsletter_mute_change` | Mudança de silenciamento de canal | `*events.NewsletterMuteChange` |
| `newsletter_live_update` | Atualização ao vivo de canal | `*events.NewsletterLiveUpdate` |

---

### 🔟 **Facebook/Meta Bridge** (`facebook`)

| Evento | Descrição | Tipo whatsmeow |
|--------|-----------|----------------|
| `fb_message` | Mensagem recebida do Facebook/Instagram | `*events.FBMessage` |

---

### ⭐ **Especiais** (`special`)

| Evento | Descrição | Uso |
|--------|-----------|-----|
| `all` | Recebe TODOS os eventos | ⚠️ Use com cuidado - muito tráfego |

---

## 🎯 Casos de Uso Comuns

### **Bot de Atendimento**
```json
{
  "events": [
    "message",
    "receipt",
    "connected",
    "disconnected",
    "logged_out"
  ]
}
```

### **Monitor de Grupos**
```json
{
  "events": [
    "message",
    "group_info",
    "joined_group",
    "picture"
  ]
}
```

### **Sistema de Presença**
```json
{
  "events": [
    "presence",
    "chat_presence",
    "connected",
    "disconnected"
  ]
}
```

### **Monitoramento Completo**
```json
{
  "events": ["all"]
}
```

---

## 📊 Estatísticas

- **Total de eventos:** 60+
- **Categorias:** 10
- **Eventos críticos:** 7
- **Eventos padrão:** 6
- **Eventos recomendados:** 10

---

## 🔗 Referências

- [Documentação oficial whatsmeow](https://pkg.go.dev/go.mau.fi/whatsmeow/types/events)
- [Código fonte zpmeow](https://github.com/seu-repo/zpmeow)

