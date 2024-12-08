package application

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Errorf("failed to generate token: %w", err))
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	return token[:length]
}
