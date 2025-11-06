package dto

type SendTextRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Message string `json:"message" binding:"required" example:"Hello, World!"`
}

type SendImageRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Image   string `json:"image" binding:"required" example:"https://example.com/image.jpg"`
	Caption string `json:"caption,omitempty" example:"Check this out!"`
}

type SendAudioRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
	Audio string `json:"audio" binding:"required" example:"https://example.com/audio.mp3"`
}

type SendVideoRequest struct {
	Phone   string `json:"phone" binding:"required" example:"5511999999999"`
	Video   string `json:"video" binding:"required" example:"https://example.com/video.mp4"`
	Caption string `json:"caption,omitempty" example:"Check this video!"`
}

type SendDocumentRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Document string `json:"document" binding:"required" example:"https://example.com/doc.pdf"`
	FileName string `json:"fileName,omitempty" example:"document.pdf"`
	Caption  string `json:"caption,omitempty" example:"Important document"`
}

type SendStickerRequest struct {
	Phone         string `json:"phone" binding:"required" example:"5511999999999"`
	Sticker       string `json:"sticker,omitempty" example:"https://example.com/sticker.webp"`
	StickerBase64 string `json:"stickerBase64,omitempty" example:"data:image/webp;base64,UklGRiQAAABXRUJQVlA4IBgAAAAwAQCdASoBAAEAAwA0JaQAA3AA/vuUAAA="`
}

type SendMediaRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Media    string `json:"media" binding:"required" example:"https://example.com/file.jpg"`
	Caption  string `json:"caption,omitempty" example:"Check this out!"`
	FileName string `json:"fileName,omitempty" example:"file.jpg"`
}

type SendLocationRequest struct {
	Phone     string  `json:"phone" binding:"required" example:"5511999999999"`
	Latitude  float64 `json:"latitude" binding:"required" example:"-23.5505"`
	Longitude float64 `json:"longitude" binding:"required" example:"-46.6333"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
}

type ContactInfo struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Phone string `json:"phone" binding:"required" example:"5511888888888"`
	Vcard string `json:"vcard,omitempty" example:"BEGIN:VCARD\\nVERSION:3.0\\nFN:John Doe\\nTEL;waid=5511888888888:+55 11 88888-8888\\nEND:VCARD"`
}

type SendContactRequest struct {
	Phone    string        `json:"phone" binding:"required" example:"5511999999999"`
	Contacts []ContactInfo `json:"contacts" binding:"required,min=1"`
}

type SendReactionRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0XXXXX"`
	Emoji     string `json:"emoji" binding:"required" example:"👍"`
}

type SendPollRequest struct {
	Phone           string   `json:"phone" binding:"required" example:"5511999999999"`
	Question        string   `json:"question" binding:"required" example:"What's your favorite color?"`
	Options         []string `json:"options" binding:"required" example:"Red,Blue,Green"`
	SelectableCount int      `json:"selectableCount,omitempty" example:"1"`
}

type MarkAsReadRequest struct {
	Phone      string   `json:"phone" binding:"required" example:"5511999999999"`
	MessageIDs []string `json:"messageIds" binding:"required" example:"3EB0XXXXX,3EB0YYYYY"`
}

type RevokeMessageRequest struct {
	Phone     string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID string `json:"messageId" binding:"required" example:"3EB0XXXXX"`
}

type EditMessageRequest struct {
	Phone      string `json:"phone" binding:"required" example:"5511999999999"`
	MessageID  string `json:"messageId" binding:"required" example:"3EB0XXXXX"`
	NewMessage string `json:"newMessage" binding:"required" example:"Updated message"`
}

type SendPresenceRequest struct {
	Phone    string `json:"phone" binding:"required" example:"5511999999999"`
	Presence string `json:"presence" binding:"required" example:"available" enums:"available,unavailable,composing,recording,paused"`
}

type MessageResponse struct {
	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"messageId" example:"3EB0XXXXX"`
	Timestamp int64  `json:"timestamp" example:"1699999999"`
	Phone     string `json:"phone" example:"5511999999999"`
}

type PollResultsResponse struct {
	Question string             `json:"question" example:"What's your favorite color?"`
	Options  []PollOptionResult `json:"options"`
	Voters   []PollVoter        `json:"voters"`
}

type PollOptionResult struct {
	Name  string `json:"name" example:"Red"`
	Votes int    `json:"votes" example:"5"`
}

type PollVoter struct {
	Phone     string   `json:"phone" example:"5511999999999"`
	Options   []string `json:"options" example:"Red"`
	Timestamp int64    `json:"timestamp" example:"1699999999"`
}
