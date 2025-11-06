package constants

// WebhookEventType representa um tipo de evento de webhook do whatsmeow
// Baseado em: https://pkg.go.dev/go.mau.fi/whatsmeow/types/events
type WebhookEventType string

// ============================================================================
// MENSAGENS E COMUNICAÇÃO
// Eventos relacionados ao recebimento e processamento de mensagens
// ============================================================================

const (
	// EventMessage - Mensagem recebida (texto, mídia, documentos, etc)
	// Tipo: *events.Message
	EventMessage WebhookEventType = "message"

	// EventUndecryptableMessage - Mensagem que não pôde ser descriptografada
	// Tipo: *events.UndecryptableMessage
	// A biblioteca automaticamente pede ao remetente para reenviar
	EventUndecryptableMessage WebhookEventType = "undecryptable_message"

	// EventReceipt - Recibo de entrega/leitura de mensagem
	// Tipo: *events.Receipt
	// Indica quando mensagens são entregues ou lidas
	EventReceipt WebhookEventType = "receipt"

	// EventMediaRetry - Resposta a uma solicitação de reenvio de mídia
	// Tipo: *events.MediaRetry
	EventMediaRetry WebhookEventType = "media_retry"

	// EventDeleteForMe - Mensagem deletada apenas para o usuário atual
	// Tipo: *events.DeleteForMe
	EventDeleteForMe WebhookEventType = "delete_for_me"
)

// ============================================================================
// GRUPOS E CONTATOS
// Eventos relacionados a grupos, contatos e listas de bloqueio
// ============================================================================

const (
	// EventGroupInfo - Metadados de grupo alterados (nome, descrição, participantes, etc)
	// Tipo: *events.GroupInfo
	EventGroupInfo WebhookEventType = "group_info"

	// EventJoinedGroup - Usuário entrou ou foi adicionado a um grupo
	// Tipo: *events.JoinedGroup
	EventJoinedGroup WebhookEventType = "joined_group"

	// EventPicture - Foto de perfil de usuário ou grupo foi alterada
	// Tipo: *events.Picture
	EventPicture WebhookEventType = "picture"

	// EventBlocklistChange - Mudança individual na lista de bloqueados
	// Tipo: *events.BlocklistChange
	EventBlocklistChange WebhookEventType = "blocklist_change"

	// EventBlocklist - Lista completa de bloqueados foi atualizada
	// Tipo: *events.Blocklist
	EventBlocklist WebhookEventType = "blocklist"

	// EventContact - Entrada na lista de contatos foi modificada
	// Tipo: *events.Contact
	EventContact WebhookEventType = "contact"

	// EventPushName - Nome de exibição de um contato mudou
	// Tipo: *events.PushName
	EventPushName WebhookEventType = "push_name"

	// EventBusinessName - Nome comercial verificado de um contato mudou
	// Tipo: *events.BusinessName
	EventBusinessName WebhookEventType = "business_name"
)

// ============================================================================
// CONEXÃO E SESSÃO
// Eventos relacionados ao estado da conexão WebSocket e autenticação
// ============================================================================

