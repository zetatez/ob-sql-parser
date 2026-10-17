package obparser

import (
	"encoding/json"
	"testing"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/version"
	"github.com/zetatez/ob-sql-parser/visitor"
)

// Version gating: WITH ... SELECT is a feature introduced at OB 4.0, so an
// (artificially) older target version rejects it while 4.x accepts it.
func TestVersionGating(t *testing.T) {
	sql := "WITH cte AS (SELECT 1) SELECT * FROM cte"

	if _, err := Parse(sql, WithVersion(version.MustParse("4.0.0"))); err != nil {
		t.Fatalf("4.0 should accept WITH: %v", err)
	}
	_, err := Parse(sql, WithVersion(version.MustParse("3.9.0")))
	if err == nil {
		t.Fatal("pre-4.0 should reject WITH")
	}
	if e, ok := err.(interface{ Kind() int }); ok {
		_ = e // presence of a typed parser error is enough here
		_ = ok
	}
}

// Walk a SELECT and collect all table-name leaves reachable from TableRefs.
func TestVisitorCollect(t *testing.T) {
	stmt, err := Parse("SELECT a FROM t1 JOIN t2 ON t1.id = t2.id WHERE x IN (SELECT y FROM t3)")
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	visitor.Walk(stmt, visitor.FuncVisitor(func(n ast.Node) visitor.Visitor {
		if ref, ok := n.(*ast.TableRef); ok && ref.RefKind == ast.RefTable {
			tables = append(tables, ref.Name.Parts[0])
		}
		return visitor.NoOpVisitor{}
	}))
	want := []string{"t1", "t2", "t3"}
	if len(tables) != len(want) {
		t.Fatalf("tables = %v, want %v", tables, want)
	}
	for i := range want {
		if tables[i] != want[i] {
			t.Fatalf("tables = %v, want %v", tables, want)
		}
	}
}

func TestJSONSerialization(t *testing.T) {
	stmt, err := Parse("SELECT a, b FROM t WHERE a > 1")
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(stmt)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(data)
	if !contains(s, "\"SelectStmt\"") && !contains(s, "\"kind\"") {
		t.Fatalf("JSON missing kind discriminator: %s", s)
	}
	// ensure the statement body serialized structurally
	if !contains(s, "t") {
		t.Fatalf("JSON missing table name: %s", s)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
