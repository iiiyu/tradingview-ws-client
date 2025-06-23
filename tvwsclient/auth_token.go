package tvwsclient

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// AuthTokenManager is a singleton that manages the auth token
type AuthTokenManager struct {
	authToken string
	mu        sync.RWMutex
	client    *TVHttpClient
}

var (
	instance *AuthTokenManager
	once     sync.Once
)

func InitAuthTokenManager(client *TVHttpClient) {
	once.Do(func() {
		instance = &AuthTokenManager{
			client: client,
		}
		token, err := instance.client.GetQuoteToken()
		if err != nil {
			panic(err)
		}
		instance.SetToken(token)
	})
}

// GetAuthTokenManager returns the singleton instance of AuthTokenManager
func GetAuthTokenManager() *AuthTokenManager {
	if instance == nil {
		panic("AuthTokenManager not initialized")
	}
	return instance
}

// SetToken sets a new auth token
func (m *AuthTokenManager) SetToken(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.authToken = token
}

// GetToken returns the current auth token, checking if it needs to be updated
func (m *AuthTokenManager) GetToken() string {
	// First check with read lock
	m.mu.RLock()
	if m.authToken != "" && !m.CheckAuthTokenExpired() {
		token := m.authToken
		m.mu.RUnlock()
		return token
	}
	m.mu.RUnlock()

	// If token is empty or expired, acquire write lock
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check condition after acquiring write lock
	if m.authToken == "" || m.CheckAuthTokenExpired() {
		token, err := m.client.GetQuoteToken()
		if err != nil {
			slog.Error("Failed to get quote token", "error", err)
			panic(err)
		}
		m.authToken = token // Direct assignment since we already have the write lock
	}
	return m.authToken
}

func (m *AuthTokenManager) CheckAuthTokenExpired() bool {
	// Split the JWT token into parts
	parts := strings.Split(m.authToken, ".")
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
	return time.Now().Add(-5*time.Minute).Unix() >= claims.Exp
}
