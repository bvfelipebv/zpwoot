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
	logger.Log.Info().
		Str("environment", config.AppConfig.Environment).
		Str("log_level", config.AppConfig.LogLevel).
		Msg("Configuration loaded")

	// Initialize database
	logger.Log.Info().Msg("Initializing database connection...")
	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()
	logger.Log.Info().Msg("✅ Database initialized successfully")

	// Initialize WhatsApp service
	logger.Log.Info().Msg("Initializing WhatsApp service...")
	whatsappSvc, err := service.NewWhatsAppService(db.DB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize WhatsApp service")
	}
	defer whatsappSvc.Close()
	logger.Log.Info().Msg("✅ WhatsApp service initialized")

	// Initialize repositories
	sessionRepo := repository.NewSessionRepository(db.DB)

	// Initialize services
	sessionManager := service.NewSessionManager(whatsappSvc, sessionRepo)
	pairingService := service.NewPairingService(whatsappSvc, sessionRepo, sessionManager)

	// Restore sessions if configured
	if config.AppConfig.AutoRestoreSessions {
		logger.Log.Info().Msg("Restoring connected sessions...")
		if err := sessionManager.RestoreAllSessions(context.Background()); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to restore sessions")
		} else {
			activeCount := sessionManager.GetActiveSessionsCount()
			logger.Log.Info().
				Int("active_sessions", activeCount).
				Msg("✅ Sessions restored")
		}
	}

	// Initialize handlers
	sessionHandler := handlers.NewSessionHandler(sessionManager, pairingService)
	messageHandler := handlers.NewMessageHandler(sessionManager)

	// Setup Gin
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// Register routes
	api.RegisterRoutes(r, sessionHandler, messageHandler)

	// Server info
	port := config.AppConfig.Port
	addr := fmt.Sprintf(":%s", port)

	logger.Log.Info().
		Str("port", port).
		Str("address", addr).
		Msg("🌐 Starting HTTP server")

	// Start server in goroutine
	go func() {
		if err := r.Run(addr); err != nil {
			logger.Log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	logger.Log.Info().Msg("✅ Server started successfully")
	logger.Log.Info().Msgf("📡 API available at http://localhost:%s", port)
	logger.Log.Info().Msgf("🏥 Health check: http://localhost:%s/health", port)
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
