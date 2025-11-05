package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"zpwoot/pkg/logger"
)

// downloadOrDecodeMedia baixa mídia de URL ou decodifica base64
// Retorna: dados, mimeType, erro
func downloadOrDecodeMedia(mediaURL string) ([]byte, string, error) {
	// Verifica se é data URL (base64)
	if strings.HasPrefix(mediaURL, "data:") {
		return decodeBase64Media(mediaURL)
	}

	// Caso contrário, faz download da URL
	return downloadFromURL(mediaURL)
}

// decodeBase64Media decodifica data URL base64
// Formato esperado: data:image/jpeg;base64,XXXXX
func decodeBase64Media(dataURL string) ([]byte, string, error) {
	// Encontrar o início dos dados base64
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("invalid data URL format")
	}

	// Extrair MIME type
	mimeType := ""
	header := parts[0]
	if strings.Contains(header, ";") {
		mimePart := strings.Split(header, ";")[0]
		mimeType = strings.TrimPrefix(mimePart, "data:")
	}

	// Decodificar base64
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Se não conseguiu extrair MIME type do header, detectar dos dados
	if mimeType == "" {
		mimeType = detectMimeType(data)
	}

	logger.Log.Debug().
		Str("mimeType", mimeType).
		Int("size", len(data)).
		Msg("Decoded base64 media")

	return data, mimeType, nil
}

// downloadFromURL faz download de URL HTTP/HTTPS
func downloadFromURL(url string) ([]byte, string, error) {
	// Validar URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, "", fmt.Errorf("invalid URL scheme, must be http or https")
	}

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Fazer requisição
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Ler dados
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response: %w", err)
	}

	// Obter MIME type do header ou detectar
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = detectMimeType(data)
	}

	logger.Log.Debug().
		Str("url", url).
		Str("mimeType", mimeType).
		Int("size", len(data)).
		Msg("Downloaded media from URL")

	return data, mimeType, nil
}

// detectMimeType detecta o tipo MIME dos dados
func detectMimeType(data []byte) string {
	if len(data) == 0 {
		return "application/octet-stream"
	}

	// Usar http.DetectContentType (detecta até 512 bytes)
	mimeType := http.DetectContentType(data)

	// Ajustes para tipos específicos do WhatsApp
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		// Manter como está
	case strings.HasPrefix(mimeType, "audio/"):
		// Para áudio, WhatsApp prefere audio/ogg; codecs=opus para PTT
		if mimeType == "audio/ogg" {
			mimeType = "audio/ogg; codecs=opus"
		}
	case strings.HasPrefix(mimeType, "video/"):
		// Manter como está
	default:
		// Para outros tipos, usar application/octet-stream
		if mimeType == "application/octet-stream" {
			// Tentar detectar tipos específicos por magic numbers
			if len(data) >= 4 {
				// PDF
				if string(data[0:4]) == "%PDF" {
					return "application/pdf"
				}
				// ZIP
				if data[0] == 0x50 && data[1] == 0x4B {
					return "application/zip"
				}
			}
		}
	}

	return mimeType
}

