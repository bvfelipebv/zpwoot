package service

import (
	"fmt"
	"time"

	"zpwoot/internal/config"
	"zpwoot/pkg/logger"
)

type MeowService struct {
	DataDir string
	running bool
}

func NewMeowService() *MeowService {
	return &MeowService{DataDir: config.AppConfig.WhatsAppDataDir}
}

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

func (m *MeowService) Stop() error {
	if !m.running {
		return nil
	}
	m.running = false
	logger.Log.Info().Msg("stopped MeowService (stub)")
	return nil
}

func (m *MeowService) Status() string {
	if m.running {
		return "running"
	}
	return "stopped"
}

func (m *MeowService) SendMessage(to string, message string) error {
	if !m.running {
		return fmt.Errorf("service not running")
	}
	logger.Log.Info().Str("to", to).Msgf("sending message: %s", message)
	return nil
}
