package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"

	"zpwoot/internal/model"
	"zpwoot/internal/repository"
	"zpwoot/pkg/logger"
)

type SessionManager struct {
	whatsappSvc *WhatsAppService
	sessionRepo *repository.SessionRepository

	// Map de clientes ativos: sessionID -> *whatsmeow.Client
	clients    map[string]*whatsmeow.Client
	clientsMux sync.RWMutex

	// Map de QR codes em memória: sessionID -> QR code string
	qrCodes    map[string]string
	qrCodesMux sync.RWMutex

	// Event handler
	eventHandler *EventHandler
}

func NewSessionManager(
	whatsappSvc *WhatsAppService,
	sessionRepo *repository.SessionRepository,
	webhookProcessor *WebhookProcessor,
	webhookFormatter *WebhookFormatter,
) *SessionManager {
	manager := &SessionManager{
		whatsappSvc: whatsappSvc,
		sessionRepo: sessionRepo,
		clients:     make(map[string]*whatsmeow.Client),
		qrCodes:     make(map[string]string),
	}

	// Criar event handler com webhook support
	manager.eventHandler = NewEventHandler(manager, sessionRepo, webhookProcessor, webhookFormatter)

	return manager
}

func (m *SessionManager) CreateSession(ctx context.Context, name, webhookURL string) (*model.Session, error) {
	session := &model.Session{
		Name:      name,
		Status:    "disconnected",
		Connected: false,
	}

	if webhookURL != "" {
		session.WebhookConfig = &model.WebhookConfig{
			Enabled: true,
			URL:     webhookURL,
			Events:  []string{"message", "qr", "connected", "disconnected"},
		}
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

func (m *SessionManager) CreateSessionWithConfig(ctx context.Context, session *model.Session) error {
	if err := m.sessionRepo.Create(ctx, session); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	logger.Log.Info().
		Str("session_id", session.ID).
		Str("name", session.Name).
		Msg("Session created with config")

	return nil
}

func (m *SessionManager) GetSession(ctx context.Context, sessionID string) (*model.Session, error) {
	return m.sessionRepo.GetByID(ctx, sessionID)
}

func (m *SessionManager) ListSessions(ctx context.Context) ([]*model.Session, error) {
	return m.sessionRepo.List(ctx)
}

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

func (m *SessionManager) ConnectSession(ctx context.Context, sessionID string) error {
	// Buscar sessão
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Verificar se já está conectado
	if m.IsClientActive(sessionID) {
		// Desconectar primeiro para reiniciar o processo
		m.DisconnectSession(ctx, sessionID)
	}

	// Limpar QR code anterior da memória
	m.qrCodesMux.Lock()
	delete(m.qrCodes, sessionID)
	m.qrCodesMux.Unlock()

	// Atualizar status para connecting
	if err := m.sessionRepo.UpdateStatus(ctx, sessionID, "connecting", false); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to update session status")
	}

	// Criar ou obter device do whatsmeow (passar JID se existir)
	device, err := m.whatsappSvc.GetOrCreateDevice(ctx, sessionID, session.DeviceJID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Criar cliente WhatsApp
	client := m.whatsappSvc.NewClient(device)

	// Registrar event handlers
	m.eventHandler.RegisterHandlers(client, sessionID)

	// Adicionar ao map de clientes ativos
	m.clientsMux.Lock()
	m.clients[sessionID] = client
	m.clientsMux.Unlock()

	// Se não está pareado, iniciar processo de QR Code
	if session.DeviceJID == "" {
		// IMPORTANTE: Usar context.Background() para não cancelar quando a requisição HTTP terminar
		qrChan, err := client.GetQRChannel(context.Background())
		if err != nil {
			// Se o erro for ErrQRStoreContainsID, significa que já está logado
			if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
				logger.Log.Info().
					Str("session_id", sessionID).
					Msg("Device already contains ID, connecting...")

				if err := client.Connect(); err != nil {
					m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
					return fmt.Errorf("failed to connect: %w", err)
				}

				m.sessionRepo.UpdateStatus(ctx, sessionID, "connected", true)
				return nil
			}

			m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
			return fmt.Errorf("failed to get QR channel: %w", err)
		}

		// Conectar DEPOIS de obter o canal (como na wuzapi)
		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("Connecting to WhatsApp to generate QR codes...")

		if err := client.Connect(); err != nil {
			m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
			return fmt.Errorf("failed to connect: %w", err)
		}

		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("Connected successfully - waiting for QR codes...")

		// Processar QR codes em goroutine (usar context.Background para não cancelar)
		go m.handleQRCodes(context.Background(), sessionID, qrChan, client)

		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("QR code handler started - scan QR code to connect")
	} else {
		// Já está pareado, apenas conectar
		if err := client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("Session connecting - already paired")
	}

	return nil
}

// handleQRCodes processa os QR codes gerados pelo whatsmeow
// Baseado no código da wuzapi: https://github.com/asternic/wuzapi
func (m *SessionManager) handleQRCodes(ctx context.Context, sessionID string, qrChan <-chan whatsmeow.QRChannelItem, client *whatsmeow.Client) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("Starting QR code handler - waiting for events...")

	// Processar TODOS os eventos do canal
	for evt := range qrChan {
		logger.Log.Info().
			Str("session_id", sessionID).
			Str("event", evt.Event).
			Msg("QR channel event received")

		if evt.Event == "code" {
			// Novo QR code gerado
			// Salvar QR code em memória
			m.qrCodesMux.Lock()
			m.qrCodes[sessionID] = evt.Code
			m.qrCodesMux.Unlock()

			// Salvar QR code no banco
			if err := m.sessionRepo.UpdateQRCode(ctx, sessionID, evt.Code); err != nil {
				logger.Log.Error().
					Err(err).
					Str("session_id", sessionID).
					Msg("Failed to save QR code to database")
			}

			// Exibir QR code no terminal
			logger.Log.Info().
				Str("session_id", sessionID).
				Dur("timeout", evt.Timeout).
				Msg("QR code generated - displaying in terminal")

			fmt.Printf("\n=== QR Code for Session: %s (timeout: %v) ===\n", sessionID, evt.Timeout)
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Printf("=== Scan with WhatsApp to connect ===\n\n")

			// Atualizar status para qr_code
			if err := m.sessionRepo.UpdateStatus(ctx, sessionID, "qr_code", false); err != nil {
				logger.Log.Warn().Err(err).Msg("Failed to update session status")
			}

		} else if evt.Event == "timeout" {
			// QR code expirou sem ser escaneado
			logger.Log.Warn().
				Str("session_id", sessionID).
				Msg("QR code timeout - no scan detected")

			// Limpar QR code da memória e banco
			m.qrCodesMux.Lock()
			delete(m.qrCodes, sessionID)
			m.qrCodesMux.Unlock()

			m.sessionRepo.UpdateQRCode(ctx, sessionID, "")
			m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)

			// Remover cliente
			m.clientsMux.Lock()
			delete(m.clients, sessionID)
			m.clientsMux.Unlock()

			return

		} else if evt.Event == "success" {
			// QR code escaneado com sucesso!
			logger.Log.Info().
				Str("session_id", sessionID).
				Msg("QR code scanned successfully - pairing completed")

			// Limpar QR code da memória e banco
			m.qrCodesMux.Lock()
			delete(m.qrCodes, sessionID)
			m.qrCodesMux.Unlock()

			m.sessionRepo.UpdateQRCode(ctx, sessionID, "")

			// Atualizar status para connected (será confirmado pelo event handler)
			m.sessionRepo.UpdateStatus(ctx, sessionID, "connected", true)

			// Salvar JID no banco
			if client.Store.ID != nil {
				jid := client.Store.ID.String()
				m.sessionRepo.UpdateDeviceJID(ctx, sessionID, jid)

				logger.Log.Info().
					Str("session_id", sessionID).
					Str("jid", jid).
					Msg("Device JID saved")
			}

			return

		} else if evt.Event == "error" {
			// Erro durante pareamento
			logger.Log.Error().
				Err(evt.Error).
				Str("session_id", sessionID).
				Msg("QR pairing error")

			// Limpar e desconectar
			m.qrCodesMux.Lock()
			delete(m.qrCodes, sessionID)
			m.qrCodesMux.Unlock()

			m.sessionRepo.UpdateQRCode(ctx, sessionID, "")
			m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)

			m.clientsMux.Lock()
			delete(m.clients, sessionID)
			m.clientsMux.Unlock()

			if client != nil {
				client.Disconnect()
			}

			return

		} else {
			// Outros eventos (err-unexpected-state, err-client-outdated, etc)
			logger.Log.Warn().
				Str("session_id", sessionID).
				Str("event", evt.Event).
				Msg("QR channel event")

			// Para eventos de erro, limpar e desconectar
			if evt.Event == "err-unexpected-state" || evt.Event == "err-client-outdated" || evt.Event == "err-scanned-without-multidevice" {
				m.qrCodesMux.Lock()
				delete(m.qrCodes, sessionID)
				m.qrCodesMux.Unlock()

				m.sessionRepo.UpdateQRCode(ctx, sessionID, "")
				m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)

				m.clientsMux.Lock()
				delete(m.clients, sessionID)
				m.clientsMux.Unlock()

				if client != nil {
					client.Disconnect()
				}

				return
			}
		}
	}

	// Canal fechado
	logger.Log.Warn().
		Str("session_id", sessionID).
		Bool("client_connected", client.IsConnected()).
		Msg("QR channel closed unexpectedly")

	// Se o cliente ainda está conectado, algo deu errado
	if client.IsConnected() {
		logger.Log.Info().
			Str("session_id", sessionID).
			Msg("Client still connected, waiting for pairing...")
		// Não fazer nada, deixar o cliente conectado
	} else {
		logger.Log.Warn().
			Str("session_id", sessionID).
			Msg("Client disconnected, updating status")

		// Limpar QR code
		m.qrCodesMux.Lock()
		delete(m.qrCodes, sessionID)
		m.qrCodesMux.Unlock()

		m.sessionRepo.UpdateQRCode(ctx, sessionID, "")
		m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)

		// Remover cliente
		m.clientsMux.Lock()
		delete(m.clients, sessionID)
		m.clientsMux.Unlock()
	}
}

