package service

import (
	"time"

	"go.mau.fi/whatsmeow/types/events"
	"zpwoot/internal/constants"
)

type WebhookFormatter struct{}

func NewWebhookFormatter() *WebhookFormatter {
	return &WebhookFormatter{}
}

type WebhookPayload struct {
	Event     string                 `json:"event"`
	SessionID string                 `json:"session_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func (f *WebhookFormatter) FormatMessage(sessionID string, evt *events.Message) *WebhookPayload {
	data := map[string]interface{}{
		"message_id": evt.Info.ID,
		"from":       evt.Info.Sender.String(),
		"from_me":    evt.Info.IsFromMe,
		"chat":       evt.Info.Chat.String(),
		"timestamp":  evt.Info.Timestamp,
	}

	// Add message content based on type
	if evt.Message.Conversation != nil {
		data["type"] = "conversation"
		data["body"] = *evt.Message.Conversation
	} else if evt.Message.ExtendedTextMessage != nil {
		data["type"] = "extended_text"
		data["body"] = *evt.Message.ExtendedTextMessage.Text
	} else if evt.Message.ImageMessage != nil {
		data["type"] = "image"
		data["caption"] = evt.Message.ImageMessage.Caption
		data["mime_type"] = evt.Message.ImageMessage.Mimetype
	} else if evt.Message.VideoMessage != nil {
		data["type"] = "video"
		data["caption"] = evt.Message.VideoMessage.Caption
		data["mime_type"] = evt.Message.VideoMessage.Mimetype
	} else if evt.Message.AudioMessage != nil {
		data["type"] = "audio"
		data["mime_type"] = evt.Message.AudioMessage.Mimetype
	} else if evt.Message.DocumentMessage != nil {
		data["type"] = "document"
		data["file_name"] = evt.Message.DocumentMessage.FileName
		data["mime_type"] = evt.Message.DocumentMessage.Mimetype
	} else {
		data["type"] = "unknown"
	}

	return &WebhookPayload{
		Event:     string(constants.EventMessage),
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      data,
	}
}

func (f *WebhookFormatter) FormatReceipt(sessionID string, evt *events.Receipt) *WebhookPayload {
	data := map[string]interface{}{
		"message_ids": evt.MessageIDs,
		"timestamp":   evt.Timestamp,
		"chat":        evt.Chat.String(),
		"sender":      evt.Sender.String(),
		"type":        string(evt.Type),
	}

	return &WebhookPayload{
		Event:     string(constants.EventReceipt),
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      data,
	}
}

func (f *WebhookFormatter) FormatConnected(sessionID string, evt *events.Connected) *WebhookPayload {
	data := map[string]interface{}{
		"status": "connected",
	}

	return &WebhookPayload{
		Event:     string(constants.EventConnected),
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      data,
	}
}

func (f *WebhookFormatter) FormatDisconnected(sessionID string, evt *events.Disconnected) *WebhookPayload {
	data := map[string]interface{}{
		"status": "disconnected",
	}

	return &WebhookPayload{
		Event:     string(constants.EventDisconnected),
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      data,
	}
}

func (f *WebhookFormatter) FormatGroupInfo(sessionID string, evt *events.GroupInfo) *WebhookPayload {
	data := map[string]interface{}{
		"jid":       evt.JID.String(),
		"name":      evt.Name,
		"topic":     evt.Topic,
		"timestamp": evt.Timestamp,
	}

	if evt.Sender != nil {
		data["sender"] = evt.Sender.String()
	}

	return &WebhookPayload{
		Event:     string(constants.EventGroupInfo),
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      data,
	}
}

func (f *WebhookFormatter) FormatPicture(sessionID string, evt *events.Picture) *WebhookPayload {
	data := map[string]interface{}{
		"jid":       evt.JID.String(),
		"timestamp": evt.Timestamp,
		"remove":    evt.Remove,
	}

	return &WebhookPayload{
		Event:     string(constants.EventPicture),
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      data,
	}
}
