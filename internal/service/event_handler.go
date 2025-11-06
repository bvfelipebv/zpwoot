package service

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"zpwoot/internal/constants"
	"zpwoot/internal/repository"
	"zpwoot/pkg/logger"
)

type EventHandler struct {
	manager          *SessionManager
	sessionRepo      *repository.SessionRepository
	webhookProcessor *WebhookProcessor
	webhookFormatter *WebhookFormatter
}

func NewEventHandler(
	manager *SessionManager,
	sessionRepo *repository.SessionRepository,
	webhookProcessor *WebhookProcessor,
	webhookFormatter *WebhookFormatter,
) *EventHandler {
	return &EventHandler{
		manager:          manager,
		sessionRepo:      sessionRepo,
		webhookProcessor: webhookProcessor,
		webhookFormatter: webhookFormatter,
	}
}

func (h *EventHandler) RegisterHandlers(client *whatsmeow.Client, sessionID string) {
	client.AddEventHandler(func(evt interface{}) {
		h.handleEvent(sessionID, evt)
	})
}

func (h *EventHandler) handleEvent(sessionID string, evt interface{}) {
	switch v := evt.(type) {
	case *events.AppStateSyncComplete:
		h.handleAppStateSyncComplete(sessionID, v)
	case *events.PairSuccess:
		h.handlePairSuccess(sessionID, v)
	case *events.Connected:
		h.handleConnected(sessionID, v)
	case *events.PushNameSetting:
		h.handlePushNameSetting(sessionID, v)
	case *events.Disconnected:
		h.handleDisconnected(sessionID, v)
	case *events.LoggedOut:
		h.handleLoggedOut(sessionID, v)
	case *events.Message:
		h.handleMessage(sessionID, v)
	case *events.Receipt:
		h.handleReceipt(sessionID, v)
	case *events.Presence:
		h.handlePresence(sessionID, v)
	case *events.HistorySync:
		h.handleHistorySync(sessionID, v)
	case *events.PushName:
		h.handlePushName(sessionID, v)
	}
}

// sendAvailablePresence envia presença disponível
func (h *EventHandler) sendAvailablePresence(ctx context.Context, sessionID string, client *whatsmeow.Client) {
	if err := client.SendPresence(ctx, types.PresenceAvailable); err != nil {
		logger.Log.Warn().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to send available presence")
	} else {
		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("Marked self as available")
	}
}

// saveDeviceJID salva o device JID no banco de dados
func (h *EventHandler) saveDeviceJID(ctx context.Context, sessionID string, deviceJID string) {
	if deviceJID == "" {
		return
	}

	if err := h.sessionRepo.UpdateDeviceJID(ctx, sessionID, deviceJID); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("device_jid", deviceJID).
			Msg("Failed to save device JID")
	} else {
		logger.Log.Info().
			Str("session_id", sessionID).
			Str("device_jid", deviceJID).
			Msg("Device JID saved successfully")
	}
}

// updateSessionStatus atualiza o status da sessão no banco
func (h *EventHandler) updateSessionStatus(ctx context.Context, sessionID string, status string, connected bool) {
	if err := h.sessionRepo.UpdateStatus(ctx, sessionID, status, connected); err != nil {
		logger.Log.Warn().
			Err(err).
			Str("session_id", sessionID).
			Str("status", status).
			Msg("Failed to update session status")
	}
}

func (h *EventHandler) handleAppStateSyncComplete(sessionID string, evt *events.AppStateSyncComplete) {
	ctx := context.Background()

	client, err := h.manager.GetClient(sessionID)
	if err != nil {
		return
	}

	// CRÍTICO (do wuzapi): Enviar presença disponível SOMENTE se tiver PushName E for WAPatchCriticalBlock
	// Isso é ESSENCIAL para o WhatsApp não deslogar a sessão
	if len(client.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
		h.sendAvailablePresence(ctx, sessionID, client)
	}
}

func (h *EventHandler) handlePushNameSetting(sessionID string, evt *events.PushNameSetting) {
	h.handleConnectedOrPushName(sessionID, evt)
}

func (h *EventHandler) handlePairSuccess(sessionID string, evt *events.PairSuccess) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Str("jid", evt.ID.String()).
		Str("business_name", evt.BusinessName).
		Str("platform", evt.Platform).
		Msg("QR Pair Success")

	ctx := context.Background()

	// Salvar device_jid no banco IMEDIATAMENTE após pairing
	h.saveDeviceJID(ctx, sessionID, evt.ID.String())

	// Limpar QR code
	if err := h.sessionRepo.UpdateQRCode(ctx, sessionID, ""); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to clear QR code after pairing")
	}

	// Atualizar status para connected
	h.updateSessionStatus(ctx, sessionID, "connected", true)
}