const (
	// EventConnected - Cliente conectado e autenticado com sucesso
	// Tipo: *events.Connected
	EventConnected WebhookEventType = "connected"

	// EventDisconnected - WebSocket foi fechado pelo servidor
	// Tipo: *events.Disconnected
	EventDisconnected WebhookEventType = "disconnected"

	// EventConnectFailure - Servidor rejeitou a conexão com código de erro desconhecido
	// Tipo: *events.ConnectFailure
	EventConnectFailure WebhookEventType = "connect_failure"

	// EventKeepAliveRestored - Pings keepalive voltaram a funcionar após timeout
	// Tipo: *events.KeepAliveRestored
	EventKeepAliveRestored WebhookEventType = "keepalive_restored"

	// EventKeepAliveTimeout - Ping keepalive para o servidor expirou
	// Tipo: *events.KeepAliveTimeout
	EventKeepAliveTimeout WebhookEventType = "keepalive_timeout"

	// EventLoggedOut - Cliente foi desconectado do telefone
	// Tipo: *events.LoggedOut
	// Pode ocorrer durante conexão ou após conectar
	EventLoggedOut WebhookEventType = "logged_out"

	// EventClientOutdated - Servidor rejeitou conexão por cliente desatualizado
	// Tipo: *events.ClientOutdated
	EventClientOutdated WebhookEventType = "client_outdated"

	// EventTemporaryBan - Conta foi temporariamente banida
	// Tipo: *events.TemporaryBan
	EventTemporaryBan WebhookEventType = "temporary_ban"

	// EventStreamError - Servidor enviou erro de stream com código desconhecido
	// Tipo: *events.StreamError
	EventStreamError WebhookEventType = "stream_error"

	// EventStreamReplaced - Cliente foi desconectado por outro cliente com mesmas chaves
	// Tipo: *events.StreamReplaced
	EventStreamReplaced WebhookEventType = "stream_replaced"

	// EventPairSuccess - QR code foi escaneado e handshake completado
	// Tipo: *events.PairSuccess
	EventPairSuccess WebhookEventType = "pair_success"

	// EventPairError - QR code escaneado mas pareamento local falhou
	// Tipo: *events.PairError
	EventPairError WebhookEventType = "pair_error"

	// EventQR - QR codes foram gerados para pareamento
	// Tipo: *events.QR
	EventQR WebhookEventType = "qr"

	// EventQRScannedWithoutMultidevice - QR escaneado mas telefone não tem multidevice
	// Tipo: *events.QRScannedWithoutMultidevice
	EventQRScannedWithoutMultidevice WebhookEventType = "qr_scanned_without_multidevice"

	// EventManualLoginReconnect - Emitido após login se DisableLoginAutoReconnect está ativo
	// Tipo: *events.ManualLoginReconnect
	EventManualLoginReconnect WebhookEventType = "manual_login_reconnect"
)

// ============================================================================
// PRIVACIDADE E CONFIGURAÇÕES
// Eventos relacionados a configurações de privacidade e perfil do usuário
// ============================================================================

const (
	// EventPrivacySettings - Usuário alterou configurações de privacidade
	// Tipo: *events.PrivacySettings
	EventPrivacySettings WebhookEventType = "privacy_settings"

	// EventPushNameSetting - Push name do usuário foi alterado de outro dispositivo
	// Tipo: *events.PushNameSetting
	EventPushNameSetting WebhookEventType = "push_name_setting"

	// EventUserAbout - Status/sobre de um usuário foi alterado
	// Tipo: *events.UserAbout
	EventUserAbout WebhookEventType = "user_about"

	// EventUserStatusMute - Usuário silenciou ou dessilenciou atualizações de status de outro usuário
	// Tipo: *events.UserStatusMute
	EventUserStatusMute WebhookEventType = "user_status_mute"
)

// ============================================================================
// SINCRONIZAÇÃO E ESTADO DA APLICAÇÃO
// Eventos relacionados à sincronização de dados e estado do app
// ============================================================================

