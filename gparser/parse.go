package gparser

import (
	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/token"
)

// spannable is implemented (via embedded ast.Base) by every concrete node.
type spannable interface {
	SetSpan(start, end token.Position)
}

// Parse parses a pre-lexed token stream into zero or more statements. The
// grammar consumes ';' itself (either as a statement separator at the top
// level or as the statement terminator inside PL blocks), so a single LALR
// run produces every statement.
func Parse(toks []token.Token, opts Options) ([]ast.Statement, error) {
	l := newLexer(toks, opts)
	if _, err := parseOne(l); err != nil {
		return nil, err
	}
	stmts := l.results
	for _, st := range stmts {
		if st != nil {
			ast.Finalize(st)
		}
	}
	return stmts, nil
}

func parseOne(l *obLexer) (st ast.Statement, err error) {
	l.result = nil
	l.err = nil
	defer func() {
		if r := recover(); r != nil {
			if b, ok := r.(bailout); ok {
				err = b.err
				st = nil
				return
			}
			panic(r)
		}
	}()
	yyNewParser().Parse(l)
	if l.err != nil {
		return nil, l.err
	}
	return l.result, nil
}
