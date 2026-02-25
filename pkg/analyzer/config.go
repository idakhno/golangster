package analyzer

import (
	"github.com/idakhno/golangster/pkg/analyzer/rules"
)

// Config holds golangster analyzer settings.
type Config struct {
	// Rules controls which checks are enabled.
	Rules RulesConfig
	// SensitiveKeywords is the list of keywords used to detect sensitive data.
	// If empty, DefaultSensitiveKeywords is used.
	SensitiveKeywords []string
}

// RulesConfig controls enabling and disabling of individual rules.
type RulesConfig struct {
	Lowercase    bool
	EnglishOnly  bool
	NoSpecialChars bool
	NoSensitive  bool
}

// DefaultConfig returns a Config with all rules enabled.
func DefaultConfig() Config {
	return Config{
		Rules: RulesConfig{
			Lowercase:      true,
			EnglishOnly:    true,
			NoSpecialChars: true,
			NoSensitive:    true,
		},
		SensitiveKeywords: rules.DefaultSensitiveKeywords,
	}
}

// effectiveKeywords returns the keyword list to use for sensitive checks.
func (c *Config) effectiveKeywords() []string {
	if len(c.SensitiveKeywords) > 0 {
		return c.SensitiveKeywords
	}
	return rules.DefaultSensitiveKeywords
}
