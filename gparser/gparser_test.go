package gparser

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/formatter"
	"github.com/zetatez/ob-sql-parser/lexer"
	"github.com/zetatez/ob-sql-parser/token"
	"github.com/zetatez/ob-sql-parser/version"
)

var lexOpts = lexer.Options{AllowDoubleQuotedIdentifiers: true}

func parseAll(t *testing.T, sql string, opts ...Options) []ast.Statement {
	t.Helper()
	toks, err := lexer.LexAll(sql, lexOpts)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	stmts, err := Parse(toks, o)
	if err != nil {
		t.Fatalf("parse %q: %v", sql, err)
	}
	return stmts
}

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
			stmts := parseAll(t, sql)
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			s1 := formatter.Format(stmts[0])
			stmts2 := parseAll(t, s1)
			if len(stmts2) != 1 {
				t.Fatalf("re-parse of formatted output gave %d statements\n  s1: %s", len(stmts2), s1)
			}
			s2 := formatter.Format(stmts2[0])
			if s1 != s2 {
				t.Fatalf("format not idempotent\n  s1: %s\n  s2: %s", s1, s2)
			}
		})
	}
}

func TestExtraRoundTrip(t *testing.T) {
	cases := []string{
		"SELECT * FROM a LEFT JOIN b ON a.id = b.id",
		"SELECT * FROM a RIGHT JOIN b ON a.id = b.id",
		"SELECT * FROM a INNER JOIN b ON a.id = b.id",
		"SELECT * FROM a JOIN b USING (id)",
		"SELECT * FROM a AS x JOIN b AS y ON x.id = y.id",
		"CREATE TABLE t (id INT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (id))",
		"CREATE TABLE t (id INT, name VARCHAR(10), KEY (name))",
		"CREATE TABLE t (id INT, UNIQUE KEY (id))",
		"CREATE INDEX idx ON t (a, b)",
		"SELECT a NOT LIKE 'x' FROM t",
		"SELECT COUNT(DISTINCT a) FROM t",
		"SELECT a IS TRUE FROM t",
		"SELECT NOT a FROM t",
		"SELECT -1, +2, ~3 FROM t",
		"SELECT @@global.x, @y",
		"SHOW FULL TABLES",
		"SHOW VARIABLES LIKE 'x'",
		"SHOW DATABASES",
		"START TRANSACTION READ ONLY",
		"ROLLBACK",
		"REPLACE INTO t VALUES (1)",
		"INSERT INTO t SET a = 1",
		"INSERT IGNORE INTO t VALUES (1)",
		"DELETE FROM t ORDER BY a LIMIT 1",
		"UPDATE t SET a = 1 ORDER BY a LIMIT 2",
		"SELECT 1 WHERE 1 IN (2,3)",
		"SELECT CASE a WHEN 1 THEN 'x' WHEN 2 THEN 'y' END FROM t",
		"SELECT a FROM t WHERE b BETWEEN 1 AND 2 AND c = 3",
		"SELECT a FROM t UNION SELECT b FROM u",
		"SELECT * FROM t WHERE EXISTS (SELECT 1 FROM u)",
		"SELECT CAST(a AS DECIMAL(10, 2)) FROM t",
		"SET @x = @y",
		"CREATE DATABASE db1 DEFAULT CHARACTER SET utf8",
		"SELECT a FROM t WHERE a = 1 OR b = 2 XOR c = 3",
		"SELECT 1 + 2 * 3 - 4 / 2",
		"SELECT a FROM t LIMIT 5 OFFSET 10",
		"SELECT * FROM t WHERE a <=> b",
		"SELECT INTERVAL 1 DAY",
		"SELECT a DIV 2 FROM t",
		"SELECT a MOD 2 FROM t",
		"SELECT a RLIKE 'x' FROM t",
		"SELECT a NOT RLIKE 'x' FROM t",
		"SELECT a SOUNDS LIKE b FROM t",
		"SELECT * FROM a NATURAL JOIN b",
		"SELECT * FROM a NATURAL INNER JOIN b",
		"SELECT * FROM a NATURAL LEFT JOIN b",
		"SELECT * FROM a NATURAL RIGHT OUTER JOIN b",
		"SHOW TABLE STATUS",
		"SHOW TABLE STATUS FROM mydb",
		"SELECT CAST(a AS BINARY) FROM t",
		"DESC t",
		"DESCRIBE mydb.t",
		// index hints
		"SELECT * FROM t USE INDEX (i)",
		"SELECT * FROM t FORCE INDEX (i, j)",
		"SELECT * FROM t IGNORE INDEX (i)",
		"SELECT * FROM t USE INDEX FOR JOIN (i)",
		"SELECT * FROM t FORCE INDEX FOR ORDER BY (i)",
		"SELECT * FROM t USE INDEX FOR GROUP BY (PRIMARY)",
		"SELECT * FROM t AS x USE INDEX (i) WHERE x.a = 1",
		"SELECT * FROM a JOIN b FORCE INDEX (i) ON a.id = b.id",
		// SHOW ... FROM
		"SHOW COLUMNS FROM t",
		"SHOW INDEX FROM t",
		"SHOW FULL COLUMNS FROM t",
		"SHOW COLUMNS FROM t FROM mydb",
		"SHOW INDEX FROM t FROM mydb",
		// ALTER TABLE
		"ALTER TABLE t ADD COLUMN c INT",
		"ALTER TABLE t ADD c INT NOT NULL",
		"ALTER TABLE t DROP COLUMN c",
		"ALTER TABLE t MODIFY COLUMN c INT NOT NULL",
		"ALTER TABLE t CHANGE COLUMN a b INT",
		"ALTER TABLE t RENAME TO t2",
		"ALTER TABLE t RENAME INDEX i TO j",
		"ALTER TABLE t ADD INDEX i (a, b)",
		"ALTER TABLE t ADD UNIQUE INDEX i (a)",
		"ALTER TABLE t ADD PRIMARY KEY (a)",
		"ALTER TABLE t DROP INDEX i",
		"ALTER TABLE t DROP PRIMARY KEY",
		"ALTER TABLE t ALTER COLUMN c SET DEFAULT 1",
		"ALTER TABLE t ALTER COLUMN c DROP DEFAULT",
		"ALTER TABLE t ADD COLUMN a INT, DROP COLUMN b",
		// GROUP BY ... WITH ROLLUP
		"SELECT a, COUNT(*) FROM t GROUP BY a WITH ROLLUP",
		// window frames
		"SELECT ROW_NUMBER() OVER (ORDER BY a ROWS UNBOUNDED PRECEDING) FROM t",
		"SELECT SUM(a) OVER (PARTITION BY b ORDER BY c RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t",
		"SELECT SUM(a) OVER (ORDER BY c ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t",
		// COLLATE
		"SELECT a FROM t ORDER BY a COLLATE utf8_bin",
		"SELECT a FROM t WHERE a = b COLLATE utf8_general_ci",
		// table options
		"CREATE TABLE t (id INT) ENGINE=InnoDB",
		"CREATE TABLE t (id INT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
		"CREATE TABLE t (id INT) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		"CREATE TABLE t (id INT) AUTO_INCREMENT=100 COMMENT=hello",
		"CREATE TABLE t (id INT) TABLEGROUP=tg1",
		"CREATE TABLE t (id INT) TABLEGROUP = tg1 COMMENT='x'",
		"CREATE TABLE t (a VARCHAR(10) COLLATE utf8_bin, b INT)",
		"CREATE DATABASE db1 COLLATE utf8_bin",
		// partitioning
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (MAXVALUE))",
		"CREATE TABLE t (id INT) PARTITION BY HASH (id) PARTITIONS 4",
		"CREATE TABLE t (id INT) PARTITION BY KEY (id) PARTITIONS 4",
		"CREATE TABLE t (id INT) PARTITION BY LIST (id) (PARTITION p0 VALUES IN (1, 2), PARTITION p1 VALUES IN (3))",
		"CREATE TABLE t (id INT) PARTITION BY RANGE COLUMNS (id) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE t (id INT) ENGINE=InnoDB PARTITION BY HASH (id) PARTITIONS 8",
		"SELECT * FROM t PARTITION (p0, p1)",
		"SELECT * FROM t PARTITION (p0) AS x WHERE x.a = 1",
		// generated columns
		"CREATE TABLE t (a INT, b INT AS (a + 1))",
		"CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a * 2) STORED)",
		"CREATE TABLE t (a INT, b INT AS (a) VIRTUAL)",
		// CTAS
		"CREATE TABLE t AS SELECT * FROM u",
		"CREATE TABLE t (a INT) AS SELECT * FROM u",
		"CREATE TABLE t (a INT) AS SELECT b FROM u WHERE b > 1",
		"CREATE TABLE t SELECT 1",
		// views
		"CREATE VIEW v AS SELECT * FROM t",
		"CREATE OR REPLACE VIEW v (a, b) AS SELECT a, b FROM t",
		"CREATE VIEW v AS SELECT a FROM t WHERE a > 1",
		"DROP VIEW v",
		"DROP VIEW IF EXISTS v1, v2",
		// ALTER DATABASE
		"ALTER DATABASE db1 DEFAULT CHARACTER SET utf8mb4",
		"ALTER DATABASE db1 COLLATE utf8mb4_bin",
		"ALTER SCHEMA db1 DEFAULT CHARACTER SET utf8",
		// sequences
		"CREATE SEQUENCE s",
		"CREATE SEQUENCE IF NOT EXISTS s START WITH 1 INCREMENT BY 1",
		"CREATE SEQUENCE s MINVALUE 1 MAXVALUE 100 CACHE 20 CYCLE ORDER",
		"CREATE SEQUENCE s NOCYCLE NOORDER",
		"DROP SEQUENCE s",
		"DROP SEQUENCE IF EXISTS s",
		// transaction control
		"SET TRANSACTION ISOLATION LEVEL READ COMMITTED",
		"SET TRANSACTION READ ONLY",
		"SET TRANSACTION READ WRITE",
		"SET GLOBAL TRANSACTION ISOLATION LEVEL REPEATABLE READ",
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE, READ WRITE",
		"SAVEPOINT s1",
		"RELEASE SAVEPOINT s1",
		"ROLLBACK TO SAVEPOINT s1",
		"ROLLBACK TO s1",
		"ROLLBACK WORK TO SAVEPOINT s1",
		// table locks
		"LOCK TABLES t READ",
		"LOCK TABLES t AS x WRITE, u READ LOCAL",
		"UNLOCK TABLES",
		// SHOW completeness
		"SHOW CREATE VIEW v",
		"SHOW CREATE PROCEDURE p",
		"SHOW CREATE FUNCTION f",
		"SHOW CREATE TRIGGER tr",
		"SHOW TRIGGERS",
		"SHOW TRIGGERS FROM mydb",
		"SHOW CHARACTER SET",
		"SHOW CHARACTER SET LIKE 'utf8%'",
		"SHOW PRIVILEGES",
		"SHOW PLUGINS",
		"SHOW PROFILES",
		"SHOW OPEN TABLES",
		"SHOW OPEN TABLES FROM mydb LIKE 'x'",
		"SHOW TABLE STATUS LIKE 'x'",
		"SHOW FULL PROCESSLIST",
		// GRANT / REVOKE
		"GRANT SELECT ON db1.* TO 'u'@'%'",
		"GRANT SELECT, INSERT, UPDATE ON db1.t TO 'u'@'localhost', 'v'@'%'",
		"GRANT ALL ON *.* TO 'u'@'%' WITH GRANT OPTION",
		"GRANT SELECT (a, b) ON db1.t TO 'u'@'%'",
		"REVOKE SELECT ON db1.* FROM 'u'@'%'",
		"REVOKE ALL PRIVILEGES ON *.* FROM 'u'@'%'",
		// CREATE/DROP USER
		"CREATE USER 'u'@'%'",
		"CREATE USER IF NOT EXISTS 'u'@'localhost' IDENTIFIED BY 'pwd'",
		"CREATE USER 'u'@'%', 'v'@'%'",
		"DROP USER 'u'@'%'",
		"DROP USER IF EXISTS 'u'@'%', 'v'@'%'",
		// PREPARE / EXECUTE / DEALLOCATE
		"PREPARE s FROM 'SELECT 1'",
		"PREPARE s FROM @q",
		"EXECUTE s",
		"EXECUTE s USING @a, @b",
		"DEALLOCATE PREPARE s",
		// "user" must remain usable as an identifier
		"SELECT user FROM t",
		"SELECT * FROM user",
		// roles
		"CREATE ROLE r1",
		"CREATE ROLE IF NOT EXISTS r1, r2",
		"DROP ROLE r1",
		"DROP ROLE IF EXISTS r1, r2",
		"SET ROLE r1, r2",
		"SET ROLE NONE",
		"SET ROLE DEFAULT",
		"SET PASSWORD = 'pwd'",
		"SET PASSWORD FOR 'u'@'%' = 'pwd'",
		// misc
		"KILL 123",
		"KILL CONNECTION 123",
		"KILL QUERY 123",
		"HELP 'SELECT'",
		"ANALYZE TABLE t",
		"OPTIMIZE TABLE t1, t2",
		"REPAIR TABLE t",
		// RENAME TABLE
		"RENAME TABLE a TO b",
		"RENAME TABLE a TO b, c TO d",
		// LOAD DATA
		"LOAD DATA INFILE 'f.csv' INTO TABLE t",
		"LOAD DATA LOCAL INFILE 'f.csv' REPLACE INTO TABLE t",
		"LOAD DATA INFILE 'f' IGNORE INTO TABLE t FIELDS TERMINATED BY ',' LINES TERMINATED BY '\\n' IGNORE 1 LINES (a, b)",
		// FLUSH
		"FLUSH TABLES",
		"FLUSH TABLES WITH READ LOCK",
		"FLUSH TABLES t1, t2 FOR EXPORT",
		"FLUSH PRIVILEGES",
		"FLUSH LOCAL TABLES",
		"FLUSH STATUS",
		// "data" must remain usable as an identifier
		"SELECT data FROM t",
		// || / &&
		"SELECT a || b FROM t",
		"SELECT a && b FROM t",
		// XA
		"XA START 'x1'",
		"XA START 'x1' JOIN",
		"XA BEGIN 'x1', 'b1', 1 RESUME",
		"XA END 'x1' SUSPEND",
		"XA END 'x1' SUSPEND FOR MIGRATE",
		"XA PREPARE 'x1'",
		"XA COMMIT 'x1' ONE PHASE",
		"XA ROLLBACK 'x1'",
		"XA RECOVER",
		"XA RECOVER CONVERT XID",
		// PURGE
		"PURGE TABLE t",
		"PURGE INDEX i ON t",
		"PURGE RECYCLEBIN",
		// MEMBER OF / MATCH AGAINST
		"SELECT a FROM t WHERE a MEMBER OF (1, 2)",
		"SELECT MATCH (a, b) AGAINST ('x') FROM t",
		"SELECT MATCH (a) AGAINST ('x' IN BOOLEAN MODE) FROM t",
		"SELECT MATCH (a) AGAINST ('x' WITH QUERY EXPANSION) FROM t",
		// GROUPING SETS / CUBE / ROLLUP
		"SELECT a FROM t GROUP BY GROUPING SETS ((a), (b))",
		"SELECT a FROM t GROUP BY CUBE (a, b)",
		"SELECT a FROM t GROUP BY ROLLUP (a, b)",
		// ALTER SYSTEM
		"ALTER SYSTEM SET enable_x = true",
		"ALTER SYSTEM SET max_connections = 100",
		"ALTER SYSTEM SET param = 'v'",
		// "system" must remain usable as an identifier
		"SELECT system FROM t",
		// MERGE
		"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN UPDATE SET t.a = u.a",
		"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN DELETE",
		"MERGE INTO t USING u ON t.id = u.id WHEN NOT MATCHED THEN INSERT (a, b) VALUES (u.a, u.b)",
		"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN UPDATE SET a = u.a WHEN NOT MATCHED THEN INSERT VALUES (u.a)",
		"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN UPDATE SET a = 1 DELETE WHERE a > 10",
		// FLASHBACK
		"FLASHBACK TABLE t TO BEFORE DROP",
		"FLASHBACK TABLE t, u TO BEFORE DROP",
		"FLASHBACK TABLE t TO BEFORE DROP RENAME TO t2",
		"FLASHBACK TABLE t TO SNAPSHOT s1",
		// materialized views
		"CREATE MATERIALIZED VIEW mv AS SELECT * FROM t",
		"CREATE MATERIALIZED VIEW mv (a, b) REFRESH FAST ON DEMAND AS SELECT a, b FROM t",
		"CREATE MATERIALIZED VIEW mv REFRESH COMPLETE ON COMMIT AS SELECT 1",
		"CREATE MATERIALIZED VIEW mv REFRESH FORCE AS SELECT 1",
		"DROP MATERIALIZED VIEW mv",
		"DROP MATERIALIZED VIEW IF EXISTS mv1, mv2",
		// subpartitions / partition options
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) SUBPARTITION BY HASH (id) SUBPARTITIONS 4 (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (MAXVALUE))",
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) SUBPARTITION BY KEY (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE t (id INT) PARTITION BY HASH (id) PARTITIONS 4 SUBPARTITION BY HASH (id) SUBPARTITIONS 2",
		"CREATE TABLE t (id INT) PARTITION BY LIST (id) (PARTITION p0 VALUES IN (1, 2) COMMENT=cmt)",
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10) TABLEGROUP = tg1, PARTITION p1 VALUES LESS THAN (MAXVALUE) TABLEGROUP tg1)",
		// SIGNAL / RESIGNAL / GET DIAGNOSTICS
		"SIGNAL SQLSTATE '45000'",
		"SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'oops', MYSQL_ERRNO = 100",
		"RESIGNAL",
		"RESIGNAL SQLSTATE '45000'",
		"RESIGNAL SET MESSAGE_TEXT = 'x'",
		"GET DIAGNOSTICS a = MESSAGE_TEXT",
		"GET CURRENT DIAGNOSTICS a = MESSAGE_TEXT, b = MYSQL_ERRNO",
		"GET STACKED DIAGNOSTICS c = RETURNED_SQLSTATE",
		// cluster admin
		"CREATE TENANT t",
		"CREATE TENANT IF NOT EXISTS t RESOURCE_POOL_LIST = ('p1', 'p2')",
		"CREATE TENANT t PRIMARY_ZONE = 'z1', LOCALITY = 'FULL{1}@z1', CHARSET = utf8mb4, COMMENT = 'x'",
		"DROP TENANT t",
		"DROP TENANT IF EXISTS t PURGE",
		"CREATE RESOURCE UNIT u MAX_CPU = 2, MEMORY_SIZE = '2G'",
		"DROP RESOURCE UNIT u",
		"CREATE RESOURCE POOL p UNIT = 'u', UNIT_NUM = 1, ZONE_LIST = ('z1', 'z2')",
		"DROP RESOURCE POOL IF EXISTS p",
		"CREATE TABLEGROUP tg",
		"CREATE TABLEGROUP IF NOT EXISTS tg",
		"DROP TABLEGROUP tg",
		"CREATE TABLESPACE ts",
		"DROP TABLESPACE IF EXISTS ts",
		"ALTER TENANT t RESOURCE_POOL_LIST = ('p1')",
		// ALTER 族 + 用户维护
		"ALTER VIEW v AS SELECT 1",
		"ALTER VIEW v (a, b) AS SELECT a, b FROM t",
		"ALTER SEQUENCE s START WITH 10 INCREMENT BY 5 CYCLE",
		"ALTER TABLEGROUP tg",
		"ALTER RESOURCE UNIT u MAX_CPU = 4",
		"ALTER RESOURCE POOL p UNIT = 'u'",
		"RENAME USER 'u'@'%' TO 'v'@'%'",
		"LOCK USER 'u'@'%'",
		"UNLOCK USER 'u'@'%', 'v'@'%'",
		// ALTER MVIEW / MLOG / OUTLINE / TABLESPACE / ROTATE MASTER KEY
		"ALTER MATERIALIZED VIEW mv REFRESH COMPLETE ON DEMAND",
		"CREATE MATERIALIZED VIEW LOG ON t",
		"CREATE MATERIALIZED VIEW LOG ON t WITH PRIMARY KEY",
		"CREATE MATERIALIZED VIEW LOG ON t WITH PRIMARY KEY INCLUDING NEW VALUES",
		"ALTER MATERIALIZED VIEW LOG ON t",
		"DROP MATERIALIZED VIEW LOG IF EXISTS ON t",
		"CREATE OUTLINE o ON SELECT * FROM t",
		"CREATE OUTLINE o ON 'SELECT 1'",
		"ALTER OUTLINE o ON SELECT 1",
		"DROP OUTLINE IF EXISTS o",
		"ALTER TABLESPACE ts ENGINE=InnoDB",
		"ROTATE MASTER KEY",
		// tenant advanced
		"ALTER SYSTEM ACTIVATE STANDBY TENANT t",
		"ALTER SYSTEM SWITCHOVER TO PRIMARY TENANT t",
		"ALTER SYSTEM SWITCHOVER TO STANDBY TENANT t",
		"ALTER SYSTEM RECOVER STANDBY TENANT t UNTIL UNLIMITED",
		"ALTER SYSTEM RECOVER STANDBY TENANT t UNTIL TIME '2020-01-01'",
		"ALTER SYSTEM RECOVER STANDBY TENANT t UNTIL SCN 100",
		"ALTER SYSTEM RECOVER STANDBY CANCEL",
		"ALTER SYSTEM CREATE LS LS_NUM = 1 TENANT t",
		"ALTER SYSTEM MODIFY LS RESOURCE_POOL_LIST = ('p')",
		"ALTER SYSTEM ROTATE INNODB MASTER KEY",
		"CREATE STANDBY TENANT t LOG_RESTORE_SOURCE = 'file://x'",
		"CREATE TENANT IF NOT EXISTS t FROM src USING SNAPSHOT snap WITH RESOURCE_POOL_LIST = ('p')",
		"CREATE RESTORE POINT rp1",
		"DROP RESTORE POINT rp1",
		"CREATE SNAPSHOT snap1 FOR TENANT t",
		"DROP SNAPSHOT snap1",
		// audit / dump / service / location / flashback standby log
		"AUDIT SELECT, INSERT ON db1.*",
		"AUDIT ALL ON *.*",
		"NOAUDIT SELECT ON db1.*",
		"DUMP ENTITY ALL",
		"DUMP CHUNK ALL",
		"DUMP ENTITY name = 't', SLOT_IDX = 1",
		"DUMP CHUNK TENANT_ID = 1, CTX_ID = 'x'",
		"ALTER SYSTEM CREATE SERVICE s1",
		"ALTER SYSTEM DELETE SERVICE s1 TENANT t",
		"CREATE LOCATION loc URL 'oss://x'",
		"CREATE OR REPLACE LOCATION loc URL = 'oss://y'",
		"REMOVE FILES IN LOCATION @path",
		"ALTER SYSTEM FLASHBACK STANDBY LOG TO SCN = 100 TENANT t",
		// routine load
		"CREATE ROUTINE LOAD job1 ON t FROM STREAM ()",
		"CREATE ROUTINE LOAD job1 ON t COLUMNS (a, b) FROM STREAM ()",
		"CREATE ROUTINE LOAD job1 ON t FROM STREAM ('k' = 'v')",
		"PAUSE ROUTINE LOAD FOR job1",
		"RESUME ROUTINE LOAD FOR job1",
		"STOP ROUTINE LOAD FOR job1",
		// ccl / sensitive rule / dblink
		"CREATE CONCURRENT_LIMITING_RULE r1 ON db1.t TO 'u'@'%' FOR CCL_DML_OPTION",
		"ALTER CONCURRENT_LIMITING_RULE r1",
		"DROP CONCURRENT_LIMITING_RULE IF EXISTS r1",
		"CREATE SENSITIVE RULE s1 ON a, b USING FULL_MASK",
		"ALTER SENSITIVE RULE s1 ON a USING MASK_ALL",
		"DROP SENSITIVE RULE s1",
		"CREATE DATABASE LINK dl CONNECT TO 'u' DATABASE 'db' IDENTIFIED BY 'pwd'",
		"CREATE DATABASE LINK dl CONNECT TO 'u' DATABASE 'db' IDENTIFIED BY 'pwd' HOST 'h' PORT 1521 CLUSTER c",
		"DROP DATABASE LINK IF EXISTS dl",
		// LATERAL / QUALIFY
		"SELECT * FROM t, LATERAL (SELECT 1) x",
		"SELECT a FROM t QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1",
		"SELECT a FROM t GROUP BY a QUALIFY COUNT(*) > 1",
		// direct load hint / auto partition / plugins / handler / analyse
		"LOAD DATA /*+ DIRECT(true, 2, 'load') */ INFILE 'f' INTO TABLE t",
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) PARTITIONS AUTO",
		"CREATE TABLE t (id INT) PARTITION BY RANGE (id) PARTITION SIZE 100 PARTITIONS AUTO",
		"INSTALL PLUGIN p1 SONAME 'p1.so'",
		"UNINSTALL PLUGIN p1",
		"HANDLER t OPEN AS a",
		"HANDLER t CLOSE",
		"HANDLER t READ FIRST (a, b) WHERE a > 1 LIMIT 10",
		"ANALYZE TABLE t METHOD_OPT 'FOR COLUMNS SIZE 1'",
		"SELECT a FROM t PROCEDURE ANALYSE(2, 10)",
		"SELECT a FROM t PROCEDURE ANALYSE()",
		// ALTER TABLE partition ops & column details
		"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (20))",
		"ALTER TABLE t DROP PARTITION p0, p1",
		"ALTER TABLE t DROP PARTITION ALL",
		"ALTER TABLE t TRUNCATE PARTITION p0",
		"ALTER TABLE t COALESCE PARTITION 2",
		"ALTER TABLE t EXCHANGE PARTITION p0 WITH TABLE t2",
		"ALTER TABLE t REORGANIZE PARTITION p0 INTO (PARTITION p0 VALUES LESS THAN (10))",
		"ALTER TABLE t REMOVE PARTITIONING",
		"ALTER TABLE t ADD COLUMN c INT FIRST",
		"ALTER TABLE t ADD COLUMN c INT AFTER b",
		"ALTER TABLE t MODIFY COLUMN c BIGINT AFTER a",
		"ALTER TABLE t ADD (c INT, d INT)",
		"ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4",
		"ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_bin",
		"ALTER TABLE t AUTO_INCREMENT = 100",
		// SELECT INTO OUTFILE / EXPLAIN / DIAGNOSTICS / FLUSH
		"SELECT a FROM t INTO OUTFILE '/tmp/x'",
		"SELECT a FROM t INTO OUTFILE '/tmp/x' FIELDS TERMINATED BY ','",
		"SELECT a FROM t INTO DUMPFILE '/tmp/x'",
		"EXPLAIN ANALYZE SELECT 1",
		"EXPLAIN FOR CONNECTION 7 SELECT 1",
		"GET DIAGNOSTICS CONDITION 1 a = MESSAGE_TEXT",
		"FLUSH LOGS",
		"FLUSH HOSTS",
		// catalog / location alter-drop
		"CREATE EXTERNAL CATALOG c PROPERTIES ('k' = 'v')",
		"CREATE EXTERNAL CATALOG IF NOT EXISTS c PROPERTIES ('endpoint' = 'oss://x')",
		"DROP CATALOG IF EXISTS c",
		"SET CATALOG c",
		"ALTER LOCATION loc SET URL 'oss://new'",
		"DROP LOCATION loc",
		"CREATE USER 'u'@'%' IDENTIFIED WITH mysql_native_password",
		"CREATE USER 'u'@'%' IDENTIFIED WITH mysql_native_password BY 'pwd'",
		// PL / stored routines
		"CREATE PROCEDURE p() BEGIN SET a = 1; END",
		"CREATE PROCEDURE p(IN x INT, OUT y VARCHAR(10)) BEGIN SET y = x; END",
		"CREATE PROCEDURE p() BEGIN DECLARE v INT DEFAULT 0; SET v = v + 1; SELECT v; END",
		"CREATE FUNCTION f(x INT) RETURNS INT BEGIN RETURN x * 2; END",
		"CREATE OR REPLACE PROCEDURE p() BEGIN IF a > 1 THEN SET b = 1; ELSEIF a = 1 THEN SET b = 2; ELSE SET b = 0; END IF; END",
		"CREATE PROCEDURE p() BEGIN WHILE a < 10 DO SET a = a + 1; END WHILE; END",
		"CREATE PROCEDURE p() BEGIN REPEAT SET a = a + 1; UNTIL a >= 10 END REPEAT; END",
		"CREATE PROCEDURE p() BEGIN LOOP SET a = a + 1; END LOOP; END",
		"CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT a FROM t; OPEN c; FETCH c INTO a, b; CLOSE c; END",
		"CREATE PROCEDURE p() BEGIN SELECT * FROM t; END",
		"CREATE PROCEDURE p() BEGIN DECLARE v INT; BEGIN SET v = 1; END; END",
		"CREATE PROCEDURE p() BEGIN SET @x = 1; END",
		"DROP PROCEDURE p",
		"DROP FUNCTION IF EXISTS f",
		"CALL p(1, 2)",
		"CALL p()",
		"ALTER SYSTEM TRANSFER PARTITION TABLE_ID = 1100611139456002, OBJECT_ID = 0 TO LS 1 TENANT t",
		"ALTER SYSTEM CANCEL TRANSFER PARTITION TABLE_ID = 1100611139456002, OBJECT_ID = 0",
		"ALTER SYSTEM CANCEL TRANSFER PARTITION ALL",
		// PL labels / handlers
		"CREATE PROCEDURE p() BEGIN lbl: BEGIN SET a = 1; LEAVE lbl; END; END",
		"CREATE PROCEDURE p() BEGIN lab: LOOP SET a = a + 1; ITERATE lab; END LOOP; END",
		"CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1; END",
		"CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '45000' ROLLBACK; END",
		"CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR SQLWARNING, SQLEXCEPTION SET a = 1; END",
		// ANALYZE is non-reserved: statement, EXPLAIN, and identifier uses coexist
		"ANALYZE TABLE t",
		"EXPLAIN ANALYZE SELECT 1",
		"USE analyze",
		"SELECT analyze FROM t",
		"CREATE TABLE analyze (id INT)",
		// engineered robustness fixes
		"SELECT a FROM t INTO OUTFILE '/tmp/f' FIELDS TERMINATED BY ',' LINES TERMINATED BY '\\n'",
		"CREATE TABLE t (id INT, CONSTRAINT fk FOREIGN KEY (id) REFERENCES u (id))",
		"CREATE TABLE t (id INT, KEY k (id), INDEX i (name), UNIQUE uq (a), UNIQUE KEY uk (b))",
		"SELECT CONVERT(b, CHAR(4)) FROM t",
		"USE analyze",
		"CREATE TABLE t (id INT, CHECK (id > 0))",
		// A: locking & show user & exchange validation
		"SELECT * FROM t FOR UPDATE NOWAIT",
		"SELECT * FROM t FOR UPDATE SKIP LOCKED",
		"SELECT * FROM t LOCK IN SHARE MODE",
		"SHOW GRANTS FOR 'u'@'%'",
		"SHOW GRANTS FOR 'u'@'localhost'",
		"SHOW CREATE USER 'u'@'%'",
		"ALTER TABLE t EXCHANGE PARTITION p0 WITH TABLE t2",
		"ALTER TABLE t EXCHANGE PARTITION p0 WITH TABLE t2 WITH VALIDATION",
		"ALTER TABLE t EXCHANGE PARTITION p0 WITH TABLE t2 WITHOUT VALIDATION",
		// 完善：SELECT INTO 变量 / PL CASE / DECLARE CONDITION / SET CHARACTER SET
		"SELECT a, b FROM t INTO @x, @y",
		"SELECT a FROM t INTO @z",
		"CREATE PROCEDURE p() BEGIN CASE WHEN a > 0 THEN SET b = 1; WHEN a = 0 THEN SET b = 0; ELSE SET b = -1; END CASE; END",
		"CREATE PROCEDURE p() BEGIN CASE x WHEN 1 THEN SET a = 1; ELSE SET a = 2; END CASE; END",
		"CREATE PROCEDURE p() BEGIN DECLARE mycond CONDITION FOR SQLSTATE '45000'; END",
		"CREATE PROCEDURE p() BEGIN DECLARE mycond CONDITION FOR 1062; END",
		"SET CHARACTER SET utf8mb4",
		"SET @x := 1",
		"SHOW VARIABLES WHERE Variable_name = 'x'",
		"SHOW STATUS LIKE 'Ssl%'",
		"SHOW TABLES FROM mydb WHERE Tables_in_mydb LIKE 'a%'",
		"SHOW COLUMNS FROM t WHERE Field LIKE 'a%'",
		"SHOW PROFILES",
		"SHOW PROFILE",
		"CREATE PROCEDURE p() BEGIN SET @x := 1; SET y := 2; END",
		"SHOW CHARSET",
		"SHOW CHARSET LIKE 'utf8%'",
		"SHOW STORAGE ENGINES",
		// 常用 SHOW 变体
		"SHOW ENGINE InnoDB STATUS",
		"SHOW ENGINE InnoDB MUTEX",
		"SHOW MASTER STATUS",
		"SHOW FUNCTION STATUS LIKE 'f%'",
		"SHOW PROCEDURE STATUS",
		"SHOW EVENTS",
		"SHOW EVENTS FROM mydb",
		"SHOW PROFILE FOR QUERY 1",
		// ON UPDATE CURRENT_TIMESTAMP 与布尔表选项（此前不支持）
		"CREATE TABLE t (id bigint NOT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id))",
		"CREATE TABLE t (updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP())",
		"CREATE TABLE t (updated_at timestamp(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3))",
		"CREATE TABLE t (updated_at timestamp DEFAULT NOW() ON UPDATE NOW())",
		"CREATE TABLE t (id bigint NOT NULL, PRIMARY KEY (id)) USE_BLOOM_FILTER = FALSE",
		"CREATE TABLE t (id bigint NOT NULL, PRIMARY KEY (id)) USE_BLOOM_FILTER = TRUE",
		"CREATE TABLE t (id bigint NOT NULL, PRIMARY KEY (id)) COMPRESSION = 'zstd' USE_BLOOM_FILTER = FALSE",
		// 常用 DDL/DML 补齐（本轮）
		"CREATE TABLE t (id INT, name VARCHAR(50), CONSTRAINT fk FOREIGN KEY (id) REFERENCES p (pid) ON DELETE CASCADE ON UPDATE RESTRICT)",
		"CREATE TABLE t (id INT, name VARCHAR(50), CONSTRAINT fk FOREIGN KEY (id) REFERENCES p (pid) ON DELETE SET NULL)",
		"CREATE TABLE t (id INT) COMMENT = 'table comment'",
		"CREATE TABLE t (id INT) ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd' REPLICA_NUM = 3",
		"CREATE TABLE t (id INT, PRIMARY KEY (id)) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (100), PARTITION p1 VALUES LESS THAN MAXVALUE)",
		"CREATE TABLE t (id INT UNSIGNED ZEROFILL)",
		"CREATE TABLE t (id INT) PARTITION BY LIST (id) SUBPARTITION BY HASH (id) SUBPARTITIONS 2",
		"CREATE TABLE t2 AS SELECT a, b FROM t1",
		"CREATE TEMPORARY TABLE t (id INT)",
		"ALTER TABLE t ADD COLUMN (c1 INT, c2 VARCHAR(10))",
		"ALTER TABLE t ADD INDEX idx (a) LOCAL",
		"ALTER TABLE t ADD INDEX idx (a) KEY_BLOCK_SIZE = 16 COMMENT 'idx'",
		"ALTER TABLE t ADD UNIQUE KEY uq (a, b)",
		"ALTER TABLE t ALTER INDEX idx INVISIBLE",
		"ALTER TABLE t ALTER INDEX idx VISIBLE",
		"ALTER TABLE t ENABLE KEYS",
		"ALTER TABLE t DISABLE KEYS",
		"ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id) ON DELETE RESTRICT ON UPDATE CASCADE",
		"UPDATE LOW_PRIORITY t SET a = 1",
		"UPDATE IGNORE t SET a = 1 WHERE b = 2",
		"UPDATE t SET a = 1 WHERE id = 2 RETURNING id, a",
		"DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id WHERE t2.x = 1",
		"DELETE t1, t2 FROM t1 JOIN t2 ON t1.id = t2.id",
		"DELETE FROM t1, t2 USING t1 JOIN t2 ON t1.id = t2.id WHERE t1.x = 1",
		"DELETE FROM t WHERE a = 1 RETURNING a",
		"DELETE QUICK FROM t WHERE a = 1",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			stmts := parseAll(t, sql)
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			s1 := formatter.Format(stmts[0])
			stmts2 := parseAll(t, s1)
			if len(stmts2) != 1 {
				t.Fatalf("re-parse gave %d statements\n  s1: %s", len(stmts2), s1)
			}
			s2 := formatter.Format(stmts2[0])
			if s1 != s2 {
				t.Fatalf("not idempotent\n  s1: %s\n  s2: %s", s1, s2)
			}
		})
	}
}

