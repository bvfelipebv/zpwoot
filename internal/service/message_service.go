package service

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	waProto "go.mau.fi/whatsmeow/binary/proto"

	"zpwoot/pkg/logger"
)

// SendTextMessage envia uma mensagem de texto
func (m *SessionManager) SendTextMessage(ctx context.Context, client *whatsmeow.Client, phone string, message string) (string, time.Time, error) {
	// Parsear JID do destinatário
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	// Criar mensagem
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Enviar mensagem
	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send message: %w", err)
	}

	logger.Log.Info().
		Str("message_id", resp.ID).
		Str("phone", phone).
		Msg("Text message sent")

	return resp.ID, resp.Timestamp, nil
}

// SendImageMessage envia uma mensagem de imagem
func (m *SessionManager) SendImageMessage(ctx context.Context, client *whatsmeow.Client, phone string, imageData []byte, caption string, mimeType string) (string, time.Time, error) {
	// Parsear JID do destinatário
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	// Upload da imagem
	uploaded, err := client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload image: %w", err)
	}

	// Criar mensagem de imagem
	msg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:        proto.String(caption),
			URL:            proto.String(uploaded.URL),
			DirectPath:     proto.String(uploaded.DirectPath),
			MediaKey:       uploaded.MediaKey,
			Mimetype:       proto.String(mimeType),
			FileEncSHA256:  uploaded.FileEncSHA256,
			FileSHA256:     uploaded.FileSHA256,
			FileLength:     proto.Uint64(uint64(len(imageData))),
		},
	}

	// Enviar mensagem
	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send image: %w", err)
	}

	logger.Log.Info().
		Str("message_id", resp.ID).
		Str("phone", phone).
		Msg("Image message sent")

	return resp.ID, resp.Timestamp, nil
}

// SendAudioMessage envia uma mensagem de áudio
func (m *SessionManager) SendAudioMessage(ctx context.Context, client *whatsmeow.Client, phone string, audioData []byte, mimeType string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	uploaded, err := client.Upload(ctx, audioData, whatsmeow.MediaAudio)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload audio: %w", err)
	}

	msg := &waProto.Message{
		AudioMessage: &waProto.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(audioData))),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send audio: %w", err)
	}

	logger.Log.Info().
		Str("message_id", resp.ID).
		Str("phone", phone).
		Msg("Audio message sent")

	return resp.ID, resp.Timestamp, nil
}

// SendImageFromURL envia imagem a partir de URL
func (m *SessionManager) SendImageFromURL(ctx context.Context, client *whatsmeow.Client, phone string, imageURL string, caption string) (string, time.Time, error) {
	// TODO: Download da imagem da URL
	// Por enquanto, retornar erro
	return "", time.Time{}, fmt.Errorf("image from URL not yet implemented")
}

// SendAudioFromURL envia áudio a partir de URL
func (m *SessionManager) SendAudioFromURL(ctx context.Context, client *whatsmeow.Client, phone string, audioURL string) (string, time.Time, error) {
	return "", time.Time{}, fmt.Errorf("audio from URL not yet implemented")
}

// SendVideoFromURL envia vídeo a partir de URL
func (m *SessionManager) SendVideoFromURL(ctx context.Context, client *whatsmeow.Client, phone string, videoURL string, caption string) (string, time.Time, error) {
	return "", time.Time{}, fmt.Errorf("video from URL not yet implemented")
}

// SendDocumentFromURL envia documento a partir de URL
func (m *SessionManager) SendDocumentFromURL(ctx context.Context, client *whatsmeow.Client, phone string, docURL string, fileName string, caption string) (string, time.Time, error) {
	return "", time.Time{}, fmt.Errorf("document from URL not yet implemented")
}

// SendPresence envia presença (digitando, gravando, etc)
func (m *SessionManager) SendPresence(ctx context.Context, client *whatsmeow.Client, phone string, presence string) error {
	recipient, err := parseJID(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	// Para presença global (available/unavailable)
	if presence == "available" {
		return client.SendPresence(ctx, types.PresenceAvailable)
	} else if presence == "unavailable" {
		return client.SendPresence(ctx, types.PresenceUnavailable)
	}

	// Para presença de chat (composing, recording, paused)
	var chatPresence types.ChatPresence
	var media types.ChatPresenceMedia

	switch presence {
	case "composing":
		chatPresence = types.ChatPresenceComposing
		media = types.ChatPresenceMediaText
	case "recording":
		chatPresence = types.ChatPresenceComposing
		media = types.ChatPresenceMediaAudio
	case "paused":
		chatPresence = types.ChatPresencePaused
		media = types.ChatPresenceMediaText
	default:
		return fmt.Errorf("invalid presence type: %s", presence)
	}

	err = client.SendChatPresence(ctx, recipient, chatPresence, media)
	if err != nil {
		return fmt.Errorf("failed to send chat presence: %w", err)
	}

	return nil
}

// parseJID converte um número de telefone em JID do WhatsApp
func parseJID(phone string) (types.JID, error) {
	// Remover caracteres não numéricos
	cleanPhone := ""
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			cleanPhone += string(c)
		}
	}

	if cleanPhone == "" {
		return types.JID{}, fmt.Errorf("invalid phone number")
	}

	// Se não contém @, adicionar @s.whatsapp.net
	if len(cleanPhone) > 0 {
		return types.NewJID(cleanPhone, types.DefaultUserServer), nil
	}

	// Se já contém @, parsear diretamente
	return types.ParseJID(phone)
}

