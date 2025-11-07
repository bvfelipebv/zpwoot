package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"

	"zpwoot/internal/api/dto"
	"zpwoot/internal/constants"
	"zpwoot/internal/model"
	"zpwoot/internal/service"
	"zpwoot/pkg/logger"
)

type SessionHandler struct {
	sessionManager *service.SessionManager
	pairingService *service.PairingService
}

func NewSessionHandler(sessionManager *service.SessionManager, pairingService *service.PairingService) *SessionHandler {
	return &SessionHandler{
		sessionManager: sessionManager,
		pairingService: pairingService,
	}
}

// @Summary Criar nova sessão
// @Description Cria uma nova sessão do WhatsApp com nome e webhook opcional
// @Tags Sessions
// @Accept json
// @Produce json
// @Param request body dto.CreateSessionRequest true "Dados da sessão"
// @Success 201 {object} dto.SessionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/create [post]
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Converter DTOs para model
	var proxyConfig *model.ProxyConfig
	if req.Proxy != nil {
		proxyConfig = &model.ProxyConfig{
			Enabled:  req.Proxy.Enabled,
			Protocol: req.Proxy.Protocol,
			Host:     req.Proxy.Host,
			Port:     req.Proxy.Port,
			Username: req.Proxy.Username,
			Password: req.Proxy.Password,
		}
	}

	var webhookConfig *model.WebhookConfig
	if req.Webhook != nil {
		webhookConfig = &model.WebhookConfig{
			Enabled: req.Webhook.Enabled,
			URL:     req.Webhook.URL,
			Events:  req.Webhook.Events,
			Token:   req.Webhook.Token,
		}
	}

	// Criar sessão
	session := &model.Session{
		Name:          req.Name,
		Status:        string(model.SessionStatusDisconnected),
		Connected:     false,
		ProxyConfig:   proxyConfig,
		WebhookConfig: webhookConfig,
		APIKey:        req.APIKey,
	}

	if err := h.sessionManager.CreateSessionWithConfig(c.Request.Context(), session); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create session")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, toSessionResponse(session))
}

// @Summary Listar sessões
// @Description Retorna a lista de todas as sessões criadas
// @Tags Sessions
// @Produce json
// @Success 200 {object} dto.SessionListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/list [get]
func (h *SessionHandler) GetSessions(c *gin.Context) {
	sessions, err := h.sessionManager.ListSessions(c.Request.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to list sessions")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "list_failed",
			Message: err.Error(),
		})
		return
	}

	responses := make([]dto.SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = toSessionResponse(session)
	}

	c.JSON(http.StatusOK, dto.SessionListResponse{
		Sessions: responses,
		Total:    len(responses),
	})
}

// @Summary Obter detalhes da sessão
// @Description Retorna informações detalhadas de uma sessão específica
// @Tags Sessions
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} dto.SessionResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/info [get]
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")

	session, err := h.sessionManager.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "session_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, toSessionResponse(session))
}

// @Summary Deletar sessão
// @Description Remove uma sessão e todos os seus dados associados
// @Tags Sessions
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/delete [delete]
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")

	if err := h.sessionManager.DeleteSession(c.Request.Context(), sessionID); err != nil {
		logger.Log.Error().Err(err).Str("session_id", sessionID).Msg("Failed to delete session")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Session deleted successfully",
	})
}

// @Summary Conectar sessão
// @Description Inicia a conexão de uma sessão com o WhatsApp
// @Tags Sessions
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/connect [post]
func (h *SessionHandler) ConnectSession(c *gin.Context) {
	sessionID := c.Param("id")

	if err := h.sessionManager.ConnectSession(c.Request.Context(), sessionID); err != nil {
		logger.Log.Error().Err(err).Str("session_id", sessionID).Msg("Failed to connect session")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "connect_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Session connecting",
	})
}

// @Summary Desconectar sessão
// @Description Desconecta uma sessão ativa do WhatsApp
// @Tags Sessions
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/disconnect [post]
func (h *SessionHandler) DisconnectSession(c *gin.Context) {
	sessionID := c.Param("id")

	if err := h.sessionManager.DisconnectSession(c.Request.Context(), sessionID); err != nil {
		logger.Log.Error().Err(err).Str("session_id", sessionID).Msg("Failed to disconnect session")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "disconnect_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Session disconnected",
	})
}

