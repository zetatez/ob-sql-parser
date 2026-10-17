package formatter_test

import (
	"strings"
	"testing"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/formatter"
)

func lit(text string) *ast.Literal { return &ast.Literal{Value: text, Text: text} }

func strlit(s string) *ast.Literal {
	return &ast.Literal{LitType: ast.LitString, Value: s, Text: "'" + s + "'"}
}

func sel() *ast.SelectStmt {
	return &ast.SelectStmt{
		Bodies: []*ast.SelectBody{{
			Items: []ast.ColumnItem{
				{Expr: &ast.Identifier{Parts: []string{"a"}}},
				{Expr: &ast.Identifier{Parts: []string{"b"}}},
			},
			From: []*ast.TableRef{{RefKind: ast.RefTable, Name: &ast.Identifier{Parts: []string{"t"}}}},
		}},
	}
}

func simpleSelect() *ast.SelectStmt {
	return &ast.SelectStmt{Bodies: []*ast.SelectBody{{
		Items: []ast.ColumnItem{{Expr: &ast.Identifier{Parts: []string{"a"}}}},
	}}}
}

func TestExpressionFormats(t *testing.T) {
	cases := []struct {
		name string
		expr ast.Expression
		want string
	}{
		{"number", lit("42"), "42"},
		{"string", &ast.Literal{LitType: ast.LitString, Value: "a'b"}, "'a''b'"},
		{"charset-string", &ast.Literal{LitType: ast.LitString, Charset: "utf8mb4", Value: "hi", Text: "_utf8mb4'hi'"}, "_utf8mb4'hi'"},
		{"null", &ast.Literal{LitType: ast.LitNull, Value: "NULL"}, "NULL"},
		{"true", &ast.Literal{LitType: ast.LitBool, Value: "true"}, "TRUE"},
		{"false", &ast.Literal{LitType: ast.LitBool, Value: "false"}, "FALSE"},
		{"hex", &ast.Literal{LitType: ast.LitHex, Text: "0xAB"}, "0xAB"},
		{"ident", &ast.Identifier{Parts: []string{"a", "b"}}, "a.b"},
		{"ident-reserved", &ast.Identifier{Parts: []string{"select"}}, "`select`"},
		{"ident-backtick", &ast.Identifier{Parts: []string{"a`b"}}, "`a``b`"},
		{"variable", &ast.Variable{Name: "x"}, "@x"},
		{"variable-system", &ast.Variable{Name: "sql_mode", IsSystem: true}, "@@sql_mode"},
		{"variable-scope", &ast.Variable{Name: "tx", IsSystem: true, Scope: "SESSION"}, "@@SESSION.tx"},
		{"unary-minus", &ast.UnaryExpr{Op: "-", Operand: lit("1")}, "-1"},
		{"unary-not", &ast.UnaryExpr{Op: "NOT", Operand: lit("1")}, "NOT 1"},
		{"binary", &ast.BinaryExpr{Op: "+", Left: lit("1"), Right: lit("2")}, "1 + 2"},
		{"func", &ast.FuncCall{Name: "COUNT", Star: true}, "count(*)"},
		{"func-distinct", &ast.FuncCall{Name: "count", Distinct: true, Args: []ast.Expression{&ast.Identifier{Parts: []string{"a"}}}}, "count(DISTINCT a)"},
		{"paren", &ast.ParenExpr{Expr: lit("1")}, "(1)"},
		{
			"case",
			&ast.CaseExpr{
				Operand: lit("x"),
				Whens:   []ast.CaseClause{{When: lit("1"), Then: lit("2")}},
				Else:    lit("3"),
			},
			"CASE x WHEN 1 THEN 2 ELSE 3 END",
		},
		{"between", &ast.BetweenExpr{Expr: lit("1"), Low: lit("2"), High: lit("3"), Not: true}, "1 NOT BETWEEN 2 AND 3"},
		{"in", &ast.InExpr{Expr: lit("1"), List: []ast.Expression{lit("2"), lit("3")}}, "1 IN (2, 3)"},
		{"like", &ast.LikeExpr{Expr: lit("a"), Pattern: strlit("a%"), Escape: strlit("!"), Not: true}, "a NOT LIKE 'a%' ESCAPE '!'"},
		{"is", &ast.IsExpr{Expr: lit("a"), Truth: "null", Not: true}, "a IS NOT NULL"},
		{"interval", &ast.IntervalExpr{Value: lit("3"), Unit: "DAY"}, "INTERVAL 3 DAY"},
		{"param", &ast.ParamMarker{}, "?"},
		{"cast", &ast.CastExpr{Expr: lit("1"), TypeName: "SIGNED"}, "CAST(1 AS SIGNED)"},
		{"subquery", &ast.SubqueryExpr{Query: simpleSelect()}, "(SELECT a)"},
		{"member-of", &ast.MemberOfExpr{Expr: lit("1"), List: []ast.Expression{lit("2")}}, "1 MEMBER OF (2)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatter.Format(tc.expr); got != tc.want {
				t.Fatalf("Format(%q) = %q, want %q", tc.name, got, tc.want)
			}
		})
	}
}

