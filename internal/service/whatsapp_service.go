package service

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"

	"zpwoot/pkg/logger"
)

type WhatsAppService struct {
	container *sqlstore.Container
	db        *sql.DB
}

func NewWhatsAppService(db *sql.DB) (*WhatsAppService, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Criar logger para whatsmeow
	waLogger := waLog.Stdout("WhatsApp", "INFO", true)

	// Criar container do whatsmeow usando a mesma conexão SQL
	container := sqlstore.NewWithDB(db, "postgres", waLogger)

	// Executar upgrade das tabelas do whatsmeow
	if err := container.Upgrade(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to upgrade whatsmeow database: %w", err)
	}

	logger.Log.Info().Msg("WhatsApp service initialized successfully")

	return &WhatsAppService{
		container: container,
		db:        db,
	}, nil
}

func (s *WhatsAppService) GetOrCreateDevice(ctx context.Context, sessionID string, jid string) (*store.Device, error) {
	// Se temos um JID, tentar obter device existente
	if jid != "" {
		parsedJID, err := types.ParseJID(jid)
		if err != nil {
			logger.Log.Warn().
				Err(err).
				Str("session_id", sessionID).
				Str("jid", jid).
				Msg("Failed to parse JID, creating new device")
			return s.container.NewDevice(), nil
		}

		device, err := s.container.GetDevice(ctx, parsedJID)
		if err != nil {
			logger.Log.Warn().
				Err(err).
				Str("session_id", sessionID).
				Str("jid", jid).
				Msg("Failed to get device, creating new one")
			return s.container.NewDevice(), nil
		}

		if device != nil {
			logger.Log.Info().
				Str("session_id", sessionID).
				Str("jid", jid).
				Msg("Retrieved existing WhatsApp device")
			return device, nil
		}
	}

	// Se não tem JID ou não encontrou device, criar novo
	device := s.container.NewDevice()

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("Created new WhatsApp device")

	return device, nil
}

func (s *WhatsAppService) GetDeviceByJID(ctx context.Context, jid types.JID) (*store.Device, error) {
	device, err := s.container.GetDevice(ctx, jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get device by JID: %w", err)
	}
	return device, nil
}

func (s *WhatsAppService) GetFirstDevice(ctx context.Context) (*store.Device, error) {
	device, err := s.container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get first device: %w", err)
	}
	return device, nil
}

func (s *WhatsAppService) SaveDevice(ctx context.Context, device *store.Device) error {
	if device == nil {
		return fmt.Errorf("device is nil")
	}

	if err := s.container.PutDevice(ctx, device); err != nil {
		return fmt.Errorf("failed to save device: %w", err)
	}

	logger.Log.Info().
		Str("jid", device.ID.String()).
		Msg("Device saved successfully")

	return nil
}

func (s *WhatsAppService) DeleteDevice(ctx context.Context, device *store.Device) error {
	if device == nil {
		return fmt.Errorf("device is nil")
	}

	if err := s.container.DeleteDevice(ctx, device); err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	logger.Log.Info().
		Str("jid", device.ID.String()).
		Msg("Device deleted successfully")

	return nil
}

func (s *WhatsAppService) GetAllDevices(ctx context.Context) ([]*store.Device, error) {
	devices, err := s.container.GetAllDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}
	return devices, nil
}

func (s *WhatsAppService) NewClient(device *store.Device) *whatsmeow.Client {
	return whatsmeow.NewClient(device, nil)
}

func (s *WhatsAppService) Close() error {
	if s.container != nil {
		return s.container.Close()
	}
	return nil
}
