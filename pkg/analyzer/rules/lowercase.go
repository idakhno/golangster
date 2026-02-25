package rules

import (
	"go/ast"
	"go/token"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

// CheckLowercase reports if a log message starts with an uppercase letter.
// A SuggestedFix is included to convert the first letter to lowercase.
func CheckLowercase(pass *analysis.Pass, msg string, lit *ast.BasicLit) {
	if len(msg) == 0 {
		return
	}
	r, _ := utf8.DecodeRuneInString(msg)
	if r == utf8.RuneError || !unicode.IsUpper(r) {
		return
	}

	// Compute the position of the first character inside the string literal.
	// lit.Pos() points to the opening quote, +1 is the first character.
	firstCharPos := lit.Pos() + 1
	// Width of the first character in bytes within the source file.
	firstCharEnd := token.Pos(int(firstCharPos) + utf8.RuneLen(r))

	lower := string(unicode.ToLower(r))

	pass.Report(analysis.Diagnostic{
		Pos:     lit.Pos(),
		End:     lit.End(),
		Message: "log message must start with a lowercase letter",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "convert first letter to lowercase",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     firstCharPos,
						End:     firstCharEnd,
						NewText: []byte(lower),
					},
				},
			},
		},
	})
}
