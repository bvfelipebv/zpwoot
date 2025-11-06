package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"zpwoot/pkg/logger"
)

// cleanPhone remove caracteres não numéricos do telefone
func cleanPhone(phone string) string {
	var cleaned strings.Builder
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			cleaned.WriteRune(c)
		}
	}
	return cleaned.String()
}

// parseJID converte um número de telefone em JID do WhatsApp
func parseJID(phone string) (types.JID, error) {
	cleaned := cleanPhone(phone)
	if cleaned == "" {
		return types.JID{}, fmt.Errorf("invalid phone number")
	}
	return types.NewJID(cleaned, types.DefaultUserServer), nil
}

// buildImageMessage cria uma mensagem de imagem
func buildImageMessage(uploaded whatsmeow.UploadResponse, imageData []byte, caption, mimeType string) *waProto.Message {
	return &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:       proto.String(caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(imageData))),
		},
	}
}

// buildAudioMessage cria uma mensagem de áudio
func buildAudioMessage(uploaded whatsmeow.UploadResponse, audioData []byte, mimeType string, isPTT bool) *waProto.Message {
	audioMsg := &waProto.AudioMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(mimeType),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(audioData))),
	}

	if isPTT {
		audioMsg.PTT = proto.Bool(true)
		audioMsg.Seconds = proto.Uint32(0)
	}

	return &waProto.Message{AudioMessage: audioMsg}
}

// buildVideoMessage cria uma mensagem de vídeo
func buildVideoMessage(uploaded whatsmeow.UploadResponse, videoData []byte, caption, mimeType string) *waProto.Message {
	return &waProto.Message{
		VideoMessage: &waProto.VideoMessage{
			Caption:       proto.String(caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(videoData))),
			Seconds:       proto.Uint32(0),
		},
	}
}

// buildDocumentMessage cria uma mensagem de documento
func buildDocumentMessage(uploaded whatsmeow.UploadResponse, docData []byte, fileName, caption, mimeType string) *waProto.Message {
	return &waProto.Message{
		DocumentMessage: &waProto.DocumentMessage{
			Caption:       proto.String(caption),
			FileName:      proto.String(fileName),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(docData))),
		},
	}
}

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

func (m *SessionManager) SendImageMessage(ctx context.Context, client *whatsmeow.Client, phone string, imageData []byte, caption string, mimeType string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	uploaded, err := client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload image: %w", err)
	}

	msg := buildImageMessage(uploaded, imageData, caption, mimeType)

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

func (m *SessionManager) SendAudioMessage(ctx context.Context, client *whatsmeow.Client, phone string, audioData []byte, mimeType string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	uploaded, err := client.Upload(ctx, audioData, whatsmeow.MediaAudio)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload audio: %w", err)
	}

	msg := buildAudioMessage(uploaded, audioData, mimeType, false)

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

func (m *SessionManager) SendImageFromURL(ctx context.Context, client *whatsmeow.Client, phone string, imageURL string, caption string) (string, time.Time, error) {
	imageData, mimeType, err := downloadOrDecodeMedia(imageURL)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get image: %w", err)
	}

	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	uploaded, err := client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload image: %w", err)
	}

	msg := buildImageMessage(uploaded, imageData, caption, mimeType)

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send image: %w", err)
	}

	logger.Log.Info().
		Str("message_id", resp.ID).
		Str("phone", phone).
		Int("size", len(imageData)).
		Str("mime", mimeType).
		Msg("Image message sent")

	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) SendAudioFromURL(ctx context.Context, client *whatsmeow.Client, phone string, audioURL string) (string, time.Time, error) {
	audioData, mimeType, err := downloadOrDecodeMedia(audioURL)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get audio: %w", err)
	}

	if mimeType == "" {
		mimeType = "audio/ogg; codecs=opus"
	}

	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	uploaded, err := client.Upload(ctx, audioData, whatsmeow.MediaAudio)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload audio: %w", err)
	}

	msg := buildAudioMessage(uploaded, audioData, mimeType, true)

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send audio: %w", err)
	}

	logger.Log.Info().
		Str("message_id", resp.ID).
		Str("phone", phone).
		Int("size", len(audioData)).
		Str("mime", mimeType).
		Msg("Audio message sent")

	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) SendVideoFromURL(ctx context.Context, client *whatsmeow.Client, phone string, videoURL string, caption string) (string, time.Time, error) {
	videoData, mimeType, err := downloadOrDecodeMedia(videoURL)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get video: %w", err)
	}

	if mimeType == "" {
		mimeType = "video/mp4"
	}

	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	uploaded, err := client.Upload(ctx, videoData, whatsmeow.MediaVideo)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload video: %w", err)
	}

	msg := buildVideoMessage(uploaded, videoData, caption, mimeType)

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send video: %w", err)
	}

	logger.Log.Info().
		Str("message_id", resp.ID).
		Str("phone", phone).
		Int("size", len(videoData)).
		Str("mime", mimeType).
		Msg("Video message sent")

	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) SendDocumentFromURL(ctx context.Context, client *whatsmeow.Client, phone string, docURL string, fileName string, caption string) (string, time.Time, error) {
	docData, mimeType, err := downloadOrDecodeMedia(docURL)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get document: %w", err)
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if fileName == "" {
		fileName = "document"
	}

	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	uploaded, err := client.Upload(ctx, docData, whatsmeow.MediaDocument)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload document: %w", err)
	}

	msg := buildDocumentMessage(uploaded, docData, fileName, caption, mimeType)

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send document: %w", err)
	}

	logger.Log.Info().Str("message_id", resp.ID).Str("phone", phone).Str("fileName", fileName).Msg("Document sent")
	return resp.ID, resp.Timestamp, nil
}

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

