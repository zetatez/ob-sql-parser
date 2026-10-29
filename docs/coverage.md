# 语法覆盖矩阵

> 本文回答"**解析器能解析什么**"。它是按**能力形态**（而非官方文法 `stmt:` 分支）组织的
> 活文档：每行能力都有一一对应的测试用例，改动能力时必须同步改表与测试。
>
> 读者：集成方做能力评估 / 语法选型；维护者做回归对照。
> 相关文档：[version-coverage](version-coverage.md) 讲"某语法在某版本是否可用"，[testing](testing.md) 讲"如何保证本表为真"。

## 口径与图例

| 图例 | 含义 |
| --- | --- |
| 支持 | 该形态完整支持（含常用子形态） |
| 部分 | 部分支持 / 仅常用子集 |
| 超集 | **超集**：OB 官方各版文法均无，作为 MySQL 兼容语法额外提供 |
| 门控 | 该能力有版本门控，`Since` 为 OB n.n（详见 version-coverage） |

验证：本表每个"支持"项至少一条 SQL 出现在 `gparser/coverage_test.go`
（supported）并通过 parse -> format -> parse 幂等；超集项同时列入 unsupported（相对 OB）。

## 1. 查询（SELECT）

| 子能力 | 状态 | 说明与示例 |
| --- | --- | --- |
| 投影 / 别名 | 支持 | `SELECT a, b AS x, t.*` |
| DISTINCT / ALL | 支持 | `SELECT DISTINCT a FROM t` |
| FROM 源 | 支持 | 表、`db.tbl`、派生表 `(SELECT ...) x`、`LATERAL (...)` |
| JOIN | 支持 | `[INNER] JOIN`、`LEFT/RIGHT [OUTER]`、`FULL`、`CROSS`、`STRAIGHT_JOIN`、`NATURAL`、`ON`/`USING` |
| 索引提示 | 支持 | `USE/FORCE/IGNORE INDEX [FOR JOIN\|ORDER BY\|GROUP BY] (idx)` |
| WHERE | 支持 | 任意布尔表达式 |
| GROUP BY | 支持 | 项、`WITH ROLLUP`、门控4.4 `GROUPING SETS / CUBE / ROLLUP` |
| HAVING / QUALIFY | 支持 / 超集 | `QUALIFY` 为窗口过滤的 MySQL 超集 |
| ORDER BY | 支持 | `ASC/DESC`、`NULLS FIRST/LAST`、`COLLATE` |
| LIMIT / OFFSET | 支持 | `LIMIT n`、`LIMIT off, n`、`LIMIT n OFFSET off` |
| 集合运算 | 支持 | `UNION [ALL]`、`EXCEPT`、`INTERSECT`、`MINUS` |
| WITH CTE | 支持 | 门控4.0 `WITH [RECURSIVE] c AS (...)` |
| 窗口函数 | 支持 | `OVER (PARTITION BY ... ORDER BY ...)`、门控4.0 ROWS/RANGE 帧、门控4.2 命名 `WINDOW` |
| 行锁 | 支持 | `FOR UPDATE [NOWAIT\|SKIP LOCKED]`、`FOR SHARE`、`LOCK IN SHARE MODE`(超集) |
| INTO | 支持 | `INTO OUTFILE/DUMPFILE '...' [FIELDS/LINES ...]`、`INTO @a, @b` |
| PROCEDURE ANALYSE | 超集 | `... PROCEDURE ANALYSE(...)` |

## 2. DML

| 语句 | 状态 | 说明与示例 |
| --- | --- | --- |
| INSERT / REPLACE | 支持 | 多行 `VALUES`、`INSERT ... SET`、`... SELECT`、`IGNORE`、`ON DUPLICATE KEY UPDATE`、`RETURNING` |
| UPDATE | 支持 | 单/多表、别名、`LOW_PRIORITY`/`IGNORE`、`WHERE`、`ORDER BY`、`LIMIT`、`RETURNING`(超集)、JOIN 源 `UPDATE t JOIN u ... SET` |
| DELETE | 支持 | 多表 `DELETE t1[,t2] FROM ...`、`USING`、`LOW_PRIORITY/QUICK/IGNORE`、`WHERE`、`ORDER BY`、`LIMIT`、`RETURNING`(超集) |
| MERGE | 支持 + 门控4.4 | `WHEN MATCHED [THEN UPDATE/DELETE]`、`WHEN NOT MATCHED THEN INSERT` |