const (
	// EventAppState - Novos dados recebidos da sincronização de app state
	// Tipo: *events.AppState
	// Use eventos de alto nível como Contact e Mute ao invés deste
	EventAppState WebhookEventType = "app_state"

	// EventAppStateSyncComplete - App state foi ressincronizado
	// Tipo: *events.AppStateSyncComplete
	EventAppStateSyncComplete WebhookEventType = "app_state_sync_complete"

	// EventHistorySync - Telefone enviou blob de mensagens históricas
	// Tipo: *events.HistorySync
	EventHistorySync WebhookEventType = "history_sync"

	// EventOfflineSyncCompleted - Servidor terminou de enviar eventos perdidos
	// Tipo: *events.OfflineSyncCompleted
	EventOfflineSyncCompleted WebhookEventType = "offline_sync_completed"

	// EventOfflineSyncPreview - Preview de eventos que serão sincronizados
	// Tipo: *events.OfflineSyncPreview
	// Emitido logo após conectar se houver eventos perdidos
	EventOfflineSyncPreview WebhookEventType = "offline_sync_preview"

	// EventArchive - Chat foi arquivado ou desarquivado de outro dispositivo
	// Tipo: *events.Archive
	EventArchive WebhookEventType = "archive"

	// EventPin - Chat foi fixado ou desfixado de outro dispositivo
	// Tipo: *events.Pin
	EventPin WebhookEventType = "pin"

	// EventMute - Chat foi silenciado ou dessilenciado de outro dispositivo
	// Tipo: *events.Mute
	EventMute WebhookEventType = "mute"

	// EventMarkChatAsRead - Chat inteiro foi marcado como lido/não lido de outro dispositivo
	// Tipo: *events.MarkChatAsRead
	EventMarkChatAsRead WebhookEventType = "mark_chat_as_read"

	// EventDeleteChat - Chat foi deletado de outro dispositivo
	// Tipo: *events.DeleteChat
	EventDeleteChat WebhookEventType = "delete_chat"

	// EventClearChat - Chat foi limpo de outro dispositivo (diferente de deletar)
	// Tipo: *events.ClearChat
	EventClearChat WebhookEventType = "clear_chat"

	// EventStar - Mensagem foi favoritada ou desfavoritada de outro dispositivo
	// Tipo: *events.Star
	EventStar WebhookEventType = "star"

	// EventUnarchiveChatsSetting - Configuração "Manter chats arquivados" foi alterada
	// Tipo: *events.UnarchiveChatsSetting
	EventUnarchiveChatsSetting WebhookEventType = "unarchive_chats_setting"

	// EventLabelEdit - Label foi editada de qualquer dispositivo
	// Tipo: *events.LabelEdit
	EventLabelEdit WebhookEventType = "label_edit"

	// EventLabelAssociationChat - Chat foi etiquetado ou desetiquetado
	// Tipo: *events.LabelAssociationChat
	EventLabelAssociationChat WebhookEventType = "label_association_chat"

	// EventLabelAssociationMessage - Mensagem foi etiquetada ou desetiquetada
	// Tipo: *events.LabelAssociationMessage
	EventLabelAssociationMessage WebhookEventType = "label_association_message"
)

// ============================================================================
// CHAMADAS
// Eventos relacionados a chamadas de voz e vídeo
// ============================================================================

const (
	// EventCallOffer - Usuário recebeu uma chamada no WhatsApp (1:1)
	// Tipo: *events.CallOffer
	EventCallOffer WebhookEventType = "call_offer"

	// EventCallAccept - Chamada foi aceita no WhatsApp
	// Tipo: *events.CallAccept
	EventCallAccept WebhookEventType = "call_accept"

	// EventCallTerminate - Chamada foi terminada
	// Tipo: *events.CallTerminate
	EventCallTerminate WebhookEventType = "call_terminate"

	// EventCallOfferNotice - Notificação de oferta de chamada (principalmente para grupos)
	// Tipo: *events.CallOfferNotice
	EventCallOfferNotice WebhookEventType = "call_offer_notice"

	// EventCallRelayLatency - Latência do relay de chamada
	// Tipo: *events.CallRelayLatency
	EventCallRelayLatency WebhookEventType = "call_relay_latency"

	// EventCallPreAccept - Chamada foi pré-aceita
	// Tipo: *events.CallPreAccept
	EventCallPreAccept WebhookEventType = "call_pre_accept"

	// EventCallReject - Chamada foi rejeitada
	// Tipo: *events.CallReject
	EventCallReject WebhookEventType = "call_reject"

	// EventCallTransport - Transporte de chamada
	// Tipo: *events.CallTransport
	EventCallTransport WebhookEventType = "call_transport"

	// EventUnknownCallEvent - Evento de chamada com conteúdo desconhecido
	// Tipo: *events.UnknownCallEvent
	EventUnknownCallEvent WebhookEventType = "unknown_call_event"
)

// ============================================================================
// PRESENÇA E ATIVIDADE
// Eventos relacionados ao status online/offline e atividade em chats
// ============================================================================

const (
	// EventPresence - Atualização de presença de usuário (online/offline/last seen)
	// Tipo: *events.Presence
	// Requer subscrição: client.SubscribePresence(userJID)
	EventPresence WebhookEventType = "presence"

	// EventChatPresence - Estado de chat (digitando, gravando áudio, pausado)
	// Tipo: *events.ChatPresence
	// Requer que você esteja online: client.SendPresence(types.PresenceAvailable)
	EventChatPresence WebhookEventType = "chat_presence"
)

// ============================================================================
// IDENTIDADE E SEGURANÇA
// Eventos relacionados a mudanças de identidade e segurança
// ============================================================================

