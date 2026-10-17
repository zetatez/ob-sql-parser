package gparser

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/formatter"
	"github.com/zetatez/ob-sql-parser/lexer"
)

// TestCoverageSupported is a living spec: one representative statement per
// supported category must parse and round-trip. Add a case here when a new
// statement kind is implemented.
func TestCoverageSupported(t *testing.T) {
	cases := []string{
		// queries / DML
		"SELECT a, b FROM t WHERE a > 1 ORDER BY a DESC LIMIT 10",
		"SELECT * FROM a JOIN b ON a.id = b.id",
		"SELECT a FROM t UNION ALL SELECT b FROM u",
		"SELECT ROW_NUMBER() OVER (PARTITION BY a ORDER BY b) FROM t",
		"INSERT INTO t (a, b) VALUES (1, 2)",
		"UPDATE t SET a = 1 WHERE id = 2",
		"DELETE FROM t WHERE id = 3",
		// DDL
		"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(10) NOT NULL)",
		"CREATE TABLE t (id INT) PARTITION BY HASH (id) PARTITIONS 4",
		"CREATE TABLE t (a INT, b INT AS (a + 1) STORED)",
		"CREATE TABLE t AS SELECT * FROM u",
		"CREATE INDEX i ON t (a)",
		"CREATE TABLE a LIKE b",
		"ALTER TABLE t ADD COLUMN c INT",
		"DROP TABLE IF EXISTS t",
		"TRUNCATE TABLE t",
		// databases / views / sequences / roles / users
		"CREATE DATABASE IF NOT EXISTS db1",
		"ALTER DATABASE db1 DEFAULT CHARACTER SET utf8mb4",
		"CREATE OR REPLACE VIEW v AS SELECT 1",
		"CREATE SEQUENCE s START WITH 1 INCREMENT BY 1",
		"CREATE USER 'u'@'%' IDENTIFIED BY 'pwd'",
		"CREATE ROLE IF NOT EXISTS r1",
		"GRANT SELECT ON db1.* TO 'u'@'%'",
		"REVOKE SELECT ON db1.* FROM 'u'@'%'",
		"SET ROLE NONE",
		"SET PASSWORD FOR 'u'@'%' = 'pwd'",
		// session / utility
		"SHOW TABLES",
		"SHOW CREATE TABLE t",
		"SET @x = 1",
		"USE mydb",
		"EXPLAIN SELECT 1",
		"DESCRIBE t",
		"PREPARE s FROM 'SELECT 1'",
		"EXECUTE s USING @a",
		"DEALLOCATE PREPARE s",
		"KILL QUERY 1",
		"HELP 'SELECT'",
		// transactions / locks
		"BEGIN",
		"COMMIT",
		"ROLLBACK",
		"SAVEPOINT s1",
		"ROLLBACK TO SAVEPOINT s1",
		"SET TRANSACTION ISOLATION LEVEL READ COMMITTED",
		"LOCK TABLES t READ",
		"UNLOCK TABLES",
		"XA START 'x1'",
		"XA COMMIT 'x1'",
		// ops
		"RENAME TABLE a TO b",
		"ANALYZE TABLE t",
		"OPTIMIZE TABLE t",
		"FLUSH TABLES WITH READ LOCK",
		"LOAD DATA INFILE 'f' INTO TABLE t",
		"PURGE RECYCLEBIN",
		"SELECT a FROM t WHERE a MEMBER OF (1, 2)",
		"SELECT MATCH(a) AGAINST ('x' IN BOOLEAN MODE) FROM t",
		"SELECT a FROM t GROUP BY CUBE (a, b)",
		"ALTER SYSTEM SET enable_x = true",
		"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN DELETE",
		"FLASHBACK TABLE t TO BEFORE DROP",
		"CREATE TENANT t RESOURCE_POOL_LIST = ('p')",
		"CREATE RESOURCE UNIT u MAX_CPU = 1",
		"CREATE TABLEGROUP g",
		"CREATE MATERIALIZED VIEW mv AS SELECT 1",
		"SIGNAL SQLSTATE '45000'",
		"AUDIT SELECT ON db1.*",
		"CREATE DATABASE LINK dl CONNECT TO 'u' DATABASE 'd' IDENTIFIED BY 'p'",
		"CREATE ROUTINE LOAD job ON t FROM STREAM ()",
		"SELECT * FROM t, LATERAL (SELECT 1) x",
		"SELECT a FROM t QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1",
		"CREATE MATERIALIZED VIEW LOG ON t WITH PRIMARY KEY",
		"CREATE OUTLINE o ON SELECT 1",
		"ALTER SYSTEM SWITCHOVER TO PRIMARY TENANT t",
		"CREATE RESTORE POINT rp",
		"CREATE PROCEDURE p() BEGIN SET a = 1; END",
		"CREATE FUNCTION f() RETURNS INT BEGIN RETURN 1; END",
		"CALL p()",
		"DROP PROCEDURE p",
	}
	opts := lexer.Options{AllowDoubleQuotedIdentifiers: true}
	for _, sql := range cases {
		toks, err := lexer.LexAll(sql, opts)
		if err != nil {
			t.Errorf("lex %q: %v", sql, err)
			continue
		}
		stmts, err := Parse(toks, Options{})
		if err != nil {
			t.Errorf("parse %q: %v", sql, err)
			continue
		}
		s1 := formatter.Format(stmts[0])
		toks2, _ := lexer.LexAll(s1, opts)
		stmts2, err := Parse(toks2, Options{})
		if err != nil {
			t.Errorf("re-parse %q (from %q): %v", s1, sql, err)
			continue
		}
		if s2 := formatter.Format(stmts2[0]); s2 != s1 {
			t.Errorf("not idempotent for %q: %q != %q", sql, s1, s2)
		}
	}
}

// TestCoverageUnsupported documents statements that are intentionally out of
// scope today. When one is implemented, move it to TestCoverageSupported.
func TestCoverageUnsupported(t *testing.T) {
	cases := []string{
		"CREATE PROCEDURE p() BEGIN NOT A PL STATEMENT END",
		"SELECT a FROM t GROUP BY a WITH GROUPING SETS ((a))",
	}
	opts := lexer.Options{AllowDoubleQuotedIdentifiers: true}
	for _, sql := range cases {
		toks, err := lexer.LexAll(sql, opts)
		if err != nil {
			continue // lexical rejection also counts as unsupported
		}
		if _, err := Parse(toks, Options{}); err == nil {
			t.Errorf("expected %q to be unsupported, but it parsed", sql)
		}
	}
}
