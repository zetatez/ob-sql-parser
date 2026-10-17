package gparser

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/formatter"
	"github.com/zetatez/ob-sql-parser/lexer"
	"github.com/zetatez/ob-sql-parser/version"
)

// TestVersionCorpus pins the "strict 4.2.5 grammar surface": statements legal
// in 4.2.5 must parse under WithVersion("4.2.5"), while post-4.2.5 features
// must be rejected there and accepted unversioned.
func TestVersionCorpus(t *testing.T) {
	v425 := Options{Version: version.MustParse("4.2.5")}
	legal := []string{
		"SELECT a FROM t",
		"CREATE TABLE t (id INT)",
		"ALTER TABLE t ADD COLUMN c INT",
		"GRANT SELECT ON db1.* TO 'u'@'%'",
		"CREATE PROCEDURE p() BEGIN SET a = 1; END",
		"WITH c AS (SELECT 1) SELECT * FROM c",
		"ALTER SYSTEM SWITCHOVER TO PRIMARY TENANT t",
		"CREATE DATABASE LINK dl CONNECT TO 'u' DATABASE 'd' IDENTIFIED BY 'p'",
		"CREATE OUTLINE o ON SELECT 1",
		"ALTER SYSTEM TRANSFER PARTITION TABLE_ID = 1, OBJECT_ID = 0 TO LS 1",
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) PARTITIONS AUTO",
	}
	post425 := []string{
		"CREATE MATERIALIZED VIEW mv AS SELECT 1",
		"CREATE MATERIALIZED VIEW LOG ON t",
		"CREATE EXTERNAL CATALOG c PROPERTIES ('k' = 'v')",
		"CREATE LOCATION loc URL 'oss://x'",
		"CREATE ROUTINE LOAD job ON t FROM STREAM ()",
		"CREATE CONCURRENT_LIMITING_RULE r1 ON db1.t TO 'u'@'%' FOR CCL_DML_OPTION",
		"CREATE SENSITIVE RULE s1 ON a USING MASK",
		"CREATE SNAPSHOT snap FOR TENANT t",
		"ALTER SYSTEM FLASHBACK STANDBY LOG TO SCN = 1",
		"SELECT a FROM t GROUP BY GROUPING SETS ((a), (b))",
	}
	opts := lexer.Options{AllowDoubleQuotedIdentifiers: true}
	for _, sql := range legal {
		toks, _ := lexer.LexAll(sql, opts)
		if _, err := Parse(toks, v425); err != nil {
			t.Errorf("4.2.5 should accept %q: %v", sql, err)
		}
	}
	for _, sql := range post425 {
		toks, _ := lexer.LexAll(sql, opts)
		if _, err := Parse(toks, v425); err == nil {
			t.Errorf("4.2.5 should REJECT %q", sql)
		}
		if stmts, err := Parse(toks, Options{}); err != nil {
			t.Errorf("unversioned should accept %q: %v", sql, err)
		} else {
			_ = formatter.Format(stmts[0])
		}
	}
}