func (m *SessionManager) SendLocation(ctx context.Context, client *whatsmeow.Client, phone string, latitude float64, longitude float64, name string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	msg := &waProto.Message{
		LocationMessage: &waProto.LocationMessage{
			DegreesLatitude:  proto.Float64(latitude),
			DegreesLongitude: proto.Float64(longitude),
			Name:             proto.String(name),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send location: %w", err)
	}

	logger.Log.Info().Str("message_id", resp.ID).Str("phone", phone).Msg("Location sent")
	return resp.ID, resp.Timestamp, nil
}

type ContactData struct {
	Name  string
	Phone string
	Vcard string
}

func (m *SessionManager) SendContact(ctx context.Context, client *whatsmeow.Client, phone string, contactName string, contactPhone string, customVcard string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	vcard := generateVcard(contactName, contactPhone, customVcard)

	msg := &waProto.Message{
		ContactMessage: &waProto.ContactMessage{
			DisplayName: proto.String(contactName),
			Vcard:       proto.String(vcard),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send contact: %w", err)
	}

	logger.Log.Info().Str("message_id", resp.ID).Str("phone", phone).Str("contact", contactName).Msg("Single contact sent")
	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) SendContactsList(ctx context.Context, client *whatsmeow.Client, phone string, contacts []ContactData) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	if len(contacts) == 0 {
		return "", time.Time{}, fmt.Errorf("contacts list cannot be empty")
	}

	// Criar array de ContactMessage
	contactMessages := make([]*waProto.ContactMessage, len(contacts))
	for i, contact := range contacts {
		vcard := generateVcard(contact.Name, contact.Phone, contact.Vcard)
		contactMessages[i] = &waProto.ContactMessage{
			DisplayName: proto.String(contact.Name),
			Vcard:       proto.String(vcard),
		}
	}

	// Usar o nome do primeiro contato como DisplayName da lista
	displayName := contacts[0].Name
	if len(contacts) > 1 {
		displayName = fmt.Sprintf("%s e mais %d", contacts[0].Name, len(contacts)-1)
	}

	msg := &waProto.Message{
		ContactsArrayMessage: &waProto.ContactsArrayMessage{
			DisplayName: proto.String(displayName),
			Contacts:    contactMessages,
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send contacts list: %w", err)
	}

	logger.Log.Info().Str("message_id", resp.ID).Str("phone", phone).Int("count", len(contacts)).Msg("Contacts list sent")
	return resp.ID, resp.Timestamp, nil
}

func formatPhoneForDisplay(phone string) string {
	if len(phone) >= 12 && phone[:2] == "55" {
		ddd := phone[2:4]
		if len(phone) == 13 {
			return fmt.Sprintf("+55 %s %s-%s", ddd, phone[4:9], phone[9:13])
		} else if len(phone) == 12 {
			return fmt.Sprintf("+55 %s %s-%s", ddd, phone[4:8], phone[8:12])
		}
	}

	if len(phone) > 0 {
		return "+" + phone
	}

	return phone
}

func generateVcard(contactName string, contactPhone string, customVcard string) string {
	if customVcard != "" {
		return customVcard
	}

	cleaned := cleanPhone(contactPhone)
	formatted := formatPhoneForDisplay(cleaned)

	return fmt.Sprintf(`BEGIN:VCARD
VERSION:3.0
FN:%s
TEL;type=CELL;type=VOICE;waid=%s:%s
END:VCARD`, contactName, cleaned, formatted)
}

func (m *SessionManager) SendSticker(ctx context.Context, client *whatsmeow.Client, phone string, stickerURL string, stickerBase64 string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	var imageData []byte
	var mimeType string

	// Determinar qual fonte usar (base64 tem prioridade)
	mediaSource := stickerURL
	if stickerBase64 != "" {
		mediaSource = stickerBase64
	}

	if mediaSource == "" {
		return "", time.Time{}, fmt.Errorf("sticker (URL) or stickerBase64 is required")
	}

	// Download ou decode do sticker
	imageData, mimeType, err = downloadOrDecodeMedia(mediaSource)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get sticker: %w", err)
	}

	// Stickers devem ser WebP
	if mimeType != "image/webp" && mimeType != "image/png" && mimeType != "image/jpeg" {
		logger.Log.Warn().Str("mime", mimeType).Msg("Sticker should be image/webp, image/png or image/jpeg")
	}

	// Upload do sticker
	uploaded, err := client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to upload sticker: %w", err)
	}

	msg := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String("image/webp"),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(imageData))),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send sticker: %w", err)
	}

	logger.Log.Info().Str("message_id", resp.ID).Str("phone", phone).Int("size", len(imageData)).Msg("Sticker sent")
	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) SendPoll(ctx context.Context, client *whatsmeow.Client, phone string, question string, options []string, selectableCount uint32) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	// Criar opções da enquete
	pollOptions := make([]*waProto.PollCreationMessage_Option, len(options))
	for i, opt := range options {
		pollOptions[i] = &waProto.PollCreationMessage_Option{
			OptionName: proto.String(opt),
		}
	}

	msg := &waProto.Message{
		PollCreationMessage: &waProto.PollCreationMessage{
			Name:                   proto.String(question),
			Options:                pollOptions,
			SelectableOptionsCount: proto.Uint32(selectableCount),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send poll: %w", err)
	}

	logger.Log.Info().Str("message_id", resp.ID).Str("phone", phone).Msg("Poll sent")
	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) SendReaction(ctx context.Context, client *whatsmeow.Client, phone string, messageID string, emoji string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	// Determinar se a mensagem é nossa (FromMe)
	// Se messageID começa com "me:", é nossa mensagem
	fromMe := false
	actualMessageID := messageID
	if len(messageID) > 3 && messageID[:3] == "me:" {
		fromMe = true
		actualMessageID = messageID[3:] // Remove o prefixo "me:"
	}

	// Se emoji vazio, remove a reação
	reactionText := emoji
	if emoji == "remove" || emoji == "" {
		reactionText = ""
	}

	msg := &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJID: proto.String(recipient.String()),
				FromMe:    proto.Bool(fromMe),
				ID:        proto.String(actualMessageID),
			},
			Text:              proto.String(reactionText),
			GroupingKey:       proto.String(reactionText),
			SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to send reaction: %w", err)
	}

	logger.Log.Info().Str("message_id", resp.ID).Str("phone", phone).Str("emoji", emoji).Bool("fromMe", fromMe).Msg("Reaction sent")
	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) MarkAsRead(ctx context.Context, client *whatsmeow.Client, phone string, messageIDs []string) error {
	recipient, err := parseJID(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	// Criar array de IDs
	ids := make([]types.MessageID, len(messageIDs))
	for i, id := range messageIDs {
		ids[i] = types.MessageID(id)
	}

	err = client.MarkRead(ctx, ids, time.Now(), recipient, recipient)
	if err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	logger.Log.Info().Str("phone", phone).Int("count", len(messageIDs)).Msg("Messages marked as read")
	return nil
}

func (m *SessionManager) RevokeMessage(ctx context.Context, client *whatsmeow.Client, phone string, messageID string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	msg := &waProto.Message{
		ProtocolMessage: &waProto.ProtocolMessage{
			Type: waProto.ProtocolMessage_REVOKE.Enum(),
			Key: &waProto.MessageKey{
				RemoteJID: proto.String(recipient.String()),
				FromMe:    proto.Bool(true),
				ID:        proto.String(messageID),
			},
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to revoke message: %w", err)
	}

	logger.Log.Info().Str("message_id", messageID).Str("phone", phone).Msg("Message revoked")
	return resp.ID, resp.Timestamp, nil
}

func (m *SessionManager) EditMessage(ctx context.Context, client *whatsmeow.Client, phone string, messageID string, newText string) (string, time.Time, error) {
	recipient, err := parseJID(phone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid phone number: %w", err)
	}

	// Usar EditMessage do whatsmeow (método correto)
	resp, err := client.SendMessage(ctx, recipient, client.BuildEdit(recipient, messageID, &waProto.Message{
		Conversation: proto.String(newText),
	}))

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to edit message: %w", err)
	}

	logger.Log.Info().Str("message_id", messageID).Str("phone", phone).Msg("Message edited")
	return resp.ID, resp.Timestamp, nil
}
