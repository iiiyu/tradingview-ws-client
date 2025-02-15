package tvwsclient

import (
	"fmt"
	"math/rand"
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