func TestRepresentativeStatements(t *testing.T) {
	ident := func(s string) *ast.Identifier { return &ast.Identifier{Parts: []string{s}} }
	tbl := func() *ast.TableRef {
		return &ast.TableRef{RefKind: ast.RefTable, Name: &ast.Identifier{Parts: []string{"t"}}}
	}
	cases := []struct {
		name string
		stmt ast.Statement
		want string
	}{
		{"select", sel(), "SELECT a, b FROM t"},
		{
			"select-full",
			&ast.SelectStmt{Bodies: []*ast.SelectBody{{
				Items: []ast.ColumnItem{
					{Expr: &ast.Identifier{Parts: []string{"a"}}},
					{Expr: &ast.Identifier{Parts: []string{"b"}}},
				},
				From:      []*ast.TableRef{tbl()},
				Where:     &ast.BinaryExpr{Op: ">", Left: ident("a"), Right: lit("1")},
				GroupBy:   []ast.Expression{ident("a")},
				Rollup:    true,
				Having:    &ast.BinaryExpr{Op: ">", Left: &ast.FuncCall{Name: "count", Star: true}, Right: lit("0")},
				OrderBy:   []ast.OrderItem{{Expr: ident("a"), Desc: true, HasNulls: true, NullsFirst: true}},
				Limit:     &ast.LimitClause{Count: lit("10"), Offset: lit("5"), HasOffset: true},
				ForUpdate: true,
				LockWait:  "NOWAIT",
			}}},
			"SELECT a, b FROM t WHERE a > 1 GROUP BY a WITH ROLLUP HAVING count(*) > 0 ORDER BY a DESC NULLS FIRST LIMIT 5, 10 FOR UPDATE NOWAIT",
		},
		{
			"insert-values",
			&ast.InsertStmt{
				Table:   ident("t"),
				Columns: []*ast.Identifier{ident("a"), ident("b")},
				Rows:    [][]ast.Expression{{lit("1"), lit("2")}, {lit("3"), lit("4")}},
			},
			"INSERT INTO t (a, b) VALUES (1, 2), (3, 4)",
		},
		{
			"replace-set",
			&ast.InsertStmt{Replace: true, Table: ident("t"), Set: []ast.Assignment{{Column: ident("a"), Expr: lit("1")}}},
			"REPLACE INTO t SET a = 1",
		},
		{
			"update",
			&ast.UpdateStmt{
				Tables: []*ast.TableRef{tbl()},
				Set:    []ast.Assignment{{Column: ident("a"), Expr: lit("1")}},
				Where:  &ast.BinaryExpr{Op: "=", Left: ident("a"), Right: lit("2")},
			},
			"UPDATE t SET a = 1 WHERE a = 2",
		},
		{"delete", &ast.DeleteStmt{From: []*ast.TableRef{tbl()}, Where: lit("1")}, "DELETE FROM t WHERE 1"},
		{
			"show-table-status",
			&ast.ShowStmt{Kind: "table status", From: ident("db"), Like: &ast.Literal{LitType: ast.LitString, Value: "t%", Text: "'t%'"}},
			"SHOW TABLE STATUS FROM db LIKE 't%'",
		},
		{"show-engine", &ast.ShowStmt{Kind: "engine", Target: ident("InnoDB"), Mode: "STATUS"}, "SHOW ENGINE InnoDB STATUS"},
		{"begin", &ast.BeginStmt{}, "BEGIN"},
		{"commit", &ast.CommitStmt{}, "COMMIT"},
		{"rollback", &ast.RollbackStmt{}, "ROLLBACK"},
		{"use", &ast.UseDbStmt{Name: "db1"}, "USE db1"},
		{"help", &ast.HelpStmt{Target: ident("select")}, "HELP `select`"},
		{"call", &ast.CallStmt{Name: "p", Args: []ast.Expression{lit("1")}}, "CALL p(1)"},
		{"do", &ast.DoStmt{Expr: lit("1")}, "DO 1"},
		{"kill", &ast.KillStmt{Kind: "QUERY", Target: lit("1")}, "KILL QUERY 1"},
		{"truncate", &ast.TruncateStmt{Table: ident("t")}, "TRUNCATE TABLE t"},
		{"lock-tables", &ast.LockTablesStmt{Tables: []ast.LockTable{{Table: ident("t"), LockType: "WRITE"}}}, "LOCK TABLES t WRITE"},
		{"unlock", &ast.UnlockTablesStmt{}, "UNLOCK TABLES"},
		{"create-db", &ast.CreateDatabaseStmt{Name: "mydb", Charset: "utf8mb4"}, "CREATE DATABASE mydb DEFAULT CHARACTER SET utf8mb4"},
		{"drop-db", &ast.DropDatabaseStmt{Name: "mydb"}, "DROP DATABASE mydb"},
		{"create-table-like", &ast.CreateTableLikeStmt{Table: ident("t2"), Like: ident("t1")}, "CREATE TABLE t2 LIKE t1"},
		{"create-view", &ast.CreateViewStmt{Name: ident("v"), Select: sel()}, "CREATE VIEW v AS SELECT a, b FROM t"},
		{"drop-table", &ast.DropTableStmt{Tables: []*ast.Identifier{ident("t")}}, "DROP TABLE t"},
		{"drop-index", &ast.DropIndexStmt{Name: ident("i"), Table: ident("t")}, "DROP INDEX i ON t"},
		{"savepoint", &ast.SavepointStmt{Name: "sp"}, "SAVEPOINT sp"},
		{"rollback-to", &ast.RollbackToSavepointStmt{Name: "sp"}, "ROLLBACK TO SAVEPOINT sp"},
		{"release", &ast.ReleaseSavepointStmt{Name: "sp"}, "RELEASE SAVEPOINT sp"},
		{"prepare", &ast.PrepareStmt{Name: "st", Source: &ast.Literal{LitType: ast.LitString, Value: "SELECT 1", Text: "'SELECT 1'"}}, "PREPARE st FROM 'SELECT 1'"},
		{"execute", &ast.ExecuteStmt{Name: "st"}, "EXECUTE st"},
		{"deallocate", &ast.DeallocateStmt{Name: "st"}, "DEALLOCATE PREPARE st"},
		{"grant", &ast.GrantStmt{Privileges: []string{"select"}, Level: "t", Users: []ast.UserSpec{{Name: "u"}}}, "GRANT select ON t TO 'u'"},
		{"revoke", &ast.RevokeStmt{Privileges: []string{"select"}, Level: "t", Users: []ast.UserSpec{{Name: "u"}}}, "REVOKE select ON t FROM 'u'"},
		{"set", &ast.SetStmt{Items: []*ast.SetItem{{Name: "a", Value: lit("1")}}}, "SET a = 1"},
		{"set-password", &ast.SetPasswordStmt{User: ast.UserSpec{Name: "u"}, HasUser: true, Password: "pw"}, "SET PASSWORD FOR 'u' = 'pw'"},
		{"xa", &ast.XAStmt{Action: "START", XID: []ast.Expression{lit("a,b,c")}}, "XA START a,b,c"},
		{"flush", &ast.FlushStmt{Kind: "PRIVILEGES"}, "FLUSH PRIVILEGES"},
		{"purge", &ast.PurgeStmt{Kind: "TABLE", Table: ident("t")}, "PURGE TABLE t"},
		{"explain", &ast.ExplainStmt{Stmt: sel()}, "EXPLAIN SELECT a, b FROM t"},
		{
			"pl-simple",
			&ast.StoredProcStmt{Name: "p1", Body: &ast.PLBlock{
				Statements: []ast.Node{&ast.PLSetStmt{Name: "x", Expr: lit("1")}},
			}},
			"CREATE PROCEDURE p1 () BEGIN SET x = 1; END",
		},
		{
			"pl-return",
			&ast.StoredProcStmt{Function: true, Name: "f1", Returns: "INT", Body: &ast.PLBlock{
				Statements: []ast.Node{&ast.PLReturnStmt{Expr: lit("1")}},
			}},
			"CREATE FUNCTION f1 () RETURNS INT BEGIN RETURN 1; END",
		},
		{"drop-routine", &ast.DropRoutineStmt{Name: "p1"}, "DROP PROCEDURE p1"},
		{"alter-system", &ast.AlterSystemStmt{Name: "x", Value: lit("1")}, "ALTER SYSTEM SET x = 1"},
		{"flashback-table", &ast.FlashbackTableStmt{Tables: []*ast.Identifier{ident("t")}}, "FLASHBACK TABLE t"},
		{"lock-user", &ast.LockUserStmt{Users: []ast.UserSpec{{Name: "u"}}}, "LOCK USER 'u'"},
		{"unlock-user", &ast.LockUserStmt{Unlock: true, Users: []ast.UserSpec{{Name: "u"}}}, "UNLOCK USER 'u'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatter.Format(tc.stmt)
			if got != tc.want {
				t.Fatalf("Format() = %q\nwant         %q", got, tc.want)
			}
		})
	}
}