const (
	// EventIdentityChange - Outro usuário mudou seu dispositivo principal
	// Tipo: *events.IdentityChange
	EventIdentityChange WebhookEventType = "identity_change"

	// EventCATRefreshError - Erro ao atualizar CAT (Client Access Token)
	// Tipo: *events.CATRefreshError
	EventCATRefreshError WebhookEventType = "cat_refresh_error"
)

// ============================================================================
// NEWSLETTER (CANAIS DO WHATSAPP)
// Eventos relacionados aos Canais do WhatsApp
// ============================================================================

const (
	// EventNewsletterJoin - Usuário entrou em um canal
	// Tipo: *events.NewsletterJoin
	EventNewsletterJoin WebhookEventType = "newsletter_join"

	// EventNewsletterLeave - Usuário saiu de um canal
	// Tipo: *events.NewsletterLeave
	EventNewsletterLeave WebhookEventType = "newsletter_leave"

	// EventNewsletterMuteChange - Mudança de silenciamento de canal
	// Tipo: *events.NewsletterMuteChange
	EventNewsletterMuteChange WebhookEventType = "newsletter_mute_change"

	// EventNewsletterLiveUpdate - Atualização ao vivo de canal
	// Tipo: *events.NewsletterLiveUpdate
	EventNewsletterLiveUpdate WebhookEventType = "newsletter_live_update"
)

// ============================================================================
// FACEBOOK/META BRIDGE
// Eventos de mensagens do Facebook/Instagram via WhatsApp
// ============================================================================

const (
	// EventFBMessage - Mensagem recebida do Facebook/Instagram
	// Tipo: *events.FBMessage
	EventFBMessage WebhookEventType = "fb_message"
)

// ============================================================================
// ESPECIAIS
// Eventos especiais e meta-eventos
// ============================================================================

const (
	// EventAll - Recebe TODOS os eventos (wildcard)
	// Use com cuidado - pode gerar muito tráfego
	EventAll WebhookEventType = "all"
)

// ============================================================================
// LISTAS E VALIDAÇÃO
// Organização e validação de eventos de webhook
// ============================================================================

// AllWebhookEvents contém todos os eventos suportados organizados por categoria
// Baseado na documentação oficial: https://pkg.go.dev/go.mau.fi/whatsmeow/types/events
var AllWebhookEvents = map[string][]WebhookEventType{
	"messages": {
		EventMessage,
		EventUndecryptableMessage,
		EventReceipt,
		EventMediaRetry,
		EventDeleteForMe,
	},
	"groups_contacts": {
		EventGroupInfo,
		EventJoinedGroup,
		EventPicture,
		EventBlocklistChange,
		EventBlocklist,
		EventContact,
		EventPushName,
		EventBusinessName,
	},
	"connection": {
		EventConnected,
		EventDisconnected,
		EventConnectFailure,
		EventKeepAliveRestored,
		EventKeepAliveTimeout,
		EventLoggedOut,
		EventClientOutdated,
		EventTemporaryBan,
		EventStreamError,
		EventStreamReplaced,
		EventPairSuccess,
		EventPairError,
		EventQR,
		EventQRScannedWithoutMultidevice,
		EventManualLoginReconnect,
	},
	"privacy": {
		EventPrivacySettings,
		EventPushNameSetting,
		EventUserAbout,
		EventUserStatusMute,
	},
	"sync": {
		EventAppState,
		EventAppStateSyncComplete,
		EventHistorySync,
		EventOfflineSyncCompleted,
		EventOfflineSyncPreview,
		EventArchive,
		EventPin,
		EventMute,
		EventMarkChatAsRead,
		EventDeleteChat,
		EventClearChat,
		EventStar,
		EventUnarchiveChatsSetting,
		EventLabelEdit,
		EventLabelAssociationChat,
		EventLabelAssociationMessage,
	},
	"calls": {
		EventCallOffer,
		EventCallAccept,
		EventCallTerminate,
		EventCallOfferNotice,
		EventCallRelayLatency,
		EventCallPreAccept,
		EventCallReject,
		EventCallTransport,
		EventUnknownCallEvent,
	},
	"presence": {
		EventPresence,
		EventChatPresence,
	},
	"identity": {
		EventIdentityChange,
		EventCATRefreshError,
	},
	"newsletter": {
		EventNewsletterJoin,
		EventNewsletterLeave,
		EventNewsletterMuteChange,
		EventNewsletterLiveUpdate,
	},
	"facebook": {
		EventFBMessage,
	},
	"special": {
		EventAll,
	},
}

