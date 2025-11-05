package middleware

import (
	"net/http"
	"strings"

	"zpwoot/internal/config"
	"zpwoot/pkg/logger"

	"github.com/gin-gonic/gin"
)

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

		// Obter API Key do header "apikey"
		apiKey := strings.TrimSpace(c.GetHeader("apikey"))

		// Validar API Key
		if apiKey == "" {
			logger.Log.Warn().
				Str("ip", c.ClientIP()).
				Str("path", c.Request.URL.Path).
				Msg("Missing API key")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "API key is required. Use header: apikey: <your_key>",
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

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, apikey, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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