func TestEmptyAllStatementKinds(t *testing.T) {
	// Breadth guard: every statement kind must format without panicking,
	// even for a zero-value instance (guard clauses run first).
	stmts := []ast.Statement{
		&ast.AlterDatabaseStmt{}, &ast.AlterLSStmt{}, &ast.AlterMViewStmt{}, &ast.AlterOutlineStmt{},
		&ast.AlterResourcePoolStmt{}, &ast.AlterResourceUnitStmt{}, &ast.AlterSequenceStmt{},
		&ast.AlterSystemStmt{}, &ast.AlterTableStmt{}, &ast.AlterTablegroupStmt{}, &ast.AlterTablespaceStmt{},
		&ast.AlterTenantStmt{}, &ast.AlterViewStmt{}, &ast.AuditStmt{}, &ast.BeginStmt{}, &ast.CallStmt{},
		&ast.CatalogStmt{}, &ast.CCLLimitRuleStmt{}, &ast.CloneTenantStmt{}, &ast.CommitStmt{},
		&ast.CreateDatabaseStmt{}, &ast.CreateIndexStmt{}, &ast.CreateLocationStmt{}, &ast.CreateMViewStmt{},
		&ast.CreateOutlineStmt{}, &ast.CreateResourcePoolStmt{}, &ast.CreateResourceUnitStmt{},
		&ast.CreateRoleStmt{}, &ast.CreateSequenceStmt{}, &ast.CreateStandbyTenantStmt{},
		&ast.CreateTablegroupStmt{}, &ast.CreateTableLikeStmt{}, &ast.CreateTablespaceStmt{},
		&ast.CreateTableStmt{}, &ast.CreateTenantStmt{}, &ast.CreateUserStmt{}, &ast.CreateViewStmt{},
		&ast.DblinkStmt{}, &ast.DeallocateStmt{}, &ast.DeleteStmt{}, &ast.DoStmt{},
		&ast.DropDatabaseStmt{}, &ast.DropIndexStmt{}, &ast.DropMViewStmt{}, &ast.DropOutlineStmt{},
		&ast.DropResourcePoolStmt{}, &ast.DropResourceUnitStmt{}, &ast.DropRoleStmt{}, &ast.DropRoutineStmt{},
		&ast.DropSequenceStmt{}, &ast.DropTablegroupStmt{}, &ast.DropTablespaceStmt{}, &ast.DropTableStmt{},
		&ast.DropTenantStmt{}, &ast.DropUserStmt{}, &ast.DropViewStmt{}, &ast.DumpMemoryStmt{},
		&ast.ExecuteStmt{}, &ast.ExplainStmt{}, &ast.FlashbackStandbyLogStmt{}, &ast.FlashbackTableStmt{},
		&ast.FlushStmt{}, &ast.GetDiagnosticsStmt{}, &ast.GrantStmt{}, &ast.HandlerStmt{}, &ast.HelpStmt{},
		&ast.InsertStmt{}, &ast.InstallPluginStmt{}, &ast.KillStmt{}, &ast.LoadDataStmt{},
		&ast.LocationUtilsStmt{}, &ast.LockTablesStmt{}, &ast.LockUserStmt{}, &ast.MergeStmt{},
		&ast.MLogStmt{}, &ast.PrepareStmt{}, &ast.PurgeStmt{}, &ast.RecoverTenantStmt{},
		&ast.ReleaseSavepointStmt{}, &ast.RenameTableStmt{}, &ast.RenameUserStmt{}, &ast.RestorePointStmt{},
		&ast.RevokeStmt{}, &ast.RollbackStmt{}, &ast.RollbackToSavepointStmt{}, &ast.RotateMasterKeyStmt{},
		&ast.RoutineLoadActionStmt{}, &ast.RoutineLoadStmt{}, &ast.SavepointStmt{}, &ast.SelectStmt{},
		&ast.SensitiveRuleStmt{}, &ast.ServiceStmt{}, &ast.SetPasswordStmt{}, &ast.SetRoleStmt{},
		&ast.SetStmt{}, &ast.SetTransactionStmt{}, &ast.ShowStmt{}, &ast.SignalStmt{}, &ast.StoredProcStmt{},
		&ast.SwitchoverTenantStmt{}, &ast.TableMaintStmt{}, &ast.TenantSnapshotStmt{},
		&ast.TransferPartitionStmt{}, &ast.TruncateStmt{}, &ast.UninstallPluginStmt{}, &ast.UnlockTablesStmt{},
		&ast.UpdateStmt{}, &ast.UseDbStmt{}, &ast.XAStmt{},
	}
	for _, s := range stmts {
		// zero-value instance must format without panicking.
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%T: Format panicked: %v", s, r)
				}
			}()
			_ = formatter.Format(s)
		}()
	}
}

