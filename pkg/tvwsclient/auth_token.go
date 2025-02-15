package tvwsclient

import (
	"encoding/base64"
	"encoding/json"
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
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.authToken == "" || m.CheckAuthTokenExpired() {
		token, err := m.client.GetQuoteToken()
		if err != nil {
			panic(err)
		}
		m.SetToken(token)
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
