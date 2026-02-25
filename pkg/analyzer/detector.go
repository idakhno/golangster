package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

// loggerPackages contains import paths of supported logger packages.
var loggerPackages = map[string]bool{
	"log":            true,
	"log/slog":       true,
	"go.uber.org/zap": true,
}

// loggerTypes maps package paths to known logger type names.
var loggerTypes = map[string]map[string]bool{
	"log/slog": {
		"Logger": true,
	},
	"go.uber.org/zap": {
		"Logger":        true,
		"SugaredLogger": true,
	},
}

// logMethods is the set of method names considered as log calls.
var logMethods = map[string]bool{
	"Info": true, "Infof": true, "Infow": true,
	"Error": true, "Errorf": true, "Errorw": true,
	"Debug": true, "Debugf": true, "Debugw": true,
	"Warn": true, "Warnf": true, "Warnw": true,
	"Fatal": true, "Fatalf": true, "Fatalw": true,
	"Panic": true, "Panicf": true, "Panicw": true,
	"Print": true, "Printf": true, "Println": true,
	"Log":   true,
}

// LogCall holds information about a detected log call.
type LogCall struct {
	// Expr is the AST expression of the first argument (the message).
	Expr ast.Expr
	// Literals are all string literals found inside the first argument.
	Literals []*ast.BasicLit
}

// FindLogCall reports whether a CallExpr is a call to a known logger.
// If so, it returns a LogCall with the message expression and its string literals.
func FindLogCall(typesInfo *types.Info, call *ast.CallExpr) (LogCall, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return LogCall{}, false
	}

	methodName := sel.Sel.Name
	if !logMethods[methodName] {
		return LogCall{}, false
	}

	if len(call.Args) == 0 {
		return LogCall{}, false
	}

	if !isLoggerReceiver(typesInfo, sel.X) {
		return LogCall{}, false
	}

	msgExpr := call.Args[0]
	lits := extractStringLiterals(msgExpr)

	return LogCall{Expr: msgExpr, Literals: lits}, true
}

// isLoggerReceiver reports whether the expression is a known logger receiver.
func isLoggerReceiver(typesInfo *types.Info, x ast.Expr) bool {
	switch v := x.(type) {
	case *ast.Ident:
		obj := typesInfo.ObjectOf(v)
		if obj == nil {
			return false
		}
		// package-level call: slog.Info(...), log.Print(...)
		if pkgName, ok := obj.(*types.PkgName); ok {
			return loggerPackages[pkgName.Imported().Path()]
		}
		// method call on a variable: logger.Info(...)
		return isLoggerType(typesInfo.TypeOf(v))

	default:
		// chained call or other expression - check the type of the left side
		return isLoggerType(typesInfo.TypeOf(x))
	}
}

// isLoggerType reports whether a type is a known logger type.
func isLoggerType(t types.Type) bool {
	if t == nil {
		return false
	}
	// dereference pointer if needed
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj.Pkg() == nil {
		return false
	}
	pkgPath := obj.Pkg().Path()
	typeName := obj.Name()

	types, exists := loggerTypes[pkgPath]
	if !exists {
		return false
	}
	return types[typeName]
}

// extractStringLiterals recursively collects all string literals from an expression.
// Supports plain literals, concatenation with "+", and parenthesized expressions.
func extractStringLiterals(expr ast.Expr) []*ast.BasicLit {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return []*ast.BasicLit{e}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			var result []*ast.BasicLit
			result = append(result, extractStringLiterals(e.X)...)
			result = append(result, extractStringLiterals(e.Y)...)
			return result
		}
	case *ast.ParenExpr:
		return extractStringLiterals(e.X)
	}
	return nil
}

// UnquoteStringLit returns the unquoted string value of a basic literal.
func UnquoteStringLit(lit *ast.BasicLit) (string, bool) {
	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return val, true
}
