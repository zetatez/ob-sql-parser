// Package formatter renders an ast back into SQL text. It performs a
// canonical reconstruction (not a source slice), producing output that can be
// re-parsed (round-trip), which doubles as a strong parser correctness check.

package formatter

import (
	"fmt"
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/token"
)

// Options control output styling.
type Options struct {
	// Pretty enables simple newline indentation for multi-item clauses.
	Pretty bool
}

// Format renders an arbitrary statement or expression to SQL text.
func Format(n ast.Node) string { return FormatWith(n, Options{}) }

// FormatWith renders n to SQL text with the given options.
func FormatWith(n ast.Node, opts Options) string {
	f := &formatter{opts: opts}
	switch v := n.(type) {
	case ast.Statement:
		f.stmt(v)
	case ast.Expression:
		f.expr(v)
	case nil:
		return ""
	default:
		f.nodeAny(n)
	}
	return f.b.String()
}

type formatter struct {
	b     strings.Builder
	opts  Options
	depth int // pretty-print indentation depth
}

func (f *formatter) printf(format string, args ...interface{}) {
	f.b.WriteString(sprintf(format, args...))
}
func (f *formatter) raw(s string) { f.b.WriteString(s) }

// nl ends the current line with the configured indentation (pretty mode).
func (f *formatter) nl() {
	if !f.opts.Pretty {
		f.raw(" ")
		return
	}
	f.raw("\n")
	for i := 0; i < f.depth; i++ {
		f.raw("  ")
	}
}
func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
func stringsLower(s string) string { return strings.ToLower(s) }
func stringsUpper(s string) string { return strings.ToUpper(s) }
func (f *formatter) ident(id *ast.Identifier) string {
	if id == nil {
		return ""
	}
	parts := make([]string, 0, len(id.Parts))
	for _, p := range id.Parts {
		parts = append(parts, f.quoteIfNeeded(p))
	}
	return strings.Join(parts, ".")
}
func (f *formatter) quoteIfNeeded(name string) string {
	if !needsQuote(name) {
		return name
	}
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
func (f *formatter) stringLiteral(s string) string {
	esc := strings.NewReplacer("\\", "\\\\", "'", "''", "\n", "\\n", "\r", "\\r", "\t", "\\t", "\x00", "\\0").Replace(s)
	return "'" + esc + "'"
}

// needsQuote reports whether an identifier must be backtick-quoted.
func needsQuote(name string) bool {
	if name == "" {
		return true
	}
	if token.IsReservedWord(name) {
		return true
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if i == 0 {
			if !(c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c >= 0x80) {
				return true
			}
			continue
		}
		if !(c == '_' || c == '$' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c >= 0x80) {
			return true
		}
	}
	return false
}
func isAlpha(s string) bool {
	if s == "" {
		return false
	}
	c := s[0]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
func (f *formatter) nodeAny(n ast.Node) {
	f.printf("<%s>", n.NodeKind())
}
