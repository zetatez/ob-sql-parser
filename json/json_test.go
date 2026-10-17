package json_test

import (
	"strings"
	"testing"

	"github.com/zetatez/ob-sql-parser"
	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/formatter"
	jsonpkg "github.com/zetatez/ob-sql-parser/json"
)

func TestRoundTripJSON(t *testing.T) {
	cases := []string{
		"SELECT a, b FROM t WHERE a > 1 AND b = 'x'",
		"SELECT t.a, COUNT(*) FROM t GROUP BY t.a HAVING COUNT(*) > 1",
		"SELECT * FROM a JOIN b ON a.id = b.aid ORDER BY a.x DESC LIMIT 10",
		"INSERT INTO t (a,b) VALUES (1,'x'),(2,'y') ON DUPLICATE KEY UPDATE a = VALUES(a)",
		"UPDATE t SET a = 1 WHERE id = 5",
		"DELETE FROM t WHERE id IN (1,2,3)",
		"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT 'x')",
		"SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t",
		"SELECT a FROM t UNION ALL SELECT b FROM u",
		"SELECT ROW_NUMBER() OVER (PARTITION BY a ORDER BY b DESC) FROM t",
		"SELECT * FROM t WHERE a IS NOT NULL AND b NOT IN (1,2)",
		"CREATE PROCEDURE p(IN x INT) BEGIN DECLARE v INT; IF x > 0 THEN SET v = x; ELSE SET v = 0; END IF; RETURN v; END",
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10))",
		"GRANT SELECT ON db1.* TO 'u'@'%'",
		"ALTER TABLE t ADD COLUMN c INT AFTER b, DROP INDEX i",
		"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN UPDATE SET a = u.a",
		"SHOW COLUMNS FROM t LIKE 'a%'",
		"SHOW CREATE TABLE t",
		"SHOW GRANTS FOR 'u'@'%'",
		"SET GLOBAL sql_mode = 'STRICT_ALL_TABLES'",
		"SET @x = 1, @y = 2",
		"LOCK TABLES t WRITE, u READ",
		"UNLOCK TABLES",
		"XA START '1,2,3'",
		"SAVEPOINT sp1",
		"ROLLBACK TO SAVEPOINT sp1",
		"RELEASE SAVEPOINT sp1",
		"CREATE SEQUENCE s START WITH 1 INCREMENT BY 2 CACHE 20",
		"DROP TABLE IF EXISTS a, b",
		"RENAME TABLE a TO b, c TO d",
		"TRUNCATE TABLE t",
		"CREATE ROLE r1",
		"SET ROLE r1",
		"CREATE USER 'u'@'%' IDENTIFIED BY 'pw'",
		"DROP USER 'u'@'%'",
		"REVOKE ALL ON *.* FROM 'u'@'%'",
		"START TRANSACTION",
		"USE mydb",
		"KILL QUERY 123",
		"DESCRIBE t",
		"ANALYZE TABLE t",
		"OPTIMIZE TABLE t",
		"CREATE PROCEDURE p() BEGIN DECLARE v INT DEFAULT 0; WHILE v < 10 DO SET v = v + 1; END WHILE; END",
		"CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT a FROM t; OPEN c; CLOSE c; END",
		"CALL p(1)",
		"DROP PROCEDURE IF EXISTS p",
		"FLUSH PRIVILEGES",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			st := obparser.MustParse(sql)
			want := formatter.Format(st)

			data, err := jsonpkg.Marshal(st)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			back, err := jsonpkg.Unmarshal(data)
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if back == nil {
				t.Fatal("nil result")
			}
			got := formatter.Format(back)
			if got != want {
				t.Fatalf("not preserved through JSON\n  want: %s\n  got:  %s", want, got)
			}
		})
	}
}

func TestDumpUnmarshal(t *testing.T) {
	sqls := []string{"SELECT 1", "INSERT INTO t VALUES (1)", "CREATE TABLE a LIKE b"}
	var nodes []ast.Node
	for _, s := range sqls {
		nodes = append(nodes, obparser.MustParse(s))
	}
	data, err := jsonpkg.Dump(nodes)
	if err != nil {
		t.Fatal(err)
	}
	back, err := jsonpkg.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	lst, ok := back.(interface{ Nodes() []ast.Node })
	if !ok {
		t.Fatalf("expected statement list, got %T", back)
	}
	if len(lst.Nodes()) != len(sqls) {
		t.Fatalf("want %d, got %d", len(sqls), len(lst.Nodes()))
	}
}

func TestJSONUnknownKind(t *testing.T) {
	if _, err := jsonpkg.Unmarshal([]byte(`{"kind":"NoSuchNode"}`)); err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestMarshalIndentAndWrite(t *testing.T) {
	st := obparser.MustParse("SELECT a FROM t")
	compact, err := jsonpkg.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	indented, err := jsonpkg.MarshalIndent(st)
	if err != nil {
		t.Fatal(err)
	}
	if len(indented) <= len(compact) {
		t.Fatalf("indented output should be larger: %d vs %d", len(indented), len(compact))
	}
	if !strings.Contains(jsonpkg.String(st), "\"kind\"") {
		t.Fatalf("String() should contain the serialized node")
	}
	var sb strings.Builder
	if err := jsonpkg.Write(&sb, st); err != nil {
		t.Fatal(err)
	}
	if sb.Len() == 0 {
		t.Fatal("Write produced no output")
	}
}

func TestDecodeAlias(t *testing.T) {
	st := obparser.MustParse("SELECT 1")
	data, err := jsonpkg.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	back, err := jsonpkg.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if back == nil {
		t.Fatal("Decode returned nil")
	}
}

func TestUnmarshalErrors(t *testing.T) {
	bad := []string{
		``,
		`{`,
		`{"kind": 123}`,
		`"string"`,
	}
	for _, b := range bad {
		if _, err := jsonpkg.Unmarshal([]byte(b)); err == nil {
			t.Errorf("Unmarshal(%q) should fail", b)
		}
	}
}

func TestDumpEmptyAndSingle(t *testing.T) {
	data, err := jsonpkg.Dump(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("Dump(nil) should still produce a list")
	}
	one, err := jsonpkg.Dump([]ast.Node{obparser.MustParse("SELECT 1")})
	if err != nil {
		t.Fatal(err)
	}
	back, err := jsonpkg.Unmarshal(one)
	if err != nil {
		t.Fatal(err)
	}
	if lst, ok := back.(interface{ Nodes() []ast.Node }); !ok || len(lst.Nodes()) != 1 {
		t.Fatalf("expected single-statement list, got %T", back)
	}
}
