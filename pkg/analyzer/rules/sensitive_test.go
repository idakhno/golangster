package rules

import (
	"testing"
)

func TestContainsSensitiveKeyword(t *testing.T) {
	tests := []struct {
		input    string
		wantBad  bool
		wantKw   string
	}{
		{"user password: ", true, "password"},
		{"api_key=", true, "api_key"},
		{"token: ", true, "token"},
		{"user login", false, ""},
		{"request id", false, ""},
		{"secret key", true, "secret"},
		{"Authorization header", true, "auth"}, // "auth" appears before "authorization" in the list
		{"bearer ", true, "bearer"},
		{"private_key path", true, "private_key"},
		{"access_key id", true, "access_key"},
		{"username and email", false, ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			kw, found := containsSensitiveKeyword(tc.input, DefaultSensitiveKeywords)
			if found != tc.wantBad {
				t.Errorf("containsSensitiveKeyword(%q) found=%v, want %v", tc.input, found, tc.wantBad)
			}
			if found && kw != tc.wantKw {
				t.Errorf("containsSensitiveKeyword(%q) keyword=%q, want %q", tc.input, kw, tc.wantKw)
			}
		})
	}
}

func TestVariableNameSensitive(t *testing.T) {
	varNames := []struct {
		name    string
		wantBad bool
	}{
		{"password", true},
		{"apiKey", true},    // matches "apikey" keyword (lowercased)
		{"userToken", true}, // matches "token" keyword
		{"userName", false},
		{"requestID", false},
		{"passwd", true},
		{"pwd", true},
		{"secretValue", true},
	}

	for _, tc := range varNames {
		t.Run(tc.name, func(t *testing.T) {
			_, found := containsSensitiveKeyword(tc.name, DefaultSensitiveKeywords)
			if found != tc.wantBad {
				t.Errorf("variable %q: found=%v, want %v", tc.name, found, tc.wantBad)
			}
		})
	}
}
