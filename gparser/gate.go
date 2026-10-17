package gparser

import (
	"fmt"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/feature"
)

// Error is a parse error carrying a source position.
type Error struct {
	Line int
	Col  int
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Col, e.Msg)
}

// bailout is panicked from a semantic action to abort parsing with an error
// (used for version gating). It is recovered by the driver.
type bailout struct{ err error }

func lset(yylex yyLexer, st ast.Statement) {
	if l, ok := yylex.(*obLexer); ok {
		l.result = st
		if st != nil && l.last.Pos.IsValid() {
			if s, ok := st.(spannable); ok {
				s.SetSpan(l.stmtStart, l.last.EndPos)
			}
		}
		l.results = append(l.results, st)
	}
}

// gateFeat aborts parsing with a version error when feat is not enabled at the
// target version. It is a no-op for an unversioned target (all features on).
func gateFeat(yylex yyLexer, feat feature.FeatureID, msg string) {
	l, ok := yylex.(*obLexer)
	if !ok {
		return
	}
	if !l.features.Enabled(feat, l.version) {
		panic(bailout{&Error{
			Line: l.last.Pos.Line,
			Col:  l.last.Pos.Col,
			Msg:  msg,
		}})
	}
}
func gateWith(yylex yyLexer) {
	gateFeat(yylex, feature.FWithCte, "WITH ... SELECT is not supported at this version")
}
