package rules

import (
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// emojiRanges is the unicode range table for emoji and decorative symbols.
var emojiRanges = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x2600, 0x27BF, 1}, // Miscellaneous Symbols, Dingbats
		{0xFE00, 0xFE0F, 1}, // Variation Selectors
	},
	R32: []unicode.Range32{
		{0x1F300, 0x1F9FF, 1}, // Emoji (Miscellaneous Symbols and Pictographs, Transport, etc.)
		{0x1FA00, 0x1FAFF, 1}, // Chess pieces, additional emoji
		{0x1F000, 0x1F02F, 1}, // Mahjong, domino
	},
}

// forbiddenChars is the set of special characters not allowed in log messages.
var forbiddenChars = map[rune]bool{
	'!': true,
	'?': true,
}

// CheckSpecialChars reports if a log message contains emoji or special characters.
func CheckSpecialChars(pass *analysis.Pass, msg string, lit *ast.BasicLit) {
	checkEmoji(pass, msg, lit)
	checkForbiddenChars(pass, msg, lit)
	checkRepeatedDots(pass, msg, lit)
}

func checkEmoji(pass *analysis.Pass, msg string, lit *ast.BasicLit) {
	for _, r := range msg {
		if unicode.Is(emojiRanges, r) {
			pass.Report(analysis.Diagnostic{
				Pos:     lit.Pos(),
				End:     lit.End(),
				Message: "log message must not contain emoji",
			})
			return
		}
	}
}

func checkForbiddenChars(pass *analysis.Pass, msg string, lit *ast.BasicLit) {
	for _, r := range msg {
		if forbiddenChars[r] {
			pass.Report(analysis.Diagnostic{
				Pos:     lit.Pos(),
				End:     lit.End(),
				Message: "log message must not contain special character '" + string(r) + "'",
			})
			return
		}
	}
}

func checkRepeatedDots(pass *analysis.Pass, msg string, lit *ast.BasicLit) {
	if strings.Contains(msg, "...") {
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message must not contain '...' (ellipsis)",
		})
	}
}