// GetQRCode retorna o QR code atual da memória
func (m *SessionManager) GetQRCode(sessionID string) (string, bool) {
	m.qrCodesMux.RLock()
	defer m.qrCodesMux.RUnlock()

	qrCode, exists := m.qrCodes[sessionID]
	return qrCode, exists
}

// GetClient retorna o cliente WhatsApp de uma sessão
func (m *SessionManager) GetClient(sessionID string) (*whatsmeow.Client, error) {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	client, exists := m.clients[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found or not connected")
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	return client, nil
}

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

	// Limpar QR code da memória
	m.qrCodesMux.Lock()
	delete(m.qrCodes, sessionID)
	m.qrCodesMux.Unlock()

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

func (m *SessionManager) IsClientActive(sessionID string) bool {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	_, exists := m.clients[sessionID]
	return exists
}

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

func (m *SessionManager) UpdateWebhook(ctx context.Context, sessionID, webhookURL, webhookEvents string) error {
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Converter para novo formato
	var events []string
	if webhookEvents != "" {
		json.Unmarshal([]byte(webhookEvents), &events)
	}

	session.WebhookConfig = &model.WebhookConfig{
		Enabled: webhookURL != "",
		URL:     webhookURL,
		Events:  events,
	}

	if err := m.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Str("webhook_url", webhookURL).
		Msg("Webhook updated")

	return nil
}

func (m *SessionManager) UpdateWebhookConfig(ctx context.Context, sessionID string, webhookConfig *model.WebhookConfig) error {
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	session.WebhookConfig = webhookConfig

	if err := m.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Bool("enabled", webhookConfig.Enabled).
		Str("url", webhookConfig.URL).
		Msg("Webhook config updated")

	return nil
}

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

func (m *SessionManager) GetActiveSessionsCount() int {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()
	return len(m.clients)
}