## 3. DDL

| 语句 | 状态 | 说明与示例 |
| --- | --- | --- |
| CREATE TABLE | 支持 | 列定义、类型参数/UNSIGNED/ZEROFILL、`NULL/NOT NULL`、`DEFAULT`、`ON UPDATE CURRENT_TIMESTAMP[(n)]`、注释、`AUTO_INCREMENT`、`CHARSET/COLLATE`、`INVISIBLE`、生成列 `AS (expr) [VIRTUAL\|STORED]`、表级约束（PK/UNIQUE/KEY/INDEX/`FOREIGN KEY ... REFERENCES ... ON DELETE/UPDATE RESTRICT\|CASCADE\|SET NULL\|NO ACTION`/`CHECK`）、列级 `CHECK (expr)`、表选项（ENGINE/CHARSET/COMMENT/BLOCK_SIZE/COMPRESSION/ROW_FORMAT/PRIMARY_ZONE/REPLICA_NUM/TABLEGROUP 等，含布尔值 `USE_BLOOM_FILTER = TRUE/FALSE`）、CTAS、`LIKE` |
| 分区 | 支持 | `PARTITION BY RANGE/LIST/HASH/KEY [COLUMNS]`、`PARTITIONS n`、`PARTITIONS AUTO`、`SUBPARTITION BY ...`、`PARTITION p VALUES LESS THAN/IN (.../MAXVALUE)`、分区项选项 |
| VIEW | 支持 | `CREATE [OR REPLACE] / ALTER / DROP VIEW` |
| 物化视图 + MLOG | 支持 + 门控4.3 | `CREATE/ALTER/DROP MATERIALIZED VIEW`、`CREATE/DROP MATERIALIZED VIEW LOG` |
| INDEX | 支持 | `CREATE [UNIQUE] INDEX [IF NOT EXISTS] name [USING BTREE\|HASH] ON t (col[(n)] [ASC\|DESC], ...)` + `STORING (cols)` 覆盖列 + `WITH (参数)` / `WITH COLUMN GROUP (ALL COLUMNS\|EACH COLUMN\|name (cols))` + 索引选项（GLOBAL/LOCAL/KEY_BLOCK_SIZE/BLOCK_SIZE/COMMENT）+ 分区、`DROP INDEX` |
| SEQUENCE / TABLESPACE / TABLEGROUP | 支持 | `CREATE/ALTER/DROP` |
| TENANT / RESOURCE | 支持 | `CREATE/ALTER/DROP TENANT`、`RESOURCE UNIT/POOL`、STANDBY/SWITCHOVER/RECOVER（见 OB 特有节） |
| OUTLINE | 支持 | `CREATE/ALTER/DROP OUTLINE` |
| DBLINK | 支持 | `CREATE/DROP DATABASE LINK` |
| LOCATION / CATALOG | 支持 + 门控4.4 | `CREATE/ALTER/DROP LOCATION`、`CREATE/DROP/SET CATALOG` |
| CCL / SENSITIVE RULE | 支持 + 门控4.4 | `CREATE/ALTER/DROP CCL RULE`、`... SENSITIVE RULE` |
| ALTER TABLE | 支持 | `ADD [COLUMN] (多列)`、`ADD/DROP/MODIFY/CHANGE [COLUMN]` + `FIRST/AFTER`、索引/约束（含索引选项 `LOCAL/GLOBAL/KEY_BLOCK_SIZE/COMMENT`）、`ALTER INDEX ... VISIBLE/INVISIBLE`、`ENABLE/DISABLE KEYS`、`RENAME`、`AUTO_INCREMENT`、`CONVERT TO CHARACTER SET`、分区操作（`ADD/DROP/TRUNCATE/COALESCE/EXCHANGE [门控4.4 WITH[OUT] VALIDATION]/REORGANIZE/REMOVE PARTITIONING`） |
| DROP / TRUNCATE | 支持 | `DROP TABLE/INDEX/DATABASE/VIEW/SEQUENCE/... [IF EXISTS]`、`TRUNCATE TABLE` |