func (h *EventHandler) handleConnected(sessionID string, evt *events.Connected) {
	h.handleConnectedOrPushName(sessionID, evt)
}

// handleConnectedOrPushName trata eventos Connected e PushNameSetting (como no wuzapi)
func (h *EventHandler) handleConnectedOrPushName(sessionID string, evt interface{}) {
	ctx := context.Background()

	client, err := h.manager.GetClient(sessionID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get client")
		return
	}

	// Log do evento
	switch evt.(type) {
	case *events.Connected:
		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("WhatsApp connected")
	case *events.PushNameSetting:
		logger.Log.Info().
			Str("session_id", sessionID).
			Str("push_name", client.Store.PushName).
			Msg("Push name setting received")
	}

	// CRÍTICO (do wuzapi): Só enviar presença se PushName estiver configurado
	if len(client.Store.PushName) == 0 {
		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("Waiting for PushName before sending presence")
		return
	}

	// Enviar presença disponível
	h.sendAvailablePresence(ctx, sessionID, client)

	// Salvar device_jid no banco
	if client.Store != nil && client.Store.ID != nil {
		h.saveDeviceJID(ctx, sessionID, client.Store.ID.String())
	}

	// Atualizar status para connected
	h.updateSessionStatus(ctx, sessionID, "connected", true)

	// Enviar webhook de conexão (apenas para Connected)
	if _, ok := evt.(*events.Connected); ok {
		payload := h.webhookFormatter.FormatConnected(sessionID, evt.(*events.Connected))
		if err := h.webhookProcessor.ProcessEvent(sessionID, constants.EventConnected, payload); err != nil {
			logger.Log.Error().
				Err(err).
				Str("session_id", sessionID).
				Msg("Failed to process connected webhook")
		}
	}
}

func (h *EventHandler) handleDisconnected(sessionID string, evt *events.Disconnected) {
	logger.Log.Warn().
		Str("session_id", sessionID).
		Msg("WhatsApp disconnected")

	ctx := context.Background()

	// Atualizar status no banco
	h.updateSessionStatus(ctx, sessionID, "disconnected", false)

	// Enviar webhook de desconexão
	payload := h.webhookFormatter.FormatDisconnected(sessionID, evt)
	if err := h.webhookProcessor.ProcessEvent(sessionID, constants.EventDisconnected, payload); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to process disconnected webhook")
	}
}

func (h *EventHandler) handleLoggedOut(sessionID string, evt *events.LoggedOut) {
	logger.Log.Warn().
		Str("session_id", sessionID).
		Str("reason", evt.Reason.String()).
		Msg("WhatsApp logged out")

	ctx := context.Background()

	// Atualizar status no banco
	h.updateSessionStatus(ctx, sessionID, "disconnected", false)

	// Limpar QR code se existir
	if err := h.sessionRepo.UpdateQRCode(ctx, sessionID, ""); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to clear QR code after logout")
	}

	// Enviar webhook de logout
	payload := h.webhookFormatter.FormatDisconnected(sessionID, &events.Disconnected{})
	if err := h.webhookProcessor.ProcessEvent(sessionID, constants.EventDisconnected, payload); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to process logout webhook")
	}
}

func (h *EventHandler) handleMessage(sessionID string, evt *events.Message) {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Str("from", evt.Info.Sender.String()).
		Str("message_id", evt.Info.ID).
		Msg("Message received")

	// Enviar webhook de mensagem
	payload := h.webhookFormatter.FormatMessage(sessionID, evt)
	if err := h.webhookProcessor.ProcessEvent(sessionID, constants.EventMessage, payload); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to process message webhook")
	}
}

func (h *EventHandler) handleReceipt(sessionID string, evt *events.Receipt) {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Str("type", string(evt.Type)).
		Msg("Receipt received")

	// Enviar webhook de recibo
	payload := h.webhookFormatter.FormatReceipt(sessionID, evt)
	if err := h.webhookProcessor.ProcessEvent(sessionID, constants.EventReceipt, payload); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to process receipt webhook")
	}
}

func (h *EventHandler) handlePresence(sessionID string, evt *events.Presence) {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Str("from", evt.From.String()).
		Msg("Presence update")

	// TODO: Processar presença e enviar webhook
}

func (h *EventHandler) handleHistorySync(sessionID string, evt *events.HistorySync) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Str("type", evt.Data.SyncType.String()).
		Msg("History sync")

	// TODO: Processar histórico
}

func (h *EventHandler) handlePushName(sessionID string, evt *events.PushName) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Str("jid", evt.JID.String()).
		Str("push_name", evt.NewPushName).
		Msg("Push name updated")

	// TODO: Atualizar push name no banco
}
