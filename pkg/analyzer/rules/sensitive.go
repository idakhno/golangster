package rules

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// DefaultSensitiveKeywords is the default list of keywords that indicate sensitive data.
var DefaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"api-key",
	"auth",
	"authorization",
	"credential",
	"credentials",
	"private_key",
	"private-key",
	"access_key",
	"access-key",
	"bearer",
}

// CheckSensitive reports if a log call may expose sensitive data.
// It checks both string literals and variable names in concatenation expressions.
func CheckSensitive(pass *analysis.Pass, call *ast.CallExpr, keywords []string) {
	if len(call.Args) == 0 {
		return
	}
	checkExprForSensitive(pass, call.Args[0], keywords)
}

// checkExprForSensitive recursively walks an expression looking for sensitive data.
func checkExprForSensitive(pass *analysis.Pass, expr ast.Expr, keywords []string) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// check the string literal for sensitive keywords
		val := strings.Trim(e.Value, `"` + "`")
		if kw, found := containsSensitiveKeyword(val, keywords); found {
			pass.Report(analysis.Diagnostic{
				Pos:     e.Pos(),
				End:     e.End(),
				Message: "log message may expose sensitive data (keyword: \"" + kw + "\")",
			})
		}

	case *ast.BinaryExpr:
		// concatenation - check both sides
		checkExprForSensitive(pass, e.X, keywords)
		checkExprForSensitive(pass, e.Y, keywords)

	case *ast.Ident:
		// identifier - check the variable name itself
		if kw, found := containsSensitiveKeyword(e.Name, keywords); found {
			pass.Report(analysis.Diagnostic{
				Pos:     e.Pos(),
				End:     e.End(),
				Message: "log message may expose sensitive data via variable \"" + e.Name + "\" (keyword: \"" + kw + "\")",
			})
		}

	case *ast.ParenExpr:
		checkExprForSensitive(pass, e.X, keywords)
	}
}

// containsSensitiveKeyword reports whether s contains any of the keywords (case-insensitive).
func containsSensitiveKeyword(s string, keywords []string) (string, bool) {
	lower := strings.ToLower(s)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return kw, true
		}
	}
	return "", false
}
