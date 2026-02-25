// Build this file as a golangci-lint plugin:
//   go build -buildmode=plugin -o golangster.so ./plugin/
//
// Then reference it in .golangci.yml:
//   linters:
//     settings:
//       custom:
//         golangster:
//           path: ./golangster.so
//           description: Checks log messages for style and security issues
//           original-url: github.com/idakhno/golangster

//go:build ignore

package main

import (
	"golang.org/x/tools/go/analysis"

	"github.com/idakhno/golangster/pkg/analyzer"
)

// New is the entry point for the golangci-lint plugin system.
// It accepts settings from .golangci.yml and returns the list of analyzers.
func New(conf any) ([]*analysis.Analyzer, error) {
	cfg := analyzer.DefaultConfig()

	// apply settings from .golangci.yml if provided
	if settings, ok := conf.(map[string]any); ok {
		if rules, ok := settings["rules"].(map[string]any); ok {
			if v, ok := rules["lowercase"].(bool); ok {
				cfg.Rules.Lowercase = v
			}
			if v, ok := rules["english_only"].(bool); ok {
				cfg.Rules.EnglishOnly = v
			}
			if v, ok := rules["no_special_chars"].(bool); ok {
				cfg.Rules.NoSpecialChars = v
			}
			if v, ok := rules["no_sensitive"].(bool); ok {
				cfg.Rules.NoSensitive = v
			}
		}
		if kws, ok := settings["sensitive_keywords"].([]any); ok {
			for _, kw := range kws {
				if s, ok := kw.(string); ok {
					cfg.SensitiveKeywords = append(cfg.SensitiveKeywords, s)
				}
			}
		}
	}

	return []*analysis.Analyzer{analyzer.NewAnalyzer(cfg)}, nil
}