func TestGating(t *testing.T) {
	sql := "WITH cte AS (SELECT 1) SELECT * FROM cte"
	if _, err := Parse(mustLex(t, sql), Options{Version: version.MustParse("4.0.0")}); err != nil {
		t.Fatalf("4.0 should accept WITH: %v", err)
	}
	if _, err := Parse(mustLex(t, sql), Options{Version: version.MustParse("3.9.0")}); err == nil {
		t.Fatal("pre-4.0 should reject WITH")
	}
}

func TestFeatureGating(t *testing.T) {
	mk := func(v string) Options { return Options{Version: version.MustParse(v)} }
	cases := []struct {
		sql     string
		version string
		ok      bool
	}{
		// WITH CTE: since 4.0
		{"WITH c AS (SELECT 1) SELECT * FROM c", "4.0.0", true},
		{"WITH c AS (SELECT 1) SELECT * FROM c", "3.9.0", false},
		// EXPLAIN EXTENDED: 4.1 <= v < 4.3
		{"EXPLAIN EXTENDED SELECT 1", "4.1.0", true},
		{"EXPLAIN EXTENDED SELECT 1", "4.2.0", true},
		{"EXPLAIN EXTENDED SELECT 1", "4.0.0", false},
		{"EXPLAIN EXTENDED SELECT 1", "4.3.0", false},
		// named WINDOW clause: since 4.2
		{"SELECT a FROM t WINDOW w AS (PARTITION BY a)", "4.2.0", true},
		{"SELECT a FROM t WINDOW w AS (PARTITION BY a)", "4.1.0", false},
		// window frame ROWS/RANGE: since 4.0
		{"SELECT ROW_NUMBER() OVER (ORDER BY a ROWS UNBOUNDED PRECEDING) FROM t", "4.0.0", true},
		{"SELECT ROW_NUMBER() OVER (ORDER BY a ROWS UNBOUNDED PRECEDING) FROM t", "3.9.0", false},
		// unversioned target accepts everything
		{"EXPLAIN EXTENDED SELECT 1", "", true},
		// post-4.2.x features: absent in 4.2.5, present later
		{"CREATE MATERIALIZED VIEW mv AS SELECT 1", "4.2.5", false},
		{"CREATE MATERIALIZED VIEW mv AS SELECT 1", "4.3.0", true},
		{"CREATE MATERIALIZED VIEW LOG ON t", "4.2.5", false},
		{"CREATE MATERIALIZED VIEW LOG ON t", "4.3.0", true},
		{"CREATE EXTERNAL CATALOG c PROPERTIES ('k' = 'v')", "4.2.5", false},
		{"CREATE EXTERNAL CATALOG c PROPERTIES ('k' = 'v')", "4.3.0", false},
		{"CREATE EXTERNAL CATALOG c PROPERTIES ('k' = 'v')", "4.4.0", true},
		{"CREATE ROUTINE LOAD job ON t FROM STREAM ()", "4.2.5", false},
		{"CREATE ROUTINE LOAD job ON t FROM STREAM ()", "5.0.0", true},
		{"CREATE CONCURRENT_LIMITING_RULE r1 ON db1.t TO 'u'@'%' FOR CCL_DML_OPTION", "4.2.5", false},
		{"CREATE CONCURRENT_LIMITING_RULE r1 ON db1.t TO 'u'@'%' FOR CCL_DML_OPTION", "4.3.0", false},
		{"CREATE CONCURRENT_LIMITING_RULE r1 ON db1.t TO 'u'@'%' FOR CCL_DML_OPTION", "4.4.0", true},
		{"CREATE SENSITIVE RULE s1 ON a USING MASK", "4.3.0", false},
		{"CREATE SENSITIVE RULE s1 ON a USING MASK", "4.4.0", true},
		{"CREATE SNAPSHOT snap FOR TENANT t", "4.2.5", false},
		{"CREATE SNAPSHOT snap FOR TENANT t", "4.3.0", true},
		{"ALTER SYSTEM FLASHBACK STANDBY LOG TO SCN = 1", "4.3.0", false},
		{"ALTER SYSTEM FLASHBACK STANDBY LOG TO SCN = 1", "4.4.0", true},
		{"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN DELETE", "4.3.0", false},
		{"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN DELETE", "4.4.0", true},
		{"CREATE TABLE t (id INT) PARTITION BY RANGE (id) PARTITIONS AUTO", "4.0.0", true},
		{"SELECT a FROM t GROUP BY GROUPING SETS ((a), (b))", "4.2.5", false},
		{"SELECT a FROM t GROUP BY GROUPING SETS ((a), (b))", "4.3.0", false},
		{"SELECT a FROM t GROUP BY GROUPING SETS ((a), (b))", "4.4.0", true},
		{"ALTER TABLE t EXCHANGE PARTITION p0 WITH TABLE t2 WITHOUT VALIDATION", "4.2.5", false},
		{"ALTER TABLE t EXCHANGE PARTITION p0 WITH TABLE t2 WITHOUT VALIDATION", "4.4.0", true},
		// 4.2.x minor-series boundaries (calibrated from official tags)
		{"ALTER SYSTEM TRANSFER PARTITION TABLE_ID = 1, OBJECT_ID = 0 TO LS 1", "4.2.1", false},
		{"ALTER SYSTEM TRANSFER PARTITION TABLE_ID = 1, OBJECT_ID = 0 TO LS 1", "4.2.2", true},
		{"CREATE ROLE r1", "4.2.1", false},
		{"CREATE ROLE r1", "4.2.4", true},
		{"ALTER SYSTEM CREATE SERVICE s", "4.2.1", false},
		{"ALTER SYSTEM CREATE SERVICE s", "4.2.4", true},
		{"XA RECOVER", "4.2.4", false},
		{"XA RECOVER", "4.2.5", true},
	}
	for _, c := range cases {
		var opts Options
		if c.version != "" {
			opts = mk(c.version)
		}
		_, err := Parse(mustLex(t, c.sql), opts)
		if c.ok && err != nil {
			t.Errorf("v%s should accept %q: %v", c.version, c.sql, err)
		}
		if !c.ok && err == nil {
			t.Errorf("v%s should reject %q", c.version, c.sql)
		}
	}
}

