// Package obparser is the top-level entry point for the OceanBase 4.x
// (MySQL mode) SQL parser. It lexes, parses, and produces a typed AST that can
// be formatted back to SQL, traversed via the visitor package, or serialized
// to JSON.
//
// Version gating: pass WithVersion for an exact OceanBase version; otherwise
// all features are enabled (assumes the latest OB 4.x).
package obparser

import (
	"fmt"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/feature"
	"github.com/zetatez/ob-sql-parser/formatter"
	"github.com/zetatez/ob-sql-parser/gparser"
	"github.com/zetatez/ob-sql-parser/lexer"
	"github.com/zetatez/ob-sql-parser/version"
)

// Option configures the parser.
type Option func(*gparser.Options)

// WithVersion sets the target OceanBase version for feature gating.
func WithVersion(v version.Version) Option {
	return func(o *gparser.Options) { o.Version = v }
}

// WithFeatures replaces the feature registry used for version gating.
func WithFeatures(f *feature.Registry) Option {
	return func(o *gparser.Options) { o.Features = f }
}

func buildOptions(opts []Option) gparser.Options {
	var o gparser.Options
	for _, fn := range opts {
		fn(&o)
	}
	return o
}

var defaultLexOpts = lexer.Options{AllowDoubleQuotedIdentifiers: true}

// Parse lexes and parses a single SQL statement.
//
// The optional trailing semicolon is accepted. An empty or whitespace-only
// input yields (nil, nil).
func Parse(sql string, opts ...Option) (ast.Statement, error) {
	toks, err := lexer.LexAll(sql, defaultLexOpts)
	if err != nil {
		return nil, err
	}
	stmts, err := gparser.Parse(toks, buildOptions(opts))
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 {
		return nil, nil
	}
	return stmts[0], nil
}

// ParseMulti parses a string that may contain multiple ';'-separated
// statements, returning all of them.
func ParseMulti(sql string, opts ...Option) ([]ast.Statement, error) {
	toks, err := lexer.LexAll(sql, defaultLexOpts)
	if err != nil {
		return nil, err
	}
	return gparser.Parse(toks, buildOptions(opts))
}

// MustParse parses a single statement, panicking on error (for tests/constants).
func MustParse(sql string, opts ...Option) ast.Statement {
	st, err := Parse(sql, opts...)
	if err != nil {
		panic(fmt.Errorf("MustParse: %w", err))
	}
	return st
}

// Format renders an AST node back to SQL text (see the formatter package).
func Format(n ast.Node) string { return formatter.Format(n) }
