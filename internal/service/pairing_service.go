package service

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"zpwoot/internal/repository"
	"zpwoot/pkg/logger"
	"zpwoot/pkg/utils"
)

type PairingService struct {
	whatsappSvc *WhatsAppService
	sessionRepo *repository.SessionRepository
	manager     *SessionManager
}

func NewPairingService(whatsappSvc *WhatsAppService, sessionRepo *repository.SessionRepository, manager *SessionManager) *PairingService {
	return &PairingService{
		whatsappSvc: whatsappSvc,
		sessionRepo: sessionRepo,
		manager:     manager,
	}
}

func (p *PairingService) GenerateQRCode(ctx context.Context, sessionID string) (string, error) {
	// Buscar sessão
	session, err := p.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("session not found: %w", err)
	}

	// Verificar se já está pareado
	if session.DeviceJID != "" {
		return "", fmt.Errorf("session already paired")
	}

	// Criar ou obter device (sem JID pois ainda não está pareado)
	device, err := p.whatsappSvc.GetOrCreateDevice(ctx, sessionID, "")
	if err != nil {
		return "", fmt.Errorf("failed to get device: %w", err)
	}

	// Criar cliente
	client := p.whatsappSvc.NewClient(device)

	// Canal para receber QR code
	qrChan, err := client.GetQRChannel(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get QR channel: %w", err)
	}

	// Conectar
	if err := client.Connect(); err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}

	// Atualizar status
	session.Status = "pairing"
	p.sessionRepo.UpdateStatus(ctx, sessionID, "pairing", false)

	// Aguardar QR code ou pareamento
	var qrCode string
	timeout := time.After(2 * time.Minute)

	for {
		select {
		case evt := <-qrChan:
			switch evt.Event {
			case "code":
				// Gerar imagem QR code
				qrCodeImage, err := utils.GenerateQRCodeImage(evt.Code)
				if err != nil {
					logger.Log.Error().Err(err).Msg("Failed to generate QR code image")
					qrCode = evt.Code // Retornar código texto se falhar
				} else {
					qrCode = qrCodeImage
				}

				// Salvar QR code no banco
				session.QRCode = qrCode
				p.sessionRepo.Update(ctx, session)

				logger.Log.Info().
					Str("session_id", sessionID).
					Msg("QR code generated")

				// Iniciar goroutine para aguardar pareamento
				go p.waitForPairing(client, sessionID)

				return qrCode, nil

			case "success":
				// Pareamento bem-sucedido
				p.handlePairingSuccess(ctx, client, sessionID)
				return "", nil
			}

		case <-timeout:
			client.Disconnect()
			p.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
			return "", fmt.Errorf("QR code generation timeout")
		}
	}
}

func (p *PairingService) PairWithPhone(ctx context.Context, sessionID, phoneNumber string) (string, error) {
	// Buscar sessão
	session, err := p.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("session not found: %w", err)
	}

	// Verificar se já está pareado
	if session.DeviceJID != "" {
		return "", fmt.Errorf("session already paired")
	}

	// Criar ou obter device (sem JID pois ainda não está pareado)
	device, err := p.whatsappSvc.GetOrCreateDevice(ctx, sessionID, "")
	if err != nil {
		return "", fmt.Errorf("failed to get device: %w", err)
	}

	// Criar cliente
	client := p.whatsappSvc.NewClient(device)

	// Conectar
	if err := client.Connect(); err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}

	// Solicitar código de pareamento
	code, err := client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		client.Disconnect()
		return "", fmt.Errorf("failed to request pairing code: %w", err)
	}

	// Atualizar status
	p.sessionRepo.UpdateStatus(ctx, sessionID, "pairing", false)

	logger.Log.Info().
		Str("session_id", sessionID).
		Str("phone", phoneNumber).
		Msg("Pairing code requested")

	// Iniciar goroutine para aguardar pareamento
	go p.waitForPairing(client, sessionID)

	return code, nil
}

func (p *PairingService) waitForPairing(client *whatsmeow.Client, sessionID string) {
	ctx := context.Background()

	// Registrar handler temporário para evento de pareamento
	eventChan := make(chan interface{}, 10)
	handlerID := client.AddEventHandler(func(evt interface{}) {
		eventChan <- evt
	})
	defer client.RemoveEventHandler(handlerID)

	// Aguardar até 5 minutos pelo pareamento
	timeout := time.After(5 * time.Minute)

	for {
		select {
		case evt := <-eventChan:
			switch evt.(type) {
			case *events.PairSuccess:
				p.handlePairingSuccess(ctx, client, sessionID)
				return
			case *events.LoggedOut:
				logger.Log.Warn().
					Str("session_id", sessionID).
					Msg("Logged out during pairing")
				client.Disconnect()
				p.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
				return
			}

		case <-timeout:
			logger.Log.Warn().
				Str("session_id", sessionID).
				Msg("Pairing timeout")
			client.Disconnect()
			p.sessionRepo.UpdateStatus(ctx, sessionID, "disconnected", false)
			return
		}
	}
}

func (p *PairingService) handlePairingSuccess(ctx context.Context, client *whatsmeow.Client, sessionID string) {
	logger.Log.Info().
		Str("session_id", sessionID).
		Str("jid", client.Store.ID.String()).
		Msg("Pairing successful")

	// Buscar sessão
	session, err := p.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get session after pairing")
		return
	}

	// Atualizar sessão com informações do device
	session.DeviceJID = client.Store.ID.String()
	session.Status = "connected"
	session.Connected = true
	session.QRCode = "" // Limpar QR code

	if err := p.sessionRepo.Update(ctx, session); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to update session after pairing")
	}

	// Salvar device no whatsmeow store
	if err := p.whatsappSvc.SaveDevice(ctx, client.Store); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to save device after pairing")
	}

	// Adicionar cliente ao manager
	p.manager.clientsMux.Lock()
	p.manager.clients[sessionID] = client
	p.manager.clientsMux.Unlock()

	// Registrar event handlers
	p.manager.eventHandler.RegisterHandlers(client, sessionID)

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("Session paired and connected successfully")
}