func TestStatementSpan(t *testing.T) {
	stmts := parseAll(t, "SELECT 1 FROM t")
	if len(stmts) != 1 {
		t.Fatalf("want 1 stmt")
	}
	sp := stmts[0].Span()
	if sp.Start.Line != 1 || sp.Start.Col != 1 {
		t.Fatalf("start = %+v, want 1:1", sp.Start)
	}
	if sp.End.Line != 1 || sp.End.Col <= sp.Start.Col {
		t.Fatalf("end = %+v should follow start", sp.End)
	}
}

func TestMulti(t *testing.T) {
	stmts := parseAll(t, "SELECT 1; SELECT 2; INSERT INTO t VALUES (3);")
	if len(stmts) != 3 {
		t.Fatalf("expected 3, got %d", len(stmts))
	}
}

func mustLex(t *testing.T, sql string) []token.Token {
	t.Helper()
	toks, err := lexer.LexAll(sql, lexOpts)
	if err != nil {
		t.Fatal(err)
	}
	return toks
}

func TestDoStmt(t *testing.T) {
	for _, s := range []string{"DO 1 + 1", "DO SLEEP(1)", "DO (SELECT 1)"} {
		st := parseAll(t, s)
		if st[0].(*ast.DoStmt) == nil {
			t.Fatalf("not a DoStmt: %s", s)
		}
	}
}

