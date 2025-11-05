package dto

// SendTextRequest representa uma requisição para enviar mensagem de texto
type SendTextRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Message string `json:"message" binding:"required" example:"Hello, World!"`
}

// SendImageRequest representa uma requisição para enviar imagem
type SendImageRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Image   string `json:"image" binding:"required" example:"https://example.com/image.jpg"`
	Caption string `json:"caption,omitempty" example:"Check this out!"`
}

// SendAudioRequest representa uma requisição para enviar áudio
type SendAudioRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	Audio string `json:"audio" binding:"required" example:"https://example.com/audio.mp3"`
}

// SendVideoRequest representa uma requisição para enviar vídeo
type SendVideoRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Video   string `json:"video" binding:"required" example:"https://example.com/video.mp4"`
	Caption string `json:"caption,omitempty" example:"Check this video!"`
}

// SendDocumentRequest representa uma requisição para enviar documento
type SendDocumentRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Document string `json:"document" binding:"required" example:"https://example.com/doc.pdf"`
	FileName string `json:"fileName,omitempty" example:"document.pdf"`
	Caption  string `json:"caption,omitempty" example:"Important document"`
}

// SendStickerRequest representa uma requisição para enviar sticker
type SendStickerRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Sticker string `json:"sticker" binding:"required" example:"https://example.com/sticker.webp"`
}

// SendMediaRequest representa uma requisição para enviar mídia genérica
type SendMediaRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Media    string `json:"media" binding:"required" example:"https://example.com/file.jpg"`
	Caption  string `json:"caption,omitempty" example:"Check this out!"`
	FileName string `json:"fileName,omitempty" example:"file.jpg"`
}

// SendLocationRequest representa uma requisição para enviar localização
type SendLocationRequest struct {
	Phone     string  `json:"phone" binding:"required" example:"5511999999999"`
	Latitude  float64 `json:"latitude" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	Address   string  `json:"address,omitempty" example:"São Paulo, Brazil"`
}

// SendContactRequest representa uma requisição para enviar contato
type SendContactRequest struct {
	Phone        string `json:"phone" binding:"required" example:"5511999999999"`
	ContactPhone string `json:"contactPhone" binding:"required" example:"5511888888888"`
	ContactName  string `json:"contactName" binding:"required" example:"John Doe"`
}

// SendReactionRequest representa uma requisição para enviar reação
type SendReactionRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0XXXXX"`
	Emoji     string `json:"emoji" binding:"required" example:"👍"`
}

// SendPollRequest representa uma requisição para enviar enquete
type SendPollRequest struct {
	Phone          string   `json:"phone" binding:"required" example:"5511999999999"`
	Question       string   `json:"question" binding:"required" example:"What's your favorite color?"`
	Options        []string `json:"options" binding:"required" example:"Red,Blue,Green"`
	SelectableCount int     `json:"selectableCount,omitempty" example:"1"`
}

// MarkAsReadRequest representa uma requisição para marcar mensagem como lida
type MarkAsReadRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"5511999999999"`
	MessageIDs []string `json:"messageIds" binding:"required" example:"3EB0XXXXX,3EB0YYYYY"`
}

// RevokeMessageRequest representa uma requisição para revogar mensagem
type RevokeMessageRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0XXXXX"`
}

// EditMessageRequest representa uma requisição para editar mensagem
type EditMessageRequest struct {
	Phone      string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID  string `json:"messageId" binding:"required" example:"3EB0XXXXX"`
	NewMessage string `json:"newMessage" binding:"required" example:"Updated message"`
}

// SendPresenceRequest representa uma requisição para enviar presença
type SendPresenceRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Presence string `json:"presence" binding:"required" example:"available" enums:"available,unavailable,composing,recording,paused"`
}

// MessageResponse representa a resposta de envio de mensagem
type MessageResponse struct {
	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"messageId" example:"3EB0XXXXX"`
	Timestamp int64  `json:"timestamp" example:"1699999999"`
	Phone     string `json:"phone" example:"5511999999999"`
}

// PollResultsResponse representa os resultados de uma enquete
type PollResultsResponse struct {
	Question string              `json:"question" example:"What's your favorite color?"`
	Options  []PollOptionResult  `json:"options"`
	Voters   []PollVoter         `json:"voters"`
}

// PollOptionResult representa o resultado de uma opção da enquete
type PollOptionResult struct {
	Name  string `json:"name" example:"Red"`
	Votes int    `json:"votes" example:"5"`
}

// PollVoter representa um votante da enquete
type PollVoter struct {
	Phone     string   `json:"phone" example:"5511999999999"`
	Options   []string `json:"options" example:"Red"`
	Timestamp int64    `json:"timestamp" example:"1699999999"`
}

