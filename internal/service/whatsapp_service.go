package service

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	"zpwoot/pkg/logger"
)

var (
	// Container global do whatsmeow (similar ao wuzapi)
	container *sqlstore.Container
	// Nome do OS para identificação no WhatsApp (baseado no wmial.bak)
	osName = "zpwoot"
)

type WhatsAppService struct {
	db *sql.DB
}

func NewWhatsAppService(db *sql.DB) (*WhatsAppService, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Criar logger para whatsmeow
	waLogger := waLog.Stdout("WhatsApp", "INFO", true)

	// Criar container global do whatsmeow usando a mesma conexão SQL
	container = sqlstore.NewWithDB(db, "postgres", waLogger)

	// Executar upgrade das tabelas do whatsmeow
	if err := container.Upgrade(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to upgrade whatsmeow database: %w", err)
	}

	// Configurar propriedades do device (IMPORTANTE para WhatsApp aceitar)
	// Baseado no wmial.bak que funciona corretamente
	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_UNKNOWN.Enum()
	store.DeviceProps.Os = proto.String(osName)



	return &WhatsAppService{
		db: db,
	}, nil
}

func (s *WhatsAppService) GetOrCreateDevice(ctx context.Context, jid string) (*store.Device, error) {
	var deviceStore *store.Device
	var err error

	// Se temos um JID, tentar obter device existente
	if jid != "" {
		parsedJID, parseErr := types.ParseJID(jid)
		if parseErr != nil {
			logger.Log.Warn().
				Err(parseErr).
				Str("jid", jid).
				Msg("Failed to parse JID, creating new device")
			deviceStore = container.NewDevice()
		} else {
			deviceStore, err = container.GetDevice(ctx, parsedJID)
			if err != nil {
				logger.Log.Warn().
					Err(err).
					Str("jid", jid).
					Msg("Failed to get device, creating new one")
				deviceStore = container.NewDevice()
			}
		}
	}

	// Se não encontrou device ou não tinha JID, criar novo
	if deviceStore == nil {
		logger.Log.Warn().Msg("No store found. Creating new one")
		deviceStore = container.NewDevice()
	}

	return deviceStore, nil
}

func (s *WhatsAppService) NewClient(device *store.Device, debug bool) *whatsmeow.Client {
	if debug {
		clientLog := waLog.Stdout("Client", "DEBUG", true)
		return whatsmeow.NewClient(device, clientLog)
	}
	return whatsmeow.NewClient(device, nil)
}

func (s *WhatsAppService) Close() error {
	if container != nil {
		return container.Close()
	}
	return nil
}