func TestColumnOnUpdate(t *testing.T) {
	// ON UPDATE CURRENT_TIMESTAMP 必须保留在 AST 与格式化输出中（此前不支持）。
	src := "CREATE TABLE t (id bigint NOT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id))"
	st := parseAll(t, src)
	ct := st[0].(*ast.CreateTableStmt)
	if len(ct.Columns) != 2 || ct.Columns[1].OnUpdate == nil {
		t.Fatalf("OnUpdate not preserved: %+v", ct.Columns)
	}
	if got := formatter.Format(st[0]); got != "CREATE TABLE t ( id bigint NOT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id))" {
		t.Fatalf("format mismatch: %s", got)
	}

	// 精度形态 CURRENT_TIMESTAMP(6)
	srcN := "CREATE TABLE t (updated_at timestamp(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6))"
	stN := parseAll(t, srcN)
	if stN[0].(*ast.CreateTableStmt).Columns[0].OnUpdate == nil {
		t.Fatal("ON UPDATE CURRENT_TIMESTAMP(6) not preserved")
	}
}

func TestTableOptionBool(t *testing.T) {
	// 布尔表选项值 TRUE/FALSE（此前 BOOL_VALUE 不在 table_option_value）。
	for _, v := range []string{"FALSE", "TRUE", "true", "false"} {
		src := "CREATE TABLE t (id bigint NOT NULL, PRIMARY KEY (id)) USE_BLOOM_FILTER = " + v
		st := parseAll(t, src)
		ct := st[0].(*ast.CreateTableStmt)
		if len(ct.Options) != 1 || ct.Options[0].Name != "use_bloom_filter" {
			t.Fatalf("option not parsed: %s", src)
		}
		// 归一化为小写布尔值，格式化后可再解析
		got := formatter.Format(st[0])
		if ct.Options[0].Value != "true" && ct.Options[0].Value != "false" {
			t.Fatalf("bool option value = %q", ct.Options[0].Value)
		}
		parseAll(t, got) // 格式化输出必须可再解析
	}
}

