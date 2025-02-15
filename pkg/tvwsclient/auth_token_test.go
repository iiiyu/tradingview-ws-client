package tvwsclient

import (
	"sync"
	"testing"
)

func TestAuthTokenManager(t *testing.T) {

	// Test token operations
	t.Run("token operations", func(t *testing.T) {
		manager := GetAuthTokenManager()

		// Test initial state
		if token := manager.GetToken(); token != "" {
			t.Errorf("Initial token should be empty, got %v", token)
		}

		// Test setting and getting token
		testToken := "test_token"
		manager.SetToken(testToken)
		if token := manager.GetToken(); token != testToken {
			t.Errorf("GetToken() = %v, want %v", token, testToken)
		}
	})

	// Test singleton pattern
	t.Run("singleton pattern", func(t *testing.T) {
		manager1 := GetAuthTokenManager()
		manager2 := GetAuthTokenManager()

		if manager1 != manager2 {
			t.Error("GetAuthTokenManager() should return the same instance")
		}
	})

	// Test concurrent access
	t.Run("concurrent access", func(t *testing.T) {
		manager := GetAuthTokenManager()
		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func(i int) {
				defer wg.Done()
				token := manager.GetToken()
				manager.SetToken("new_token")
				_ = token // Use token to prevent compiler optimization
			}(i)
		}

		wg.Wait()
	})

	// Test CheckAuthTokenExpired
	t.Run("check token expiration", func(t *testing.T) {
		tests := []struct {
			name        string
			token       string
			want        bool
			description string
		}{
			{
				name:        "valid token not expired",
				token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4Nzc4MzY4MDB9.QESXXx9ErawCjaMDrf51kakqBhpcBerH9Xtq6Tcc8Ko",
				want:        false,
				description: "Token expires on 2029-06-04",
			},
			{
				name:        "valid token expired",
				token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4MzY4MDB9.eLbv9IVChXwHHX0zw3-5yG2aMEwJKmkGhpYJJFKLGEk",
				want:        true,
				description: "Token expired on 2020-01-01",
			},
			{
				name:        "invalid token format",
				token:       "invalid.token.format",
				want:        true,
				description: "Malformed token should be considered expired",
			},
			{
				name:        "empty token",
				token:       "",
				want:        true,
				description: "Empty token should be considered expired",
			},
			{
				name:        "malformed base64",
				token:       "header.@#$%^&*.signature",
				want:        true,
				description: "Invalid base64 encoding should be considered expired",
			},
		}

		manager := GetAuthTokenManager()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				manager.SetToken(tt.token)
				got := manager.CheckAuthTokenExpired()
				if got != tt.want {
					t.Errorf("CheckAuthTokenExpired() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
