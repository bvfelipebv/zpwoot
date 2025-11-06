package service

import (
	"context"

	"go.mau.fi/whatsmeow"
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
	case *events.Connected:
		h.handleConnected(sessionID, v)
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

func (h *EventHandler) handleConnected(sessionID string, evt *events.Connected) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("WhatsApp connected")

	ctx := context.Background()

	// Atualizar status no banco
	if err := h.sessionRepo.UpdateStatus(ctx, sessionID, "connected", true); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to update session status to connected")
	}

	// Enviar webhook de conexão
	payload := h.webhookFormatter.FormatConnected(sessionID, evt)
	if err := h.webhookProcessor.ProcessEvent(sessionID, constants.EventConnected, payload); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to process connected webhook")
	}
}

func (h *EventHandler) handleDisconnected(sessionID string, evt *events.Disconnected) {
	logger.Log.Warn().
		Str("session_id", sessionID).
		Msg("WhatsApp disconnected")

	ctx := context.Background()

	// Atualizar status no banco
	if err := h.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to update session status to disconnected")
	}

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
	if err := h.sessionRepo.UpdateStatus(ctx, sessionID, "logged_out", false); err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to update session status to logged_out")
	}

	// Limpar device JID
	session, err := h.sessionRepo.GetByID(ctx, sessionID)
	if err == nil {
		session.DeviceJID = ""
		session.QRCode = ""
		h.sessionRepo.Update(ctx, session)
	}

	// TODO: Enviar webhook de logout
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
