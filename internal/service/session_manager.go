package service

import (
	"context"
	"fmt"
	"sync"

	"go.mau.fi/whatsmeow"

	"zpmeow/internal/model"
	"zpmeow/internal/repository"
	"zpmeow/pkg/logger"
)

// SessionManager gerencia sessões WhatsApp ativas
type SessionManager struct {
	whatsappSvc *WhatsAppService
	sessionRepo *repository.SessionRepository

	// Map de clientes ativos: sessionID -> *whatsmeow.Client
	clients     map[string]*whatsmeow.Client
	clientsMux  sync.RWMutex

	// Event handler
	eventHandler *EventHandler
}

// NewSessionManager cria um novo gerenciador de sessões
func NewSessionManager(whatsappSvc *WhatsAppService, sessionRepo *repository.SessionRepository) *SessionManager {
	manager := &SessionManager{
		whatsappSvc: whatsappSvc,
		sessionRepo: sessionRepo,
		clients:     make(map[string]*whatsmeow.Client),
	}

	// Criar event handler
	manager.eventHandler = NewEventHandler(manager, sessionRepo)

	return manager
}

// CreateSession cria uma nova sessão no banco
func (m *SessionManager) CreateSession(ctx context.Context, name, webhookURL string) (*model.Session, error) {
	session := &model.Session{
		Name:       name,
		Status:     "disconnected",
		Connected:  false,
		WebhookURL: webhookURL,
	}

	if err := m.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	logger.Log.Info().
		Str("session_id", session.ID).
		Str("name", name).
		Msg("Session created")

	return session, nil
}

// GetSession busca uma sessão por ID
func (m *SessionManager) GetSession(ctx context.Context, sessionID string) (*model.Session, error) {
	return m.sessionRepo.GetByID(ctx, sessionID)
}

// ListSessions lista todas as sessões
func (m *SessionManager) ListSessions(ctx context.Context) ([]*model.Session, error) {
	return m.sessionRepo.List(ctx)
}

// DeleteSession deleta uma sessão e desconecta se estiver ativa
func (m *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	// Desconectar se estiver ativo
	if err := m.DisconnectSession(ctx, sessionID); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to disconnect session before delete")
	}

	// Deletar do banco
	if err := m.sessionRepo.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("Session deleted")

	return nil
}

// ConnectSession conecta uma sessão ao WhatsApp
func (m *SessionManager) ConnectSession(ctx context.Context, sessionID string) error {
	// Buscar sessão
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Verificar se já está conectado
	if m.IsClientActive(sessionID) {
		return fmt.Errorf("session already connected")
	}

	// Verificar se tem device JID (já foi pareado)
	if session.DeviceJID == "" {
		return fmt.Errorf("session not paired yet")
	}

	// Criar ou obter device do whatsmeow
	device, err := m.whatsappSvc.GetOrCreateDevice(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Criar cliente WhatsApp
	client := m.whatsappSvc.NewClient(device)

	// Registrar event handlers
	m.eventHandler.RegisterHandlers(client, sessionID)

	// Conectar
	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Adicionar ao map de clientes ativos
	m.clientsMux.Lock()
	m.clients[sessionID] = client
	m.clientsMux.Unlock()

	// Atualizar status no banco
	if err := m.sessionRepo.UpdateStatus(ctx, sessionID, "connecting", false); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to update session status")
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("Session connecting")

	return nil
}

// DisconnectSession desconecta uma sessão
func (m *SessionManager) DisconnectSession(ctx context.Context, sessionID string) error {
	m.clientsMux.Lock()
	client, exists := m.clients[sessionID]
	if exists {
		delete(m.clients, sessionID)
	}
	m.clientsMux.Unlock()

	if !exists {
		return fmt.Errorf("session not connected")
	}

	// Desconectar cliente
	client.Disconnect()

	// Atualizar status no banco
	if err := m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to update session status")
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("Session disconnected")

	return nil
}

// GetClient retorna o cliente WhatsApp ativo para uma sessão
func (m *SessionManager) GetClient(sessionID string) (*whatsmeow.Client, error) {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	client, exists := m.clients[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not connected")
	}

	return client, nil
}

// IsClientActive verifica se um cliente está ativo
func (m *SessionManager) IsClientActive(sessionID string) bool {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	_, exists := m.clients[sessionID]
	return exists
}

// GetSessionStatus retorna o status detalhado de uma sessão
func (m *SessionManager) GetSessionStatus(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	// Buscar sessão do banco
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Verificar se tem cliente ativo
	isActive := m.IsClientActive(sessionID)

	status := map[string]interface{}{
		"id":            session.ID,
		"name":          session.Name,
		"status":        session.Status,
		"connected":     session.Connected,
		"device_jid":    session.DeviceJID,
		"is_active":     isActive,
		"needs_pairing": session.NeedsPairing(),
		"can_connect":   session.CanConnect(),
	}

	// Se tem cliente ativo, adicionar informações do cliente
	if isActive {
		client, _ := m.GetClient(sessionID)
		if client != nil {
			status["is_logged_in"] = client.IsLoggedIn()
			status["is_connected"] = client.IsConnected()
			if client.Store != nil && client.Store.ID != nil {
				status["push_name"] = client.Store.PushName
				status["platform"] = client.Store.Platform
			}
		}
	}

	return status, nil
}

// UpdateWebhook atualiza as configurações de webhook de uma sessão
func (m *SessionManager) UpdateWebhook(ctx context.Context, sessionID, webhookURL, webhookEvents string) error {
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	session.WebhookURL = webhookURL
	session.WebhookEvents = webhookEvents

	if err := m.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Str("webhook_url", webhookURL).
		Msg("Webhook updated")

	return nil
}

// RestoreAllSessions restaura todas as sessões conectadas ao iniciar
func (m *SessionManager) RestoreAllSessions(ctx context.Context) error {
	// Buscar sessões conectadas
	sessions, err := m.sessionRepo.ListConnected(ctx)
	if err != nil {
		return fmt.Errorf("failed to list connected sessions: %w", err)
	}

	logger.Log.Info().
		Int("count", len(sessions)).
		Msg("Restoring connected sessions")

	// Reconectar cada sessão
	for _, session := range sessions {
		if err := m.ConnectSession(ctx, session.ID); err != nil {
			logger.Log.Error().
				Err(err).
				Str("session_id", session.ID).
				Msg("Failed to restore session")
			continue
		}

		logger.Log.Info().
			Str("session_id", session.ID).
			Str("name", session.Name).
			Msg("Session restored")
	}

	return nil
}

// Shutdown desconecta todas as sessões ativas
func (m *SessionManager) Shutdown(ctx context.Context) error {
	m.clientsMux.Lock()
	defer m.clientsMux.Unlock()

	logger.Log.Info().
		Int("count", len(m.clients)).
		Msg("Shutting down all sessions")

	// Desconectar todos os clientes
	for sessionID, client := range m.clients {
		client.Disconnect()

		// Atualizar status no banco
		if err := m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false); err != nil {
			logger.Log.Warn().
				Err(err).
				Str("session_id", sessionID).
				Msg("Failed to update session status on shutdown")
		}

		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("Session disconnected on shutdown")
	}

	// Limpar map
	m.clients = make(map[string]*whatsmeow.Client)

	return nil
}

// GetActiveSessionsCount retorna o número de sessões ativas
func (m *SessionManager) GetActiveSessionsCount() int {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()
	return len(m.clients)
}

