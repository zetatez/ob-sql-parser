package visitor_test

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/visitor"
)

func TestWalkCollectsAllNodes(t *testing.T) {
	// SELECT a + b FROM t WHERE a > 1  (hand-built tree)
	ident := func(s string) ast.Expression { return &ast.Identifier{Parts: []string{s}} }
	sel := &ast.SelectStmt{Bodies: []*ast.SelectBody{{
		Items: []ast.ColumnItem{
			{Expr: &ast.BinaryExpr{Op: "+", Left: ident("a"), Right: ident("b")}},
		},
		From:  []*ast.TableRef{{RefKind: ast.RefTable, Name: &ast.Identifier{Parts: []string{"t"}}}},
		Where: &ast.BinaryExpr{Op: ">", Left: ident("a"), Right: &ast.Literal{Value: "1"}},
	}}}
	ast.Finalize(sel)
	var kinds []string
	visitor.Walk(sel, visitor.FuncVisitor(func(n ast.Node) visitor.Visitor {
		kinds = append(kinds, string(n.NodeKind()))
		return visitor.NoOpVisitor{}
	}))
	for _, want := range []string{"SelectStmt", "SelectBody", "BinaryExpr", "Identifier", "Literal", "TableRef"} {
		if !contains(kinds, want) {
			t.Errorf("kinds %v missing %s", kinds, want)
		}
	}
}

func TestWalkPruneSubtree(t *testing.T) {
	total := 0
	visitor.Walk(&ast.SelectStmt{}, visitor.FuncVisitor(func(n ast.Node) visitor.Visitor {
		total++
		return nil // prune everything below the root
	}))
	if total != 1 {
		t.Fatalf("expected 1 visit after pruning, got %d", total)
	}
}

func TestWalkHandlesNil(t *testing.T) {
	// Must not panic on typed-nil children (absent optional fields).
	visitor.Walk(&ast.StoredProcStmt{}, visitor.NoOpVisitor{})
	// nil root is a no-op.
	visitor.Walk(nil, visitor.NoOpVisitor{})
	// nil visitor is a no-op.
	visitor.Walk(&ast.SelectStmt{}, nil)
}

func TestWalkStatements(t *testing.T) {
	var count int
	stmts := []ast.Statement{&ast.SelectStmt{}, &ast.CommitStmt{}}
	visitor.WalkStatements(stmts, visitor.FuncVisitor(func(n ast.Node) visitor.Visitor {
		count++
		return visitor.NoOpVisitor{}
	}))
	if count != 2 {
		t.Fatalf("expected 2 visits, got %d", count)
	}
}

func TestFuncVisitorAndNoOp(t *testing.T) {
	v := visitor.FuncVisitor(func(n ast.Node) visitor.Visitor { return visitor.NoOpVisitor{} })
	var ran bool
	visitor.Walk(&ast.SelectStmt{}, v)
	ran = v != nil
	if !ran {
		t.Fatal("unreachable")
	}
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
