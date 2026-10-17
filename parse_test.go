package obparser

import (
	"strings"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	cases := []string{
		"SELECT * FROM t",
		"SELECT a, b FROM t WHERE a > 1 AND b = 'x'",
		"SELECT t.a, COUNT(*) FROM t GROUP BY t.a HAVING COUNT(*) > 1",
		"SELECT * FROM a JOIN b ON a.id = b.aid WHERE a.x = 2 ORDER BY a.x DESC LIMIT 10",
		"INSERT INTO t (a,b) VALUES (1,'x'),(2,'y')",
		"INSERT INTO t VALUES (1,2) ON DUPLICATE KEY UPDATE a = VALUES(a)",
		"UPDATE t SET a = 1, b = 'z' WHERE id = 5",
		"DELETE FROM t WHERE id IN (1,2,3)",
		"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL)",
		"CREATE TABLE a LIKE b",
		"DROP TABLE IF EXISTS t",
		"SHOW TABLES",
		"SHOW CREATE TABLE t",
		"USE mydb",
		"SET @x = 1",
		"EXPLAIN SELECT * FROM t",
		"SELECT a BETWEEN 1 AND 5 FROM t",
		"SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t",
		"SELECT DATE_ADD('2020-01-01', INTERVAL 1 DAY)",
		"SELECT a FROM t UNION ALL SELECT b FROM u",
		"SELECT * FROM t WHERE a IS NOT NULL",
		"SELECT * FROM t LIMIT 10, 5",
		"SELECT ROW_NUMBER() OVER (PARTITION BY a ORDER BY b DESC) FROM t",
		"TRUNCATE TABLE t",
		"DROP INDEX idx ON t",
		"CREATE UNIQUE INDEX idx ON t (a)",
		"CREATE DATABASE IF NOT EXISTS db1",
		"DROP DATABASE db1",
		"SELECT * FROM (SELECT a FROM t) sub WHERE sub.a = 1",
		"SELECT a FROM t WHERE a NOT IN (1,2)",
		"SELECT * FROM t WHERE name LIKE '%foo%'",
		"SELECT x FROM t ORDER BY x DESC NULLS FIRST",
		"SELECT CAST(a AS CHAR(3)) FROM t",
		"BEGIN",
		"COMMIT",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			x, err := Parse(sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			s1 := Format(x)
			y, err := Parse(s1)
			if err != nil {
				t.Fatalf("re-parse of formatted output failed: %v\n  sql: %s\n  fmt: %s", err, sql, s1)
			}
			s2 := Format(y)
			if s1 != s2 {
				t.Fatalf("format not idempotent\n  s1: %s\n  s2: %s", s1, s2)
			}
		})
	}
}

func TestParseBasics(t *testing.T) {
	st, err := Parse("SELECT a, b FROM t WHERE a > 1")
	if err != nil {
		t.Fatal(err)
	}
	_ = st
	if got := Format(st); got != "SELECT a, b FROM t WHERE a > 1" {
		t.Fatalf("unexpected format: %q", got)
	}
}

func TestParseMulti(t *testing.T) {
	stmts, err := ParseMulti("SELECT 1; SELECT 2; INSERT INTO t VALUES (3);")
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(stmts))
	}
}

func TestComments(t *testing.T) {
	st, err := Parse("SELECT 1 -- trailing comment\n")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(Format(st)) != "SELECT 1" {
		t.Fatalf("got %q", Format(st))
	}
}

func TestParseError(t *testing.T) {
	_, err := Parse("SELECT FROM")
	if err == nil {
		t.Fatal("expected an error")
	}
	e := err.Error()
	if !strings.Contains(e, ":") {
		t.Fatalf("error should carry position, got %q", e)
	}
}