## 4. 表达式

支持 全套：算术/位/移位/比较/逻辑（含 `DIV/MOD`、`\|\|/&&`、`RLIKE/REGEXP/SOUNDS LIKE`）、
`COLLATE`、`IN/BETWEEN/LIKE ... ESCAPE/IS [NOT] NULL|TRUE|FALSE`、`CAST`/`CONVERT`、
`CASE`、`INTERVAL`、`EXISTS`/标量子查询、`VALUES()`、`MATCH ... AGAINST`、
`MEMBER OF`(超集)、字面量（含 `_charset` 引入符、hex/bit/字符串转义）、`@var`/`@@var`、`?` 占位。

```sql
SELECT a + b * 2, a DIV 3, a || b, a & b, a << 2, a COLLATE utf8mb4_bin
     , a IN (1,2), a BETWEEN 1 AND 9, a LIKE 'x%' ESCAPE '!', a IS NOT NULL
     , CAST(a AS SIGNED), CONVERT(a, CHAR), CASE WHEN a>0 THEN 1 ELSE 0 END
     , INTERVAL 3 DAY, EXISTS (SELECT 1 FROM u WHERE u.id = t.id)
     , MATCH (title) AGAINST ('ocean'), _utf8mb4'字串', 0xFF, b'101', @x, @@sql_mode, ?
FROM t;
```

## 5. 权限 / 账号

支持 `GRANT/REVOKE`（列级 priv、`ALL [PRIVILEGES]`、`WITH GRANT OPTION`）、
`CREATE/DROP USER`（`IDENTIFIED BY` / `IDENTIFIED WITH <plugin>`）、
`CREATE/DROP/SET ROLE`、`SET PASSWORD`、`RENAME/LOCK/UNLOCK USER`。

```sql
GRANT SELECT (a), INSERT ON db.* TO 'u'@'%' WITH GRANT OPTION;
CREATE USER 'u'@'%' IDENTIFIED BY 'pw';
SET ROLE r1;  -- 账号须带引号（'u'@'h'）；裸 u@localhost 暂不支持
```

## 6. 事务 / 锁 / 会话

支持 `BEGIN / START TRANSACTION / COMMIT / ROLLBACK [TO SAVEPOINT]`、
`SAVEPOINT/RELEASE SAVEPOINT`、`SET [GLOBAL|SESSION] TRANSACTION ...`、
`LOCK/UNLOCK TABLES`、XA 全动作（门控4.2.5 `XA RECOVER`）、
`SET @var` / `SET NAMES` / `SET CHARACTER SET` / `SET ROLE` / `SET CATALOG`。

## 7. SHOW 族

支持/部分：`TABLES`、`DATABASES/SCHEMAS`、`VARIABLES`、`STATUS`、`CREATE TABLE/DATABASE/VIEW/
PROCEDURE/FUNCTION/TRIGGER/USER`、`COLUMNS/FIELDS`、`INDEX/INDEXES/KEYS`、`TABLE STATUS`、
`TRIGGERS`、`CHARACTER SET/CHARSET`、`PRIVILEGES`、`PLUGINS`、`PROFILES/PROFILE`、
`OPEN TABLES`、`WARNINGS`、`ERRORS`、`GRANTS [FOR ...]`、`PROCESSLIST`、`COLLATION`、
`ENGINES/STORAGE ENGINES`、`ENGINE name STATUS|MUTEX`、`FUNCTION/PROCEDURE STATUS`、
`EVENTS`、`MASTER STATUS`；各形态按官方支持 `LIKE/FROM/WHERE`。
（部分：SHOW 子参数存在无穷变体，按需扩。）