func TestDDLCommonExtras(t *testing.T) {
	t.Run("fk-reference-actions", func(t *testing.T) {
		st := parseAll(t, "CREATE TABLE t (id INT, CONSTRAINT fk FOREIGN KEY (id) REFERENCES p (pid) ON DELETE CASCADE ON UPDATE RESTRICT)")
		ct := st[0].(*ast.CreateTableStmt)
		if len(ct.Constraints) != 1 {
			t.Fatalf("constraints = %d", len(ct.Constraints))
		}
		c := ct.Constraints[0]
		if c.Kind != "FOREIGN KEY" || c.RefTable == nil || len(c.RefTable.Parts) == 0 ||
			len(c.RefColumns) != 1 || c.OnDelete != "CASCADE" || c.OnUpdate != "RESTRICT" {
			t.Fatalf("fk not preserved: %+v", c)
		}
	})

	t.Run("zerofill", func(t *testing.T) {
		st := parseAll(t, "CREATE TABLE t (id INT UNSIGNED ZEROFILL)")
		ct := st[0].(*ast.CreateTableStmt)
		if !ct.Columns[0].Zerofill || !ct.Columns[0].Unsigned {
			t.Fatalf("zerofill/unsigned not preserved: %+v", ct.Columns[0])
		}
	})

	t.Run("maxvalue-no-paren", func(t *testing.T) {
		st := parseAll(t, "CREATE TABLE t (id INT, PRIMARY KEY (id)) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (100), PARTITION p1 VALUES LESS THAN MAXVALUE)")
		ct := st[0].(*ast.CreateTableStmt)
		pds := ct.Partition.Partitions
		if len(pds) != 2 || !pds[1].MaxValue || pds[1].ValueKind != "LESS THAN" {
			t.Fatalf("maxvalue partition not preserved: %+v", pds)
		}
	})

	t.Run("comment-table-option-quoted", func(t *testing.T) {
		st := parseAll(t, "CREATE TABLE t (id INT) COMMENT = 'table comment'")
		ct := st[0].(*ast.CreateTableStmt)
		if len(ct.Options) != 1 || ct.Options[0].Name != "comment" {
			t.Fatalf("option missing: %+v", ct.Options)
		}
		if got := formatter.Format(st[0]); got != "CREATE TABLE t ( id INT) COMMENT='table comment'" {
			t.Fatalf("comment option not quoted: %s", got)
		}
	})
}

