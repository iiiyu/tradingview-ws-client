package tvwsclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateSession generates a random session ID with the given prefix
func GenerateSession(prefix string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return prefix + string(b)
}

func wrappedMessage(msg string) string {
	wrappedMsg := fmt.Sprintf("~m~%d~m~%s", len(msg), msg)
	return wrappedMsg
}

func CheckAuthTokenExpired(authToken string) bool {
	// Split the JWT token into parts
	parts := strings.Split(authToken, ".")
	if len(parts) != 3 {
		return true
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return true
	}

	// Parse the payload
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true
	}

	// Check if token is expired
	return time.Now().Unix() >= claims.Exp
}
