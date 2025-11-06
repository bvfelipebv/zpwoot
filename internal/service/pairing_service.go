package service

import (
	"context"
	"fmt"

	"zpwoot/internal/repository"
	"zpwoot/pkg/logger"
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

	// Iniciar conexão (que vai gerar QR code automaticamente)
	if err := p.manager.ConnectSession(ctx, sessionID); err != nil {
		return "", fmt.Errorf("failed to start connection: %w", err)
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("QR code generation started - check session QR code")

	return "QR code generation started", nil
}

func (p *PairingService) PairWithPhone(ctx context.Context, sessionID, phoneNumber string) (string, error) {
	// Por enquanto, não implementado - usar QR code
	return "", fmt.Errorf("phone pairing not implemented yet - use QR code")
}

