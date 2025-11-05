package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"zpwoot/internal/api/handlers"
	"zpwoot/internal/api/middleware"
)

func RegisterRoutes(r *gin.Engine, sessionHandler *handlers.SessionHandler) {
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

		// POST /sessions/:id/pair - Parear com telefone
		sessions.POST("/:id/pair", sessionHandler.PairPhone)

		// GET /sessions/:id/status - Obter status da sessão
		sessions.GET("/:id/status", sessionHandler.GetSessionStatus)

		// PUT /sessions/:id/webhook - Atualizar webhook
		sessions.PUT("/:id/webhook", sessionHandler.UpdateSessionWebhook)
	}
}