// SupportedEventTypes lista plana de todos os eventos suportados
var SupportedEventTypes []string

// EventTypeMap mapa para validação rápida de eventos
var EventTypeMap map[string]bool

// DefaultWebhookEvents eventos padrão quando nenhum é especificado
// Inclui os eventos mais comuns que a maioria dos usuários precisa
var DefaultWebhookEvents = []string{
	string(EventMessage),           // Mensagens recebidas
	string(EventReceipt),            // Confirmações de entrega/leitura
	string(EventQR),                 // QR Code para pareamento
	string(EventConnected),          // Conexão estabelecida
	string(EventDisconnected),       // Desconectado
	string(EventLoggedOut),          // Deslogado
}

// CriticalEvents eventos críticos relacionados ao estado da conexão
// Estes eventos devem sempre ser monitorados para manter a sessão saudável
var CriticalEvents = []string{
	string(EventConnected),          // Conexão estabelecida
	string(EventDisconnected),       // Desconectado
	string(EventLoggedOut),          // Deslogado do telefone
	string(EventStreamReplaced),     // Login em outro lugar
	string(EventTemporaryBan),       // Conta banida temporariamente
	string(EventClientOutdated),     // Cliente desatualizado
	string(EventKeepAliveTimeout),   // Timeout de keepalive
}

// RecommendedEvents eventos recomendados para a maioria dos casos de uso
var RecommendedEvents = []string{
	// Mensagens
	string(EventMessage),
	string(EventReceipt),

	// Conexão
	string(EventConnected),
	string(EventDisconnected),
	string(EventLoggedOut),
	string(EventQR),

	// Grupos
	string(EventGroupInfo),
	string(EventJoinedGroup),

	// Presença
	string(EventPresence),
	string(EventChatPresence),
}

// MessageEvents eventos relacionados apenas a mensagens
var MessageEvents = []string{
	string(EventMessage),
	string(EventUndecryptableMessage),
	string(EventReceipt),
	string(EventMediaRetry),
	string(EventDeleteForMe),
}

// ConnectionEvents eventos relacionados apenas a conexão
var ConnectionEvents = []string{
	string(EventConnected),
	string(EventDisconnected),
	string(EventConnectFailure),
	string(EventKeepAliveRestored),
	string(EventKeepAliveTimeout),
	string(EventLoggedOut),
	string(EventClientOutdated),
	string(EventTemporaryBan),
	string(EventStreamError),
	string(EventStreamReplaced),
	string(EventPairSuccess),
	string(EventPairError),
	string(EventQR),
	string(EventQRScannedWithoutMultidevice),
	string(EventManualLoginReconnect),
}

func init() {
	// Inicializar lista plana de eventos
	SupportedEventTypes = make([]string, 0)
	EventTypeMap = make(map[string]bool)

	// Percorrer todas as categorias e adicionar eventos
	for _, events := range AllWebhookEvents {
		for _, event := range events {
			eventStr := string(event)
			SupportedEventTypes = append(SupportedEventTypes, eventStr)
			EventTypeMap[eventStr] = true
		}
	}
}

// IsValidEventType verifica se um tipo de evento é válido
func IsValidEventType(eventType string) bool {
	return EventTypeMap[eventType]
}

// GetEventsByCategory retorna eventos de uma categoria específica
func GetEventsByCategory(category string) []WebhookEventType {
	if events, ok := AllWebhookEvents[category]; ok {
		return events
	}
	return []WebhookEventType{}
}

// GetAllCategories retorna todas as categorias de eventos
func GetAllCategories() []string {
	categories := make([]string, 0, len(AllWebhookEvents))
	for category := range AllWebhookEvents {
		categories = append(categories, category)
	}
	return categories
}

// IsCriticalEvent verifica se um evento é crítico
func IsCriticalEvent(eventType string) bool {
	for _, critical := range CriticalEvents {
		if critical == eventType {
			return true
		}
	}
	return false
}