func TestAlterTableExtras(t *testing.T) {
	cases := []struct {
		src  string
		kind ast.AlterActionKind
	}{
		{"ALTER TABLE t ADD COLUMN (c1 INT, c2 VARCHAR(10))", ast.AlterAddColumns},
		{"ALTER TABLE t ADD INDEX idx (a) LOCAL", ast.AlterAddIndex},
		{"ALTER TABLE t ALTER INDEX idx INVISIBLE", ast.AlterInvisibleIndex},
		{"ALTER TABLE t ALTER INDEX idx VISIBLE", ast.AlterVisibleIndex},
		{"ALTER TABLE t ENABLE KEYS", ast.AlterEnableKeys},
		{"ALTER TABLE t DISABLE KEYS", ast.AlterDisableKeys},
	}
	for _, tc := range cases {
		st := parseAll(t, tc.src)
		at := st[0].(*ast.AlterTableStmt)
		if len(at.Actions) != 1 || at.Actions[0].Kind != tc.kind {
			t.Fatalf("%s: kind = %q, want %q", tc.src, at.Actions[0].Kind, tc.kind)
		}
	}
	if st := parseAll(t, "ALTER TABLE t ADD INDEX idx (a) LOCAL"); st[0].(*ast.AlterTableStmt).Actions[0].IndexOptions[0] != "LOCAL" {
		t.Fatal("index option LOCAL not preserved")
	}
}

