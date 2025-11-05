package service

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"zpwoot/internal/repository"
	"zpwoot/pkg/logger"
)

type EventHandler struct {
	manager     *SessionManager
	sessionRepo *repository.SessionRepository
}

func NewEventHandler(manager *SessionManager, sessionRepo *repository.SessionRepository) *EventHandler {
	return &EventHandler{
		manager:     manager,
		sessionRepo: sessionRepo,
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

func (h *EventHandler) handleConnected(sessionID string, _ *events.Connected) {
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

	// TODO: Enviar webhook de conexão
}

func (h *EventHandler) handleDisconnected(sessionID string, _ *events.Disconnected) {
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

	// TODO: Enviar webhook de desconexão
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

	// TODO: Processar mensagem e enviar webhook
}

func (h *EventHandler) handleReceipt(sessionID string, evt *events.Receipt) {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Str("type", string(evt.Type)).
		Msg("Receipt received")

	// TODO: Processar recibo e enviar webhook
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
