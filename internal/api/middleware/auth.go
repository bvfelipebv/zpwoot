package middleware

import (
	"net/http"
	"strings"

	"zpwoot/internal/config"
	"zpwoot/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthenticateGlobal middleware para autenticação via API Key
// Suporta 3 métodos:
// 1. Authorization: Bearer <token>
// 2. X-API-Key: <token>
// 3. Query param: ?api_key=<token>
func AuthenticateGlobal() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter API Key configurada
		expectedAPIKey := config.AppConfig.APIKey

		// Se não houver API Key configurada, permitir acesso
		if expectedAPIKey == "" {
			logger.Log.Warn().Msg("No API key configured - authentication disabled")
			c.Next()
			return
		}

		// Tentar obter API Key de diferentes fontes
		apiKey := extractAPIKey(c)

		// Validar API Key
		if apiKey == "" {
			logger.Log.Warn().
				Str("ip", c.ClientIP()).
				Str("path", c.Request.URL.Path).
				Msg("Missing API key")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "API key is required",
			})
			c.Abort()
			return
		}

		if apiKey != expectedAPIKey {
			logger.Log.Warn().
				Str("ip", c.ClientIP()).
				Str("path", c.Request.URL.Path).
				Msg("Invalid API key")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		// API Key válida - continuar
		logger.Log.Debug().
			Str("ip", c.ClientIP()).
			Str("path", c.Request.URL.Path).
			Msg("Request authenticated")

		c.Next()
	}
}

// extractAPIKey extrai a API Key de diferentes fontes
func extractAPIKey(c *gin.Context) string {
	// 1. Tentar Authorization Bearer token
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// Formato: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return strings.TrimSpace(parts[1])
		}
	}

	// 2. Tentar X-API-Key header
	apiKeyHeader := c.GetHeader("X-API-Key")
	if apiKeyHeader != "" {
		return strings.TrimSpace(apiKeyHeader)
	}

	// 3. Tentar query parameter
	apiKeyQuery := c.Query("api_key")
	if apiKeyQuery != "" {
		return strings.TrimSpace(apiKeyQuery)
	}

	return ""
}

// CORS middleware para permitir requisições cross-origin
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLogger middleware para logging de requisições
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Incoming request")

		c.Next()

		logger.Log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Msg("Request completed")
	}
}
