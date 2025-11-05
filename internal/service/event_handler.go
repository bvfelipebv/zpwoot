package service

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"zpmeow/internal/repository"
	"zpmeow/pkg/logger"
)

// EventHandler gerencia eventos do WhatsApp
type EventHandler struct {
	manager     *SessionManager
	sessionRepo *repository.SessionRepository
}

// NewEventHandler cria um novo handler de eventos
func NewEventHandler(manager *SessionManager, sessionRepo *repository.SessionRepository) *EventHandler {
	return &EventHandler{
		manager:     manager,
		sessionRepo: sessionRepo,
	}
}

// RegisterHandlers registra os handlers de eventos para um cliente
func (h *EventHandler) RegisterHandlers(client *whatsmeow.Client, sessionID string) {
	client.AddEventHandler(func(evt interface{}) {
		h.handleEvent(sessionID, evt)
	})
}

// handleEvent processa eventos do WhatsApp
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

// handleConnected processa evento de conexão
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

// handleDisconnected processa evento de desconexão
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

// handleLoggedOut processa evento de logout
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

// handleMessage processa mensagens recebidas
func (h *EventHandler) handleMessage(sessionID string, evt *events.Message) {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Str("from", evt.Info.Sender.String()).
		Str("message_id", evt.Info.ID).
		Msg("Message received")
	
	// TODO: Processar mensagem e enviar webhook
}

// handleReceipt processa recibos de mensagem
func (h *EventHandler) handleReceipt(sessionID string, evt *events.Receipt) {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Str("type", string(evt.Type)).
		Msg("Receipt received")
	
	// TODO: Processar recibo e enviar webhook
}

// handlePresence processa eventos de presença
func (h *EventHandler) handlePresence(sessionID string, evt *events.Presence) {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Str("from", evt.From.String()).
		Msg("Presence update")
	
	// TODO: Processar presença e enviar webhook
}

// handleHistorySync processa sincronização de histórico
func (h *EventHandler) handleHistorySync(sessionID string, evt *events.HistorySync) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Str("type", evt.Data.SyncType.String()).
		Msg("History sync")
	
	// TODO: Processar histórico
}

// handlePushName processa atualização de nome
func (h *EventHandler) handlePushName(sessionID string, evt *events.PushName) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Str("jid", evt.JID.String()).
		Str("push_name", evt.NewPushName).
		Msg("Push name updated")
	
	// TODO: Atualizar push name no banco
}

