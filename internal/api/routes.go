package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"zpwoot/internal/api/handlers"
	"zpwoot/internal/api/middleware"
)

func RegisterRoutes(r *gin.Engine, sessionHandler *handlers.SessionHandler, messageHandler *handlers.MessageHandler) {
	// Middlewares globais
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// Swagger documentation (sem autenticação)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check (sem autenticação)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "zpwoot",
		})
	})

	// Grupo de rotas de sessões com autenticação
	sessions := r.Group("/sessions")
	sessions.Use(middleware.AuthenticateGlobal())
	{
		// === ROTAS DE EVENTOS DE WEBHOOK (GLOBAIS) ===
		// GET /sessions/webhook/events - Listar todos os eventos suportados
		sessions.GET("/webhook/events", sessionHandler.ListWebhookEvents)

		// POST /sessions/create - Criar nova sessão
		sessions.POST("/create", sessionHandler.CreateSession)

		// GET /sessions/list - Listar todas as sessões
		sessions.GET("/list", sessionHandler.GetSessions)

		// GET /sessions/:id/info - Obter detalhes da sessão
		sessions.GET("/:id/info", sessionHandler.GetSession)

		// DELETE /sessions/:id/delete - Deletar sessão
		sessions.DELETE("/:id/delete", sessionHandler.DeleteSession)

		// POST /sessions/:id/connect - Conectar sessão
		sessions.POST("/:id/connect", sessionHandler.ConnectSession)

		// POST /sessions/:id/disconnect - Desconectar sessão
		sessions.POST("/:id/disconnect", sessionHandler.DisconnectSession)

		// GET /sessions/:id/qr - Obter QR Code atual
		sessions.GET("/:id/qr", sessionHandler.GetQRCode)

		// POST /sessions/:id/pair - Parear com telefone
		sessions.POST("/:id/pair", sessionHandler.PairPhone)

		// GET /sessions/:id/status - Obter status da sessão
		sessions.GET("/:id/status", sessionHandler.GetSessionStatus)

		// === ROTAS DE WEBHOOK ===
		webhook := sessions.Group("/:id/webhook")
		{
			// POST /sessions/:id/webhook/set - Configurar/Atualizar webhook
			webhook.POST("/set", sessionHandler.SetWebhook)

			// GET /sessions/:id/webhook/find - Obter configuração de webhook
			webhook.GET("/find", sessionHandler.FindWebhook)
		}

		// === ROTAS DE MENSAGENS ===
		messages := sessions.Group("/:id/message")
		{
			// POST /sessions/:id/message/text - Enviar mensagem de texto
			messages.POST("/text", messageHandler.SendText)

			// POST /sessions/:id/message/image - Enviar imagem
			messages.POST("/image", messageHandler.SendImage)

			// POST /sessions/:id/message/audio - Enviar áudio
			messages.POST("/audio", messageHandler.SendAudio)

			// POST /sessions/:id/message/video - Enviar vídeo
			messages.POST("/video", messageHandler.SendVideo)

			// POST /sessions/:id/message/document - Enviar documento
			messages.POST("/document", messageHandler.SendDocument)

			// POST /sessions/:id/message/sticker - Enviar sticker
			messages.POST("/sticker", messageHandler.SendSticker)

			// POST /sessions/:id/message/media - Enviar mídia genérica (auto-detect)
			messages.POST("/media", messageHandler.SendMedia)

			// POST /sessions/:id/message/contact - Enviar contato
			messages.POST("/contact", messageHandler.SendContact)

			// POST /sessions/:id/message/location - Enviar localização
			messages.POST("/location", messageHandler.SendLocation)

			// POST /sessions/:id/message/poll - Enviar enquete
			messages.POST("/poll", messageHandler.SendPoll)

			// POST /sessions/:id/message/reaction - Enviar reação
			messages.POST("/reaction", messageHandler.SendReaction)

			// POST /sessions/:id/message/presence - Enviar presença (digitando, gravando, etc)
			messages.POST("/presence", messageHandler.SendPresence)

			// POST /sessions/:id/message/read - Marcar como lida
			messages.POST("/read", messageHandler.MarkAsRead)

			// DELETE /sessions/:id/message/revoke - Revogar mensagem
			messages.DELETE("/revoke", messageHandler.RevokeMessage)

			// PUT /sessions/:id/message/edit - Editar mensagem
			messages.PUT("/edit", messageHandler.EditMessage)
		}
	}
}
