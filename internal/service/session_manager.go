package service

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"golang.org/x/net/proxy"

	"zpwoot/internal/model"
	"zpwoot/internal/repository"
	"zpwoot/pkg/logger"
)

var (
	// Map de canais kill para controlar sessões: sessionID -> chan bool
	killChannels = make(map[string]chan bool)
	killMux      sync.RWMutex
)

type SessionManager struct {
	whatsappSvc *WhatsAppService
	sessionRepo *repository.SessionRepository

	// Map de clientes ativos: sessionID -> *whatsmeow.Client
	clients    map[string]*whatsmeow.Client
	clientsMux sync.RWMutex

	// Map de HTTP clients para cada sessão (para proxy support)
	httpClients    map[string]*resty.Client
	httpClientsMux sync.RWMutex

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
		httpClients: make(map[string]*resty.Client),
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

func (m *SessionManager) GetQRCode(sessionID string) (string, bool) {
	// Buscar QR code do banco
	ctx := context.Background()
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", false
	}

	if session.QRCode == "" {
		return "", false
	}

	return session.QRCode, true
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

func (m *SessionManager) IsClientActive(sessionID string) bool {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	_, exists := m.clients[sessionID]
	return exists
}

func (m *SessionManager) GetActiveSessionsCount() int {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()
	return len(m.clients)
}

// ConnectSession inicia conexão com WhatsApp (baseado em wuzapi startClient)
func (m *SessionManager) ConnectSession(ctx context.Context, sessionID string) error {
	logger.Log.Info().Str("session_id", sessionID).Msg("Starting WhatsApp connection")

	// Buscar sessão
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Verificar se já está conectado
	if m.IsClientActive(sessionID) {
		logger.Log.Warn().Str("session_id", sessionID).Msg("Session already active, disconnecting first")
		m.DisconnectSession(ctx, sessionID)
	}

	// Criar canal kill para esta sessão
	killMux.Lock()
	killChannels[sessionID] = make(chan bool)
	killMux.Unlock()

	// Iniciar cliente em goroutine (como no wuzapi)
	go m.startClient(sessionID, session.DeviceJID)

	return nil
}

// startClient é a função principal de conexão (baseada no wuzapi)
func (m *SessionManager) startClient(sessionID string, textJID string) {
	logger.Log.Info().Str("session_id", sessionID).Str("jid", textJID).Msg("Starting websocket connection to WhatsApp")

	ctx := context.Background()

	// Constantes de retry (do wuzapi)
	const maxConnectionRetries = 3
	const connectionRetryBaseWait = 5 * time.Second

	// Obter ou criar device
	deviceStore, err := m.whatsappSvc.GetOrCreateDevice(ctx, textJID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get device")
		m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
		return
	}

	// Criar cliente WhatsApp
	client := m.whatsappSvc.NewClient(deviceStore, false)

	// Adicionar ao map de clientes
	m.clientsMux.Lock()
	m.clients[sessionID] = client
	m.clientsMux.Unlock()

	// Configurar HTTP client com proxy support (do wuzapi)
	httpClient := resty.New()
	httpClient.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	httpClient.SetTimeout(30 * time.Second)
	httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// Buscar configuração de proxy do banco
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err == nil && session.ProxyConfig != nil && session.ProxyConfig.Enabled {
		proxyURL := buildProxyURL(session.ProxyConfig)
		if proxyURL != "" {
			m.configureProxy(client, httpClient, proxyURL)
		}
	}

	// Salvar HTTP client
	m.httpClientsMux.Lock()
	m.httpClients[sessionID] = httpClient
	m.httpClientsMux.Unlock()

	// Registrar event handlers
	m.eventHandler.RegisterHandlers(client, sessionID)

	// Verificar se precisa fazer pairing (QR code)
	if client.Store.ID == nil {
		// No ID stored, new login - precisa QR code
		qrChan, err := client.GetQRChannel(context.Background())
		if err != nil {
			if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
				logger.Log.Error().Err(err).Msg("Failed to get QR channel")
				m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
				return
			}
		} else {
			// Conectar ANTES de processar QR codes (IMPORTANTE!)
			err = client.Connect()
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to connect client")
				m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
				return
			}

			// Processar QR codes SÍNCRONAMENTE (como no wuzapi!)
			// O loop só termina quando o canal fecha (após pareamento ou timeout)
			for evt := range qrChan {
				if evt.Event == "code" {
					// Novo QR code gerado
					logger.Log.Info().
						Str("session_id", sessionID).
						Msg("QR code received")

					// Gerar imagem QR code base64
					image, _ := qrcode.Encode(evt.Code, qrcode.Medium, 256)
					base64QRCode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)

					// Salvar no banco
					if err := m.sessionRepo.UpdateQRCode(ctx, sessionID, base64QRCode); err != nil {
						logger.Log.Error().Err(err).Msg("Failed to save QR code")
					}

					// Atualizar status
					m.sessionRepo.UpdateStatus(ctx, sessionID, "qr_code", false)

					// Exibir QR code no terminal
					fmt.Printf("\n========================================\n")
					fmt.Printf("QR CODE for Session: %s\n", sessionID)
					fmt.Printf("========================================\n")
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					fmt.Printf("========================================\n\n")

				} else if evt.Event == "timeout" {
					logger.Log.Warn().Str("session_id", sessionID).Msg("QR code timeout - killing channel")

					// Limpar QR code
					m.sessionRepo.UpdateQRCode(ctx, sessionID, "")
					m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)

					// Cleanup
					m.cleanupSession(sessionID)

					// Enviar kill signal
					killMux.Lock()
					if ch, exists := killChannels[sessionID]; exists {
						ch <- true
					}
					killMux.Unlock()
					return

				} else if evt.Event == "success" {
					logger.Log.Info().Str("session_id", sessionID).Msg("QR pairing ok!")

					// Limpar QR code e atualizar status
					m.sessionRepo.UpdateQRCode(ctx, sessionID, "")
					m.sessionRepo.UpdateStatus(ctx, sessionID, "connected", true)

					// Salvar JID
					if client.Store.ID != nil {
						m.sessionRepo.UpdateDeviceJID(ctx, sessionID, client.Store.ID.String())
					}

				} else {
					logger.Log.Info().Str("session_id", sessionID).Str("event", evt.Event).Msg("Login event")
				}
			}
			// Canal QR fechado - continuar para keep-alive loop
		}
	} else {
		// Já está pareado, apenas conectar com retry
		logger.Log.Info().Msg("Already logged in, connecting with retry logic")

		var lastErr error
		for attempt := 0; attempt < maxConnectionRetries; attempt++ {
			if attempt > 0 {
				waitTime := time.Duration(attempt) * connectionRetryBaseWait
				logger.Log.Warn().
					Int("attempt", attempt+1).
					Int("max_retries", maxConnectionRetries).
					Dur("wait_time", waitTime).
					Msg("Retrying connection after delay")
				time.Sleep(waitTime)
			}

			err = client.Connect()
			if err == nil {
				logger.Log.Info().
					Int("attempt", attempt+1).
					Msg("Successfully connected to WhatsApp")
				break
			}

			lastErr = err
			logger.Log.Warn().
				Err(err).
				Int("attempt", attempt+1).
				Int("max_retries", maxConnectionRetries).
				Msg("Failed to connect to WhatsApp")
		}

		if lastErr != nil {
			logger.Log.Error().
				Err(lastErr).
				Str("session_id", sessionID).
				Int("attempts", maxConnectionRetries).
				Msg("Failed to connect after all retry attempts")

			m.cleanupSession(sessionID)
			m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
			return
		}
	}

	// Keep alive loop (do wuzapi) - mantém a conexão viva
	m.keepAliveLoop(sessionID, client)
}

