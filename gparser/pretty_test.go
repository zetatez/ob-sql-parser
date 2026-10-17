package gparser

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/formatter"
	"github.com/zetatez/ob-sql-parser/lexer"
)

// TestPrettyFormat renders a selection of statements with pretty-printing and
// asserts the output still parses and is idempotent under the compact form.
func TestPrettyFormat(t *testing.T) {
	cases := []string{
		"SELECT a, b, c FROM t WHERE a > 1 GROUP BY a HAVING COUNT(*) > 1 ORDER BY a LIMIT 10",
		"SELECT * FROM a JOIN b ON a.id = b.id LEFT JOIN c ON c.aid = b.id WHERE a.x = 2",
		"INSERT INTO t (a, b) VALUES (1, 'x'), (2, 'y'), (3, 'z')",
		"UPDATE t SET a = 1, b = 2 WHERE id = 5",
		"CREATE TABLE t (id INT NOT NULL, name VARCHAR(20), PRIMARY KEY (id))",
		"SELECT a FROM t UNION ALL SELECT b FROM u ORDER BY 1",
		"DELETE FROM t WHERE x IN (1, 2, 3)",
	}
	opts := lexer.Options{AllowDoubleQuotedIdentifiers: true}
	for _, sql := range cases {
		toks, _ := lexer.LexAll(sql, opts)
		stmts, err := Parse(toks, Options{})
		if err != nil {
			t.Fatalf("parse %q: %v", sql, err)
		}
		pretty := formatter.FormatWith(stmts[0], formatter.Options{Pretty: true})
		toks2, _ := lexer.LexAll(pretty, opts)
		again, err := Parse(toks2, Options{})
		if err != nil {
			t.Fatalf("pretty re-parse %q failed: %v", pretty, err)
		}
		if compact := formatter.Format(again[0]); compact != formatter.Format(stmts[0]) {
			t.Fatalf("pretty output does not round-trip to same AST\n  sql:    %s\n  pretty: %q", sql, pretty)
		}
	}
}
