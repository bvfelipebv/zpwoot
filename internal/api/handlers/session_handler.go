package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zpwoot/internal/api/dto"
	"zpwoot/internal/model"
	"zpwoot/internal/service"
	"zpwoot/pkg/logger"
)

// SessionHandler gerencia requisições HTTP de sessões
type SessionHandler struct {
	sessionManager *service.SessionManager
	pairingService *service.PairingService
}

// NewSessionHandler cria um novo handler de sessões
func NewSessionHandler(sessionManager *service.SessionManager, pairingService *service.PairingService) *SessionHandler {
	return &SessionHandler{
		sessionManager: sessionManager,
		pairingService: pairingService,
	}
}

// CreateSession cria uma nova sessão
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

// GetSessions lista todas as sessões
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

// GetSession obtém detalhes de uma sessão
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

// DeleteSession deleta uma sessão
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

// ConnectSession conecta uma sessão ao WhatsApp
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

// DisconnectSession desconecta uma sessão
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

// PairPhone inicia pareamento com número de telefone
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

// GetSessionStatus obtém status detalhado da sessão
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

// UpdateSessionWebhook atualiza configurações de webhook
// @Summary Atualizar webhook
// @Description Atualiza a URL e eventos do webhook de uma sessão
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

// toSessionResponse converte model.Session para dto.SessionResponse
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
