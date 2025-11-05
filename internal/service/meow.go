package service

import (
	"fmt"
	"time"

	"zpmeow/internal/config"
	"zpmeow/pkg/logger"
)

// MeowService é um stub para encapsular a lógica do cliente whatsmeow.
// No momento é um stub que armazena o dataDir e possui Start/Stop simples.
type MeowService struct {
	DataDir string
	running bool
}

// NewMeowService cria uma instância do MeowService
func NewMeowService() *MeowService {
	return &MeowService{DataDir: config.AppConfig.WhatsAppDataDir}
}

// Start inicia o MeowService (stub)
func (m *MeowService) Start() error {
	if m.running {
		return nil
	}
	logger.Log.Info().Str("data_dir", m.DataDir).Msg("starting MeowService (stub)")
	// Simular inicialização leve
	go func() {
		m.running = true
		// Simular algum trabalho
		time.Sleep(100 * time.Millisecond)
		logger.Log.Debug().Msg("MeowService stub ready")
	}()
	return nil
}

// Stop para o MeowService
func (m *MeowService) Stop() error {
	if !m.running {
		return nil
	}
	m.running = false
	logger.Log.Info().Msg("stopped MeowService (stub)")
	return nil
}

// Status retorna se o serviço está rodando
func (m *MeowService) Status() string {
	if m.running {
		return "running"
	}
	return "stopped"
}

// Example method to send a message (stub)
func (m *MeowService) SendMessage(to string, message string) error {
	if !m.running {
		return fmt.Errorf("service not running")
	}
	logger.Log.Info().Str("to", to).Msgf("sending message: %s", message)
	return nil
}
