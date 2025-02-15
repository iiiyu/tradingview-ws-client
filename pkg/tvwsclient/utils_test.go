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