func TestPrettyMultiLine(t *testing.T) {
	got := formatter.FormatWith(sel(), formatter.Options{Pretty: true})
	if !strings.Contains(got, "\n") {
		t.Fatalf("pretty output has no newlines: %q", got)
	}
	if !strings.HasPrefix(got, "SELECT a,") {
		t.Fatalf("pretty output wrong: %q", got)
	}

	complex := &ast.SelectStmt{Bodies: []*ast.SelectBody{{
		Items: []ast.ColumnItem{
			{Expr: &ast.Identifier{Parts: []string{"a"}}},
			{Expr: &ast.FuncCall{Name: "count", Star: true}},
		},
		From:  []*ast.TableRef{{RefKind: ast.RefTable, Name: &ast.Identifier{Parts: []string{"t"}}, Alias: "x", HasAlias: true}},
		Where: &ast.BinaryExpr{Op: "AND", Left: lit("1"), Right: lit("2")},
	}}}
	pretty := formatter.FormatWith(complex, formatter.Options{Pretty: true})
	if !strings.Contains(pretty, "\n  FROM") || !strings.Contains(pretty, "\n  WHERE") {
		t.Fatalf("pretty multi-clause output missing lines: %q", pretty)
	}
}

func TestStringLiteralEscaping(t *testing.T) {
	s := &ast.Literal{LitType: ast.LitString, Value: "it's\na\n\x00"}
	out := formatter.Format(s)
	if out != "'it''s\\na\\n\\0'" {
		t.Fatalf("escape output = %q", out)
	}
}