func TestIndexOptionsAndColumnGroup(t *testing.T) {
	// CREATE TABLE 约束级索引选项 + WITH COLUMN GROUP
	st := parseAll(t, "CREATE TABLE t (a INT, KEY k (a) LOCAL BLOCK_SIZE = 16384 WITH COLUMN GROUP (ALL COLUMNS))")
	ct := st[0].(*ast.CreateTableStmt)
	if len(ct.Constraints) != 1 || len(ct.Constraints[0].Options) != 2 || ct.Constraints[0].ColumnGroup == "" {
		t.Fatalf("constraint options: %+v", ct.Constraints[0])
	}
	// 命名 UNIQUE / UNIQUE KEY / INDEX 组合
	st = parseAll(t, "CREATE TABLE t (a INT, UNIQUE uq (a), UNIQUE KEY uk (b) LOCAL, INDEX i (c) WITH COLUMN GROUP (cg (c)))")
	if len(st[0].(*ast.CreateTableStmt).Constraints) != 3 {
		t.Fatal("constraint kinds")
	}
	// CREATE INDEX / ALTER ADD INDEX 的 WITH COLUMN GROUP
	st = parseAll(t, "CREATE INDEX idx ON t (a) LOCAL WITH COLUMN GROUP (EACH COLUMN)")
	if st[0].(*ast.CreateIndexStmt).ColumnGroup == "" {
		t.Fatal("create index colgroup lost")
	}
	st = parseAll(t, "ALTER TABLE t ADD INDEX idx (a) LOCAL WITH COLUMN GROUP (all columns)")
	if st[0].(*ast.AlterTableStmt).Actions[0].ColumnGroup == "" {
		t.Fatal("alter index colgroup lost")
	}
	// ALTER CHANGE COLUMN + ON UPDATE CURRENT_TIMESTAMP
	st = parseAll(t, "ALTER TABLE t CHANGE COLUMN old new timestamp DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()")
	col := st[0].(*ast.AlterTableStmt).Actions[0].Column
	if col == nil || col.OnUpdate == nil {
		t.Fatal("change column on update lost")
	}
	// 索引选项 STRONG/WEAK 值
	if st := parseAll(t, "CREATE TABLE t (id INT) read_consistency = STRONG"); st[0].(*ast.CreateTableStmt).Options[0].Value != "STRONG" {
		t.Fatal("strong value lost")
	}
}

func TestIndexStoringWithAndCheck(t *testing.T) {
	// STORING(...) 覆盖索引列（建索引/建表约束/ALTER 三处）
	st := parseAll(t, "CREATE INDEX idx ON t (a) STORING (b, c)")
	if len(st[0].(*ast.CreateIndexStmt).Storing) != 2 {
		t.Fatal("create index storing lost")
	}
	st = parseAll(t, "CREATE TABLE t (a INT, b INT, KEY k (a) STORING (b))")
	if len(st[0].(*ast.CreateTableStmt).Constraints[0].Storing) != 1 {
		t.Fatal("constraint storing lost")
	}
	st = parseAll(t, "ALTER TABLE t ADD INDEX idx (a) STORING (b, c)")
	if len(st[0].(*ast.AlterTableStmt).Actions[0].Storing) != 2 {
		t.Fatal("alter storing lost")
	}
	// WITH (...): 索引参数
	st = parseAll(t, "CREATE UNIQUE INDEX idx ON t (a) STORING (b) WITH (cf = 'default')")
	if len(st[0].(*ast.CreateIndexStmt).WithOpts) != 1 {
		t.Fatal("create index with lost")
	}
	// USING HASH
	if st := parseAll(t, "CREATE INDEX idx USING HASH ON t (a)"); st[0].(*ast.CreateIndexStmt).Method != "HASH" {
		t.Fatal("using hash lost")
	}
	// 列级 CHECK
	st = parseAll(t, "CREATE TABLE t (a INT CHECK (a > 0))")
	if st[0].(*ast.CreateTableStmt).Columns[0].Check == nil {
		t.Fatal("column check lost")
	}
	// CHANGE COLUMN ... FIRST/AFTER
	st = parseAll(t, "ALTER TABLE t CHANGE COLUMN old new INT AFTER b")
	a := st[0].(*ast.AlterTableStmt).Actions[0]
	if a.Kind != ast.AlterChangeColumn || a.Position != "AFTER" || a.RefCol != "b" {
		t.Fatalf("change position: %+v", a)
	}
	// DROP CHECK / DROP CONSTRAINT
	for _, tc := range []struct {
		src  string
		kind ast.AlterActionKind
	}{
		{"ALTER TABLE t DROP CHECK c1", ast.AlterDropCheck},
		{"ALTER TABLE t DROP CONSTRAINT c1", ast.AlterDropConstraint},
		{"ALTER TABLE t DROP FOREIGN KEY fk1", ast.AlterDropForeignKey},
	} {
		st = parseAll(t, tc.src)
		if st[0].(*ast.AlterTableStmt).Actions[0].Kind != tc.kind {
			t.Fatalf("%s: kind = %q", tc.src, st[0].(*ast.AlterTableStmt).Actions[0].Kind)
		}
	}

	// 表级 PRIMARY KEY 索引选项
	st = parseAll(t, "CREATE TABLE t (a INT, PRIMARY KEY (a) LOCAL BLOCK_SIZE 8192)")
	pk := st[0].(*ast.CreateTableStmt).Constraints[0]
	if pk.Kind != "PRIMARY KEY" || len(pk.Options) != 2 {
		t.Fatalf("pk options: %+v", pk)
	}

	// HASH/KEY 分区显式命名分区
	st = parseAll(t, "CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0, PARTITION p1)")
	ps := st[0].(*ast.CreateTableStmt).Partition
	if ps.Kind != ast.PartitionHash || len(ps.Partitions) != 2 || ps.Partitions[0].Name != "p0" {
		t.Fatalf("hash named partitions: %+v", ps)
	}
}

