package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/idakhno/golangster/pkg/analyzer/rules"
)

const name = "golangster"
const doc = `golangster checks log messages for style and security issues.

Rules:
  - log messages must start with a lowercase letter
  - log messages must be written in English only
  - log messages must not contain special characters or emoji
  - log messages must not expose sensitive data (passwords, tokens, etc.)

Supported loggers: log, log/slog, go.uber.org/zap.`

// Analyzer is the public instance used in plugins and tests.
var Analyzer = newAnalyzer(DefaultConfig())

// NewAnalyzer creates an analyzer with a custom configuration.
// Used by the golangci-lint plugin via the New function.
func NewAnalyzer(cfg Config) *analysis.Analyzer {
	return newAnalyzer(cfg)
}

func newAnalyzer(cfg Config) *analysis.Analyzer {
	r := &runner{cfg: cfg}
	a := &analysis.Analyzer{
		Name:             name,
		Doc:              doc,
		Run:              r.run,
		RunDespiteErrors: true,
		Requires:         []*analysis.Analyzer{inspect.Analyzer},
	}
	// flags for standalone mode (go vet -vettool)
	a.Flags.BoolVar(&r.cfg.Rules.Lowercase, "lowercase", cfg.Rules.Lowercase,
		"check that log messages start with a lowercase letter")
	a.Flags.BoolVar(&r.cfg.Rules.EnglishOnly, "english", cfg.Rules.EnglishOnly,
		"check that log messages are in English only")
	a.Flags.BoolVar(&r.cfg.Rules.NoSpecialChars, "special-chars", cfg.Rules.NoSpecialChars,
		"check that log messages contain no special characters or emoji")
	a.Flags.BoolVar(&r.cfg.Rules.NoSensitive, "sensitive", cfg.Rules.NoSensitive,
		"check that log messages do not expose sensitive data")
	return a
}

type runner struct {
	cfg Config
}

func (r *runner) run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		logCall, ok := FindLogCall(pass.TypesInfo, call)
		if !ok {
			return
		}

		// apply rules to each string literal found in the message
		for _, lit := range logCall.Literals {
			msg, ok := UnquoteStringLit(lit)
			if !ok {
				continue
			}

			if r.cfg.Rules.Lowercase {
				rules.CheckLowercase(pass, msg, lit)
			}
			if r.cfg.Rules.EnglishOnly {
				rules.CheckEnglish(pass, msg, lit)
			}
			if r.cfg.Rules.NoSpecialChars {
				rules.CheckSpecialChars(pass, msg, lit)
			}
		}

		// sensitive rule inspects the full call including variable names
		if r.cfg.Rules.NoSensitive {
			rules.CheckSensitive(pass, call, r.cfg.effectiveKeywords())
		}
	})

	return nil, nil
}
