package middleware

import (
	"net/http"
	"strings"
	"time"

	"zpwoot/internal/config"
	"zpwoot/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
				Str(logger.FieldIP, c.ClientIP()).
				Str(logger.FieldPath, c.Request.URL.Path).
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
				Str(logger.FieldIP, c.ClientIP()).
				Str(logger.FieldPath, c.Request.URL.Path).
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
			Str(logger.FieldIP, c.ClientIP()).
			Str(logger.FieldPath, c.Request.URL.Path).
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
		// Generate unique request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Record start time
		startTime := time.Now()

		// Log incoming request
		logger.Log.Info().
			Str(logger.FieldRequestID, requestID).
			Str(logger.FieldMethod, c.Request.Method).
			Str(logger.FieldPath, c.Request.URL.Path).
			Str(logger.FieldIP, c.ClientIP()).
			Str(logger.FieldUserAgent, c.Request.UserAgent()).
			Str("query", c.Request.URL.RawQuery).
			Msg("→ Incoming request")

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)
		status := c.Writer.Status()

		// Determine log level based on status code
		logEvent := logger.Log.Info()
		if status >= 500 {
			logEvent = logger.Log.Error()
		} else if status >= 400 {
			logEvent = logger.Log.Warn()
		}

		// Log completed request
		logEvent.
			Str(logger.FieldRequestID, requestID).
			Str(logger.FieldMethod, c.Request.Method).
			Str(logger.FieldPath, c.Request.URL.Path).
			Int(logger.FieldStatus, status).
			Dur(logger.FieldDuration, duration).
			Int("response_size", c.Writer.Size()).
			Msg("← Request completed")
	}
}

func RequestLoggerWithSkip(skipPaths ...string) gin.HandlerFunc {
	skipMap := make(map[string]bool)
	for _, path := range skipPaths {
		skipMap[path] = true
	}

	return func(c *gin.Context) {
		// Skip logging for certain paths
		if skipMap[c.Request.URL.Path] {
			c.Next()
			return
		}

		// Use regular request logger
		RequestLogger()(c)
	}
}