## 8. 运维 / OB 特有

支持 `KILL [QUERY|CONNECTION]`、`HELP`、`DO`、`ANALYZE/OPTIMIZE/REPAIR TABLE`（含 `METHOD_OPT`）、
`LOAD DATA [/*+ DIRECT */]`、`FLUSH`、`PURGE`、`FLASHBACK TABLE`、`AUDIT/NOAUDIT`、
`DUMP ENTITY|CHUNK`、`SERVICE`、门控5.0 `ROUTINE LOAD`、`ALTER SYSTEM SET`、
`ROTATE MASTER KEY`、门控4.2.2 `TRANSFER PARTITION`、`INSTALL/UNINSTALL PLUGIN`、
`HANDLER`、`EXPLAIN [EXTENDED|PRETTY|FORMAT=|ANALYZE|FOR CONNECTION]`。

## 9. PL（存储过程 / 函数）

支持 `CREATE [OR REPLACE] PROCEDURE/FUNCTION`（参数/返回类型 + `BEGIN...END` 复合块）：

```sql
CREATE PROCEDURE p(IN x INT) BEGIN
  DECLARE v INT DEFAULT 0;
  DECLARE done INT DEFAULT FALSE;
  DECLARE c CURSOR FOR SELECT a FROM t;
  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
  IF x > 0 THEN
    SET v = x;
  ELSEIF x = 0 THEN
    SET v = -1;
  ELSE
    SET v = 0;
  END IF;
  loop1: LOOP
    FETCH c INTO v;
    IF done THEN LEAVE loop1; END IF;
  END LOOP;
  RETURN v;
END;
```

覆盖：`DECLARE`（变量 / `CURSOR FOR` / `CONDITION`）、`SET`、`IF...ELSEIF...ELSE`、
`WHILE/REPEAT/LOOP`、`CASE` 语句、`RETURN`、标签 `label:` + `LEAVE/ITERATE`、
`DECLARE ... HANDLER`（`CONTINUE/EXIT/UNDO` + `SQLSTATE/NOT FOUND`/命名条件）、
`SIGNAL/RESIGNAL`、`GET DIAGNOSTICS [CONDITION n]`、嵌套块、体内任意 SQL；
`DROP ROUTINE`、`CALL`。`;` 由文法消费以支撑体内多语句。

## 10. 超集清单（MySQL 兼容，OB 各版均无）

超集 `QUALIFY`、`PROCEDURE ANALYSE`、`MEMBER OF`、`LOCK IN SHARE MODE`、
`GENERATED ALWAYS AS (expr)`、`WINDOW` 命名窗（OB 4.2 已支持则视为支持，按 version-coverage）。

## 11. 未实现 / 边界

- `MOCK`（OB 内部测试专用集合语法；其构成语句均已独立覆盖）。
- SHOW 无穷子参数（`SHOW PROFILE ... OFFSET` 等，按需）。
- 裸写 `u@localhost` 账号（与 `@` 变量词法冲突）。
- 残留 24 处 shift/reduce 冲突（官方文法固有，默认移进优先 = 期望语义）。

## 12. 数据来源与验证

| 能力类别 | 断言所在测试 |
| --- | --- |
| 支持/部分 面 | `gparser/coverage_test.go::TestCoverageSupported`（约 50 条） |
| 超集 / 不支持 | `gparser/coverage_test.go::TestCoverageUnsupported` |
| 语句分支全量（官方 91/91 + 5.x 26） | `gparser/corpus_test.go::TestOfficialCorpus`（82 条三重校验） |
| 版本门控（门控 标注） | `gparser/version_corpus_test.go`、`gparser/gparser_test.go::TestFeatureGating` |

改动能力时：**先加测试用例**（round-trip 断言通过），再更新本表对应行与图例。