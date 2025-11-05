package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCodeImage gera uma imagem QR code em formato base64 data URL
func GenerateQRCodeImage(content string) (string, error) {
	// Gerar QR code como PNG
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	
	// Converter para base64
	base64Str := base64.StdEncoding.EncodeToString(png)
	
	// Retornar como data URL
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64Str)
	
	return dataURL, nil
}

// GenerateQRCodePNG gera uma imagem QR code como bytes PNG
func GenerateQRCodePNG(content string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256
	}
	
	png, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}
	
	return png, nil
}

