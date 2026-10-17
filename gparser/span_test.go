package gparser

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/ast"
)

// TestExpressionSpans verifies that expression-level source spans are stamped
// by the generated parser's injected position stack (goyacc has no %locations;
// see genpos.go).
func TestExpressionSpans(t *testing.T) {
	st := parseAll(t, "SELECT a + b FROM t WHERE a > 1 AND c = 'x' ORDER BY d")
	sel := st[0].(*ast.SelectStmt)
	body := sel.Bodies[0]

	where := body.Where.(*ast.BinaryExpr)
	if where.Span().Start.Line != 1 || where.Span().Start.Col != 27 || where.Span().End.Col != 44 {
		t.Fatalf("WHERE span = %+v, want [1:27 - 1:44]", where.Span())
	}
	left := where.Left.(*ast.BinaryExpr)
	if left.Span().Start.Col != 27 || left.Span().End.Col != 32 {
		t.Fatalf("a>1 span = %+v, want [1:27 - 1:32]", left.Span())
	}
	// projection expr a + b
	proj := body.Items[0].Expr.(*ast.BinaryExpr)
	if proj.Span().Start.Col != 8 || proj.Span().End.Col != 13 {
		t.Fatalf("a+b span = %+v, want [1:8 - 1:13]", proj.Span())
	}
	// order item
	order := sel.OrderBy[0].Expr.(*ast.Identifier)
	if !order.Span().Start.IsValid() || order.Span().Start.Col != 54 {
		t.Fatalf("order id span = %+v, want start 1:54", order.Span())
	}
}