func TestDDLDMLExtras3(t *testing.T) {
	// CREATE INDEX 完整形态
	st := parseAll(t, "CREATE INDEX idx USING BTREE ON t (a(10) DESC, b)")
	ci := st[0].(*ast.CreateIndexStmt)
	if ci.Method != "BTREE" || len(ci.Columns) != 2 || !ci.Columns[0].Desc || ci.Columns[0].Length == nil {
		t.Fatalf("create index: %+v", ci)
	}
	if got := formatter.Format(st[0]); got != "CREATE INDEX idx USING BTREE ON t (a(10) DESC, b)" {
		t.Fatalf("create index format: %s", got)
	}
	if st := parseAll(t, "CREATE INDEX IF NOT EXISTS idx ON t (a)"); !st[0].(*ast.CreateIndexStmt).IfNotExists {
		t.Fatal("IF NOT EXISTS not preserved")
	}

	// GRANT 单选表 priv_level（此前只支持 db.t / *.*）
	st = parseAll(t, "GRANT SELECT ON t TO 'u'@'%', 'v'@'host'")
	g := st[0].(*ast.GrantStmt)
	if g.Level != "t" || len(g.Users) != 2 {
		t.Fatalf("grant: %+v", g)
	}

	// ALTER USER / LOCK USER ACCOUNT
	st = parseAll(t, "ALTER USER 'u'@'%' IDENTIFIED BY 'newpw'")
	if st[0].(*ast.AlterUserStmt).Password != "newpw" {
		t.Fatal("alter user password lost")
	}
	st = parseAll(t, "LOCK USER 'u'@'%' ACCOUNT LOCK")
	if st[0].(*ast.LockUserStmt).AccountLock != "LOCK" {
		t.Fatal("ACCOUNT LOCK lost")
	}

	// 序列 NOCACHE / NOCYCLE / NOORDER
	st = parseAll(t, "CREATE SEQUENCE s NOCACHE NOCYCLE NOORDER")
	if got := formatter.Format(st[0]); got != "CREATE SEQUENCE s NOCACHE NOCYCLE NOORDER" {
		t.Fatalf("sequence format: %s", got)
	}

	// 列级约束（PK/UNIQUE/REFERENCES）
	st = parseAll(t, "CREATE TABLE t (a INT PRIMARY KEY, b INT UNIQUE, c INT REFERENCES p (x) ON DELETE CASCADE)")
	cols := st[0].(*ast.CreateTableStmt).Columns
	if !cols[0].PrimaryKey || !cols[1].UniqueKey || cols[2].RefTable == nil || cols[2].OnDelete != "CASCADE" {
		t.Fatalf("column constraints: %+v", cols)
	}

	// SET ROLE ALL / DEFAULT / ALL EXCEPT
	st = parseAll(t, "SET ROLE ALL")
	if st[0].(*ast.SetRoleStmt).Roles[0] != "ALL" {
		t.Fatal("SET ROLE ALL")
	}

	// INSERT 修饰符
	st = parseAll(t, "INSERT LOW_PRIORITY INTO t (a) VALUES (1)")
	if st[0].(*ast.InsertStmt).Mods != "LOW_PRIORITY" {
		t.Fatal("insert mod lost")
	}

	// COMMIT / ROLLBACK AND CHAIN|RELEASE
	st = parseAll(t, "COMMIT WORK AND CHAIN")
	if !st[0].(*ast.CommitStmt).Chain {
		t.Fatal("commit chain")
	}
	st = parseAll(t, "ROLLBACK WORK AND RELEASE")
	if !st[0].(*ast.RollbackStmt).Release {
		t.Fatal("rollback release")
	}

	// WINDOW 命名引用 + 定义 双向
	st = parseAll(t, "SELECT ROW_NUMBER() OVER w FROM t WINDOW w AS (PARTITION BY a)")
	sel := st[0].(*ast.SelectStmt)
	if len(sel.Bodies[0].NamedWindow) != 1 || sel.Bodies[0].NamedWindow[0].Name != "w" {
		t.Fatalf("named window: %+v", sel.Bodies[0].NamedWindow)
	}

	// FOR SHARE / STRAIGHT_JOIN
	st = parseAll(t, "SELECT * FROM t FOR SHARE")
	if got := formatter.Format(st[0]); got != "SELECT * FROM t FOR SHARE" {
		t.Fatalf("for share: %s", got)
	}
	st = parseAll(t, "SELECT * FROM t STRAIGHT_JOIN u ON t.id = u.id")
	if st[0].(*ast.SelectStmt).Bodies[0].From[0].RefKind != ast.RefJoin {
		t.Fatal("straight join not parsed as join")
	}

	// ROW() 构造器
	st = parseAll(t, "INSERT INTO t (a, b) VALUES (ROW(1, 2))")
	if st[0].(*ast.InsertStmt) == nil {
		t.Fatal("row constructor")
	}

	// SHOW ... WHERE（关键字列名）
	st = parseAll(t, "SHOW COLLATION WHERE Charset = 'utf8mb4'")
	if st[0].(*ast.ShowStmt).Where == nil {
		t.Fatal("show where lost")
	}
}

func TestUpdateDeleteExtras(t *testing.T) {
	// UPDATE LOW_PRIORITY + RETURNING
	st := parseAll(t, "UPDATE LOW_PRIORITY t SET a = 1 WHERE id = 2 RETURNING id, a")
	up := st[0].(*ast.UpdateStmt)
	if !up.LowPriority || len(up.Returning) != 2 {
		t.Fatalf("update extras: %+v", up)
	}
	if got := formatter.Format(st[0]); got != "UPDATE LOW_PRIORITY t SET a = 1 WHERE id = 2 RETURNING id, a" {
		t.Fatalf("update format: %s", got)
	}

	// 多表删除（表列表形式）+ RETURNING
	st = parseAll(t, "DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id WHERE t2.x = 1")
	del := st[0].(*ast.DeleteStmt)
	if len(del.Tables) != 1 || len(del.From) != 1 {
		t.Fatalf("delete list form: %+v", del)
	}
	if got := formatter.Format(st[0]); got != "DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id WHERE t2.x = 1" {
		t.Fatalf("delete format: %s", got)
	}

	// DELETE RETURNING
	st = parseAll(t, "DELETE FROM t WHERE a = 1 RETURNING a")
	if len(st[0].(*ast.DeleteStmt).Returning) != 1 {
		t.Fatal("delete returning not preserved")
	}

	// DELETE mods
	st = parseAll(t, "DELETE QUICK FROM t WHERE a = 1")
	if st[0].(*ast.DeleteStmt).Mods != "QUICK" {
		t.Fatal("delete mod not preserved")
	}
}