// IsMessageEvent verifica se um evento é relacionado a mensagens
func IsMessageEvent(eventType string) bool {
	for _, msgEvent := range MessageEvents {
		if msgEvent == eventType {
			return true
		}
	}
	return false
}

// IsConnectionEvent verifica se um evento é relacionado a conexão
func IsConnectionEvent(eventType string) bool {
	for _, connEvent := range ConnectionEvents {
		if connEvent == eventType {
			return true
		}
	}
	return false
}

// GetEventDescription retorna uma descrição amigável do evento
func GetEventDescription(eventType string) string {
	descriptions := map[string]string{
		// Messages
		string(EventMessage):               "Mensagem recebida (texto, mídia, documentos, etc)",
		string(EventUndecryptableMessage):  "Mensagem que não pôde ser descriptografada",
		string(EventReceipt):               "Confirmação de entrega/leitura de mensagem",
		string(EventMediaRetry):            "Resposta a solicitação de reenvio de mídia",
		string(EventDeleteForMe):           "Mensagem deletada apenas para o usuário",

		// Groups & Contacts
		string(EventGroupInfo):             "Metadados de grupo alterados",
		string(EventJoinedGroup):           "Entrou ou foi adicionado a um grupo",
		string(EventPicture):               "Foto de perfil alterada",
		string(EventBlocklistChange):       "Mudança individual na lista de bloqueados",
		string(EventBlocklist):             "Lista completa de bloqueados atualizada",
		string(EventContact):               "Entrada na lista de contatos modificada",
		string(EventPushName):              "Nome de exibição de contato mudou",
		string(EventBusinessName):          "Nome comercial verificado mudou",

		// Connection
		string(EventConnected):             "Conectado e autenticado com sucesso",
		string(EventDisconnected):          "WebSocket fechado pelo servidor",
		string(EventConnectFailure):        "Servidor rejeitou a conexão",
		string(EventKeepAliveRestored):     "Pings keepalive restaurados",
		string(EventKeepAliveTimeout):      "Ping keepalive expirou",
		string(EventLoggedOut):             "Desconectado do telefone",
		string(EventClientOutdated):        "Cliente desatualizado",
		string(EventTemporaryBan):          "Conta temporariamente banida",
		string(EventStreamError):           "Erro de stream desconhecido",
		string(EventStreamReplaced):        "Login em outro dispositivo",
		string(EventPairSuccess):           "QR code escaneado com sucesso",
		string(EventPairError):             "Erro no pareamento",
		string(EventQR):                    "QR code gerado",
		string(EventQRScannedWithoutMultidevice): "QR escaneado sem multidevice",
		string(EventManualLoginReconnect):  "Reconexão manual necessária após login",

		// Privacy
		string(EventPrivacySettings):       "Configurações de privacidade alteradas",
		string(EventPushNameSetting):       "Push name alterado de outro dispositivo",
		string(EventUserAbout):             "Status/sobre de usuário alterado",
		string(EventUserStatusMute):        "Atualizações de status silenciadas/dessilenciadas",

		// Calls
		string(EventCallOffer):             "Chamada recebida",
		string(EventCallAccept):            "Chamada aceita",
		string(EventCallTerminate):         "Chamada terminada",
		string(EventCallOfferNotice):       "Notificação de oferta de chamada",
		string(EventCallRelayLatency):      "Latência do relay de chamada",

		// Presence
		string(EventPresence):              "Atualização de presença (online/offline)",
		string(EventChatPresence):          "Estado de chat (digitando, gravando)",

		// Special
		string(EventAll):                   "Recebe TODOS os eventos",
	}

	if desc, ok := descriptions[eventType]; ok {
		return desc
	}
	return "Evento sem descrição disponível"
}

// ValidateEventList valida uma lista de eventos e retorna eventos inválidos
func ValidateEventList(events []string) (valid []string, invalid []string) {
	valid = make([]string, 0)
	invalid = make([]string, 0)

	for _, event := range events {
		if IsValidEventType(event) {
			valid = append(valid, event)
		} else {
			invalid = append(invalid, event)
		}
	}

	return valid, invalid
}

// GetEventCategory retorna a categoria de um evento
func GetEventCategory(eventType string) string {
	for category, events := range AllWebhookEvents {
		for _, event := range events {
			if string(event) == eventType {
				return category
			}
		}
	}
	return "unknown"
}

