package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"zpwoot/internal/api"
	"zpwoot/internal/api/handlers"
	"zpwoot/internal/config"
	"zpwoot/internal/db"
	natsclient "zpwoot/internal/nats"
	"zpwoot/internal/repository"
	"zpwoot/internal/service"
	"zpwoot/pkg/logger"

	_ "zpwoot/docs" // Importa a documentação gerada pelo swag
)

// @title           ZPWoot - WhatsApp Multi-Device API
// @version         1.0
// @description     API REST para gerenciamento de múltiplas sessões do WhatsApp usando whatsmeow
// @description     Permite criar, conectar e gerenciar sessões do WhatsApp via API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name apikey
// @description Insira sua API Key (exemplo: sldkfjsldkflskdfjlsd)

func main() {
	// Initialize logger
	logger.Init("info")
	logger.Log.Info().Msg("🚀 Starting zpwoot - WhatsApp Multi-Device API")

	// Load configuration
	if err := config.Load(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Reconfigure logger with config level
	logger.Init(config.AppConfig.LogLevel)

	// Initialize database
	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()

	// Initialize WhatsApp service
	whatsappSvc, err := service.NewWhatsAppService(db.DB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize WhatsApp service")
	}
	defer whatsappSvc.Close()

	// Initialize NATS
	natsClient := natsclient.NewClient(natsclient.Config{
		URL:           config.AppConfig.NATSURL,
		MaxReconnect:  config.AppConfig.NATSMaxReconnect,
		ReconnectWait: config.AppConfig.NATSReconnectWait,
	})
	if err := natsClient.Connect(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	defer natsClient.Close()

	// Initialize repositories
	sessionRepo := repository.NewSessionRepository(db.DB)

	// Initialize webhook services
	webhookFormatter := service.NewWebhookFormatter()
	webhookProcessor := service.NewWebhookProcessor(natsClient, webhookFormatter, sessionRepo)
	webhookDelivery := service.NewWebhookDelivery(config.AppConfig.WebhookTimeout)

	// Initialize services
	sessionManager := service.NewSessionManager(whatsappSvc, sessionRepo, webhookProcessor, webhookFormatter)
	pairingService := service.NewPairingService(whatsappSvc, sessionRepo, sessionManager)

	// Start webhook workers
	webhookWorkers := make([]*service.WebhookWorker, config.AppConfig.WebhookWorkers)
	for i := 0; i < config.AppConfig.WebhookWorkers; i++ {
		worker := service.NewWebhookWorker(
			i+1,
			natsClient,
			webhookDelivery,
			config.AppConfig.WebhookMaxRetries,
			config.AppConfig.WebhookRetryBaseDelay,
		)
		if err := worker.Start(); err != nil {
			logger.Log.Fatal().
				Err(err).
				Int("worker_id", i+1).
				Msg("Failed to start webhook worker")
		}
		webhookWorkers[i] = worker
	}
	logger.Log.Info().
		Int("workers", config.AppConfig.WebhookWorkers).
		Msg("✅ Webhook workers started")

	// Restore sessions if configured
	if config.AppConfig.AutoRestoreSessions {
		if err := sessionManager.RestoreAllSessions(context.Background()); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to restore sessions")
		} else {
			activeCount := sessionManager.GetActiveSessionsCount()
			if activeCount > 0 {
				logger.Log.Info().
					Int("active_sessions", activeCount).
					Msg("✅ Sessions restored")
			}
		}
	}

	// Initialize handlers
	sessionHandler := handlers.NewSessionHandler(sessionManager, pairingService)
	messageHandler := handlers.NewMessageHandler(sessionManager)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	// Register routes
	api.RegisterRoutes(r, sessionHandler, messageHandler)

	// Server info
	port := config.AppConfig.Port
	addr := fmt.Sprintf(":%s", port)

	// Start server in goroutine
	go func() {
		if err := r.Run(addr); err != nil {
			logger.Log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	logger.Log.Info().Msgf("✅ Server started at http://localhost:%s", port)
	logger.Log.Info().Msg("Press Ctrl+C to shutdown")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info().Msg("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown all sessions
	logger.Log.Info().Msg("Disconnecting all sessions...")
	if err := sessionManager.Shutdown(ctx); err != nil {
		logger.Log.Error().Err(err).Msg("Error during session shutdown")
	}

	logger.Log.Info().Msg("✅ Server shutdown complete")
}