// configureProxy configura proxy para o cliente (do wuzapi)
func (m *SessionManager) configureProxy(client *whatsmeow.Client, httpClient *resty.Client, proxyURL string) {
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		logger.Log.Warn().Err(err).Str("proxy", proxyURL).Msg("Invalid proxy URL")
		return
	}

	logger.Log.Info().Str("proxy", proxyURL).Msg("Configuring proxy")

	if parsed.Scheme == "socks5" || parsed.Scheme == "socks5h" {
		dialer, err := proxy.FromURL(parsed, nil)
		if err != nil {
			logger.Log.Warn().Err(err).Msg("Failed to build SOCKS proxy dialer")
		} else {
			httpClient.SetProxy(proxyURL)
			client.SetSOCKSProxy(dialer, whatsmeow.SetProxyOptions{})
			logger.Log.Info().Msg("SOCKS proxy configured successfully")
		}
	} else {
		httpClient.SetProxy(proxyURL)
		client.SetProxyAddress(parsed.String(), whatsmeow.SetProxyOptions{})
		logger.Log.Info().Msg("HTTP/HTTPS proxy configured successfully")
	}
}

// keepAliveLoop mantém a conexão ativa (do wuzapi)
func (m *SessionManager) keepAliveLoop(sessionID string, client *whatsmeow.Client) {
	killMux.RLock()
	killChan, exists := killChannels[sessionID]
	killMux.RUnlock()

	if !exists {
		logger.Log.Error().Str("session_id", sessionID).Msg("Kill channel not found")
		return
	}

	for {
		select {
		case <-killChan:
			logger.Log.Info().Str("session_id", sessionID).Msg("Received kill signal")
			client.Disconnect()
			m.cleanupSession(sessionID)
			m.sessionRepo.UpdateStatus(context.Background(), sessionID, "disconnected", false)
			return
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

// cleanupSession limpa recursos da sessão
func (m *SessionManager) cleanupSession(sessionID string) {
	m.clientsMux.Lock()
	delete(m.clients, sessionID)
	m.clientsMux.Unlock()

	m.httpClientsMux.Lock()
	delete(m.httpClients, sessionID)
	m.httpClientsMux.Unlock()

	killMux.Lock()
	delete(killChannels, sessionID)
	killMux.Unlock()
}

// DisconnectSession desconecta uma sessão
func (m *SessionManager) DisconnectSession(ctx context.Context, sessionID string) error {
	logger.Log.Info().Str("session_id", sessionID).Msg("Disconnecting session")

	// Enviar sinal kill
	killMux.Lock()
	if killChan, exists := killChannels[sessionID]; exists {
		select {
		case killChan <- true:
		default:
		}
	}
	killMux.Unlock()

	// Aguardar um pouco para cleanup
	time.Sleep(500 * time.Millisecond)

	// Garantir cleanup
	m.cleanupSession(sessionID)

	// Atualizar status no banco
	if err := m.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to update session status")
	}

	return nil
}

// RestoreAllSessions reconecta sessões que estavam conectadas (do wuzapi connectOnStartup)
func (m *SessionManager) RestoreAllSessions(ctx context.Context) error {
	sessions, err := m.sessionRepo.ListConnected(ctx)
	if err != nil {
		return fmt.Errorf("failed to list connected sessions: %w", err)
	}

	logger.Log.Info().
		Int("count", len(sessions)).
		Msg("Restoring connected sessions")

	for _, session := range sessions {
		logger.Log.Info().
			Str("session_id", session.ID).
			Str("name", session.Name).
			Str("jid", session.DeviceJID).
			Msg("Attempting to restore session")

		// Criar canal kill
		killMux.Lock()
		killChannels[session.ID] = make(chan bool)
		killMux.Unlock()

		// Iniciar cliente em goroutine
		go m.startClient(session.ID, session.DeviceJID)
	}

	return nil
}

// Shutdown desconecta todas as sessões
func (m *SessionManager) Shutdown(ctx context.Context) error {
	m.clientsMux.Lock()
	sessionIDs := make([]string, 0, len(m.clients))
	for sessionID := range m.clients {
		sessionIDs = append(sessionIDs, sessionID)
	}
	m.clientsMux.Unlock()

	logger.Log.Info().
		Int("count", len(sessionIDs)).
		Msg("Shutting down all sessions")

	for _, sessionID := range sessionIDs {
		m.DisconnectSession(ctx, sessionID)
	}

	return nil
}

// GetSessionStatus retorna status detalhado da sessão
func (m *SessionManager) GetSessionStatus(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	session, err := m.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

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

	if isActive {
		if client, _ := m.GetClient(sessionID); client != nil {
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

// UpdateWebhookConfig atualiza configuração de webhook
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

func buildProxyURL(config *model.ProxyConfig) string {
	if config == nil || !config.Enabled {
		return ""
	}

	auth := ""
	if config.Username != "" {
		auth = config.Username
		if config.Password != "" {
			auth += ":" + config.Password
		}
		auth += "@"
	}

	return fmt.Sprintf("%s://%s%s:%d", config.Protocol, auth, config.Host, config.Port)
}

