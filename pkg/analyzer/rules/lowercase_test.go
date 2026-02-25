package rules

import (
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestCheckLowercaseLogic(t *testing.T) {
	tests := []struct {
		msg     string
		wantBad bool
	}{
		{"starting server", false},
		{"Starting server", true},
		{"ERROR: something", true},
		{"error: something", false},
		{"", false},
		{"123 digits first", false},
		{"Ошибка подключения", true},
	}

	for _, tc := range tests {
		t.Run(tc.msg, func(t *testing.T) {
			got := isUpperStart(tc.msg)
			if got != tc.wantBad {
				t.Errorf("isUpperStart(%q) = %v, want %v", tc.msg, got, tc.wantBad)
			}
		})
	}
}

// isUpperStart is a helper for tests that checks uppercase start without a pass.
func isUpperStart(msg string) bool {
	if len(msg) == 0 {
		return false
	}
	r, _ := utf8.DecodeRuneInString(msg)
	return r != utf8.RuneError && unicode.IsUpper(r)
}