// @Summary Parear com telefone
// @Description Gera um código de pareamento para conectar o WhatsApp usando número de telefone
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.PairPhoneRequest true "Número de telefone"
// @Success 200 {object} dto.PairPhoneResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/pair [post]
func (h *SessionHandler) PairPhone(c *gin.Context) {
	sessionID := c.Param("id")

	var req dto.PairPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	code, err := h.pairingService.PairWithPhone(c.Request.Context(), sessionID, req.PhoneNumber)
	if err != nil {
		logger.Log.Error().Err(err).Str("session_id", sessionID).Msg("Failed to pair with phone")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "pairing_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.PairPhoneResponse{
		SessionID:   sessionID,
		PhoneNumber: req.PhoneNumber,
		PairingCode: code,
		Message:     "Enter the pairing code on your phone",
	})
}

// @Summary Obter status da sessão
// @Description Retorna informações detalhadas sobre o status de conexão da sessão
// @Tags Sessions
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} dto.SessionStatusResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/status [get]
func (h *SessionHandler) GetSessionStatus(c *gin.Context) {
	sessionID := c.Param("id")

	status, err := h.sessionManager.GetSessionStatus(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "session_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// @Summary Atualizar webhook (DEPRECATED)
// @Description Atualiza a URL e eventos do webhook de uma sessão (use /webhook/set)
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.UpdateWebhookRequest true "Configurações de webhook"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/webhook [put]
// @Deprecated
func (h *SessionHandler) UpdateSessionWebhook(c *gin.Context) {
	sessionID := c.Param("id")

	var req dto.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Converter DTO para model
	webhookConfig := &model.WebhookConfig{
		Enabled: req.Webhook.Enabled,
		URL:     req.Webhook.URL,
		Events:  req.Webhook.Events,
		Token:   req.Webhook.Token,
	}

	if err := h.sessionManager.UpdateWebhookConfig(c.Request.Context(), sessionID, webhookConfig); err != nil {
		logger.Log.Error().Err(err).Str("session_id", sessionID).Msg("Failed to update webhook")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Webhook updated successfully",
	})
}

// @Summary Configurar webhook
// @Description Configura ou atualiza o webhook de uma sessão específica com seus eventos subscritos
// @Tags Webhook
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SetWebhookRequest true "Configurações de webhook"
// @Success 200 {object} dto.WebhookConfigResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/webhook/set [post]
func (h *SessionHandler) SetWebhook(c *gin.Context) {
	sessionID := c.Param("id")

	var req dto.SetWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Validar eventos se fornecidos
	if len(req.Events) > 0 {
		for _, event := range req.Events {
			if !constants.IsValidEventType(event) {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "invalid_event",
					Message: fmt.Sprintf("Invalid event type: %s. Use /webhook/events to see supported events", event),
				})
				return
			}
		}
	}

	// Se enabled=true mas URL vazia, retornar erro
	if req.Enabled && req.URL == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "URL is required when webhook is enabled",
		})
		return
	}

	// Converter DTO para model
	webhookConfig := &model.WebhookConfig{
		Enabled: req.Enabled,
		URL:     req.URL,
		Events:  req.Events,
		Token:   req.Token,
	}

	// Se eventos não fornecidos e webhook habilitado, usar eventos padrão
	if req.Enabled && len(webhookConfig.Events) == 0 {
		webhookConfig.Events = constants.DefaultWebhookEvents
	}

	// Atualizar configuração
	if err := h.sessionManager.UpdateWebhookConfig(c.Request.Context(), sessionID, webhookConfig); err != nil {
		logger.Log.Error().Err(err).Str("session_id", sessionID).Msg("Failed to set webhook")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	// Buscar sessão atualizada para retornar
	session, err := h.sessionManager.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "fetch_failed",
			Message: "Webhook configured but failed to fetch updated session",
		})
		return
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Bool("enabled", req.Enabled).
		Str("url", req.URL).
		Int("events_count", len(req.Events)).
		Msg("Webhook configuration updated")

	// Retornar configuração atualizada
	response := dto.WebhookConfigResponse{
		SessionID: sessionID,
		Enabled:   webhookConfig.Enabled,
		URL:       webhookConfig.URL,
		Events:    webhookConfig.Events,
		Token:     webhookConfig.Token,
		UpdatedAt: session.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Obter configuração de webhook
// @Description Retorna a configuração atual de webhook e seus eventos subscritos
// @Tags Webhook
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} dto.WebhookConfigResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/webhook/find [get]
func (h *SessionHandler) FindWebhook(c *gin.Context) {
	sessionID := c.Param("id")

	// Buscar sessão
	session, err := h.sessionManager.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "session_not_found",
			Message: fmt.Sprintf("Session not found: %s", sessionID),
		})
		return
	}

	// Preparar resposta
	response := dto.WebhookConfigResponse{
		SessionID: sessionID,
		Enabled:   false,
		URL:       "",
		Events:    []string{},
		Token:     "",
		UpdatedAt: session.UpdatedAt,
	}

	// Se tem configuração de webhook, preencher
	if session.WebhookConfig != nil {
		response.Enabled = session.WebhookConfig.Enabled
		response.URL = session.WebhookConfig.URL
		response.Events = session.WebhookConfig.Events
		response.Token = session.WebhookConfig.Token
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Bool("enabled", response.Enabled).
		Msg("Webhook configuration retrieved")

	c.JSON(http.StatusOK, response)
}

