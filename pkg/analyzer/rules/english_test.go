package rules

import (
	"testing"
	"unicode"
)

func TestIsNonLatin(t *testing.T) {
	tests := []struct {
		msg     string
		wantBad bool
	}{
		{"starting server on port 8080", false},
		{"failed to connect to database", false},
		{"error reading file", false},
		{"Запуск сервера", true}, // Cyrillic
		{"服务器启动失败", true},        // Han (Chinese)
		{"خطأ في الاتصال", true}, // Arabic
		{"שגיאת חיבור", true},    // Hebrew
		{"エラーが発生しました", true},     // Katakana/Hiragana
		{"user created: john@example.com", false},
		{"request timeout after 30s", false},
	}

	for _, tc := range tests {
		t.Run(tc.msg, func(t *testing.T) {
			got := hasNonLatinScript(tc.msg)
			if got != tc.wantBad {
				t.Errorf("hasNonLatinScript(%q) = %v, want %v", tc.msg, got, tc.wantBad)
			}
		})
	}
}

func hasNonLatinScript(msg string) bool {
	for _, r := range msg {
		for _, script := range nonLatinScripts {
			if unicode.Is(script, r) {
				return true
			}
		}
	}
	return false
}
