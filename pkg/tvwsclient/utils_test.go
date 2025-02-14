package tvwsclient

import (
	"regexp"
	"testing"
)

func TestGenerateSession(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{
			name:   "with empty prefix",
			prefix: "",
		},
		{
			name:   "with custom prefix",
			prefix: "test_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate multiple sessions to test randomness
			sessions := make(map[string]bool)
			for i := 0; i < 100; i++ {
				result := GenerateSession(tt.prefix)

				// Test prefix
				if !regexp.MustCompile("^" + tt.prefix).MatchString(result) {
					t.Errorf("generateSession() = %v, should start with prefix %v", result, tt.prefix)
				}

				// Test length (prefix + 12 random characters)
				expectedLen := len(tt.prefix) + 12
				if len(result) != expectedLen {
					t.Errorf("generateSession() length = %v, want %v", len(result), expectedLen)
				}

				// Test character set
				randomPart := result[len(tt.prefix):]
				if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(randomPart) {
					t.Errorf("generateSession() random part = %v, should only contain alphanumeric characters", randomPart)
				}

				// Test uniqueness
				if sessions[result] {
					t.Errorf("generateSession() generated duplicate session: %v", result)
				}
				sessions[result] = true
			}
		})
	}
}

func TestCheckAuthTokenExpired(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		want        bool
		description string
	}{
		{
			name: "valid token not expired",
			// This is a valid JWT token with exp set to 2024-12-31
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4Nzc4MzY4MDB9.QESXXx9ErawCjaMDrf51kakqBhpcBerH9Xtq6Tcc8Ko",
			want:        false,
			description: "Token expires on 2029-06-04",
		},
		{
			name: "valid token expired",
			// This is a valid JWT token with exp set to 2020-01-01
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckAuthTokenExpired(tt.token)
			if got != tt.want {
				t.Errorf("CheckAuthTokenExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