// @Summary Listar eventos de webhook suportados
// @Description Retorna lista completa de todos os eventos de webhook suportados, organizados por categoria
// @Tags Webhook
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/webhook/events [get]
func (h *SessionHandler) ListWebhookEvents(c *gin.Context) {
	response := gin.H{
		"total":              len(constants.SupportedEventTypes),
		"events":             constants.SupportedEventTypes,
		"events_by_category": constants.AllWebhookEvents,
		"categories":         constants.GetAllCategories(),
		"default_events":     constants.DefaultWebhookEvents,
		"critical_events":    constants.CriticalEvents,
		"description": map[string]string{
			"total":              "Total de eventos suportados",
			"events":             "Lista plana de todos os eventos",
			"events_by_category": "Eventos organizados por categoria",
			"categories":         "Lista de categorias disponíveis",
			"default_events":     "Eventos padrão quando nenhum é especificado",
			"critical_events":    "Eventos críticos que sempre devem ser monitorados",
		},
	}

	c.JSON(http.StatusOK, response)
}

func toSessionResponse(session *model.Session) dto.SessionResponse {
	var webhookURL string
	var webhookEvents []string

	if session.WebhookConfig != nil {
		webhookURL = session.WebhookConfig.URL
		webhookEvents = session.WebhookConfig.Events
	}

	return dto.SessionResponse{
		ID:            session.ID,
		Name:          session.Name,
		JID:           session.DeviceJID,
		Status:        session.Status,
		WebhookURL:    webhookURL,
		WebhookEvents: webhookEvents,
		CreatedAt:     session.CreatedAt,
		UpdatedAt:     session.UpdatedAt,
	}
}

// @Summary Obter QR Code
// @Description Retorna o QR Code atual da sessão (em memória) com imagem base64 decodificada
// @Tags Sessions
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/qr [get]
func (h *SessionHandler) GetQRCode(c *gin.Context) {
	sessionID := c.Param("id")

	// Obter QR code da memória
	qrCode, exists := h.sessionManager.GetQRCode(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "qr_not_found",
			Message: "No QR code available. Please connect the session first.",
		})
		return
	}

	// Gerar imagem PNG do QR code usando go-qrcode
	qrImage, err := qrcode.Encode(qrCode, qrcode.Medium, 512)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to generate QR code image")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "qr_generation_failed",
			Message: "Failed to generate QR code image",
		})
		return
	}

	// Converter para base64
	qrBase64 := base64.StdEncoding.EncodeToString(qrImage)

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"qr_code":    qrCode,
		"qrbase64":   qrBase64,
		"message":    "Scan this QR code with WhatsApp",
	})
}
