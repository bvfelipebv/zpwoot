package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"zpwoot/internal/api/dto"
	"zpwoot/internal/service"
	"zpwoot/pkg/logger"
)

type MessageHandler struct {
	sessionManager *service.SessionManager
}

func NewMessageHandler(sessionManager *service.SessionManager) *MessageHandler {
	return &MessageHandler{
		sessionManager: sessionManager,
	}
}

// @Summary Enviar mensagem de texto
// @Description Envia uma mensagem de texto para um contato ou grupo
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendTextRequest true "Dados da mensagem"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/messages/send/text [post]
func (h *MessageHandler) SendText(c *gin.Context) {
	sessionID := c.Param("id")

	var req dto.SendTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Validar campos obrigatórios
	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "phone is required",
		})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "message is required",
		})
		return
	}

	// Obter cliente WhatsApp
	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "session_not_found",
			Message: "Session not found or not connected",
		})
		return
	}

	// Enviar mensagem
	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendTextMessage(ctx, client, req.Phone, req.Message)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("phone", req.Phone).
			Msg("Failed to send text message")

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "send_failed",
			Message: err.Error(),
		})
		return
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Str("phone", req.Phone).
		Str("message_id", messageID).
		Msg("Text message sent successfully")

	c.JSON(http.StatusOK, dto.MessageResponse{
		Success:   true,
		MessageID: messageID,
		Timestamp: timestamp.Unix(),
		Phone:     req.Phone,
	})
}

// @Summary Enviar imagem
// @Description Envia uma imagem para um contato ou grupo
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendImageRequest true "Dados da imagem"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/image [post]
func (h *MessageHandler) SendImage(c *gin.Context) {
	sessionID := c.Param("id")

	var req dto.SendImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	if req.Phone == "" || req.Image == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "phone and image are required",
		})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "session_not_found",
			Message: "Session not found or not connected",
		})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendImageFromURL(ctx, client, req.Phone, req.Image, req.Caption)
	if err != nil {
		logger.Log.Error().Err(err).Str("session_id", sessionID).Msg("Failed to send image")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "send_failed",
			Message: err.Error(),
		})
		return
	}

	logger.Log.Info().Str("session_id", sessionID).Str("message_id", messageID).Msg("Image sent")

	c.JSON(http.StatusOK, dto.MessageResponse{
		Success:   true,
		MessageID: messageID,
		Timestamp: timestamp.Unix(),
		Phone:     req.Phone,
	})
}

// @Summary Enviar áudio
// @Description Envia um áudio para um contato ou grupo
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendAudioRequest true "Dados do áudio"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/audio [post]
func (h *MessageHandler) SendAudio(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendAudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.Audio == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone and audio are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendAudioFromURL(ctx, client, req.Phone, req.Audio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Enviar vídeo
// @Description Envia um vídeo para um contato ou grupo
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendVideoRequest true "Dados do vídeo"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/video [post]
func (h *MessageHandler) SendVideo(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.Video == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone and video are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendVideoFromURL(ctx, client, req.Phone, req.Video, req.Caption)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Enviar documento
// @Description Envia um documento para um contato ou grupo
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendDocumentRequest true "Dados do documento"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/document [post]
func (h *MessageHandler) SendDocument(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.Document == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone and document are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendDocumentFromURL(ctx, client, req.Phone, req.Document, req.FileName, req.Caption)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Enviar sticker
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendStickerRequest true "Dados do sticker"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/sticker [post]
func (h *MessageHandler) SendSticker(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "not_implemented", Message: "Sticker not yet implemented"})
}

// @Summary Enviar mídia genérica
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendMediaRequest true "Dados da mídia"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/media [post]
func (h *MessageHandler) SendMedia(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "not_implemented", Message: "Generic media not yet implemented"})
}

// @Summary Enviar contato
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendContactRequest true "Dados do contato"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/contact [post]
func (h *MessageHandler) SendContact(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.ContactPhone == "" || req.ContactName == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone, contactPhone and contactName are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendContact(ctx, client, req.Phone, req.ContactName, req.ContactPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Enviar localização
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendLocationRequest true "Dados da localização"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/location [post]
func (h *MessageHandler) SendLocation(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone, latitude and longitude are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendLocation(ctx, client, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Enviar enquete
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendPollRequest true "Dados da enquete"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/poll [post]
func (h *MessageHandler) SendPoll(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendPollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.Question == "" || len(req.Options) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone, question and options are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	selectableCount := uint32(req.SelectableCount)
	if selectableCount == 0 {
		selectableCount = 1
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendPoll(ctx, client, req.Phone, req.Question, req.Options, selectableCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Enviar reação
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendReactionRequest true "Dados da reação"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/reaction [post]
func (h *MessageHandler) SendReaction(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.MessageID == "" || req.Emoji == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone, messageId and emoji are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.SendReaction(ctx, client, req.Phone, req.MessageID, req.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Enviar presença
// @Description Envia presença (digitando, gravando, disponível, etc)
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.SendPresenceRequest true "Dados da presença"
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/presence [post]
func (h *MessageHandler) SendPresence(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.SendPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	err = h.sessionManager.SendPresence(ctx, client, req.Phone, req.Presence)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	logger.Log.Info().Str("session_id", sessionID).Str("presence", req.Presence).Msg("Presence sent")
	c.JSON(http.StatusOK, gin.H{"success": true, "presence": req.Presence})
}

// @Summary Marcar mensagens como lidas
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.MarkAsReadRequest true "Dados das mensagens"
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/read [post]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || len(req.MessageIDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone and messageIds are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	err = h.sessionManager.MarkAsRead(ctx, client, req.Phone, req.MessageIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "count": len(req.MessageIDs)})
}

// @Summary Revogar mensagem
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.RevokeMessageRequest true "Dados da mensagem"
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/revoke [delete]
func (h *MessageHandler) RevokeMessage(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.RevokeMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.MessageID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone and messageId are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.RevokeMessage(ctx, client, req.Phone, req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

// @Summary Editar mensagem
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body dto.EditMessageRequest true "Dados da mensagem"
// @Success 200 {object} dto.MessageResponse
// @Security ApiKeyAuth
// @Router /sessions/{id}/message/edit [put]
func (h *MessageHandler) EditMessage(c *gin.Context) {
	sessionID := c.Param("id")
	var req dto.EditMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if req.Phone == "" || req.MessageID == "" || req.NewMessage == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: "phone, messageId and newMessage are required"})
		return
	}

	client, err := h.sessionManager.GetClient(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "session_not_found", Message: err.Error()})
		return
	}

	ctx := context.Background()
	messageID, timestamp, err := h.sessionManager.EditMessage(ctx, client, req.Phone, req.MessageID, req.NewMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "send_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Success: true, MessageID: messageID, Timestamp: timestamp.Unix(), Phone: req.Phone})
}

