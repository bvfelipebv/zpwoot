package api

import (
	"github.com/gin-gonic/gin"

	"zpwoot/internal/api/handlers"
	"zpwoot/internal/api/middleware"
)

// RegisterRoutes registra todas as rotas da API
func RegisterRoutes(r *gin.Engine, sessionHandler *handlers.SessionHandler) {
	// Middlewares globais
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// Health check (sem autenticação)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "zpwoot",
		})
	})

	// Grupo de rotas da API com autenticação
	api := r.Group("/api")
	api.Use(middleware.AuthenticateGlobal())
	{
		// Grupo de rotas de sessões
		sessions := api.Group("/sessions")
		{
			// POST /api/sessions/create - Criar nova sessão
			sessions.POST("/create", sessionHandler.CreateSession)

			// GET /api/sessions/list - Listar todas as sessões
			sessions.GET("/list", sessionHandler.GetSessions)

			// GET /api/sessions/:id/info - Obter detalhes da sessão
			sessions.GET("/:id/info", sessionHandler.GetSession)

			// DELETE /api/sessions/:id/delete - Deletar sessão
			sessions.DELETE("/:id/delete", sessionHandler.DeleteSession)

			// POST /api/sessions/:id/connect - Conectar sessão
			sessions.POST("/:id/connect", sessionHandler.ConnectSession)

			// POST /api/sessions/:id/disconnect - Desconectar sessão
			sessions.POST("/:id/disconnect", sessionHandler.DisconnectSession)

			// POST /api/sessions/:id/pair - Parear com telefone
			sessions.POST("/:id/pair", sessionHandler.PairPhone)

			// GET /api/sessions/:id/status - Obter status da sessão
			sessions.GET("/:id/status", sessionHandler.GetSessionStatus)

			// PUT /api/sessions/:id/webhook - Atualizar webhook
			sessions.PUT("/:id/webhook", sessionHandler.UpdateSessionWebhook)
		}
	}
}
