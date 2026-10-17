package gparser

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/formatter"
	jsonpkg "github.com/zetatez/ob-sql-parser/json"
	"github.com/zetatez/ob-sql-parser/lexer"
)

// TestOfficialCorpus samples every statement branch of the official OB 4.2.5
// MySQL-mode grammar (92 stmt: alternatives from v4.2.5_CE
// sql_parser_mysql_mode.y) with one representative SQL each, asserting
// parse + canonical round-trip + JSON round-trip.
func TestOfficialCorpus(t *testing.T) {
	cases := map[string]string{
		"select":                "SELECT a, b FROM t WHERE a > 1 ORDER BY a LIMIT 2",
		"insert":                "INSERT INTO t (a) VALUES (1)",
		"update":                "UPDATE t SET a = 1 WHERE id = 1",
		"delete":                "DELETE FROM t WHERE id = 1",
		"create_table":          "CREATE TABLE t (id INT)",
		"create_table_like":     "CREATE TABLE a LIKE b",
		"alter_table":           "ALTER TABLE t ADD COLUMN c INT",
		"drop_table":            "DROP TABLE IF EXISTS t",
		"truncate_table":        "TRUNCATE TABLE t",
		"rename_table":          "RENAME TABLE a TO b",
		"create_index":          "CREATE INDEX i ON t (a)",
		"drop_index":            "DROP INDEX i ON t",
		"create_database":       "CREATE DATABASE db1",
		"alter_database":        "ALTER DATABASE db1 DEFAULT CHARACTER SET utf8mb4",
		"drop_database":         "DROP DATABASE db1",
		"use_database":          "USE mydb",
		"create_view":           "CREATE VIEW v AS SELECT 1",
		"drop_view":             "DROP VIEW v",
		"create_sequence":       "CREATE SEQUENCE s START WITH 1",
		"alter_sequence":        "ALTER SEQUENCE s INCREMENT BY 2",
		"drop_sequence":         "DROP SEQUENCE s",
		"create_user":           "CREATE USER 'u'@'%' IDENTIFIED BY 'p'",
		"drop_user":             "DROP USER 'u'@'%'",
		"grant":                 "GRANT SELECT ON db1.* TO 'u'@'%'",
		"revoke":                "REVOKE SELECT ON db1.* FROM 'u'@'%'",
		"create_role":           "CREATE ROLE r1",
		"drop_role":             "DROP ROLE r1",
		"set_role":              "SET ROLE r1",
		"set_password":          "SET PASSWORD FOR 'u'@'%' = 'p'",
		"rename_user":           "RENAME USER 'a'@'%' TO 'b'@'%'",
		"lock_user":             "LOCK USER 'u'@'%'",
		"set_transaction":       "SET TRANSACTION ISOLATION LEVEL READ COMMITTED",
		"begin":                 "BEGIN",
		"commit":                "COMMIT",
		"rollback":              "ROLLBACK",
		"create_savepoint":      "SAVEPOINT s1",
		"rollback_savepoint":    "ROLLBACK TO SAVEPOINT s1",
		"release_savepoint":     "RELEASE SAVEPOINT s1",
		"lock_tables":           "LOCK TABLES t READ",
		"unlock_tables":         "UNLOCK TABLES",
		"prepare":               "PREPARE s FROM 'SELECT 1'",
		"execute":               "EXECUTE s USING @a",
		"deallocate_prepare":    "DEALLOCATE PREPARE s",
		"show":                  "SHOW TABLES",
		"set_var":               "SET @x = 1",
		"explain":               "EXPLAIN SELECT 1",
		"kill":                  "KILL QUERY 1",
		"help":                  "HELP 'SELECT'",
		"use":                   "USE db1",
		"analyze":               "ANALYZE TABLE t",
		"optimize":              "OPTIMIZE TABLE t",
		"purge":                 "PURGE RECYCLEBIN",
		"flashback":             "FLASHBACK TABLE t TO BEFORE DROP",
		"audit":                 "AUDIT SELECT ON db1.*",
		"dump_memory":           "DUMP ENTITY ALL",
		"load_data":             "LOAD DATA INFILE 'f' INTO TABLE t",
		"get_diagnostics":       "GET DIAGNOSTICS a = MESSAGE_TEXT",
		"create_outline":        "CREATE OUTLINE o ON SELECT 1",
		"alter_outline":         "ALTER OUTLINE o ON SELECT 1",
		"drop_outline":          "DROP OUTLINE o",
		"create_tablespace":     "CREATE TABLESPACE ts",
		"alter_tablespace":      "ALTER TABLESPACE ts",
		"drop_tablespace":       "DROP TABLESPACE ts",
		"create_tablegroup":     "CREATE TABLEGROUP tg",
		"alter_tablegroup":      "ALTER TABLEGROUP tg",
		"drop_tablegroup":       "DROP TABLEGROUP tg",
		"create_tenant":         "CREATE TENANT t RESOURCE_POOL_LIST = ('p')",
		"alter_tenant":          "ALTER TENANT t",
		"drop_tenant":           "DROP TENANT t",
		"create_standby_tenant": "CREATE STANDBY TENANT t LOG_RESTORE_SOURCE = 's'",
		"switchover_tenant":     "ALTER SYSTEM SWITCHOVER TO PRIMARY TENANT t",
		"recover_tenant":        "ALTER SYSTEM RECOVER STANDBY TENANT t UNTIL UNLIMITED",
		"rotate_master_key":     "ALTER SYSTEM ROTATE INNODB MASTER KEY",
		"create_dblink":         "CREATE DATABASE LINK dl CONNECT TO 'u' DATABASE 'd' IDENTIFIED BY 'p'",
		"drop_dblink":           "DROP DATABASE LINK dl",
		"create_restore_point":  "CREATE RESTORE POINT rp",
		"drop_restore_point":    "DROP RESTORE POINT rp",
		"service_name":          "ALTER SYSTEM CREATE SERVICE s",
		"transfer_partition":    "ALTER SYSTEM TRANSFER PARTITION TABLE_ID = 1, OBJECT_ID = 0 TO LS 1",
		"xa_begin":              "XA START 'x'",
		"xa_end":                "XA END 'x'",
		"xa_prepare":            "XA PREPARE 'x'",
		"xa_commit":             "XA COMMIT 'x' ONE PHASE",
		"xa_rollback":           "XA ROLLBACK 'x'",
		"xa_recover":            "XA RECOVER",
	}
	opts := lexer.Options{AllowDoubleQuotedIdentifiers: true}
	for name, sql := range cases {
		t.Run(name, func(t *testing.T) {
			toks, err := lexer.LexAll(sql, opts)
			if err != nil {
				t.Fatalf("lex %q: %v", sql, err)
			}
			stmts, err := Parse(toks, Options{})
			if err != nil || len(stmts) != 1 {
				t.Fatalf("parse %q: %v (%d stmts)", sql, err, len(stmts))
			}
			s1 := formatter.Format(stmts[0])
			// canonical re-parse
			toks2, _ := lexer.LexAll(s1, opts)
			again, err := Parse(toks2, Options{})
			if err != nil || len(again) != 1 {
				t.Fatalf("re-parse %q (from %q): %v", s1, sql, err)
			}
			if s2 := formatter.Format(again[0]); s2 != s1 {
				t.Fatalf("not idempotent: %q != %q", s1, s2)
			}
			// JSON round-trip
			data, err := jsonpkg.Marshal(stmts[0])
			if err != nil {
				t.Fatalf("json marshal: %v", err)
			}
			back, err := jsonpkg.Unmarshal(data)
			if err != nil {
				t.Fatalf("json unmarshal: %v", err)
			}
			if s3 := formatter.Format(back); s3 != s1 {
				t.Fatalf("json not preserved: %q != %q", s1, s3)
			}
		})
	}
}
