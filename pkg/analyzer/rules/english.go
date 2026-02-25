package rules

import (
	"go/ast"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// nonLatinScripts is a list of unicode range tables for non-Latin scripts.
var nonLatinScripts = []*unicode.RangeTable{
	unicode.Cyrillic,
	unicode.Han,
	unicode.Arabic,
	unicode.Hebrew,
	unicode.Devanagari,
	unicode.Hiragana,
	unicode.Katakana,
	unicode.Hangul,
	unicode.Thai,
	unicode.Georgian,
	unicode.Armenian,
	unicode.Greek,
	unicode.Ethiopic,
	unicode.Myanmar,
	unicode.Khmer,
}

// CheckEnglish reports if a log message contains non-Latin script characters.
func CheckEnglish(pass *analysis.Pass, msg string, lit *ast.BasicLit) {
	for _, r := range msg {
		for _, script := range nonLatinScripts {
			if unicode.Is(script, r) {
				pass.Report(analysis.Diagnostic{
					Pos:     lit.Pos(),
					End:     lit.End(),
					Message: "log message must be in English only",
				})
				return
			}
		}
	}
}
