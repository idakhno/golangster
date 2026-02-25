package rules

import (
	"testing"
	"unicode"
)

func TestIsEmoji(t *testing.T) {
	tests := []struct {
		char    rune
		isEmoji bool
	}{
		{'🚀', true},
		{'✅', true},
		{'⚠', true},
		{'❌', true},
		{'a', false},
		{'Z', false},
		{'!', false},
		{'.', false},
		{' ', false},
	}
	for _, tc := range tests {
		got := unicode.Is(emojiRanges, tc.char)
		if got != tc.isEmoji {
			t.Errorf("isEmoji(%q) = %v, want %v", tc.char, got, tc.isEmoji)
		}
	}
}

func TestForbiddenChars(t *testing.T) {
	tests := []struct {
		msg     string
		wantBad bool
	}{
		{"server started", false},
		{"server started!", true},
		{"what happened?", true},
		{"connection failed!!!", true},
		{"request processed", false},
	}
	for _, tc := range tests {
		t.Run(tc.msg, func(t *testing.T) {
			got := hasForbiddenChar(tc.msg)
			if got != tc.wantBad {
				t.Errorf("hasForbiddenChar(%q) = %v, want %v", tc.msg, got, tc.wantBad)
			}
		})
	}
}

func TestRepeatedDots(t *testing.T) {
	tests := []struct {
		msg     string
		wantBad bool
	}{
		{"something went wrong...", true},
		{"loading...", true},
		{"v1.2.3 deployed", false},
		{"error at line 1.2", false},
		{"connecting to db", false},
	}
	for _, tc := range tests {
		t.Run(tc.msg, func(t *testing.T) {
			got := hasEllipsis(tc.msg)
			if got != tc.wantBad {
				t.Errorf("hasEllipsis(%q) = %v, want %v", tc.msg, got, tc.wantBad)
			}
		})
	}
}

func hasForbiddenChar(msg string) bool {
	for _, r := range msg {
		if forbiddenChars[r] {
			return true
		}
	}
	return false
}

func hasEllipsis(msg string) bool {
	for i := 0; i+2 < len(msg); i++ {
		if msg[i] == '.' && msg[i+1] == '.' && msg[i+2] == '.' {
			return true
		}
	}
	return false
}
