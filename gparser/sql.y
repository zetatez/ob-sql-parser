// Code generated from sql.y by goyacc. DO NOT EDIT.

%{
package gparser

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/feature"
)
%}

%union {
	node interface{}
}

%token <node> NAME_OB STRING_VALUE INTNUM DECIMAL_VAL APPROXNUM BOOL_VALUE NULLX HEXNUM BITNUM HEX_STRING_VALUE BIN_STRING_VALUE USER_VARIABLE SYSTEM_VARIABLE QUESTIONMARK

%token SELECT INSERT REPLACE UPDATE DELETE CREATE DROP TRUNCATE SHOW SET USE EXPLAIN DESCRIBE DESC
%token BEGIN START COMMIT ROLLBACK WORK TRANSACTION READ WRITE ONLY CHAIN
%token WITH RECURSIVE VALUES VALUE INTO FROM WHERE GROUP BY HAVING ORDER LIMIT OFFSET
%token UNION EXCEPT INTERSECT MINUS ALL DISTINCT AS AND OR XOR NOT IS LIKE RLIKE REGEXP IN BETWEEN
%token CASE WHEN THEN ELSE END EXISTS INTERVAL ASC NULLS FIRST LAST
%token OVER PARTITION ROWS RANGE CURRENT ROW UNBOUNDED PRECEDING FOLLOWING WINDOW
%token JOIN INNER OUTER LEFT RIGHT FULL CROSS NATURAL STRAIGHT_JOIN ON USING DUAL FOR SHARE
%token TABLE DATABASE SCHEMA INDEX KEY PRIMARY UNIQUE FOREIGN REFERENCES CONSTRAINT CHECK
%token IF TEMPORARY IGNORE LOW_PRIORITY HIGH_PRIORITY DELAYED DUPLICATE RETURNING QUICK
%token <node> GLOBAL SESSION LOCAL DATABASES SCHEMAS TABLES VARIABLES STATUS WARNINGS ERRORS GRANTS PROCESSLIST
%token ENGINES STORAGE COLLATION COLUMNS FIELDS INDEXES KEYS VIEW TRIGGERS PROCEDURE FUNCTION NAMES
%token CHARACTER EXTENDED PRETTY FORMAT CASCADE RESTRICT UNSIGNED DEFAULT COMMENT AUTO_INCREMENT CHARSET COLLATE INVISIBLE BINARY CAST
%token COMP_EQ COMP_NE COMP_LT COMP_GT COMP_LE COMP_GE NSEQ SHIFT_LEFT SHIFT_RIGHT ESCAPE ILLEGAL_TOK
%token DIV MOD SOUNDS
%token ALTER ADD MODIFY CHANGE RENAME COLUMN TO FORCE ROLLUP
%token LIST HASH MAXVALUE LESS THAN PARTITIONS
%token GENERATED ALWAYS VIRTUAL STORED
%token SEQUENCE INCREMENT CACHE MINVALUE CYCLE NOCYCLE NOORDER
%token ISOLATION LEVEL UNCOMMITTED COMMITTED REPEATABLE SERIALIZABLE SAVEPOINT RELEASE LOCK UNLOCK
%token PLUGINS PROFILES PROFILE PRIVILEGES TRIGGER OPEN EVENTS MUTEX
%token GRANT REVOKE OPTION AT
%token USER ROLE IDENTIFIED PREPARE EXECUTE DEALLOCATE
%token <node> NONE PASSWORD KILL HELP ANALYZE OPTIMIZE REPAIR CONNECTION QUERY
%token LOAD DATA INFILE TERMINATED LINES FLUSH EXPORT
%token OROR ANDAND
%token XA RESUME SUSPEND MIGRATE ONE PHASE RECOVER CONVERT XID
%token PURGE RECYCLEBIN
%token MEMBER AGAINST MATCH GROUPING SETS CUBE BOOLEAN MODE EXPANSION OF
%token SYSTEM
%token MERGE MATCHED FLASHBACK BEFORE SNAPSHOT
%token MATERIALIZED REFRESH FAST COMPLETE FORCE DEMAND
%token SUBPARTITION SUBPARTITIONS
%token SIGNAL RESIGNAL SQLSTATE GET DIAGNOSTICS STACKED
%token TENANT RESOURCE UNIT POOL TABLEGROUP TABLESPACE
%token OUTLINE MASTER ROTATE INCLUDING NEW LOG
%token STANDBY ACTIVATE SWITCHOVER UNTIL UNLIMITED CANCEL SCN RESTORE RECOVER
%token POINT LS TIME
%token <node> AUDIT NOAUDIT DUMP ENTITY CHUNK SERVICE REMOVE FILES LOCATION URL
%token ROUTINE STREAM PAUSE RESUME STOP
%token <node> CONCURRENT_LIMITING_RULE SENSITIVE RULE USING LINK CONNECT HOST PORT CLUSTER
%token LATERAL QUALIFY AUTO SIZE
%token PLUGIN SONAME INSTALL UNINSTALL HANDLER METHOD_OPT ANALYSE NEXT PREV CLOSE
%token REORGANIZE AFTER COALESCE EXCHANGE
%token DUMPFILE LOGS HOSTS OUTFILE CONDITION EXTERNAL CATALOG PROPERTIES TRANSFER
%token RETURNS DECLARE CURSOR FETCH LOOP WHILE REPEAT UNTIL RETURN INOUT ELSEIF
%token CALL OUT DO CONTINUE EXIT UNDO LEAVE ITERATE FOUND WITHOUT LOCKED VALIDATION NOWAIT
%token SKIP
%token <node> DIRECT_HINT

%type <node> stmt stmt_inner inner_stmt
%type <node> call_stmt do_stmt proc_params proc_param_list proc_param proc_body opt_proc_mode
%type <node> sp_stmts sp_stmt sp_elseifs opt_sp_else
%type <node> opt_decl_default handler_action handler_conditions handler_condition
%type <node> opt_pl_case_operand pl_case_when_list pl_case_when opt_pl_case_else
%type <node> select_stmt select_no_parens select_rhs simple_select
%type <node> opt_with cte_list cte cte_cols_opt opt_recursive
%type <node> opt_order order_list order_item opt_asc_desc opt_nulls
%type <node> opt_limit limit_clause
%type <node> distinct_opt select_list select_item alias_opt
%type <node> from_opt from_list table_ref table_primary table_name join_cond_opt opt_outer
%type <node> where_opt group_opt group_list group_item opt_rollup having_opt opt_qualify window_opt window_list named_window for_opt
%type <node> grouping_expr grouping_sets grouping_set opt_match_modifier
%type <node> expr or_expr xor_expr and_expr not_expr bool_pri predicate bit_expr
%type <node> bit_or bit_and shift_expr additive multiplicative power unary simple_expr
%type <node> expr_list in_expr comp_op
%type <node> column_ref func_call call_args over_opt window_spec window_part_opt window_order_opt window_frame_opt frame_extent frame_bound
%type <node> date_unit interval_expr cast_expr case_expr case_operand when_list else_opt literal
%type <node> row row_list
%type <node> insert_stmt ignore_opt col_list_opt insert_source on_dup_opt assignment_list assignment returning_opt
%type <node> update_stmt delete_stmt delete_mods using_opt table_ref_list
%type <node> create_stmt create_table_tail create_def_list create_def column_def column_attrs_opt column_attr
%type <node> opt_stored opt_as_select as_select sequence_options sequence_option
%type <node> table_constraint constraint_body table_options_opt table_option table_option_value
%type <node> opt_if_not_exists opt_if_exists opt_unique index_col_list index_col db_opts_opt db_opt opt_into
%type <node> drop_stmt opt_on_table table_name_list truncate_stmt opt_table_kw
%type <node> alter_stmt alter_action_list alter_action opt_column
%type <node> column_def_list opt_col_position opt_convert_collate part_name_or_all opt_exchange_validation create_view_stmt drop_view_stmt
%type <node> opt_index_hint_list index_hint_list index_hint index_hint_for_opt index_list index_name
%type <node> opt_partition partition_by opt_partitions_count opt_partition_defs partition_def_list partition_def
%type <node> opt_subpartition opt_subpartitions_count
%type <node> opt_auto_partition
%type <node> opt_partition_values part_value_list part_value opt_partition_clause name_list
%type <node> show_stmt full_opt show_body show_scope_opt show_from_opt like_opt
%type <node> set_stmt set_body set_item_list set_item set_lhs opt_scope names_arg
%type <node> use_stmt explain_stmt opt_explain_type begin_stmt commit_stmt rollback_stmt describe_stmt opt_work opt_start_read
%type <node> transaction_characteristics transaction_characteristic isolation_level savepoint_kw
%type <node> savepoint_stmt release_savepoint_stmt lock_tables_stmt unlock_tables_stmt
%type <node> lock_table_list lock_table opt_lock_alias lock_type
%type <node> grant_stmt revoke_stmt grant_priv_list grant_priv opt_priv_cols opt_privileges_kw
%type <node> priv_level user_spec_list user_spec user_name user_host opt_with_grant
%type <node> create_user_list create_user_spec opt_identified
%type <node> prepare_stmt execute_stmt deallocate_stmt prep_source execute_opt_using user_var_list
%type <node> role_name_list set_role_target opt_for_user
%type <node> kill_stmt help_stmt opt_kill_kind table_maint_stmt
%type <node> rename_table_stmt rename_pair_list rename_pair
%type <node> load_data_stmt opt_local opt_load_conflict load_fields_opt load_lines_opt load_ignore_lines_opt load_cols_opt
%type <node> flush_stmt flush_body flush_tables_opt flush_opt
%type <node> xa_stmt xid xa_join_opt xa_suspend_opt xa_one_phase_opt xa_recover_opt purge_stmt
%type <node> merge_stmt merge_when_list merge_when opt_delete_where opt_merge_where flashback_table_stmt opt_flashback_rename
%type <node> mview_opts mview_opt
%type <node> signal_stmt opt_sqlstate_clause opt_signal_set signal_item_list signal_item
%type <node> get_diagnostics_stmt opt_diag_scope diag_item_list diag_item
%type <node> opt_cluster_options cluster_option_list cluster_option cluster_option_value cluster_opt_name opt_purge
%type <node> rename_user_stmt lock_user_stmt
%type <node> part_info
%type <node> mlog_stmt opt_mlog_opts mlog_opt
%type <node> create_outline_stmt alter_outline_stmt drop_outline_stmt outline_target
%type <node> rotate_master_key_stmt
%type <node> opt_tenant_kw recover_until opt_log_restore_source opt_clone_snapshot opt_for_tenant opt_innodb
%type <node> audit_stmt dump_memory_stmt location_utils_stmt
%type <node> opt_or_replace opt_equal
%type <node> opt_load_properties opt_job_properties load_property_list load_property routine_load_action_stmt opt_direct_hint
%type <node> opt_dblink_tail
%type <node> install_plugin_stmt uninstall_plugin_stmt handler_stmt
%type <node> opt_handler_alias handler_direction opt_handler_cols opt_handler_where opt_handler_limit
%type <node> opt_procedure_analyse
%type <node> opt_into_file opt_into_fields_opt opt_diag_condition into_var_list opt_show_where show_engine_mode
%type <node> ident_list data_type type_args_opt unsigned_opt set_op opt_all

%start stmt

%%

stmt:
	stmt_inner
	{
		lset(yylex, asStmt($1))
		$$ = $1
	}
	| stmt ';'
	{
		$$ = $1
	}
	| stmt ';' stmt_inner
	{
		lset(yylex, asStmt($3))
		$$ = $3
	}

// inner_stmt is a statement used where no result should be recorded (EXPLAIN).
inner_stmt:
	stmt_inner
	{
		$$ = $1
	}

stmt_inner:
	select_stmt
	| insert_stmt
	| update_stmt
	| delete_stmt
	| create_stmt
	| create_view_stmt
	| alter_stmt
	| drop_stmt
	| drop_view_stmt
	| truncate_stmt
	| show_stmt
	| set_stmt
	| use_stmt
	| explain_stmt
	| begin_stmt
	| commit_stmt
	| rollback_stmt
	| savepoint_stmt
	| release_savepoint_stmt
	| lock_tables_stmt
	| unlock_tables_stmt
	| grant_stmt
	| revoke_stmt
	| prepare_stmt
	| execute_stmt
	| deallocate_stmt
	| kill_stmt
	| help_stmt
	| table_maint_stmt
	| rename_table_stmt
	| load_data_stmt
	| flush_stmt
	| xa_stmt
	| purge_stmt
	| merge_stmt
	| flashback_table_stmt
	| signal_stmt
	| get_diagnostics_stmt
	| rename_user_stmt
	| lock_user_stmt
	| mlog_stmt
	| create_outline_stmt
	| alter_outline_stmt
	| drop_outline_stmt
	| rotate_master_key_stmt
	| audit_stmt
	| dump_memory_stmt
	| location_utils_stmt
	| routine_load_action_stmt
	| install_plugin_stmt
	| uninstall_plugin_stmt
	| handler_stmt
	| call_stmt
	| do_stmt
	| describe_stmt

/* ---------------- SELECT ---------------- */

select_stmt:
	opt_with select_no_parens opt_order opt_limit opt_into_file
	{
		sel := asSelect($2)
		if sel != nil {
			sel.OrderBy = asOrderList($3)
			sel.Limit = asLimit($4)
			if $5 != nil {
				fi := $5.(intoFile)
				sel.IntoKind = fi.kind
				sel.IntoFile = fi.file
				sel.IntoVars = fi.vars
				sel.IntoFields = fi.fields
				sel.IntoLines = fi.lines
			}
		}
		$$ = sel
	}

opt_into_file:
	/*empty*/
	{
		$$ = nil
	}
	| INTO OUTFILE STRING_VALUE opt_into_fields_opt
	{
		fo := $4.(intoFields)
		$$ = intoFile{kind: "OUTFILE", file: asLitValue($3), fields: fo.fields, lines: fo.lines}
	}
	| INTO DUMPFILE STRING_VALUE
	{
		$$ = intoFile{kind: "DUMPFILE", file: asLitValue($3)}
	}
	| INTO into_var_list
	{
		$$ = intoFile{kind: "@vars", vars: asVariableList($2)}
	}

into_var_list:
	USER_VARIABLE
	{
		$$ = []*ast.Variable{$1.(*ast.Variable)}
	}
	| into_var_list ',' USER_VARIABLE
	{
		$$ = append(asVariableList($1), $3.(*ast.Variable))
	}

opt_into_fields_opt:
	/*empty*/
	{
		$$ = intoFields{}
	}
	| FIELDS TERMINATED BY STRING_VALUE
	{
		$$ = intoFields{fields: asLitValue($4)}
	}
	| FIELDS TERMINATED BY STRING_VALUE LINES TERMINATED BY STRING_VALUE
	{
		$$ = intoFields{fields: asLitValue($4), lines: asLitValue($8)}
	}
	| LINES TERMINATED BY STRING_VALUE
	{
		$$ = intoFields{lines: asLitValue($4)}
	}

opt_with:
	/*empty*/
	{
		$$ = nil
	}
	| WITH opt_recursive cte_list
	{
		gateWith(yylex)
		$$ = nil
	}

opt_recursive:
	/*empty*/
	{
		$$ = nil
	}
	| RECURSIVE
	{
		$$ = nil
	}

cte_list:
	cte
	| cte_list ',' cte

cte:
	NAME_OB cte_cols_opt AS '(' select_stmt ')'

cte_cols_opt:
	/*empty*/
	{
		$$ = nil
	}
	| '(' ident_list ')'
	{
		$$ = nil
	}

select_no_parens:
	simple_select
	{
		$$ = &ast.SelectStmt{Bodies: []*ast.SelectBody{asBody($1)}}
	}
	| select_no_parens set_op select_rhs
	{
		sel := asSelect($1)
		op := $2.(setOp)
		right := asBody($3)
		sel.Bodies = append(sel.Bodies, right)
		sel.SetOps = append(sel.SetOps, &ast.SetOp{Op: op.op, All: op.all, Right: right})
		$$ = sel
	}

set_op:
	UNION opt_all
	{
		$$ = setOp{op: "UNION", all: asBool($2)}
	}
	| EXCEPT opt_all
	{
		$$ = setOp{op: "EXCEPT", all: asBool($2)}
	}
	| INTERSECT opt_all
	{
		$$ = setOp{op: "INTERSECT", all: asBool($2)}
	}
	| MINUS
	{
		$$ = setOp{op: "MINUS"}
	}

opt_all:
	/*empty*/
	{
		$$ = false
	}
	| ALL
	{
		$$ = true
	}

select_rhs:
	simple_select

simple_select:
	SELECT distinct_opt select_list from_opt where_opt group_opt having_opt opt_qualify window_opt for_opt opt_procedure_analyse
	{
		b := &ast.SelectBody{Distinct: asBool($2), Items: asItems($3)}
		if $4 != nil {
			b.From = asTableRefList($4)
		}
		b.Where = asExpr($5)
		g := asGroupOpt($6)
		b.GroupBy = g.exprs
		b.Rollup = g.rollup
		b.Having = asExpr($7)
		b.Qualify = asExpr($8)
		b.NamedWindow = asNamedWindows($9)
		lo := asLockOpt($10)
		b.ForUpdate = lo.update
		b.LockWait = lo.wait
		b.InShareMode = lo.share
		pa := $11.(procAnalyse)
		b.HasProcAnalyse = pa.has
		b.ProcAnalyse = pa.args
		$$ = b
	}

opt_procedure_analyse:
	/*empty*/
	{
		$$ = procAnalyse{}
	}
	| PROCEDURE ANALYSE '(' expr_list ')'
	{
		$$ = procAnalyse{has: true, args: asExprList($4)}
	}
	| PROCEDURE ANALYSE '(' ')'
	{
		$$ = procAnalyse{has: true}
	}

opt_qualify:
	/*empty*/
	{
		$$ = nil
	}
	| QUALIFY expr
	{
		$$ = $2
	}

distinct_opt:
	/*empty*/
	{
		$$ = false
	}
	| DISTINCT
	{
		$$ = true
	}
	| ALL
	{
		$$ = false
	}

select_list:
	select_item
	{
		$$ = []ast.ColumnItem{asItem($1)}
	}
	| select_list ',' select_item
	{
		$$ = append(asItems($1), asItem($3))
	}

select_item:
	expr alias_opt
	{
		ci := ast.ColumnItem{Expr: asExpr($1)}
		if $2 != nil {
			ci.Alias = asStr($2)
			ci.HasAlias = true
		}
		$$ = ci
	}

alias_opt:
	/*empty*/
	{
		$$ = nil
	}
	| AS NAME_OB
	{
		$$ = asIdentName($2)
	}
	| NAME_OB
	{
		$$ = asIdentName($1)
	}

from_opt:
	/*empty*/
	{
		$$ = nil
	}
	| FROM from_list
	{
		$$ = $2
	}

from_list:
	table_ref
	{
		$$ = []*ast.TableRef{asTableRef($1)}
	}
	| from_list ',' table_ref
	{
		$$ = append(asTableRefList($1), asTableRef($3))
	}

table_ref:
	table_primary
	| table_ref JOIN table_primary join_cond_opt
	{
		$$ = mkJoin(asTableRef($1), asTableRef($3), "JOIN", $4)
	}
	| table_ref INNER JOIN table_primary join_cond_opt
	{
		$$ = mkJoin(asTableRef($1), asTableRef($4), "INNER JOIN", $5)
	}
	| table_ref CROSS JOIN table_primary
	{
		$$ = mkJoin(asTableRef($1), asTableRef($4), "CROSS JOIN", nil)
	}
	| table_ref LEFT opt_outer JOIN table_primary join_cond_opt
	{
		$$ = mkJoin(asTableRef($1), asTableRef($5), "LEFT JOIN", $6)
	}
	| table_ref RIGHT opt_outer JOIN table_primary join_cond_opt
	{
		$$ = mkJoin(asTableRef($1), asTableRef($5), "RIGHT JOIN", $6)
	}
	| table_ref FULL opt_outer JOIN table_primary join_cond_opt
	{
		$$ = mkJoin(asTableRef($1), asTableRef($5), "FULL JOIN", $6)
	}
	| table_ref STRAIGHT_JOIN table_primary
	{
		$$ = mkJoin(asTableRef($1), asTableRef($3), "STRAIGHT_JOIN", nil)
	}
	| table_ref NATURAL JOIN table_primary
	{
		$$ = mkJoin(asTableRef($1), asTableRef($4), "NATURAL JOIN", nil)
	}
	| table_ref NATURAL INNER JOIN table_primary
	{
		$$ = mkJoin(asTableRef($1), asTableRef($5), "NATURAL INNER JOIN", nil)
	}
	| table_ref NATURAL LEFT opt_outer JOIN table_primary
	{
		$$ = mkJoin(asTableRef($1), asTableRef($6), "NATURAL LEFT JOIN", nil)
	}
	| table_ref NATURAL RIGHT opt_outer JOIN table_primary
	{
		$$ = mkJoin(asTableRef($1), asTableRef($6), "NATURAL RIGHT JOIN", nil)
	}

opt_outer:
	/*empty*/
	{
		$$ = nil
	}
	| OUTER
	{
		$$ = nil
	}

table_primary:
	table_name opt_partition_clause alias_opt opt_index_hint_list
	{
		r := &ast.TableRef{RefKind: ast.RefTable, Name: asIdent($1)}
		r.Partitions = asStringList($2)
		if $3 != nil {
			r.Alias = asStr($3)
			r.HasAlias = true
		}
		r.IndexHints = asIndexHints($4)
		$$ = r
	}
	| '(' select_stmt ')' alias_opt
	{
		r := &ast.TableRef{RefKind: ast.RefSubquery, Subquery: asSelect($2)}
		if $4 != nil {
			r.Alias = asStr($4)
			r.HasAlias = true
		}
		$$ = r
	}
	| LATERAL '(' select_stmt ')' alias_opt
	{
		r := &ast.TableRef{RefKind: ast.RefSubquery, Lateral: true, Subquery: asSelect($3)}
		if $5 != nil {
			r.Alias = asStr($5)
			r.HasAlias = true
		}
		$$ = r
	}

table_name:
	NAME_OB
	{
		$$ = newIdent(asIdentName($1))
	}
	| NAME_OB '.' NAME_OB
	{
		$$ = idParts(asIdentName($1), asIdentName($3))
	}
	| NAME_OB '.' NAME_OB '.' NAME_OB
	{
		$$ = idParts(asIdentName($1), asIdentName($3), asIdentName($5))
	}

opt_index_hint_list:
	/*empty*/
	{
		$$ = nil
	}
	| index_hint_list

index_hint_list:
	index_hint
	{
		$$ = []ast.IndexHint{$1.(ast.IndexHint)}
	}
	| index_hint_list index_hint
	{
		$$ = append(asIndexHints($1), $2.(ast.IndexHint))
	}

index_hint:
	USE INDEX index_hint_for_opt '(' index_list ')'
	{
		$$ = ast.IndexHint{Type: ast.IndexHintUse, For: asStr($3), Indexes: asStringList($5)}
	}
	| FORCE INDEX index_hint_for_opt '(' index_list ')'
	{
		$$ = ast.IndexHint{Type: ast.IndexHintForce, For: asStr($3), Indexes: asStringList($5)}
	}
	| IGNORE INDEX index_hint_for_opt '(' index_list ')'
	{
		$$ = ast.IndexHint{Type: ast.IndexHintIgnore, For: asStr($3), Indexes: asStringList($5)}
	}

index_hint_for_opt:
	/*empty*/
	{
		$$ = ""
	}
	| FOR JOIN
	{
		$$ = "JOIN"
	}
	| FOR ORDER BY
	{
		$$ = "ORDER BY"
	}
	| FOR GROUP BY
	{
		$$ = "GROUP BY"
	}

index_list:
	index_name
	{
		$$ = []string{$1.(string)}
	}
	| index_list ',' index_name
	{
		$$ = append(asStringList($1), $3.(string))
	}

index_name:
	NAME_OB
	{
		$$ = asIdentName($1)
	}
	| PRIMARY
	{
		$$ = "PRIMARY"
	}

join_cond_opt:
	/*empty*/
	{
		$$ = nil
	}
	| ON expr
	{
		$$ = joinCond{on: asExpr($2)}
	}
	| USING '(' ident_list ')'
	{
		$$ = joinCond{using: asIdentList($3)}
	}

where_opt:
	/*empty*/
	{
		$$ = nil
	}
	| WHERE expr
	{
		$$ = $2
	}

group_opt:
	/*empty*/
	{
		$$ = groupOpt{}
	}
	| GROUP BY group_list opt_rollup
	{
		$$ = groupOpt{exprs: asExprList($3), rollup: asBool($4)}
	}

opt_rollup:
	/*empty*/
	{
		$$ = false
	}
	| WITH ROLLUP
	{
		$$ = true
	}

group_list:
	group_item
	{
		$$ = []ast.Expression{asExpr($1)}
	}
	| group_list ',' group_item
	{
		$$ = append(asExprList($1), asExpr($3))
	}

group_item:
	expr opt_asc_desc
	{
		$$ = $1
	}
	| grouping_expr

grouping_expr:
	GROUPING SETS '(' grouping_sets ')'
	{
		gateFeat(yylex, feature.FGroupingSets, "GROUPING SETS requires OceanBase 4.4+")
		$$ = &ast.GroupingSetExpr{Kind: "GROUPING SETS", Sets: asExprSets($4)}
	}
	| CUBE '(' expr_list ')'
	{
		gateFeat(yylex, feature.FGroupingSets, "GROUPING CUBE requires OceanBase 4.4+")
		$$ = &ast.GroupingSetExpr{Kind: "CUBE", Exprs: asExprList($3)}
	}
	| ROLLUP '(' expr_list ')'
	{
		gateFeat(yylex, feature.FGroupingSets, "GROUPING ROLLUP requires OceanBase 4.4+")
		$$ = &ast.GroupingSetExpr{Kind: "ROLLUP", Exprs: asExprList($3)}
	}

grouping_sets:
	grouping_set
	{
		$$ = [][]ast.Expression{asExprList($1)}
	}
	| grouping_sets ',' grouping_set
	{
		$$ = append(asExprSets($1), asExprList($3))
	}

grouping_set:
	'(' expr_list ')'
	{
		$$ = $2
	}
	| '(' ')'
	{
		$$ = nil
	}

opt_match_modifier:
	/*empty*/
	{
		$$ = ""
	}
	| IN BOOLEAN MODE
	{
		$$ = "IN BOOLEAN MODE"
	}
	| WITH QUERY EXPANSION
	{
		$$ = "WITH QUERY EXPANSION"
	}

having_opt:
	/*empty*/
	{
		$$ = nil
	}
	| HAVING expr
	{
		$$ = $2
	}

window_opt:
	/*empty*/
	{
		$$ = nil
	}
	| WINDOW window_list
	{
		gateFeat(yylex, feature.FWithWindows, "named WINDOW clause requires OceanBase 4.2+")
		$$ = $2
	}

window_list:
	named_window
	{
		$$ = []*ast.NamedWindow{$1.(*ast.NamedWindow)}
	}
	| window_list ',' named_window
	{
		$$ = append(asNamedWindows($1), $3.(*ast.NamedWindow))
	}

named_window:
	NAME_OB AS window_spec
	{
		$$ = &ast.NamedWindow{Name: strings.ToLower(asIdentName($1)), Spec: $3.(*ast.WindowSpec)}
	}

for_opt:
	/*empty*/
	{
		$$ = lockOpt{}
	}
	| FOR UPDATE
	{
		$$ = lockOpt{update: true}
	}
	| FOR UPDATE NOWAIT
	{
		$$ = lockOpt{update: true, wait: "NOWAIT"}
	}
	| FOR UPDATE SKIP LOCKED
	{
		$$ = lockOpt{update: true, wait: "SKIP LOCKED"}
	}
	| FOR SHARE
	{
		$$ = lockOpt{update: true, wait: "SHARE"}
	}
	| LOCK IN SHARE MODE
	{
		$$ = lockOpt{share: true}
	}

/* ---------------- ORDER / LIMIT ---------------- */

opt_order:
	/*empty*/
	{
		$$ = nil
	}
	| ORDER BY order_list
	{
		$$ = $3
	}

order_list:
	order_item
	{
		$$ = []ast.OrderItem{$1.(ast.OrderItem)}
	}
	| order_list ',' order_item
	{
		$$ = append(asOrderList($1), $3.(ast.OrderItem))
	}

order_item:
	expr opt_asc_desc opt_nulls
	{
		o := ast.OrderItem{Expr: asExpr($1), Desc: asBool($2)}
		n := $3.(nullsOpt)
		o.HasNulls = n.has
		o.NullsFirst = n.first
		$$ = o
	}

opt_asc_desc:
	/*empty*/
	{
		$$ = false
	}
	| ASC
	{
		$$ = false
	}
	| DESC
	{
		$$ = true
	}

opt_nulls:
	/*empty*/
	{
		$$ = nullsOpt{}
	}
	| NULLS FIRST
	{
		$$ = nullsOpt{has: true, first: true}
	}
	| NULLS LAST
	{
		$$ = nullsOpt{has: true, first: false}
	}

opt_limit:
	/*empty*/
	{
		$$ = nil
	}
	| limit_clause
	{
		$$ = $1
	}

limit_clause:
	LIMIT expr
	{
		$$ = &ast.LimitClause{Count: asExpr($2)}
	}
	| LIMIT expr ',' expr
	{
		$$ = &ast.LimitClause{Offset: asExpr($2), Count: asExpr($4), HasOffset: true}
	}
	| LIMIT expr OFFSET expr
	{
		$$ = &ast.LimitClause{Count: asExpr($2), Offset: asExpr($4), HasOffset: true}
	}

/* ---------------- expressions ---------------- */

expr:
	or_expr

or_expr:
	xor_expr
	| or_expr OR xor_expr
	{
		$$ = bin("OR", asExpr($1), asExpr($3))
	}
	| or_expr OROR xor_expr
	{
		$$ = bin("||", asExpr($1), asExpr($3))
	}

xor_expr:
	and_expr
	| xor_expr XOR and_expr
	{
		$$ = bin("XOR", asExpr($1), asExpr($3))
	}

and_expr:
	not_expr
	| and_expr AND not_expr
	{
		$$ = bin("AND", asExpr($1), asExpr($3))
	}
	| and_expr ANDAND not_expr
	{
		$$ = bin("&&", asExpr($1), asExpr($3))
	}

not_expr:
	bool_pri
	| NOT not_expr
	{
		$$ = un("NOT", asExpr($2))
	}

bool_pri:
	predicate
	| bool_pri comp_op predicate
	{
		$$ = bin(asStr($2), asExpr($1), asExpr($3))
	}
	| bool_pri IS NULLX
	{
		$$ = &ast.IsExpr{Expr: asExpr($1), Truth: "null"}
	}
	| bool_pri IS NOT NULLX
	{
		$$ = &ast.IsExpr{Expr: asExpr($1), Not: true, Truth: "null"}
	}
	| bool_pri IS BOOL_VALUE
	{
		$$ = &ast.IsExpr{Expr: asExpr($1), Truth: asLitValue($3)}
	}
	| bool_pri IS NOT BOOL_VALUE
	{
		$$ = &ast.IsExpr{Expr: asExpr($1), Not: true, Truth: asLitValue($4)}
	}

predicate:
	bit_expr
	| bit_expr IN in_expr
	{
		$$ = mkIn(asExpr($1), $3, false)
	}
	| bit_expr NOT IN in_expr
	{
		$$ = mkIn(asExpr($1), $4, true)
	}
	| bit_expr BETWEEN bit_expr AND predicate
	{
		$$ = &ast.BetweenExpr{Expr: asExpr($1), Low: asExpr($3), High: asExpr($5)}
	}
	| bit_expr NOT BETWEEN bit_expr AND predicate
	{
		$$ = &ast.BetweenExpr{Expr: asExpr($1), Low: asExpr($4), High: asExpr($6), Not: true}
	}
	| bit_expr LIKE bit_expr
	{
		$$ = &ast.LikeExpr{Expr: asExpr($1), Pattern: asExpr($3)}
	}
	| bit_expr LIKE bit_expr ESCAPE bit_expr
	{
		$$ = &ast.LikeExpr{Expr: asExpr($1), Pattern: asExpr($3), Escape: asExpr($5)}
	}
	| bit_expr NOT LIKE bit_expr
	{
		$$ = &ast.LikeExpr{Expr: asExpr($1), Pattern: asExpr($4), Not: true}
	}
	| bit_expr NOT LIKE bit_expr ESCAPE bit_expr
	{
		$$ = &ast.LikeExpr{Expr: asExpr($1), Pattern: asExpr($4), Escape: asExpr($6), Not: true}
	}
	| bit_expr REGEXP bit_expr
	{
		$$ = bin("REGEXP", asExpr($1), asExpr($3))
	}
	| bit_expr NOT REGEXP bit_expr
	{
		$$ = un("NOT", bin("REGEXP", asExpr($1), asExpr($4)))
	}
	| bit_expr RLIKE bit_expr
	{
		$$ = bin("RLIKE", asExpr($1), asExpr($3))
	}
	| bit_expr NOT RLIKE bit_expr
	{
		$$ = un("NOT", bin("RLIKE", asExpr($1), asExpr($4)))
	}
	| bit_expr SOUNDS LIKE bit_expr
	{
		$$ = bin("SOUNDS LIKE", asExpr($1), asExpr($4))
	}
	| bit_expr MEMBER OF '(' expr_list ')'
	{
		$$ = &ast.MemberOfExpr{Expr: asExpr($1), List: asExprList($5)}
	}

in_expr:
	'(' expr_list ')'
	{
		$$ = inList{list: asExprList($2)}
	}
	| '(' select_stmt ')'
	{
		$$ = inList{sub: asSelect($2)}
	}
	| '(' ')'
	{
		$$ = inList{}
	}

/* in_expr alternatives intentionally start with '(' (untyped); actions set $$. */

comp_op:
	COMP_EQ
	{
		$$ = "="
	}
	| COMP_NE
	{
		$$ = "<>"
	}
	| COMP_LT
	{
		$$ = "<"
	}
	| COMP_GT
	{
		$$ = ">"
	}
	| COMP_LE
	{
		$$ = "<="
	}
	| COMP_GE
	{
		$$ = ">="
	}
	| NSEQ
	{
		$$ = "<=>"
	}

bit_expr:
	bit_or

bit_or:
	bit_and
	| bit_or '|' bit_and
	{
		$$ = bin("|", asExpr($1), asExpr($3))
	}

bit_and:
	shift_expr
	| bit_and '&' shift_expr
	{
		$$ = bin("&", asExpr($1), asExpr($3))
	}

shift_expr:
	additive
	| shift_expr SHIFT_LEFT additive
	{
		$$ = bin("<<", asExpr($1), asExpr($3))
	}
	| shift_expr SHIFT_RIGHT additive
	{
		$$ = bin(">>", asExpr($1), asExpr($3))
	}

additive:
	multiplicative
	| additive '+' multiplicative
	{
		$$ = bin("+", asExpr($1), asExpr($3))
	}
	| additive '-' multiplicative
	{
		$$ = bin("-", asExpr($1), asExpr($3))
	}

multiplicative:
	power
	| multiplicative '*' power
	{
		$$ = bin("*", asExpr($1), asExpr($3))
	}
	| multiplicative '/' power
	{
		$$ = bin("/", asExpr($1), asExpr($3))
	}
	| multiplicative '%' power
	{
		$$ = bin("%", asExpr($1), asExpr($3))
	}
	| multiplicative DIV power
	{
		$$ = bin("DIV", asExpr($1), asExpr($3))
	}
	| multiplicative MOD power
	{
		$$ = bin("MOD", asExpr($1), asExpr($3))
	}

power:
	unary
	| power '^' unary
	{
		$$ = bin("^", asExpr($1), asExpr($3))
	}

unary:
	simple_expr
	| unary COLLATE NAME_OB
	{
		$$ = bin("COLLATE", asExpr($1), newIdent(asIdentName($3)))
	}
	| '-' unary
	{
		$$ = un("-", asExpr($2))
	}
	| '+' unary
	{
		$$ = un("+", asExpr($2))
	}
	| '~' unary
	{
		$$ = un("~", asExpr($2))
	}

simple_expr:
	literal
	| column_ref
	| func_call
	| '(' expr ')'
	{
		$$ = &ast.ParenExpr{Expr: asExpr($2)}
	}
	| '(' select_stmt ')'
	{
		$$ = &ast.SubqueryExpr{Query: asSelect($2)}
	}
	| EXISTS '(' select_stmt ')'
	{
		$$ = &ast.ExistsExpr{Query: asSelect($3)}
	}
	| MATCH '(' ident_list ')' AGAINST '(' bit_expr opt_match_modifier ')'
	{
		$$ = &ast.MatchAgainstExpr{Columns: asIdentList($3), Against: asExpr($7), Modifier: asStr($8)}
	}
	| CONVERT '(' expr ',' data_type ')'
	{
		d := asDType($5)
		name := d.name
		if d.unsigned {
			name += " UNSIGNED"
		}
		$$ = &ast.CastExpr{Expr: asExpr($3), TypeName: name, TypeArgs: d.args, IsConvert: true}
	}
	| case_expr
	| cast_expr
	| interval_expr

literal:
	INTNUM
	| DECIMAL_VAL
	| APPROXNUM
	| STRING_VALUE
	| BOOL_VALUE
	| NULLX
	| HEXNUM
	| BITNUM
	| HEX_STRING_VALUE
	| BIN_STRING_VALUE
	| QUESTIONMARK
	| USER_VARIABLE
	| SYSTEM_VARIABLE

column_ref:
	NAME_OB
	{
		$$ = newIdent(asIdentName($1))
	}
	| NAME_OB '.' NAME_OB
	{
		$$ = idParts(asIdentName($1), asIdentName($3))
	}
	| NAME_OB '.' NAME_OB '.' NAME_OB
	{
		$$ = idParts(asIdentName($1), asIdentName($3), asIdentName($5))
	}
	| NAME_OB '.' '*'
	{
		$$ = &ast.StarLiteral{Qualifier: asIdentName($1)}
	}
	| '*'
	{
		$$ = &ast.StarLiteral{}
	}

func_call:
	NAME_OB '(' call_args ')' over_opt
	{
		ca := $3.(callArgs)
		fc := &ast.FuncCall{Name: strings.ToLower(asIdentName($1)), Distinct: ca.distinct, Star: ca.star, Args: ca.args}
		normalizeFuncStar(fc)
		if $5 != nil {
			fc.Over = $5.(*ast.WindowSpec)
		}
		$$ = fc
	}
	| VALUES '(' expr_list ')' over_opt
	{
		fc := &ast.FuncCall{Name: "values", Args: asExprList($3)}
		if $5 != nil {
			fc.Over = $5.(*ast.WindowSpec)
		}
		$$ = fc
	}

call_args:
	/*empty*/
	{
		$$ = callArgs{}
	}
	| DISTINCT expr_list
	{
		$$ = callArgs{distinct: true, args: asExprList($2)}
	}
	| expr_list
	{
		$$ = callArgs{args: asExprList($1)}
	}

over_opt:
	/*empty*/
	{
		$$ = nil
	}
	| OVER window_spec
	{
		$$ = $2
	}

window_spec:
	'(' window_part_opt window_order_opt window_frame_opt ')'
	{
		ws := &ast.WindowSpec{}
		ws.PartitionBy = asExprList($2)
		ws.OrderBy = asOrderList($3)
		fs := $4.(frameSpec)
		ws.Frame = fs.frame
		ws.FrameStart = fs.start
		ws.FrameEnd = fs.end
		$$ = ws
	}

window_part_opt:
	/*empty*/
	{
		$$ = nil
	}
	| PARTITION BY expr_list
	{
		$$ = $3
	}

window_order_opt:
	/*empty*/
	{
		$$ = nil
	}
	| ORDER BY order_list
	{
		$$ = $3
	}

window_frame_opt:
	/*empty*/
	{
		$$ = frameSpec{}
	}
	| ROWS frame_extent
	{
		gateFeat(yylex, feature.FRowWindowFrame, "window frame (ROWS) requires OceanBase 4.0+")
		fs := $2.(frameSpec)
		fs.frame = "ROWS"
		$$ = fs
	}
	| RANGE frame_extent
	{
		gateFeat(yylex, feature.FRowWindowFrame, "window frame (RANGE) requires OceanBase 4.0+")
		fs := $2.(frameSpec)
		fs.frame = "RANGE"
		$$ = fs
	}

frame_extent:
	frame_bound
	{
		$$ = frameSpec{start: asStr($1)}
	}
	| BETWEEN frame_bound AND frame_bound
	{
		$$ = frameSpec{start: asStr($2), end: asStr($4)}
	}

frame_bound:
	CURRENT ROW
	{
		$$ = "CURRENT ROW"
	}
	| UNBOUNDED PRECEDING
	{
		$$ = "UNBOUNDED PRECEDING"
	}
	| UNBOUNDED FOLLOWING
	{
		$$ = "UNBOUNDED FOLLOWING"
	}
	| expr PRECEDING
	{
		$$ = exprText(asExpr($1)) + " PRECEDING"
	}
	| expr FOLLOWING
	{
		$$ = exprText(asExpr($1)) + " FOLLOWING"
	}

case_expr:
	CASE case_operand when_list else_opt END
	{
		$$ = &ast.CaseExpr{Operand: asExpr($2), Whens: asCaseWhens($3), Else: asExpr($4)}
	}

case_operand:
	/*empty*/
	{
		$$ = nil
	}
	| expr

when_list:
	WHEN expr THEN expr
	{
		$$ = []ast.CaseClause{{When: asExpr($2), Then: asExpr($4)}}
	}
	| when_list WHEN expr THEN expr
	{
		$$ = append(asCaseWhens($1), ast.CaseClause{When: asExpr($3), Then: asExpr($5)})
	}

else_opt:
	/*empty*/
	{
		$$ = nil
	}
	| ELSE expr
	{
		$$ = $2
	}

cast_expr:
	CAST '(' expr AS data_type ')'
	{
		d := asDType($5)
		name := d.name
		if d.unsigned {
			name += " UNSIGNED"
		}
		$$ = &ast.CastExpr{Expr: asExpr($3), TypeName: name, TypeArgs: d.args}
	}

interval_expr:
	INTERVAL expr date_unit
	{
		$$ = &ast.IntervalExpr{Value: asExpr($2), Unit: asStr($3)}
	}

date_unit:
	NAME_OB
	{
		$$ = strings.ToLower(asIdentName($1))
	}

expr_list:
	expr
	{
		$$ = []ast.Expression{asExpr($1)}
	}
	| expr_list ',' expr
	{
		$$ = append(asExprList($1), asExpr($3))
	}

ident_list:
	NAME_OB
	{
		$$ = []*ast.Identifier{asIdent($1)}
	}
	| ident_list ',' NAME_OB
	{
		$$ = append(asIdentList($1), asIdent($3))
	}

data_type:
	NAME_OB type_args_opt unsigned_opt
	{
		$$ = dtype{name: asIdentName($1), args: asExprList($2), unsigned: asBool($3)}
	}

type_args_opt:
	/*empty*/
	{
		$$ = nil
	}
	| '(' expr_list ')'
	{
		$$ = $2
	}

unsigned_opt:
	/*empty*/
	{
		$$ = false
	}
	| UNSIGNED
	{
		$$ = true
	}

/* ---------------- INSERT ---------------- */

insert_stmt:
	INSERT ignore_opt opt_into table_name col_list_opt insert_source on_dup_opt returning_opt
	{
		ins := &ast.InsertStmt{Ignore: asBool($2), Table: asIdent($4), Columns: asIdentList($5)}
		src := $6.(insertSource)
		ins.Rows = src.rows
		ins.Set = src.set
		ins.Source = src.sel
		ins.OnDuplicate = asAssignList($7)
		ins.Returning = asIdentList($8)
		$$ = ins
	}
	| REPLACE ignore_opt opt_into table_name col_list_opt insert_source
	{
		ins := &ast.InsertStmt{Replace: true, Ignore: asBool($2), Table: asIdent($4), Columns: asIdentList($5)}
		src := $6.(insertSource)
		ins.Rows = src.rows
		ins.Set = src.set
		ins.Source = src.sel
		$$ = ins
	}

opt_into:
	/*empty*/
	{
		$$ = nil
	}
	| INTO
	{
		$$ = nil
	}

ignore_opt:
	/*empty*/
	{
		$$ = false
	}
	| IGNORE
	{
		$$ = true
	}

col_list_opt:
	/*empty*/
	{
		$$ = nil
	}
	| '(' ident_list ')'
	{
		$$ = $2
	}

insert_source:
	VALUES row_list
	{
		$$ = insertSource{rows: asRows($2)}
	}
	| VALUE row_list
	{
		$$ = insertSource{rows: asRows($2)}
	}
	| SET assignment_list
	{
		$$ = insertSource{set: asAssignList($2)}
	}
	| select_stmt
	{
		$$ = insertSource{sel: asSelect($1)}
	}

row_list:
	row
	{
		$$ = [][]ast.Expression{asExprList($1)}
	}
	| row_list ',' row
	{
		$$ = append(asRows($1), asExprList($3))
	}

row:
	'(' expr_list ')'
	{
		$$ = $2
	}

on_dup_opt:
	/*empty*/
	{
		$$ = nil
	}
	| ON DUPLICATE KEY UPDATE assignment_list
	{
		$$ = $5
	}

assignment_list:
	assignment
	{
		$$ = []ast.Assignment{$1.(ast.Assignment)}
	}
	| assignment_list ',' assignment
	{
		$$ = append(asAssignList($1), $3.(ast.Assignment))
	}

assignment:
	column_ref COMP_EQ expr
	{
		$$ = ast.Assignment{Column: asIdent($1), Expr: asExpr($3)}
	}

returning_opt:
	/*empty*/
	{
		$$ = nil
	}
	| RETURNING ident_list
	{
		$$ = $2
	}

/* ---------------- UPDATE / DELETE ---------------- */

update_stmt:
	UPDATE ignore_opt table_ref_list SET assignment_list where_opt opt_order opt_limit
	{
		up := &ast.UpdateStmt{Ignore: asBool($2), Tables: asTableRefList($3), Set: asAssignList($5)}
		up.Where = asExpr($6)
		up.OrderBy = asOrderList($7)
		up.Limit = asLimit($8)
		$$ = up
	}

table_ref_list:
	table_ref
	{
		$$ = []*ast.TableRef{asTableRef($1)}
	}
	| table_ref_list ',' table_ref
	{
		$$ = append(asTableRefList($1), asTableRef($3))
	}

delete_stmt:
	DELETE delete_mods FROM table_ref_list using_opt where_opt opt_order opt_limit
	{
		d := &ast.DeleteStmt{From: asTableRefList($4)}
		d.Using = asTableRefList($5)
		d.Where = asExpr($6)
		d.OrderBy = asOrderList($7)
		d.Limit = asLimit($8)
		$$ = d
	}

delete_mods:
	/*empty*/
	{
		$$ = nil
	}
	| LOW_PRIORITY
	{
		$$ = nil
	}
	| QUICK
	{
		$$ = nil
	}
	| IGNORE
	{
		$$ = nil
	}

using_opt:
	/*empty*/
	{
		$$ = nil
	}
	| USING table_ref_list
	{
		$$ = $2
	}

/* ---------------- CREATE ---------------- */

create_stmt:
	CREATE TABLE opt_if_not_exists table_name create_table_tail
	{
		body := asTableBody($5)
		if body.isLike {
			$$ = &ast.CreateTableLikeStmt{Table: asIdent($4), Like: body.like}
		} else {
			$$ = &ast.CreateTableStmt{IfNotExists: asBool($3), Table: asIdent($4), Columns: body.defs.cols, Constraints: body.defs.cons, Options: body.opts, Partition: body.part, Select: body.sel}
		}
	}
	| CREATE TEMPORARY TABLE opt_if_not_exists table_name create_table_tail
	{
		body := asTableBody($6)
		if body.isLike {
			$$ = &ast.CreateTableLikeStmt{Table: asIdent($5), Like: body.like}
		} else {
			$$ = &ast.CreateTableStmt{Temp: true, IfNotExists: asBool($4), Table: asIdent($5), Columns: body.defs.cols, Constraints: body.defs.cons, Options: body.opts, Partition: body.part, Select: body.sel}
		}
	}
	| CREATE opt_unique INDEX table_name ON table_name '(' index_col_list ')'
	{
		$$ = &ast.CreateIndexStmt{Unique: asBool($2), Name: asIdent($4), Table: asIdent($6), Columns: asIdentList($8)}
	}
	| CREATE DATABASE opt_if_not_exists NAME_OB db_opts_opt
	{
		o := asDBOpts($5)
		$$ = &ast.CreateDatabaseStmt{IfNotExists: asBool($4), Name: asIdentName($4), Charset: o.charset, Collate: o.collate}
	}
	| CREATE SCHEMA opt_if_not_exists NAME_OB db_opts_opt
	{
		o := asDBOpts($5)
		$$ = &ast.CreateDatabaseStmt{IfNotExists: asBool($4), Name: asIdentName($4), Charset: o.charset, Collate: o.collate}
	}
	| CREATE SEQUENCE opt_if_not_exists NAME_OB sequence_options
	{
		so := asSeqOpts($5)
		$$ = &ast.CreateSequenceStmt{IfNotExists: asBool($3), Name: asIdentName($4),
			StartWith: so.start, IncrementBy: so.incr, MinValue: so.min, MaxValue: so.max,
			Cache: so.cache, Cycle: so.cycle, Order: so.order}
	}
	| CREATE USER opt_if_not_exists create_user_list
	{
		$$ = &ast.CreateUserStmt{IfNotExists: asBool($3), Users: asUserSpecs($4)}
	}
	| CREATE ROLE opt_if_not_exists role_name_list
	{
		gateFeat(yylex, feature.FRoleManagement, "ROLE management requires OceanBase 4.2.4+")

		$$ = &ast.CreateRoleStmt{IfNotExists: asBool($3), Roles: asStringList($4)}
	}
	| CREATE MATERIALIZED VIEW table_name col_list_opt mview_opts AS select_stmt
	{
		gateFeat(yylex, feature.FMaterializedView, "materialized view requires OceanBase 4.3+")

		o := asMViewOpts($6)
		$$ = &ast.CreateMViewStmt{Name: asIdent($4), Columns: asIdentList($5), Refresh: o.refresh, Mode: o.mode, Select: asSelect($8)}
	}
	| CREATE TENANT opt_if_not_exists NAME_OB opt_cluster_options
	{
		$$ = &ast.CreateTenantStmt{IfNotExists: asBool($3), Name: asIdentName($4), Options: asOptions($5)}
	}
	| CREATE RESOURCE UNIT opt_if_not_exists NAME_OB opt_cluster_options
	{
		$$ = &ast.CreateResourceUnitStmt{IfNotExists: asBool($4), Name: asIdentName($5), Options: asOptions($6)}
	}
	| CREATE RESOURCE POOL opt_if_not_exists NAME_OB opt_cluster_options
	{
		$$ = &ast.CreateResourcePoolStmt{IfNotExists: asBool($4), Name: asIdentName($5), Options: asOptions($6)}
	}
	| CREATE TABLEGROUP opt_if_not_exists NAME_OB opt_cluster_options
	{
		$$ = &ast.CreateTablegroupStmt{IfNotExists: asBool($3), Name: asIdentName($4), Options: asOptions($5)}
	}
	| CREATE TABLESPACE NAME_OB opt_cluster_options
	{
		$$ = &ast.CreateTablespaceStmt{Name: asIdentName($3), Options: asOptions($4)}
	}
	| CREATE STANDBY TENANT opt_if_not_exists NAME_OB opt_log_restore_source opt_cluster_options
	{
		$$ = &ast.CreateStandbyTenantStmt{IfNotExists: asBool($4), Name: asIdentName($5), RestoreSource: asOptPtr($6), Options: asOptions($7)}
	}
	| CREATE TENANT opt_if_not_exists NAME_OB FROM NAME_OB opt_clone_snapshot WITH cluster_option_list
	{
		gateFeat(yylex, feature.FCloneTenant, "CLONE TENANT requires OceanBase 4.3+")

		$$ = &ast.CloneTenantStmt{IfNotExists: asBool($3), Name: asIdentName($4), From: asIdentName($6), Snapshot: asStr($7), Options: asOptions($9)}
	}
	| CREATE RESTORE POINT NAME_OB
	{
		$$ = &ast.RestorePointStmt{Action: "CREATE", Name: asIdentName($4)}
	}
	| CREATE SNAPSHOT NAME_OB opt_for_tenant
	{
		gateFeat(yylex, feature.FTenantSnapshot, "tenant snapshot requires OceanBase 4.3+")

		$$ = &ast.TenantSnapshotStmt{Action: "CREATE", Name: asIdentName($3), Tenant: asStr($4), HasTenant: asStr($4) != ""}
	}
	| CREATE opt_or_replace LOCATION NAME_OB URL opt_equal STRING_VALUE
	{
		gateFeat(yylex, feature.FExternalLocation, "LOCATION requires OceanBase 4.3+")

		$$ = &ast.CreateLocationStmt{Action: "CREATE", OrReplace: asBool($2), Name: asIdentName($4), URL: asLitValue($7)}
	}
	| CREATE EXTERNAL CATALOG opt_if_not_exists NAME_OB PROPERTIES '(' load_property_list ')'
	{
		gateFeat(yylex, feature.FExternalCatalog, "EXTERNAL CATALOG requires OceanBase 4.3+")

		$$ = &ast.CatalogStmt{Action: "CREATE", IfNotExists: asBool($4), Name: asIdentName($5), Options: asOptions($8)}
	}
	| CREATE opt_or_replace PROCEDURE NAME_OB proc_params proc_body
	{
		$$ = &ast.StoredProcStmt{OrReplace: asBool($2), Name: asIdentName($4), Parameters: asProcParams($5), Body: asPLBlock($6)}
	}
	| CREATE opt_or_replace FUNCTION NAME_OB proc_params RETURNS data_type proc_body
	{
		d := asDType($7)
		$$ = &ast.StoredProcStmt{Function: true, OrReplace: asBool($2), Name: asIdentName($4), Parameters: asProcParams($5),
			Returns: d.name, ReturnsArgs: d.args, Body: asPLBlock($8)}
	}
	| CREATE ROUTINE LOAD NAME_OB ON table_name opt_load_properties FROM STREAM '(' opt_job_properties ')'
	{
		gateFeat(yylex, feature.FRoutineLoad, "ROUTINE LOAD requires OceanBase 4.3+")

		$$ = &ast.RoutineLoadStmt{Job: asIdentName($4), Table: asIdent($6), LoadProps: asOptions($7), JobProps: asOptions($11)}
	}
	| CREATE CONCURRENT_LIMITING_RULE opt_if_not_exists NAME_OB ON table_name TO user_spec FOR NAME_OB opt_cluster_options
	{
		gateFeat(yylex, feature.FCCLRule, "CONCURRENT_LIMITING_RULE requires OceanBase 4.3+")

		$$ = &ast.CCLLimitRuleStmt{Action: "CREATE", Name: asIdentName($4), Table: asIdent($6), User: $8.(ast.UserSpec), HasUser: true, ForName: asIdentName($10), Options: asOptions($11)}
	}
	| CREATE SENSITIVE RULE NAME_OB ON ident_list USING NAME_OB
	{
		gateFeat(yylex, feature.FSensitiveRule, "SENSITIVE RULE requires OceanBase 4.4+")

		$$ = &ast.SensitiveRuleStmt{Action: "CREATE", Name: asIdentName($4), Fields: idNamesList(asIdentList($6)), Using: asIdentName($8)}
	}
	| CREATE DATABASE LINK NAME_OB CONNECT TO STRING_VALUE DATABASE STRING_VALUE IDENTIFIED BY STRING_VALUE opt_dblink_tail
	{
		dt := $13.(dblinkTail)
		$$ = &ast.DblinkStmt{Action: "CREATE", Name: asIdentName($4), User: asLitValue($7), Database: asLitValue($9),
			Password: asLitValue($12), Host: dt.host, Port: dt.port, Cluster: dt.cluster}
	}

opt_dblink_tail:
	/*empty*/
	{
		$$ = dblinkTail{}
	}
	| HOST STRING_VALUE PORT expr
	{
		$$ = dblinkTail{host: asLitValue($2), port: asExpr($4)}
	}
	| HOST STRING_VALUE PORT expr CLUSTER NAME_OB
	{
		$$ = dblinkTail{host: asLitValue($2), port: asExpr($4), cluster: asIdentName($6)}
	}

opt_load_properties:
	/*empty*/
	{
		$$ = nil
	}
	| load_property_list

opt_job_properties:
	/*empty*/
	{
		$$ = nil
	}
	| load_property_list

load_property_list:
	load_property
	{
		$$ = []ast.Option{$1.(ast.Option)}
	}
	| load_property_list ',' load_property
	{
		$$ = append(asOptions($1), $3.(ast.Option))
	}

load_property:
	NAME_OB COMP_EQ cluster_option_value
	{
		cv := $3.(optionValue)
		$$ = ast.Option{Name: strings.ToLower(asIdentName($1)), Value: cv.value, Values: cv.values}
	}
	| NAME_OB '(' expr_list ')'
	{
		$$ = ast.Option{Name: strings.ToLower(asIdentName($1)), Values: asExprList($3)}
	}
	| COLUMNS '(' expr_list ')'
	{
		$$ = ast.Option{Name: "columns", Values: asExprList($3)}
	}
	| FIELDS '(' expr_list ')'
	{
		$$ = ast.Option{Name: "fields", Values: asExprList($3)}
	}
	| STRING_VALUE COMP_EQ STRING_VALUE
	{
		$$ = ast.Option{Name: asLitValue($1), Value: asExpr($3)}
	}

routine_load_action_stmt:
	PAUSE ROUTINE LOAD FOR NAME_OB
	{
		gateFeat(yylex, feature.FRoutineLoad, "ROUTINE LOAD requires OceanBase 4.3+")

		$$ = &ast.RoutineLoadActionStmt{Action: "PAUSE", Job: asIdentName($5)}
	}
	| RESUME ROUTINE LOAD FOR NAME_OB
	{
		gateFeat(yylex, feature.FRoutineLoad, "ROUTINE LOAD requires OceanBase 4.3+")

		$$ = &ast.RoutineLoadActionStmt{Action: "RESUME", Job: asIdentName($5)}
	}
	| STOP ROUTINE LOAD FOR NAME_OB
	{
		gateFeat(yylex, feature.FRoutineLoad, "ROUTINE LOAD requires OceanBase 4.3+")

		$$ = &ast.RoutineLoadActionStmt{Action: "STOP", Job: asIdentName($5)}
	}

opt_or_replace:
	/*empty*/
	{
		$$ = false
	}
	| OR REPLACE
	{
		$$ = true
	}

opt_equal:
	/*empty*/
	{
		$$ = nil
	}
	| COMP_EQ
	{
		$$ = nil
	}

opt_log_restore_source:
	/*empty*/
	{
		$$ = nil
	}
	| NAME_OB COMP_EQ cluster_option_value
	{
		cv := $3.(optionValue)
		$$ = &ast.Option{Name: strings.ToLower(asIdentName($1)), Value: cv.value, Values: cv.values}
	}

opt_clone_snapshot:
	/*empty*/
	{
		$$ = ""
	}
	| USING SNAPSHOT NAME_OB
	{
		$$ = asIdentName($3)
	}

opt_for_tenant:
	/*empty*/
	{
		$$ = ""
	}
	| FOR TENANT NAME_OB
	{
		$$ = asIdentName($3)
	}

opt_cluster_options:
	/*empty*/
	{
		$$ = nil
	}
	| cluster_option_list

cluster_option_list:
	cluster_option
	{
		$$ = []ast.Option{$1.(ast.Option)}
	}
	| cluster_option_list ',' cluster_option
	{
		$$ = append(asOptions($1), $3.(ast.Option))
	}

cluster_option:
	cluster_opt_name COMP_EQ cluster_option_value
	{
		cv := $3.(optionValue)
		$$ = ast.Option{Name: strings.ToLower(asStr($1)), Value: cv.value, Values: cv.values}
	}

cluster_opt_name:
	NAME_OB
	{
		$$ = asIdentName($1)
	}
	| UNIT
	{
		$$ = "unit"
	}
	| CHARSET
	{
		$$ = "charset"
	}
	| COLLATE
	{
		$$ = "collate"
	}
	| COMMENT
	{
		$$ = "comment"
	}
	| PRIMARY
	{
		$$ = "primary"
	}

cluster_option_value:
	expr
	{
		$$ = optionValue{value: asExpr($1)}
	}
	| '(' expr_list ')'
	{
		$$ = optionValue{values: asExprList($2)}
	}

mview_opts:
	/*empty*/
	{
		$$ = mviewOpts{}
	}
	| mview_opts mview_opt
	{
		a := asMViewOpts($1)
		b := asMViewOpts($2)
		if b.refresh != "" {
			a.refresh = b.refresh
		}
		if b.mode != "" {
			a.mode = b.mode
		}
		$$ = a
	}

mview_opt:
	REFRESH FAST
	{
		$$ = mviewOpts{refresh: "FAST"}
	}
	| REFRESH COMPLETE
	{
		$$ = mviewOpts{refresh: "COMPLETE"}
	}
	| REFRESH FORCE
	{
		$$ = mviewOpts{refresh: "FORCE"}
	}
	| ON DEMAND
	{
		$$ = mviewOpts{mode: "ON DEMAND"}
	}
	| ON COMMIT
	{
		$$ = mviewOpts{mode: "ON COMMIT"}
	}

role_name_list:
	user_name
	{
		$$ = []string{asStr($1)}
	}
	| role_name_list ',' user_name
	{
		$$ = append(asStringList($1), asStr($3))
	}

create_user_list:
	create_user_spec
	{
		$$ = []ast.UserSpec{$1.(ast.UserSpec)}
	}
	| create_user_list ',' create_user_spec
	{
		$$ = append(asUserSpecs($1), $3.(ast.UserSpec))
	}

create_user_spec:
	user_name AT user_host opt_identified
	{
		us := ast.UserSpec{Name: asStr($1), Host: asStr($3), HasHost: true}
		if $4 != nil {
			ia := $4.(identAuth)
			us.Password = ia.pwd
			us.HasPassword = ia.pwd != ""
			us.Plugin = ia.plugin
		}
		$$ = us
	}
	| user_name opt_identified
	{
		us := ast.UserSpec{Name: asStr($1)}
		if $2 != nil {
			ia := $2.(identAuth)
			us.Password = ia.pwd
			us.HasPassword = ia.pwd != ""
			us.Plugin = ia.plugin
		}
		$$ = us
	}

opt_identified:
	/*empty*/
	{
		$$ = nil
	}
	| IDENTIFIED BY STRING_VALUE
	{
		$$ = identAuth{pwd: asLitValue($3)}
	}
	| IDENTIFIED WITH NAME_OB
	{
		$$ = identAuth{plugin: asIdentName($3)}
	}
	| IDENTIFIED WITH NAME_OB BY STRING_VALUE
	{
		$$ = identAuth{plugin: asIdentName($3), pwd: asLitValue($5)}
	}

sequence_options:
	/*empty*/
	{
		$$ = seqOpts{}
	}
	| sequence_options sequence_option
	{
		$$ = mergeSeqOpts($1, $2)
	}

sequence_option:
	START WITH expr
	{
		$$ = seqOpts{start: asExpr($3)}
	}
	| INCREMENT BY expr
	{
		$$ = seqOpts{incr: asExpr($3)}
	}
	| MINVALUE expr
	{
		$$ = seqOpts{min: asExpr($2)}
	}
	| MAXVALUE expr
	{
		$$ = seqOpts{max: asExpr($2)}
	}
	| CACHE expr
	{
		$$ = seqOpts{cache: asExpr($2)}
	}
	| CYCLE
	{
		$$ = seqOpts{cycle: "CYCLE"}
	}
	| NOCYCLE
	{
		$$ = seqOpts{cycle: "NOCYCLE"}
	}
	| ORDER
	{
		$$ = seqOpts{order: "ORDER"}
	}
	| NOORDER
	{
		$$ = seqOpts{order: "NOORDER"}
	}

opt_unique:
	/*empty*/
	{
		$$ = false
	}
	| UNIQUE
	{
		$$ = true
	}

create_table_tail:
	'(' create_def_list ')' table_options_opt opt_partition opt_as_select
	{
		$$ = tableBody{defs: asDefs($2), opts: asTableOptions($4), part: asPartition($5), sel: asSelect($6)}
	}
	| LIKE table_name
	{
		$$ = tableBody{isLike: true, like: asIdent($2)}
	}
	| as_select
	{
		$$ = tableBody{sel: asSelect($1)}
	}

as_select:
	AS select_stmt
	{
		$$ = $2
	}
	| select_stmt

opt_as_select:
	/*empty*/
	{
		$$ = nil
	}
	| AS select_stmt
	{
		$$ = $2
	}
	| select_stmt

opt_partition:
	/*empty*/
	{
		$$ = nil
	}
	| PARTITION BY partition_by
	{
		$$ = $3
	}

partition_by:
	RANGE '(' expr ')' opt_auto_partition opt_partitions_count opt_subpartition opt_partition_defs
	{
		ap := asAutoPart($5)
		ss := asSubSpec($7)
		$$ = &ast.PartitionSpec{Kind: ast.PartitionRange, Expr: asExpr($3), Count: asExpr($6), Auto: ap.auto, Size: ap.size,
			SubKind: ss.kind, SubExpr: ss.expr, SubColumns: ss.cols, SubCount: ss.count, Partitions: asPartitionDefs($8)}
	}
	| RANGE COLUMNS '(' ident_list ')' opt_auto_partition opt_subpartition opt_partition_defs
	{
		ap := asAutoPart($6)
		ss := asSubSpec($7)
		$$ = &ast.PartitionSpec{Kind: ast.PartitionRangeColumns, Columns: asIdentList($4), Auto: ap.auto, Size: ap.size,
			SubKind: ss.kind, SubExpr: ss.expr, SubColumns: ss.cols, SubCount: ss.count, Partitions: asPartitionDefs($8)}
	}
	| LIST '(' expr ')' opt_auto_partition opt_subpartition opt_partition_defs
	{
		ap := asAutoPart($5)
		ss := asSubSpec($6)
		$$ = &ast.PartitionSpec{Kind: ast.PartitionList, Expr: asExpr($3), Auto: ap.auto, Size: ap.size,
			SubKind: ss.kind, SubExpr: ss.expr, SubColumns: ss.cols, SubCount: ss.count, Partitions: asPartitionDefs($7)}
	}
	| LIST COLUMNS '(' ident_list ')' opt_auto_partition opt_subpartition opt_partition_defs
	{
		ap := asAutoPart($6)
		ss := asSubSpec($7)
		$$ = &ast.PartitionSpec{Kind: ast.PartitionListColumns, Columns: asIdentList($4), Auto: ap.auto, Size: ap.size,
			SubKind: ss.kind, SubExpr: ss.expr, SubColumns: ss.cols, SubCount: ss.count, Partitions: asPartitionDefs($8)}
	}
	| HASH '(' expr ')' opt_auto_partition opt_partitions_count opt_subpartition
	{
		ap := asAutoPart($5)
		ss := asSubSpec($7)
		$$ = &ast.PartitionSpec{Kind: ast.PartitionHash, Expr: asExpr($3), Count: asExpr($6), Auto: ap.auto, Size: ap.size,
			SubKind: ss.kind, SubExpr: ss.expr, SubColumns: ss.cols, SubCount: ss.count}
	}
	| KEY '(' ident_list ')' opt_auto_partition opt_partitions_count opt_subpartition
	{
		ap := asAutoPart($5)
		ss := asSubSpec($7)
		$$ = &ast.PartitionSpec{Kind: ast.PartitionKey, Columns: asIdentList($3), Count: asExpr($6), Auto: ap.auto, Size: ap.size,
			SubKind: ss.kind, SubExpr: ss.expr, SubColumns: ss.cols, SubCount: ss.count}
	}

opt_auto_partition:
	/*empty*/
	{
		$$ = autoPart{}
	}
	| PARTITION SIZE expr PARTITIONS AUTO
	{
		gateFeat(yylex, feature.FAutoPartition, "PARTITIONS AUTO requires OceanBase 4.3+")

		$$ = autoPart{auto: true, size: asExpr($3)}
	}

opt_partitions_count:
	/*empty*/
	{
		$$ = nil
	}
	| PARTITIONS expr
	{
		$$ = $2
	}
	| PARTITIONS AUTO
	{
		gateFeat(yylex, feature.FAutoPartition, "PARTITIONS AUTO requires OceanBase 4.3+")

		$$ = newIdent("AUTO")
	}

opt_subpartition:
	/*empty*/
	{
		$$ = subSpec{}
	}
	| SUBPARTITION BY HASH '(' expr ')' opt_subpartitions_count
	{
		$$ = subSpec{kind: ast.PartitionHash, expr: asExpr($5), count: asExpr($7)}
	}
	| SUBPARTITION BY KEY '(' ident_list ')' opt_subpartitions_count
	{
		$$ = subSpec{kind: ast.PartitionKey, cols: asIdentList($5), count: asExpr($7)}
	}

opt_subpartitions_count:
	/*empty*/
	{
		$$ = nil
	}
	| SUBPARTITIONS expr
	{
		$$ = $2
	}

opt_partition_defs:
	/*empty*/
	{
		$$ = nil
	}
	| '(' partition_def_list ')'
	{
		$$ = $2
	}

partition_def_list:
	partition_def
	{
		$$ = []ast.PartitionDef{$1.(ast.PartitionDef)}
	}
	| partition_def_list ',' partition_def
	{
		$$ = append(asPartitionDefs($1), $3.(ast.PartitionDef))
	}

partition_def:
	PARTITION NAME_OB opt_partition_values table_options_opt
	{
		pv := asPartitionValues($3)
		$$ = ast.PartitionDef{Name: asIdentName($2), ValueKind: pv.kind, Values: pv.values, Options: asTableOptions($4)}
	}

opt_partition_values:
	/*empty*/
	{
		$$ = partitionValues{}
	}
	| VALUES LESS THAN '(' part_value_list ')'
	{
		$$ = partitionValues{kind: "LESS THAN", values: asExprList($5)}
	}
	| VALUES IN '(' expr_list ')'
	{
		$$ = partitionValues{kind: "IN", values: asExprList($4)}
	}

part_value_list:
	part_value
	{
		$$ = []ast.Expression{asExpr($1)}
	}
	| part_value_list ',' part_value
	{
		$$ = append(asExprList($1), asExpr($3))
	}

part_value:
	expr
	| MAXVALUE
	{
		$$ = newIdent("MAXVALUE")
	}

opt_partition_clause:
	/*empty*/
	{
		$$ = nil
	}
	| PARTITION '(' name_list ')'
	{
		$$ = $3
	}

name_list:
	NAME_OB
	{
		$$ = []string{asIdentName($1)}
	}
	| name_list ',' NAME_OB
	{
		$$ = append(asStringList($1), asIdentName($3))
	}

create_def_list:
	create_def
	{
		$$ = appendDef(nil, $1)
	}
	| create_def_list ',' create_def
	{
		$$ = appendDef($1, $3)
	}

create_def:
	column_def
	| table_constraint
	| KEY NAME_OB '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "KEY", Name: asIdentName($2), Columns: asIdentList($4)}
	}
	| INDEX NAME_OB '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "INDEX", Name: asIdentName($2), Columns: asIdentList($4)}
	}
	| UNIQUE NAME_OB '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "UNIQUE", Name: asIdentName($2), Columns: asIdentList($4)}
	}
	| UNIQUE KEY NAME_OB '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "UNIQUE", Name: asIdentName($3), Columns: asIdentList($5)}
	}

column_def:
	NAME_OB data_type column_attrs_opt
	{
		$$ = buildColumnDef($1, $2, $3)
	}

column_attrs_opt:
	/*empty*/
	{
		$$ = &colAttrs{}
	}
	| column_attrs_opt column_attr
	{
		$$ = mergeAttrs($1.(*colAttrs), $2.(*colAttrs))
	}

column_attr:
	NOT NULLX
	{
		$$ = &colAttrs{nullable: ast.NullableNo}
	}
	| NULLX
	{
		$$ = &colAttrs{nullable: ast.NullableYes}
	}
	| DEFAULT expr
	{
		$$ = &colAttrs{def: asExpr($2)}
	}
	| ON UPDATE expr
	{
		$$ = &colAttrs{onUpdate: asExpr($3)}
	}
	| AUTO_INCREMENT
	{
		$$ = &colAttrs{autoIncrement: true}
	}
	| COMMENT STRING_VALUE
	{
		$$ = &colAttrs{comment: asLitValue($2)}
	}
	| CHARACTER SET NAME_OB
	{
		$$ = &colAttrs{charset: asIdentName($3)}
	}
	| CHARSET NAME_OB
	{
		$$ = &colAttrs{charset: asIdentName($2)}
	}
	| COLLATE NAME_OB
	{
		$$ = &colAttrs{collate: asIdentName($2)}
	}
	| INVISIBLE
	{
		$$ = &colAttrs{invisible: true}
	}
	| PRIMARY KEY
	{
		$$ = &colAttrs{}
	}
	| UNIQUE
	{
		$$ = &colAttrs{}
	}
	| UNIQUE KEY
	{
		$$ = &colAttrs{}
	}
	| GENERATED ALWAYS AS '(' expr ')' opt_stored
	{
		$$ = &colAttrs{generated: asExpr($5), genStored: asBool($7), hasGenerated: true}
	}
	| AS '(' expr ')' opt_stored
	{
		$$ = &colAttrs{generated: asExpr($3), genStored: asBool($5), hasGenerated: true}
	}

opt_stored:
	/*empty*/
	{
		$$ = false
	}
	| VIRTUAL
	{
		$$ = false
	}
	| STORED
	{
		$$ = true
	}

table_constraint:
	CONSTRAINT NAME_OB constraint_body
	{
		cd := $3.(*ast.ConstraintDef)
		cd.Name = asIdentName($2)
		$$ = cd
	}
	| constraint_body

constraint_body:
	PRIMARY KEY '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "PRIMARY KEY", Columns: asIdentList($4)}
	}
	| UNIQUE '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "UNIQUE", Columns: asIdentList($3)}
	}
	| UNIQUE KEY '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "UNIQUE", Columns: asIdentList($4)}
	}
	| KEY '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "KEY", Columns: asIdentList($3)}
	}
	| INDEX '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "INDEX", Columns: asIdentList($3)}
	}
	| FOREIGN KEY '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "FOREIGN KEY", Columns: asIdentList($4)}
	}
	| FOREIGN KEY '(' ident_list ')' REFERENCES table_name
	{
		$$ = &ast.ConstraintDef{Kind: "FOREIGN KEY", Columns: asIdentList($4)}
	}
	| FOREIGN KEY '(' ident_list ')' REFERENCES table_name '(' ident_list ')'
	{
		$$ = &ast.ConstraintDef{Kind: "FOREIGN KEY", Columns: asIdentList($4)}
	}
	| CHECK '(' expr ')'
	{
		$$ = &ast.ConstraintDef{Kind: "CHECK", Check: asExpr($3)}
	}

table_options_opt:
	/*empty*/
	{
		$$ = nil
	}
	| table_options_opt table_option
	{
		$$ = appendOpts($1, $2.(ast.TableOption))
	}

table_option:
	NAME_OB COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: strings.ToLower(asIdentName($1)), Value: asStr($3)}
	}
	| DEFAULT NAME_OB COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "default " + strings.ToLower(asIdentName($2)), Value: asStr($4)}
	}
	| DEFAULT CHARSET COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "default charset", Value: asStr($4)}
	}
	| DEFAULT CHARACTER SET COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "default character set", Value: asStr($5)}
	}
	| DEFAULT COLLATE COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "default collate", Value: asStr($4)}
	}
	| CHARSET COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "charset", Value: asStr($3)}
	}
	| CHARACTER SET COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "character set", Value: asStr($4)}
	}
	| COLLATE COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "collate", Value: asStr($3)}
	}
	| AUTO_INCREMENT COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "auto_increment", Value: asStr($3)}
	}
	| COMMENT COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "comment", Value: asStr($3)}
	}
	| TABLEGROUP COMP_EQ table_option_value
	{
		$$ = ast.TableOption{Name: "tablegroup", Value: asStr($3)}
	}
	| TABLEGROUP table_option_value
	{
		$$ = ast.TableOption{Name: "tablegroup", Value: asStr($2)}
	}

table_option_value:
	NAME_OB
	{
		$$ = asIdentName($1)
	}
	| STRING_VALUE
	{
		$$ = asLitValue($1)
	}
	| INTNUM
	{
		$$ = asLitValue($1)
	}
	| DECIMAL_VAL
	{
		$$ = asLitValue($1)
	}
	| HEXNUM
	{
		$$ = asLitValue($1)
	}
	| BOOL_VALUE
	{
		$$ = asLitValue($1)
	}

opt_if_not_exists:
	/*empty*/
	{
		$$ = false
	}
	| IF NOT EXISTS
	{
		$$ = true
	}

opt_if_exists:
	/*empty*/
	{
		$$ = false
	}
	| IF EXISTS
	{
		$$ = true
	}

opt_purge:
	/*empty*/
	{
		$$ = false
	}
	| PURGE
	{
		$$ = true
	}

index_col_list:
	index_col
	{
		$$ = []*ast.Identifier{asIdent($1)}
	}
	| index_col_list ',' index_col
	{
		$$ = append(asIdentList($1), asIdent($3))
	}

index_col:
	NAME_OB opt_asc_desc
	{
		$$ = $1
	}

db_opts_opt:
	/*empty*/
	{
		$$ = dbOpts{}
	}
	| db_opts_opt db_opt
	{
		a := asDBOpts($1)
		b := asDBOpts($2)
		if b.charset != "" {
			a.charset = b.charset
		}
		if b.collate != "" {
			a.collate = b.collate
		}
		$$ = a
	}

db_opt:
	DEFAULT CHARSET NAME_OB
	{
		$$ = dbOpts{charset: asIdentName($3)}
	}
	| DEFAULT CHARACTER SET NAME_OB
	{
		$$ = dbOpts{charset: asIdentName($4)}
	}
	| DEFAULT COLLATE NAME_OB
	{
		$$ = dbOpts{collate: asIdentName($3)}
	}
	| CHARSET NAME_OB
	{
		$$ = dbOpts{charset: asIdentName($2)}
	}
	| CHARACTER SET NAME_OB
	{
		$$ = dbOpts{charset: asIdentName($3)}
	}
	| COLLATE NAME_OB
	{
		$$ = dbOpts{collate: asIdentName($2)}
	}

/* ---------------- DROP / TRUNCATE ---------------- */

drop_stmt:
	DROP TABLE opt_if_exists table_name_list
	{
		$$ = &ast.DropTableStmt{IfExists: asBool($3), Tables: asIdentList($4)}
	}
	| DROP INDEX table_name opt_on_table
	{
		$$ = &ast.DropIndexStmt{Name: asIdent($3), Table: asTableName($4)}
	}
	| DROP DATABASE opt_if_exists NAME_OB
	{
		$$ = &ast.DropDatabaseStmt{IfExists: asBool($3), Name: asIdentName($4)}
	}
	| DROP SCHEMA opt_if_exists NAME_OB
	{
		$$ = &ast.DropDatabaseStmt{IfExists: asBool($3), Name: asIdentName($4)}
	}
	| DROP SEQUENCE opt_if_exists NAME_OB
	{
		$$ = &ast.DropSequenceStmt{IfExists: asBool($3), Name: asIdentName($4)}
	}
	| DROP USER opt_if_exists user_spec_list
	{
		$$ = &ast.DropUserStmt{IfExists: asBool($3), Users: asUserSpecs($4)}
	}
	| DROP ROLE opt_if_exists role_name_list
	{
		gateFeat(yylex, feature.FRoleManagement, "ROLE management requires OceanBase 4.2.4+")

		$$ = &ast.DropRoleStmt{IfExists: asBool($3), Roles: asStringList($4)}
	}
	| DROP MATERIALIZED VIEW opt_if_exists table_name_list
	{
		gateFeat(yylex, feature.FMaterializedView, "materialized view requires OceanBase 4.3+")

		$$ = &ast.DropMViewStmt{IfExists: asBool($4), Names: asIdentList($5)}
	}
	| DROP TENANT opt_if_exists NAME_OB opt_purge
	{
		$$ = &ast.DropTenantStmt{IfExists: asBool($3), Name: asIdentName($4), Purge: asBool($5)}
	}
	| DROP RESOURCE UNIT opt_if_exists NAME_OB
	{
		$$ = &ast.DropResourceUnitStmt{IfExists: asBool($4), Name: asIdentName($5)}
	}
	| DROP RESOURCE POOL opt_if_exists NAME_OB
	{
		$$ = &ast.DropResourcePoolStmt{IfExists: asBool($4), Name: asIdentName($5)}
	}
	| DROP TABLEGROUP opt_if_exists NAME_OB
	{
		$$ = &ast.DropTablegroupStmt{IfExists: asBool($3), Name: asIdentName($4)}
	}
	| DROP TABLESPACE opt_if_exists NAME_OB
	{
		$$ = &ast.DropTablespaceStmt{IfExists: asBool($3), Name: asIdentName($4)}
	}
	| DROP RESTORE POINT NAME_OB
	{
		$$ = &ast.RestorePointStmt{Action: "DROP", Name: asIdentName($4)}
	}
	| DROP SNAPSHOT NAME_OB opt_for_tenant
	{
		gateFeat(yylex, feature.FTenantSnapshot, "tenant snapshot requires OceanBase 4.3+")

		$$ = &ast.TenantSnapshotStmt{Action: "DROP", Name: asIdentName($3), Tenant: asStr($4), HasTenant: asStr($4) != ""}
	}
	| DROP CONCURRENT_LIMITING_RULE opt_if_exists NAME_OB
	{
		gateFeat(yylex, feature.FCCLRule, "CONCURRENT_LIMITING_RULE requires OceanBase 4.3+")

		$$ = &ast.CCLLimitRuleStmt{Action: "DROP", IfExists: asBool($3), Name: asIdentName($4)}
	}
	| DROP SENSITIVE RULE opt_if_exists NAME_OB
	{
		gateFeat(yylex, feature.FSensitiveRule, "SENSITIVE RULE requires OceanBase 4.4+")

		$$ = &ast.SensitiveRuleStmt{Action: "DROP", Name: asIdentName($4)}
	}
	| DROP DATABASE LINK opt_if_exists NAME_OB
	{
		$$ = &ast.DblinkStmt{Action: "DROP", Name: asIdentName($4)}
	}
	| DROP CATALOG opt_if_exists NAME_OB
	{
		gateFeat(yylex, feature.FExternalCatalog, "CATALOG requires OceanBase 4.3+")

		$$ = &ast.CatalogStmt{Action: "DROP", IfExists: asBool($3), Name: asIdentName($4)}
	}
	| DROP LOCATION NAME_OB
	{
		gateFeat(yylex, feature.FExternalLocation, "LOCATION requires OceanBase 4.3+")

		$$ = &ast.CreateLocationStmt{Action: "DROP", Name: asIdentName($3)}
	}
	| DROP PROCEDURE opt_if_exists NAME_OB
	{
		$$ = &ast.DropRoutineStmt{Name: asIdentName($4)}
	}
	| DROP FUNCTION opt_if_exists NAME_OB
	{
		$$ = &ast.DropRoutineStmt{Function: true, Name: asIdentName($4)}
	}

table_name_list:
	table_name
	{
		$$ = []*ast.Identifier{asIdent($1)}
	}
	| table_name_list ',' table_name
	{
		$$ = append(asIdentList($1), asIdent($3))
	}

opt_on_table:
	/*empty*/
	{
		$$ = nil
	}
	| ON table_name
	{
		$$ = $2
	}

truncate_stmt:
	TRUNCATE opt_table_kw table_name
	{
		$$ = &ast.TruncateStmt{Table: asIdent($3)}
	}

opt_table_kw:
	/*empty*/
	{
		$$ = nil
	}
	| TABLE
	{
		$$ = nil
	}

/* ---------------- ALTER TABLE ---------------- */

alter_stmt:
	ALTER TABLE table_name alter_action_list
	{
		$$ = &ast.AlterTableStmt{Table: asIdent($3), Actions: asAlterActions($4)}
	}
	| ALTER DATABASE NAME_OB db_opts_opt
	{
		o := asDBOpts($4)
		$$ = &ast.AlterDatabaseStmt{Name: asIdentName($3), Charset: o.charset, Collate: o.collate}
	}
	| ALTER SCHEMA NAME_OB db_opts_opt
	{
		o := asDBOpts($4)
		$$ = &ast.AlterDatabaseStmt{Name: asIdentName($3), Charset: o.charset, Collate: o.collate}
	}
	| ALTER SYSTEM SET NAME_OB COMP_EQ expr
	{
		$$ = &ast.AlterSystemStmt{Name: asIdentName($4), Value: asExpr($6)}
	}
	| ALTER TENANT NAME_OB opt_cluster_options
	{
		$$ = &ast.AlterTenantStmt{Name: asIdentName($3), Options: asOptions($4)}
	}
	| ALTER VIEW table_name col_list_opt AS select_stmt
	{
		$$ = &ast.AlterViewStmt{Name: asIdent($3), Columns: asIdentList($4), Select: asSelect($6)}
	}
	| ALTER SEQUENCE NAME_OB sequence_options
	{
		so := asSeqOpts($4)
		$$ = &ast.AlterSequenceStmt{Name: asIdentName($3),
			StartWith: so.start, IncrementBy: so.incr, MinValue: so.min, MaxValue: so.max,
			Cache: so.cache, Cycle: so.cycle, Order: so.order}
	}
	| ALTER TABLEGROUP NAME_OB opt_cluster_options
	{
		$$ = &ast.AlterTablegroupStmt{Name: asIdentName($3), Options: asOptions($4)}
	}
	| ALTER RESOURCE UNIT NAME_OB opt_cluster_options
	{
		$$ = &ast.AlterResourceUnitStmt{Name: asIdentName($4), Options: asOptions($5)}
	}
	| ALTER RESOURCE POOL NAME_OB opt_cluster_options
	{
		$$ = &ast.AlterResourcePoolStmt{Name: asIdentName($4), Options: asOptions($5)}
	}
	| ALTER MATERIALIZED VIEW table_name mview_opts
	{
		gateFeat(yylex, feature.FMaterializedView, "materialized view requires OceanBase 4.3+")

		o := asMViewOpts($5)
		$$ = &ast.AlterMViewStmt{Name: asIdent($4), Refresh: o.refresh, Mode: o.mode}
	}
	| ALTER TABLESPACE NAME_OB opt_cluster_options
	{
		$$ = &ast.AlterTablespaceStmt{Name: asIdentName($3), Options: asOptions($4)}
	}
	| ALTER SYSTEM ACTIVATE STANDBY opt_tenant_kw
	{
		$$ = &ast.SwitchoverTenantStmt{Clause: "ACTIVATE STANDBY", Name: asStr($5)}
	}
	| ALTER SYSTEM SWITCHOVER TO PRIMARY opt_tenant_kw
	{
		$$ = &ast.SwitchoverTenantStmt{Clause: "SWITCHOVER TO PRIMARY", Name: asStr($6)}
	}
	| ALTER SYSTEM SWITCHOVER TO STANDBY opt_tenant_kw
	{
		$$ = &ast.SwitchoverTenantStmt{Clause: "SWITCHOVER TO STANDBY", Name: asStr($6)}
	}
	| ALTER SYSTEM RECOVER STANDBY opt_tenant_kw recover_until
	{
		ru := $6.(recUntil)
		$$ = &ast.RecoverTenantStmt{Tenant: asStr($5), Until: ru.until, UntilValue: ru.value, Cancel: ru.cancel}
	}
	| ALTER SYSTEM CREATE LS opt_cluster_options opt_tenant_kw
	{
		gateFeat(yylex, feature.FAlterLS, "ALTER LS requires OceanBase 4.4+")

		$$ = &ast.AlterLSStmt{Op: "CREATE", Options: asOptions($5), Tenant: asStr($6)}
	}
	| ALTER SYSTEM MODIFY LS opt_cluster_options opt_tenant_kw
	{
		gateFeat(yylex, feature.FAlterLS, "ALTER LS requires OceanBase 4.4+")

		$$ = &ast.AlterLSStmt{Op: "MODIFY", Options: asOptions($5), Tenant: asStr($6)}
	}
	| ALTER SYSTEM ROTATE opt_innodb MASTER KEY
	{
		$$ = &ast.RotateMasterKeyStmt{Innodb: asBool($4)}
	}
	| ALTER SYSTEM CREATE SERVICE NAME_OB opt_tenant_kw
	{
		gateFeat(yylex, feature.FServiceName, "SERVICE requires OceanBase 4.2.4+")

		$$ = &ast.ServiceStmt{Create: true, Name: asIdentName($4), Tenant: asStr($5)}
	}
	| ALTER SYSTEM DELETE SERVICE NAME_OB opt_tenant_kw
	{
		gateFeat(yylex, feature.FServiceName, "SERVICE requires OceanBase 4.2.4+")

		$$ = &ast.ServiceStmt{Create: false, Name: asIdentName($4), Tenant: asStr($5)}
	}
	| ALTER SYSTEM FLASHBACK STANDBY NAME_OB TO SCN COMP_EQ expr opt_tenant_kw
	{
		gateFeat(yylex, feature.FFlashbackStandbyLog, "FLASHBACK STANDBY LOG requires OceanBase 5.0+")

		$$ = &ast.FlashbackStandbyLogStmt{SCN: asExpr($9), Tenant: asStr($10)}
	}
	| ALTER CONCURRENT_LIMITING_RULE NAME_OB opt_cluster_options
	{
		gateFeat(yylex, feature.FCCLRule, "CONCURRENT_LIMITING_RULE requires OceanBase 4.3+")

		$$ = &ast.CCLLimitRuleStmt{Action: "ALTER", Name: asIdentName($3), Options: asOptions($4)}
	}
	| ALTER SENSITIVE RULE NAME_OB ON ident_list USING NAME_OB
	{
		gateFeat(yylex, feature.FSensitiveRule, "SENSITIVE RULE requires OceanBase 4.4+")

		$$ = &ast.SensitiveRuleStmt{Action: "ALTER", Name: asIdentName($4), Fields: idNamesList(asIdentList($6)), Using: asIdentName($8)}
	}
	| ALTER LOCATION NAME_OB SET URL opt_equal STRING_VALUE
	{
		gateFeat(yylex, feature.FExternalLocation, "LOCATION requires OceanBase 4.3+")

		$$ = &ast.CreateLocationStmt{Action: "ALTER", Name: asIdentName($3), URL: asLitValue($7)}
	}
	| ALTER SYSTEM TRANSFER PARTITION part_info TO LS expr opt_tenant_kw
	{
		gateFeat(yylex, feature.FTransferPartition, "TRANSFER PARTITION requires OceanBase 4.2.2+")

		pi := $5.(partInfo)
		$$ = &ast.TransferPartitionStmt{TableID: pi.table, ObjectID: pi.obj, LS: asExpr($8), Tenant: asStr($9)}
	}
	| ALTER SYSTEM CANCEL TRANSFER PARTITION part_info opt_tenant_kw
	{
		gateFeat(yylex, feature.FTransferPartition, "TRANSFER PARTITION requires OceanBase 4.2.2+")

		pi := $6.(partInfo)
		$$ = &ast.TransferPartitionStmt{Cancel: true, TableID: pi.table, ObjectID: pi.obj, Tenant: asStr($7)}
	}
	| ALTER SYSTEM CANCEL TRANSFER PARTITION ALL opt_tenant_kw
	{
		gateFeat(yylex, feature.FTransferPartition, "TRANSFER PARTITION requires OceanBase 4.2.2+")

		$$ = &ast.TransferPartitionStmt{Cancel: true, All: true, Tenant: asStr($7)}
	}

part_info:
	NAME_OB opt_equal INTNUM ',' NAME_OB opt_equal INTNUM
	{
		$$ = partInfo{table: asExpr($3), obj: asExpr($7)}
	}

opt_innodb:
	/*empty*/
	{
		$$ = false
	}
	| NAME_OB
	{
		$$ = true
	}

opt_tenant_kw:
	/*empty*/
	{
		$$ = ""
	}
	| TENANT NAME_OB
	{
		$$ = asIdentName($2)
	}

recover_until:
	/*empty*/
	{
		$$ = recUntil{}
	}
	| UNTIL UNLIMITED
	{
		$$ = recUntil{until: "UNLIMITED"}
	}
	| UNTIL TIME STRING_VALUE
	{
		$$ = recUntil{until: "TIME", value: asExpr($3)}
	}
	| UNTIL SCN expr
	{
		$$ = recUntil{until: "SCN", value: asExpr($3)}
	}
	| CANCEL
	{
		$$ = recUntil{cancel: true}
	}

mlog_stmt:
	CREATE MATERIALIZED VIEW LOG ON table_name opt_mlog_opts
	{
		gateFeat(yylex, feature.FMaterializedView, "MVIEW LOG requires OceanBase 4.3+")

		o := asMLogOpts($7)
		$$ = &ast.MLogStmt{Action: "CREATE", Table: asIdent($6), WithPrimary: o.withPrimary, IncludingNew: o.includingNew}
	}
	| ALTER MATERIALIZED VIEW LOG ON table_name opt_mlog_opts
	{
		gateFeat(yylex, feature.FMaterializedView, "MVIEW LOG requires OceanBase 4.3+")

		o := asMLogOpts($7)
		$$ = &ast.MLogStmt{Action: "ALTER", Table: asIdent($6), WithPrimary: o.withPrimary, IncludingNew: o.includingNew}
	}
	| DROP MATERIALIZED VIEW LOG opt_if_exists ON table_name
	{
		gateFeat(yylex, feature.FMaterializedView, "MVIEW LOG requires OceanBase 4.3+")

		$$ = &ast.MLogStmt{Action: "DROP", IfExists: asBool($5), Table: asIdent($7)}
	}

opt_mlog_opts:
	/*empty*/
	{
		$$ = mlogOpts{}
	}
	| opt_mlog_opts mlog_opt
	{
		a := asMLogOpts($1)
		b := asMLogOpts($2)
		if b.withPrimary {
			a.withPrimary = true
		}
		if b.includingNew {
			a.includingNew = true
		}
		$$ = a
	}

mlog_opt:
	WITH PRIMARY KEY
	{
		$$ = mlogOpts{withPrimary: true}
	}
	| INCLUDING NEW VALUES
	{
		$$ = mlogOpts{includingNew: true}
	}

create_outline_stmt:
	CREATE OUTLINE NAME_OB ON outline_target
	{
		co := &ast.CreateOutlineStmt{Name: asIdent($3)}
		if sel := asSelect($5); sel != nil {
			co.Select = sel
		} else {
			co.Text = asLitValue($5)
		}
		$$ = co
	}

alter_outline_stmt:
	ALTER OUTLINE NAME_OB ON outline_target
	{
		co := &ast.AlterOutlineStmt{Name: asIdent($3)}
		if sel := asSelect($5); sel != nil {
			co.Select = sel
		} else {
			co.Text = asLitValue($5)
		}
		$$ = co
	}

drop_outline_stmt:
	DROP OUTLINE opt_if_exists NAME_OB
	{
		$$ = &ast.DropOutlineStmt{IfExists: asBool($3), Name: asIdent($4)}
	}

outline_target:
	select_stmt
	| STRING_VALUE

rotate_master_key_stmt:
	ROTATE MASTER KEY
	{
		$$ = &ast.RotateMasterKeyStmt{}
	}

audit_stmt:
	AUDIT grant_priv_list ON priv_level
	{
		$$ = &ast.AuditStmt{Audit: true, Privileges: asStringList($2), Level: asStr($4)}
	}
	| NOAUDIT grant_priv_list ON priv_level
	{
		$$ = &ast.AuditStmt{Audit: false, Privileges: asStringList($2), Level: asStr($4)}
	}

dump_memory_stmt:
	DUMP ENTITY ALL
	{
		$$ = &ast.DumpMemoryStmt{Kind: "ENTITY", All: true}
	}
	| DUMP CHUNK ALL
	{
		$$ = &ast.DumpMemoryStmt{Kind: "CHUNK", All: true}
	}
	| DUMP ENTITY NAME_OB COMP_EQ STRING_VALUE ',' NAME_OB COMP_EQ expr
	{
		$$ = &ast.DumpMemoryStmt{Kind: "ENTITY", FirstName: asIdentName($3), FirstVal: asExpr($5), SecondName: asIdentName($7), SecondVal: asExpr($9)}
	}
	| DUMP CHUNK NAME_OB COMP_EQ expr ',' NAME_OB COMP_EQ expr
	{
		$$ = &ast.DumpMemoryStmt{Kind: "CHUNK", FirstName: asIdentName($3), FirstVal: asExpr($5), SecondName: asIdentName($7), SecondVal: asExpr($9)}
	}

location_utils_stmt:
	REMOVE FILES IN LOCATION USER_VARIABLE
	{
		gateFeat(yylex, feature.FExternalLocation, "LOCATION requires OceanBase 4.3+")

		$$ = &ast.LocationUtilsStmt{Path: variableText($5.(*ast.Variable))}
	}

install_plugin_stmt:
	INSTALL PLUGIN NAME_OB SONAME STRING_VALUE
	{
		$$ = &ast.InstallPluginStmt{Name: asIdentName($3), Soname: asLitValue($5)}
	}

uninstall_plugin_stmt:
	UNINSTALL PLUGIN NAME_OB
	{
		$$ = &ast.UninstallPluginStmt{Name: asIdentName($3)}
	}

handler_stmt:
	HANDLER table_name OPEN opt_handler_alias
	{
		$$ = &ast.HandlerStmt{Op: "OPEN", Table: asIdent($2), Alias: asStr($4)}
	}
	| HANDLER table_name CLOSE
	{
		$$ = &ast.HandlerStmt{Op: "CLOSE", Table: asIdent($2)}
	}
	| HANDLER table_name READ handler_direction opt_handler_cols opt_handler_where opt_handler_limit
	{
		$$ = &ast.HandlerStmt{Op: "READ", Table: asIdent($2), Direction: asStr($4), Columns: asIdentList($5), Where: asExpr($6), Limit: asLimit($7)}
	}

opt_handler_alias:
	/*empty*/
	{
		$$ = ""
	}
	| AS NAME_OB
	{
		$$ = asIdentName($2)
	}

handler_direction:
	FIRST
	{
		$$ = "FIRST"
	}
	| NEXT
	{
		$$ = "NEXT"
	}
	| PREV
	{
		$$ = "PREV"
	}
	| LAST
	{
		$$ = "LAST"
	}

opt_handler_cols:
	/*empty*/
	{
		$$ = nil
	}
	| '(' ident_list ')'
	{
		$$ = $2
	}

opt_handler_where:
	/*empty*/
	{
		$$ = nil
	}
	| WHERE expr
	{
		$$ = $2
	}

opt_handler_limit:
	/*empty*/
	{
		$$ = nil
	}
	| LIMIT expr
	{
		$$ = &ast.LimitClause{Count: asExpr($2)}
	}

create_view_stmt:
	CREATE VIEW table_name col_list_opt AS select_stmt
	{
		$$ = &ast.CreateViewStmt{Name: asIdent($3), Columns: asIdentList($4), Select: asSelect($6)}
	}
	| CREATE OR REPLACE VIEW table_name col_list_opt AS select_stmt
	{
		$$ = &ast.CreateViewStmt{OrReplace: true, Name: asIdent($5), Columns: asIdentList($6), Select: asSelect($8)}
	}

drop_view_stmt:
	DROP VIEW opt_if_exists table_name_list
	{
		$$ = &ast.DropViewStmt{IfExists: asBool($3), Names: asIdentList($4)}
	}

grant_stmt:
	GRANT grant_priv_list ON priv_level TO user_spec_list opt_with_grant
	{
		$$ = &ast.GrantStmt{Privileges: asStringList($2), Level: asStr($4), Users: asUserSpecs($6), WithGrant: asBool($7)}
	}

revoke_stmt:
	REVOKE grant_priv_list ON priv_level FROM user_spec_list
	{
		$$ = &ast.RevokeStmt{Privileges: asStringList($2), Level: asStr($4), Users: asUserSpecs($6)}
	}

grant_priv_list:
	grant_priv
	{
		$$ = []string{$1.(string)}
	}
	| grant_priv_list ',' grant_priv
	{
		$$ = append(asStringList($1), $3.(string))
	}

grant_priv:
	NAME_OB opt_priv_cols
	{
		$$ = asIdentName($1) + asStr($2)
	}
	| ALL opt_privileges_kw
	{
		$$ = "ALL"
	}

opt_privileges_kw:
	/*empty*/
	{
		$$ = nil
	}
	| PRIVILEGES
	{
		$$ = nil
	}

opt_priv_cols:
	/*empty*/
	{
		$$ = ""
	}
	| '(' ident_list ')'
	{
		$$ = " (" + joinIdentNames(asIdentList($2)) + ")"
	}

priv_level:
	'*'
	{
		$$ = "*"
	}
	| '*' '.' '*'
	{
		$$ = "*.*"
	}
	| NAME_OB '.' '*'
	{
		$$ = asIdentName($1) + ".*"
	}
	| NAME_OB '.' NAME_OB
	{
		$$ = asIdentName($1) + "." + asIdentName($3)
	}

user_spec_list:
	user_spec
	{
		$$ = []ast.UserSpec{$1.(ast.UserSpec)}
	}
	| user_spec_list ',' user_spec
	{
		$$ = append(asUserSpecs($1), $3.(ast.UserSpec))
	}

user_spec:
	user_name AT user_host
	{
		$$ = ast.UserSpec{Name: asStr($1), Host: asStr($3), HasHost: true}
	}
	| user_name
	{
		$$ = ast.UserSpec{Name: asStr($1)}
	}

user_name:
	STRING_VALUE
	{
		$$ = asLitValue($1)
	}
	| NAME_OB
	{
		$$ = asIdentName($1)
	}

user_host:
	STRING_VALUE
	{
		$$ = asLitValue($1)
	}
	| NAME_OB
	{
		$$ = asIdentName($1)
	}

opt_with_grant:
	/*empty*/
	{
		$$ = false
	}
	| WITH GRANT OPTION
	{
		$$ = true
	}

rename_user_stmt:
	RENAME USER user_spec TO user_spec
	{
		$$ = &ast.RenameUserStmt{Old: $3.(ast.UserSpec), New: $5.(ast.UserSpec)}
	}

lock_user_stmt:
	LOCK USER user_spec_list
	{
		$$ = &ast.LockUserStmt{Users: asUserSpecs($3)}
	}
	| UNLOCK USER user_spec_list
	{
		$$ = &ast.LockUserStmt{Unlock: true, Users: asUserSpecs($3)}
	}

prepare_stmt:
	PREPARE NAME_OB FROM prep_source
	{
		$$ = &ast.PrepareStmt{Name: asIdentName($2), Source: asExpr($4)}
	}

prep_source:
	STRING_VALUE
	| USER_VARIABLE

execute_stmt:
	EXECUTE NAME_OB execute_opt_using
	{
		$$ = &ast.ExecuteStmt{Name: asIdentName($2), Using: asVariableList($3)}
	}

execute_opt_using:
	/*empty*/
	{
		$$ = nil
	}
	| USING user_var_list
	{
		$$ = $2
	}

user_var_list:
	USER_VARIABLE
	{
		$$ = []*ast.Variable{$1.(*ast.Variable)}
	}
	| user_var_list ',' USER_VARIABLE
	{
		$$ = append(asVariableList($1), $3.(*ast.Variable))
	}

deallocate_stmt:
	DEALLOCATE PREPARE NAME_OB
	{
		$$ = &ast.DeallocateStmt{Name: asIdentName($3)}
	}

kill_stmt:
	KILL opt_kill_kind expr
	{
		$$ = &ast.KillStmt{Kind: asStr($2), Target: asExpr($3)}
	}

opt_kill_kind:
	/*empty*/
	{
		$$ = ""
	}
	| CONNECTION
	{
		$$ = "CONNECTION"
	}
	| QUERY
	{
		$$ = "QUERY"
	}

help_stmt:
	HELP expr
	{
		$$ = &ast.HelpStmt{Target: asExpr($2)}
	}

table_maint_stmt:
	ANALYZE TABLE table_name_list
	{
		$$ = &ast.TableMaintStmt{Kind: "ANALYZE", Tables: asIdentList($3)}
	}
	| ANALYZE TABLE table_name_list METHOD_OPT STRING_VALUE
	{
		$$ = &ast.TableMaintStmt{Kind: "ANALYZE", Tables: asIdentList($3), MethodOpt: asLitValue($5)}
	}
	| OPTIMIZE TABLE table_name_list
	{
		$$ = &ast.TableMaintStmt{Kind: "OPTIMIZE", Tables: asIdentList($3)}
	}
	| REPAIR TABLE table_name_list
	{
		$$ = &ast.TableMaintStmt{Kind: "REPAIR", Tables: asIdentList($3)}
	}

rename_table_stmt:
	RENAME TABLE rename_pair_list
	{
		$$ = &ast.RenameTableStmt{Pairs: asRenamePairs($3)}
	}

rename_pair_list:
	rename_pair
	{
		$$ = []ast.RenamePair{$1.(ast.RenamePair)}
	}
	| rename_pair_list ',' rename_pair
	{
		$$ = append(asRenamePairs($1), $3.(ast.RenamePair))
	}

rename_pair:
	table_name TO table_name
	{
		$$ = ast.RenamePair{From: asIdent($1), To: asIdent($3)}
	}

load_data_stmt:
	LOAD DATA opt_direct_hint opt_local INFILE STRING_VALUE opt_load_conflict INTO TABLE table_name load_fields_opt load_lines_opt load_ignore_lines_opt load_cols_opt
	{
		conflict := asStr($7)
		$$ = &ast.LoadDataStmt{
			Hint:               asStr($3),
			Local:              asBool($4),
			Infile:             asLitValue($6),
			Replace:            conflict == "REPLACE",
			Ignore:             conflict == "IGNORE",
			Table:              asIdent($10),
			FieldsTerminatedBy: asStr($11),
			LinesTerminatedBy:  asStr($12),
			IgnoreLines:        asExpr($13),
			Columns:            asIdentList($14),
		}
	}

opt_direct_hint:
	/*empty*/
	{
		$$ = ""
	}
	| DIRECT_HINT
	{
		$$ = asLitValue($1)
	}

opt_local:
	/*empty*/
	{
		$$ = false
	}
	| LOCAL
	{
		$$ = true
	}

opt_load_conflict:
	/*empty*/
	{
		$$ = ""
	}
	| REPLACE
	{
		$$ = "REPLACE"
	}
	| IGNORE
	{
		$$ = "IGNORE"
	}

load_fields_opt:
	/*empty*/
	{
		$$ = ""
	}
	| FIELDS TERMINATED BY STRING_VALUE
	{
		$$ = asLitValue($4)
	}

load_lines_opt:
	/*empty*/
	{
		$$ = ""
	}
	| LINES TERMINATED BY STRING_VALUE
	{
		$$ = asLitValue($4)
	}

load_ignore_lines_opt:
	/*empty*/
	{
		$$ = nil
	}
	| IGNORE expr LINES
	{
		$$ = $2
	}

load_cols_opt:
	/*empty*/
	{
		$$ = nil
	}
	| '(' ident_list ')'
	{
		$$ = $2
	}

flush_stmt:
	FLUSH opt_local flush_body
	{
		fs := $3.(*ast.FlushStmt)
		fs.Local = asBool($2)
		$$ = fs
	}

flush_body:
	TABLES flush_tables_opt flush_opt
	{
		$$ = &ast.FlushStmt{Kind: "TABLES", Tables: asIdentList($2), Option: asStr($3)}
	}
	| PRIVILEGES
	{
		$$ = &ast.FlushStmt{Kind: "PRIVILEGES"}
	}
	| STATUS
	{
		$$ = &ast.FlushStmt{Kind: "STATUS"}
	}
	| LOGS
	{
		$$ = &ast.FlushStmt{Kind: "LOGS"}
	}
	| HOSTS
	{
		$$ = &ast.FlushStmt{Kind: "HOSTS"}
	}

flush_tables_opt:
	/*empty*/
	{
		$$ = nil
	}
	| table_name_list

flush_opt:
	/*empty*/
	{
		$$ = ""
	}
	| WITH READ LOCK
	{
		$$ = "WITH READ LOCK"
	}
	| FOR EXPORT
	{
		$$ = "FOR EXPORT"
	}

xa_stmt:
	XA START xid xa_join_opt
	{
		$$ = &ast.XAStmt{Action: "START", XID: asExprList($3), Option: asStr($4)}
	}
	| XA BEGIN xid xa_join_opt
	{
		$$ = &ast.XAStmt{Action: "START", XID: asExprList($3), Option: asStr($4)}
	}
	| XA END xid xa_suspend_opt
	{
		$$ = &ast.XAStmt{Action: "END", XID: asExprList($3), Option: asStr($4)}
	}
	| XA PREPARE xid
	{
		$$ = &ast.XAStmt{Action: "PREPARE", XID: asExprList($3)}
	}
	| XA COMMIT xid xa_one_phase_opt
	{
		$$ = &ast.XAStmt{Action: "COMMIT", XID: asExprList($3), Option: asStr($4)}
	}
	| XA ROLLBACK xid
	{
		$$ = &ast.XAStmt{Action: "ROLLBACK", XID: asExprList($3)}
	}
	| XA RECOVER xa_recover_opt
	{
		gateFeat(yylex, feature.FXaRecover, "XA RECOVER requires OceanBase 4.2.5+")

		$$ = &ast.XAStmt{Action: "RECOVER", Option: asStr($3)}
	}

xid:
	expr
	{
		$$ = []ast.Expression{asExpr($1)}
	}
	| expr ',' expr
	{
		$$ = []ast.Expression{asExpr($1), asExpr($3)}
	}
	| expr ',' expr ',' expr
	{
		$$ = []ast.Expression{asExpr($1), asExpr($3), asExpr($5)}
	}

xa_join_opt:
	/*empty*/
	{
		$$ = ""
	}
	| JOIN
	{
		$$ = "JOIN"
	}
	| RESUME
	{
		$$ = "RESUME"
	}

xa_suspend_opt:
	/*empty*/
	{
		$$ = ""
	}
	| SUSPEND
	{
		$$ = "SUSPEND"
	}
	| SUSPEND FOR MIGRATE
	{
		$$ = "SUSPEND FOR MIGRATE"
	}

xa_one_phase_opt:
	/*empty*/
	{
		$$ = ""
	}
	| ONE PHASE
	{
		$$ = "ONE PHASE"
	}

xa_recover_opt:
	/*empty*/
	{
		$$ = ""
	}
	| CONVERT XID
	{
		$$ = "CONVERT XID"
	}

purge_stmt:
	PURGE TABLE table_name
	{
		$$ = &ast.PurgeStmt{Kind: "TABLE", Table: asIdent($3)}
	}
	| PURGE INDEX table_name ON table_name
	{
		$$ = &ast.PurgeStmt{Kind: "INDEX", Index: asIdent($3), Table: asIdent($5)}
	}
	| PURGE RECYCLEBIN
	{
		$$ = &ast.PurgeStmt{Kind: "RECYCLEBIN"}
	}

merge_stmt:
	MERGE INTO table_ref USING table_ref ON expr merge_when_list
	{
		gateFeat(yylex, feature.FMerge, "MERGE requires OceanBase 4.4+")
		$$ = &ast.MergeStmt{Target: asTableRef($3), Source: asTableRef($5), On: asExpr($7), Whens: asMergeWhens($8)}
	}

merge_when_list:
	merge_when
	{
		$$ = []ast.MergeWhen{$1.(ast.MergeWhen)}
	}
	| merge_when_list merge_when
	{
		$$ = append(asMergeWhens($1), $2.(ast.MergeWhen))
	}

merge_when:
	WHEN MATCHED THEN UPDATE SET assignment_list opt_delete_where
	{
		$$ = ast.MergeWhen{Matched: true, Action: "UPDATE", Set: asAssignList($6), DeleteCond: asExpr($7)}
	}
	| WHEN MATCHED THEN DELETE opt_merge_where
	{
		$$ = ast.MergeWhen{Matched: true, Action: "DELETE", DeleteCond: asExpr($5)}
	}
	| WHEN NOT MATCHED THEN INSERT '(' ident_list ')' VALUES '(' expr_list ')'
	{
		$$ = ast.MergeWhen{Action: "INSERT", Columns: asIdentList($7), Values: asExprList($11)}
	}
	| WHEN NOT MATCHED THEN INSERT VALUES '(' expr_list ')'
	{
		$$ = ast.MergeWhen{Action: "INSERT", Values: asExprList($8)}
	}

opt_delete_where:
	/*empty*/
	{
		$$ = nil
	}
	| DELETE WHERE expr
	{
		$$ = $3
	}

opt_merge_where:
	/*empty*/
	{
		$$ = nil
	}
	| WHERE expr
	{
		$$ = $2
	}

flashback_table_stmt:
	FLASHBACK TABLE table_name_list TO BEFORE DROP opt_flashback_rename
	{
		$$ = &ast.FlashbackTableStmt{Tables: asIdentList($3), ToBeforeDrop: true, RenameTo: asIdent($7)}
	}
	| FLASHBACK TABLE table_name_list TO SNAPSHOT NAME_OB
	{
		$$ = &ast.FlashbackTableStmt{Tables: asIdentList($3), Snapshot: asIdentName($6)}
	}

opt_flashback_rename:
	/*empty*/
	{
		$$ = nil
	}
	| RENAME TO table_name
	{
		$$ = $3
	}

signal_stmt:
	SIGNAL SQLSTATE STRING_VALUE opt_signal_set
	{
		$$ = &ast.SignalStmt{SqlState: asLitValue($3), HasSqlState: true, Items: asSignalItems($4)}
	}
	| RESIGNAL opt_sqlstate_clause opt_signal_set
	{
		$$ = &ast.SignalStmt{IsResignal: true, SqlState: asStr($2), HasSqlState: asStr($2) != "", Items: asSignalItems($3)}
	}

opt_sqlstate_clause:
	/*empty*/
	{
		$$ = ""
	}
	| SQLSTATE STRING_VALUE
	{
		$$ = asLitValue($2)
	}

opt_signal_set:
	/*empty*/
	{
		$$ = nil
	}
	| SET signal_item_list
	{
		$$ = $2
	}

signal_item_list:
	signal_item
	{
		$$ = []ast.SignalItem{$1.(ast.SignalItem)}
	}
	| signal_item_list ',' signal_item
	{
		$$ = append(asSignalItems($1), $3.(ast.SignalItem))
	}

signal_item:
	NAME_OB COMP_EQ expr
	{
		$$ = ast.SignalItem{Name: asIdentName($1), Value: asExpr($3)}
	}

get_diagnostics_stmt:
	GET opt_diag_scope DIAGNOSTICS opt_diag_condition diag_item_list
	{
		$$ = &ast.GetDiagnosticsStmt{Scope: asStr($2), Condition: asExpr($4), HasCondition: $4 != nil, Items: asDiagItems($5)}
	}

opt_diag_condition:
	/*empty*/
	{
		$$ = nil
	}
	| CONDITION expr
	{
		$$ = $2
	}

opt_diag_scope:
	/*empty*/
	{
		$$ = ""
	}
	| CURRENT
	{
		$$ = "CURRENT"
	}
	| STACKED
	{
		$$ = "STACKED"
	}

diag_item_list:
	diag_item
	{
		$$ = []ast.DiagItem{$1.(ast.DiagItem)}
	}
	| diag_item_list ',' diag_item
	{
		$$ = append(asDiagItems($1), $3.(ast.DiagItem))
	}

diag_item:
	NAME_OB COMP_EQ NAME_OB
	{
		$$ = ast.DiagItem{Target: asIdentName($1), Source: asIdentName($3)}
	}

alter_action_list:
	alter_action
	{
		$$ = []ast.AlterAction{$1.(ast.AlterAction)}
	}
	| alter_action_list ',' alter_action
	{
		$$ = append(asAlterActions($1), $3.(ast.AlterAction))
	}

alter_action:
	ADD opt_column column_def opt_col_position
	{
		a := ast.AlterAction{Kind: ast.AlterAddColumn, Column: asColumnDef($3)}
		cp := $4.(colPos)
		a.Position = cp.position
		a.RefCol = cp.ref
		$$ = a
	}
	| ADD '(' column_def_list ')'
	{
		$$ = ast.AlterAction{Kind: ast.AlterAddColumns, ColumnDefs: asColumnDefs($3)}
	}
	| ADD PARTITION '(' partition_def_list ')'
	{
		$$ = ast.AlterAction{Kind: ast.AlterAddPartition, PartitionDefs: asPartitionDefs($4)}
	}
	| ADD INDEX NAME_OB '(' ident_list ')'
	{
		$$ = ast.AlterAction{Kind: ast.AlterAddIndex, Name: asIdentName($3), Columns: asIdentList($5)}
	}
	| ADD KEY NAME_OB '(' ident_list ')'
	{
		$$ = ast.AlterAction{Kind: ast.AlterAddIndex, Name: asIdentName($3), Columns: asIdentList($5)}
	}
	| ADD UNIQUE INDEX NAME_OB '(' ident_list ')'
	{
		$$ = ast.AlterAction{Kind: ast.AlterAddIndex, Unique: true, Name: asIdentName($4), Columns: asIdentList($6)}
	}
	| ADD UNIQUE KEY NAME_OB '(' ident_list ')'
	{
		$$ = ast.AlterAction{Kind: ast.AlterAddIndex, Unique: true, Name: asIdentName($4), Columns: asIdentList($6)}
	}
	| ADD CONSTRAINT NAME_OB constraint_body
	{
		cd := $4.(*ast.ConstraintDef)
		cd.Name = asIdentName($3)
		$$ = ast.AlterAction{Kind: ast.AlterAddConstraint, Constraint: cd}
	}
	| ADD constraint_body
	{
		$$ = ast.AlterAction{Kind: ast.AlterAddConstraint, Constraint: $2.(*ast.ConstraintDef)}
	}
	| DROP opt_column NAME_OB
	{
		$$ = ast.AlterAction{Kind: ast.AlterDropColumn, Name: asIdentName($3)}
	}
	| DROP INDEX NAME_OB
	{
		$$ = ast.AlterAction{Kind: ast.AlterDropIndex, Name: asIdentName($3)}
	}
	| DROP PRIMARY KEY
	{
		$$ = ast.AlterAction{Kind: ast.AlterDropPrimaryKey}
	}
	| MODIFY opt_column column_def opt_col_position
	{
		a := ast.AlterAction{Kind: ast.AlterModifyColumn, Column: asColumnDef($3)}
		cp := $4.(colPos)
		a.Position = cp.position
		a.RefCol = cp.ref
		$$ = a
	}
	| CHANGE opt_column NAME_OB column_def
	{
		$$ = ast.AlterAction{Kind: ast.AlterChangeColumn, OldName: asIdentName($3), Column: asColumnDef($4)}
	}
	| DROP PARTITION part_name_or_all
	{
		na := $3.(nameOrAll)
		$$ = ast.AlterAction{Kind: ast.AlterDropPartition, Names: na.names, All: na.all}
	}
	| TRUNCATE PARTITION part_name_or_all
	{
		na := $3.(nameOrAll)
		$$ = ast.AlterAction{Kind: ast.AlterTruncatePartition, Names: na.names, All: na.all}
	}
	| COALESCE PARTITION expr
	{
		$$ = ast.AlterAction{Kind: ast.AlterCoalescePartition, Count: asExpr($3)}
	}
	| EXCHANGE PARTITION name_list WITH TABLE table_name opt_exchange_validation
	{
		a := ast.AlterAction{Kind: ast.AlterExchangePartition, Names: asStringList($3), With: asIdent($6)}
		a.Validation = asStr($7)
		$$ = a
	}
	| REORGANIZE PARTITION part_name_or_all INTO '(' partition_def_list ')'
	{
		na := $3.(nameOrAll)
		$$ = ast.AlterAction{Kind: ast.AlterReorganizePartition, Names: na.names, All: na.all, PartitionDefs: asPartitionDefs($6)}
	}
	| REMOVE NAME_OB
	{
		$$ = ast.AlterAction{Kind: ast.AlterRemovePartitioning}
	}
	| CONVERT TO CHARACTER SET NAME_OB opt_convert_collate
	{
		a := ast.AlterAction{Kind: ast.AlterConvertCharset, Charset: asIdentName($5)}
		a.Collate = asStr($6)
		$$ = a
	}
	| AUTO_INCREMENT COMP_EQ expr
	{
		$$ = ast.AlterAction{Kind: ast.AlterAutoIncrement, Count: asExpr($3)}
	}
	| RENAME TO table_name
	{
		$$ = ast.AlterAction{Kind: ast.AlterRenameTable, Target: asIdent($3)}
	}
	| RENAME table_name
	{
		$$ = ast.AlterAction{Kind: ast.AlterRenameTable, Target: asIdent($2)}
	}
	| RENAME COLUMN NAME_OB TO NAME_OB
	{
		$$ = ast.AlterAction{Kind: ast.AlterRenameColumn, Name: asIdentName($3), NewName: asIdentName($5)}
	}
	| RENAME INDEX NAME_OB TO NAME_OB
	{
		$$ = ast.AlterAction{Kind: ast.AlterRenameIndex, Name: asIdentName($3), NewName: asIdentName($5)}
	}
	| ALTER opt_column NAME_OB SET DEFAULT expr
	{
		$$ = ast.AlterAction{Kind: ast.AlterSetDefault, Name: asIdentName($3), Default: asExpr($6)}
	}
	| ALTER opt_column NAME_OB DROP DEFAULT
	{
		$$ = ast.AlterAction{Kind: ast.AlterDropDefault, Name: asIdentName($3)}
	}

column_def_list:
	column_def
	{
		$$ = []*ast.ColumnDef{asColumnDef($1)}
	}
	| column_def_list ',' column_def
	{
		$$ = append(asColumnDefs($1), asColumnDef($3))
	}

opt_col_position:
	/*empty*/
	{
		$$ = colPos{}
	}
	| FIRST
	{
		$$ = colPos{position: "FIRST"}
	}
	| AFTER NAME_OB
	{
		$$ = colPos{position: "AFTER", ref: asIdentName($2)}
	}

opt_convert_collate:
	/*empty*/
	{
		$$ = ""
	}
	| COLLATE NAME_OB
	{
		$$ = asIdentName($2)
	}

part_name_or_all:
	name_list
	{
		$$ = nameOrAll{names: asStringList($1)}
	}
	| ALL
	{
		$$ = nameOrAll{all: true}
	}

opt_exchange_validation:
	/*empty*/
	{
		$$ = ""
	}
	| WITH VALIDATION
	{
		gateFeat(yylex, feature.FExchangeValidation, "EXCHANGE ... WITH VALIDATION requires OceanBase 4.4+")
		$$ = "WITH VALIDATION"
	}
	| WITHOUT VALIDATION
	{
		gateFeat(yylex, feature.FExchangeValidation, "EXCHANGE ... WITHOUT VALIDATION requires OceanBase 4.4+")
		$$ = "WITHOUT VALIDATION"
	}

opt_column:
	/*empty*/
	{
		$$ = nil
	}
	| COLUMN
	{
		$$ = nil
	}

/* ---------------- PL / stored routines ---------------- */

do_stmt:
	DO expr
	{
		$$ = &ast.DoStmt{Expr: asExpr($2)}
	}

call_stmt:
	CALL NAME_OB '(' expr_list ')'
	{
		$$ = &ast.CallStmt{Name: asIdentName($2), Args: asExprList($4)}
	}
	| CALL NAME_OB '(' ')'
	{
		$$ = &ast.CallStmt{Name: asIdentName($2)}
	}
	| CALL NAME_OB
	{
		$$ = &ast.CallStmt{Name: asIdentName($2)}
	}

proc_params:
	/*empty*/
	{
		$$ = nil
	}
	| '(' proc_param_list ')'
	{
		$$ = $2
	}
	| '(' ')'
	{
		$$ = nil
	}

proc_param_list:
	proc_param
	{
		$$ = []ast.ProcParam{$1.(ast.ProcParam)}
	}
	| proc_param_list ',' proc_param
	{
		$$ = append(asProcParams($1), $3.(ast.ProcParam))
	}

proc_param:
	opt_proc_mode NAME_OB data_type
	{
		d := asDType($3)
		$$ = ast.ProcParam{Mode: asStr($1), Name: asIdentName($2), Type: d.name, TypeArgs: d.args}
	}

opt_proc_mode:
	/*empty*/
	{
		$$ = ""
	}
	| IN
	{
		$$ = ""
	}
	| OUT
	{
		$$ = "OUT"
	}
	| INOUT
	{
		$$ = "INOUT"
	}

proc_body:
	BEGIN sp_stmts END
	{
		$$ = &ast.PLBlock{Statements: asNodes($2)}
	}

sp_stmts:
	/*empty*/
	{
		$$ = nil
	}
	| sp_stmts sp_stmt ';'
	{
		$$ = append(asNodes($1), asNode($2))
	}

sp_stmt:
	inner_stmt
	{
		$$ = $1
	}
	| NAME_OB ':' BEGIN sp_stmts END
	{
		$$ = &ast.PLBlock{Statements: asNodes($4), Label: asIdentName($1)}
	}
	| BEGIN sp_stmts END
	{
		$$ = &ast.PLBlock{Statements: asNodes($2)}
	}
	| LEAVE NAME_OB
	{
		$$ = &ast.PLLeaveStmt{Name: asIdentName($2)}
	}
	| ITERATE NAME_OB
	{
		$$ = &ast.PLIterateStmt{Name: asIdentName($2)}
	}
	| DECLARE handler_action HANDLER FOR handler_conditions sp_stmt
	{
		$$ = &ast.PLConditionHandlerStmt{Action: asStr($2), Conditions: asHandlerConditions($5), Body: asNode($6)}
	}
	| DECLARE ident_list data_type opt_decl_default
	{
		d := asDType($3)
		$$ = &ast.PLDeclareStmt{Names: idNamesList(asIdentList($2)), Type: d.name, TypeArgs: d.args, Default: asExpr($4)}
	}
	| DECLARE NAME_OB CURSOR FOR select_stmt
	{
		$$ = &ast.PLCursorDeclStmt{Name: asIdentName($2), Select: asSelect($5)}
	}
	| SET set_lhs COMP_EQ expr
	{
		$$ = &ast.PLSetStmt{Name: asStr($2), IsUserVar: strings.HasPrefix(asStr($2), "@"), Expr: asExpr($4)}
	}
	| IF expr THEN sp_stmts sp_elseifs opt_sp_else END IF
	{
		branches := append([]ast.PLIfBranch{{Cond: asExpr($2), Body: &ast.PLBlock{Statements: asNodes($4)}}}, asPLBranches($5)...)
		$$ = &ast.PLIfStmt{Branches: branches, Else: asPLBlock($6)}
	}
	| WHILE expr DO sp_stmts END WHILE
	{
		$$ = &ast.PLWhileStmt{Cond: asExpr($2), Body: &ast.PLBlock{Statements: asNodes($4)}}
	}
	| NAME_OB ':' WHILE expr DO sp_stmts END WHILE
	{
		$$ = &ast.PLWhileStmt{Label: asIdentName($1), Cond: asExpr($4), Body: &ast.PLBlock{Statements: asNodes($6)}}
	}
	| REPEAT sp_stmts UNTIL expr END REPEAT
	{
		$$ = &ast.PLRepeatStmt{Body: &ast.PLBlock{Statements: asNodes($2)}, Until: asExpr($4)}
	}
	| NAME_OB ':' REPEAT sp_stmts UNTIL expr END REPEAT
	{
		$$ = &ast.PLRepeatStmt{Label: asIdentName($1), Body: &ast.PLBlock{Statements: asNodes($4)}, Until: asExpr($6)}
	}
	| LOOP sp_stmts END LOOP
	{
		$$ = &ast.PLLoopStmt{Body: &ast.PLBlock{Statements: asNodes($2)}}
	}
	| NAME_OB ':' LOOP sp_stmts END LOOP
	{
		$$ = &ast.PLLoopStmt{Label: asIdentName($1), Body: &ast.PLBlock{Statements: asNodes($4)}}
	}
	| RETURN expr
	{
		$$ = &ast.PLReturnStmt{Expr: asExpr($2)}
	}
	| OPEN NAME_OB
	{
		$$ = &ast.PLCursorStmt{Op: "OPEN", Name: asIdentName($2)}
	}
	| CLOSE NAME_OB
	{
		$$ = &ast.PLCursorStmt{Op: "CLOSE", Name: asIdentName($2)}
	}
	| FETCH NAME_OB INTO ident_list
	{
		$$ = &ast.PLCursorStmt{Op: "FETCH", Name: asIdentName($2), Into: idNamesList(asIdentList($4))}
	}
	| CASE opt_pl_case_operand pl_case_when_list opt_pl_case_else END CASE
	{
		$$ = &ast.PLCaseStmt{Operand: asExpr($2), Branches: asPLCaseWhens($3), Else: asPLBlock($4)}
	}
	| DECLARE NAME_OB CONDITION FOR SQLSTATE STRING_VALUE
	{
		$$ = &ast.PLConditionDeclStmt{Name: asIdentName($2), Kind: "sqlstate", Value: asLitValue($6)}
	}
	| DECLARE NAME_OB CONDITION FOR expr
	{
		$$ = &ast.PLConditionDeclStmt{Name: asIdentName($2), Kind: "code", Value: exprText(asExpr($5))}
	}

opt_pl_case_operand:
	/*empty*/
	{
		$$ = nil
	}
	| expr

pl_case_when_list:
	pl_case_when
	{
		$$ = []ast.PLCaseWhen{$1.(ast.PLCaseWhen)}
	}
	| pl_case_when_list pl_case_when
	{
		$$ = append(asPLCaseWhens($1), $2.(ast.PLCaseWhen))
	}

pl_case_when:
	WHEN expr THEN sp_stmts
	{
		$$ = ast.PLCaseWhen{When: asExpr($2), Body: &ast.PLBlock{Statements: asNodes($4)}}
	}

opt_pl_case_else:
	/*empty*/
	{
		$$ = nil
	}
	| ELSE sp_stmts
	{
		$$ = &ast.PLBlock{Statements: asNodes($2)}
	}

sp_elseifs:
	/*empty*/
	{
		$$ = nil
	}
	| sp_elseifs ELSEIF expr THEN sp_stmts
	{
		b := ast.PLIfBranch{Cond: asExpr($3), Body: &ast.PLBlock{Statements: asNodes($5)}}
		$$ = append(asPLBranches($1), b)
	}

opt_sp_else:
	/*empty*/
	{
		$$ = nil
	}
	| ELSE sp_stmts
	{
		$$ = &ast.PLBlock{Statements: asNodes($2)}
	}

opt_decl_default:
	/*empty*/
	{
		$$ = nil
	}
	| DEFAULT expr
	{
		$$ = $2
	}

handler_action:
	CONTINUE
	{
		$$ = "CONTINUE"
	}
	| EXIT
	{
		$$ = "EXIT"
	}
	| UNDO
	{
		$$ = "UNDO"
	}

handler_conditions:
	handler_condition
	{
		$$ = []ast.HandlerCondition{$1.(ast.HandlerCondition)}
	}
	| handler_conditions ',' handler_condition
	{
		$$ = append(asHandlerConditions($1), $3.(ast.HandlerCondition))
	}

handler_condition:
	SQLSTATE STRING_VALUE
	{
		$$ = ast.HandlerCondition{Kind: "sqlstate", Value: asLitValue($2)}
	}
	| NOT FOUND
	{
		$$ = ast.HandlerCondition{Kind: "notfound"}
	}
	| NAME_OB
	{
		$$ = ast.HandlerCondition{Kind: "name", Value: asIdentName($1)}
	}

/* ---------------- SHOW ---------------- */

show_stmt:
	SHOW full_opt show_body
	{
		st := $3.(*ast.ShowStmt)
		st.Full = asBool($2)
		$$ = st
	}

full_opt:
	/*empty*/
	{
		$$ = false
	}
	| FULL
	{
		$$ = true
	}

show_body:
	show_scope_opt TABLES show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "tables", From: asIdent($3), Like: asExpr($4), Where: asExpr($5)}
	}
	| show_scope_opt DATABASES like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "databases", Like: asExpr($3), Where: asExpr($4)}
	}
	| show_scope_opt SCHEMAS like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "schemas", Like: asExpr($3), Where: asExpr($4)}
	}
	| show_scope_opt VARIABLES like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "variables", Like: asExpr($3), Where: asExpr($4)}
	}
	| show_scope_opt STATUS like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "status", Like: asExpr($3), Where: asExpr($4)}
	}
	| show_scope_opt PROFILES
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "profiles"}
	}
	| show_scope_opt PROFILE like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "profile", Like: asExpr($3), Where: asExpr($4)}
	}
	| show_scope_opt PROFILE FOR QUERY expr
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "profile", Mode: "FOR QUERY " + exprText(asExpr($5))}
	}
	| show_scope_opt NAME_OB table_name show_engine_mode
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "engine", Target: asIdent($3), Mode: asStr($4)}
	}
	| show_scope_opt MASTER STATUS
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "master status"}
	}
	| show_scope_opt FUNCTION STATUS like_opt
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "function status", Like: asExpr($4)}
	}
	| show_scope_opt PROCEDURE STATUS like_opt
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "procedure status", Like: asExpr($4)}
	}
	| show_scope_opt EVENTS show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Scope: asStr($1), Kind: "events", From: asIdent($3), Like: asExpr($4), Where: asExpr($5)}
	}

	| CREATE TABLE table_name
	{
		$$ = &ast.ShowStmt{Kind: "create table", Target: asIdent($3)}
	}
	| CREATE DATABASE table_name
	{
		$$ = &ast.ShowStmt{Kind: "create database", Target: asIdent($3)}
	}
	| CREATE VIEW table_name
	{
		$$ = &ast.ShowStmt{Kind: "create view", Target: asIdent($3)}
	}
	| CREATE PROCEDURE table_name
	{
		$$ = &ast.ShowStmt{Kind: "create procedure", Target: asIdent($3)}
	}
	| CREATE FUNCTION table_name
	{
		$$ = &ast.ShowStmt{Kind: "create function", Target: asIdent($3)}
	}
	| CREATE TRIGGER table_name
	{
		$$ = &ast.ShowStmt{Kind: "create trigger", Target: asIdent($3)}
	}
	| COLUMNS FROM table_name show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "columns", Target: asIdent($3), From: asIdent($4), Like: asExpr($5), Where: asExpr($6)}
	}
	| FIELDS FROM table_name show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "columns", Target: asIdent($3), From: asIdent($4), Like: asExpr($5), Where: asExpr($6)}
	}
	| INDEX FROM table_name show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "index", Target: asIdent($3), From: asIdent($4), Like: asExpr($5), Where: asExpr($6)}
	}
	| INDEXES FROM table_name show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "index", Target: asIdent($3), From: asIdent($4), Like: asExpr($5), Where: asExpr($6)}
	}
	| KEYS FROM table_name show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "index", Target: asIdent($3), From: asIdent($4), Like: asExpr($5), Where: asExpr($6)}
	}
	| WARNINGS
	{
		$$ = &ast.ShowStmt{Kind: "warnings"}
	}
	| ERRORS
	{
		$$ = &ast.ShowStmt{Kind: "errors"}
	}
	| GRANTS
	{
		$$ = &ast.ShowStmt{Kind: "grants"}
	}
	| GRANTS FOR user_spec
	{
		u := $3.(ast.UserSpec)
		$$ = &ast.ShowStmt{Kind: "grants", User: &u}
	}
	| CREATE NAME_OB user_spec
	{
		u := $3.(ast.UserSpec)
		$$ = &ast.ShowStmt{Kind: "create user", User: &u}
	}
	| PROCESSLIST
	{
		$$ = &ast.ShowStmt{Kind: "processlist"}
	}
	| COLLATION
	{
		$$ = &ast.ShowStmt{Kind: "collation"}
	}
	| ENGINES
	{
		$$ = &ast.ShowStmt{Kind: "engines"}
	}
	| STORAGE ENGINES
	{
		$$ = &ast.ShowStmt{Kind: "engines"}
	}
	| CHARSET like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "character set", Like: asExpr($2), Where: asExpr($3)}
	}
	| TABLE STATUS show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "table status", From: asIdent($3), Like: asExpr($4), Where: asExpr($5)}
	}
	| TRIGGERS show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "triggers", From: asIdent($2), Like: asExpr($3), Where: asExpr($4)}
	}
	| CHARACTER SET like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "character set", Like: asExpr($3), Where: asExpr($4)}
	}
	| PRIVILEGES
	{
		$$ = &ast.ShowStmt{Kind: "privileges"}
	}
	| PLUGINS
	{
		$$ = &ast.ShowStmt{Kind: "plugins"}
	}
	| PROFILES
	{
		$$ = &ast.ShowStmt{Kind: "profiles"}
	}
	| OPEN TABLES show_from_opt like_opt opt_show_where
	{
		$$ = &ast.ShowStmt{Kind: "open tables", From: asIdent($3), Like: asExpr($4), Where: asExpr($5)}
	}

show_engine_mode:
	STATUS
	{
		$$ = "STATUS"
	}
	| MUTEX
	{
		$$ = "MUTEX"
	}

opt_show_where:
	/*empty*/
	{
		$$ = nil
	}
	| WHERE expr
	{
		$$ = $2
	}
show_scope_opt:
	/*empty*/
	{
		$$ = ""
	}
	| GLOBAL
	{
		$$ = "GLOBAL"
	}
	| SESSION
	{
		$$ = "SESSION"
	}

show_from_opt:
	/*empty*/
	{
		$$ = nil
	}
	| FROM table_name
	{
		$$ = $2
	}

like_opt:
	/*empty*/
	{
		$$ = nil
	}
	| LIKE expr
	{
		$$ = $2
	}

/* ---------------- SET / USE ---------------- */

set_stmt:
	SET set_body
	{
		$$ = $2
	}

set_body:
	NAMES names_arg
	{
		$$ = &ast.SetStmt{Items: []*ast.SetItem{{Name: "NAMES", Value: asExpr($2)}}}
	}
	| opt_scope TRANSACTION transaction_characteristics
	{
		$$ = &ast.SetTransactionStmt{Scope: asStr($1), Characteristics: asStringList($3)}
	}
	| ROLE set_role_target
	{
		gateFeat(yylex, feature.FRoleManagement, "ROLE management requires OceanBase 4.2.4+")

		$$ = &ast.SetRoleStmt{Roles: asStringList($2)}
	}
	| PASSWORD opt_for_user COMP_EQ STRING_VALUE
	{
		sp := &ast.SetPasswordStmt{Password: asLitValue($4)}
		if $2 != nil {
			sp.User = $2.(ast.UserSpec)
			sp.HasUser = true
		}
		$$ = sp
	}
	| CATALOG NAME_OB
	{
		gateFeat(yylex, feature.FExternalCatalog, "CATALOG requires OceanBase 4.3+")

		$$ = &ast.CatalogStmt{Action: "SET", Name: asIdentName($2)}
	}
	| CHARACTER SET NAME_OB
	{
		$$ = &ast.SetStmt{Items: []*ast.SetItem{{Name: "CHARACTER SET " + asIdentName($3), NoValue: true}}}
	}
	| set_item_list
	{
		$$ = $1
	}

set_role_target:
	NONE
	{
		$$ = []string{"NONE"}
	}
	| DEFAULT
	{
		$$ = []string{"DEFAULT"}
	}
	| role_name_list

opt_for_user:
	/*empty*/
	{
		$$ = nil
	}
	| FOR user_spec
	{
		$$ = $2
	}

transaction_characteristics:
	transaction_characteristic
	{
		$$ = []string{$1.(string)}
	}
	| transaction_characteristics ',' transaction_characteristic
	{
		$$ = append(asStringList($1), $3.(string))
	}

transaction_characteristic:
	ISOLATION LEVEL isolation_level
	{
		$$ = "ISOLATION LEVEL " + asStr($3)
	}
	| READ WRITE
	{
		$$ = "READ WRITE"
	}
	| READ ONLY
	{
		$$ = "READ ONLY"
	}

isolation_level:
	READ UNCOMMITTED
	{
		$$ = "READ UNCOMMITTED"
	}
	| READ COMMITTED
	{
		$$ = "READ COMMITTED"
	}
	| REPEATABLE READ
	{
		$$ = "REPEATABLE READ"
	}
	| SERIALIZABLE
	{
		$$ = "SERIALIZABLE"
	}

set_item_list:
	set_item
	{
		$$ = &ast.SetStmt{Items: []*ast.SetItem{$1.(*ast.SetItem)}}
	}
	| set_item_list ',' set_item
	{
		st := $1.(*ast.SetStmt)
		st.Items = append(st.Items, $3.(*ast.SetItem))
		$$ = st
	}

set_item:
	opt_scope set_lhs COMP_EQ expr
	{
		it := &ast.SetItem{Name: asStr($2), Value: asExpr($4)}
		if s := asStr($1); s != "" {
			it.Scope = strings.ToLower(s)
		}
		$$ = it
	}

set_lhs:
	USER_VARIABLE
	{
		$$ = variableText($1.(*ast.Variable))
	}
	| SYSTEM_VARIABLE
	{
		$$ = variableText($1.(*ast.Variable))
	}
	| NAME_OB
	{
		$$ = asIdentName($1)
	}

opt_scope:
	/*empty*/
	{
		$$ = ""
	}
	| GLOBAL
	{
		$$ = "GLOBAL"
	}
	| SESSION
	{
		$$ = "SESSION"
	}
	| LOCAL
	{
		$$ = "LOCAL"
	}

names_arg:
	NAME_OB
	{
		$$ = &ast.Literal{LitType: ast.LitString, Value: asIdentName($1), Text: asIdentName($1)}
	}
	| STRING_VALUE
	{
		$$ = $1
	}

use_stmt:
	USE NAME_OB
	{
		$$ = &ast.UseDbStmt{Name: asIdentName($2)}
	}

/* ---------------- EXPLAIN / TXN / DESCRIBE ---------------- */

explain_stmt:
	EXPLAIN opt_explain_type inner_stmt
	{
		$$ = &ast.ExplainStmt{ExplainType: asStr($2), Stmt: asStmt($3)}
	}

opt_explain_type:
	/*empty*/
	{
		$$ = ""
	}
	| EXTENDED
	{
		gateFeat(yylex, feature.FExplainExtended, "EXPLAIN EXTENDED is supported from OceanBase 4.1 until 4.3")
		$$ = "EXTENDED"
	}
	| PRETTY
	{
		$$ = "PRETTY"
	}
	| FORMAT COMP_EQ NAME_OB
	{
		$$ = "FORMAT=" + strings.ToLower(asIdentName($3))
	}
	| NAME_OB
	{
		$$ = strings.ToUpper(asIdentName($1))
	}
	| FOR CONNECTION expr
	{
		$$ = "FOR CONNECTION " + exprText(asExpr($3))
	}

begin_stmt:
	BEGIN opt_work
	{
		$$ = &ast.BeginStmt{Work: asBool($2)}
	}
	| START TRANSACTION opt_start_read
	{
		$$ = &ast.BeginStmt{ReadOnly: asBool($3)}
	}

opt_work:
	/*empty*/
	{
		$$ = false
	}
	| WORK
	{
		$$ = true
	}

opt_start_read:
	/*empty*/
	{
		$$ = false
	}
	| READ ONLY
	{
		$$ = true
	}
	| READ WRITE
	{
		$$ = false
	}

commit_stmt:
	COMMIT opt_work
	{
		$$ = &ast.CommitStmt{Work: asBool($2)}
	}

rollback_stmt:
	ROLLBACK opt_work
	{
		$$ = &ast.RollbackStmt{Work: asBool($2)}
	}
	| ROLLBACK opt_work TO savepoint_kw NAME_OB
	{
		$$ = &ast.RollbackToSavepointStmt{Work: asBool($2), Name: asIdentName($5)}
	}

savepoint_kw:
	/*empty*/
	{
		$$ = nil
	}
	| SAVEPOINT
	{
		$$ = nil
	}

savepoint_stmt:
	SAVEPOINT NAME_OB
	{
		$$ = &ast.SavepointStmt{Name: asIdentName($2)}
	}

release_savepoint_stmt:
	RELEASE SAVEPOINT NAME_OB
	{
		$$ = &ast.ReleaseSavepointStmt{Name: asIdentName($3)}
	}

lock_tables_stmt:
	LOCK TABLES lock_table_list
	{
		$$ = &ast.LockTablesStmt{Tables: asLockTables($3)}
	}

lock_table_list:
	lock_table
	{
		$$ = []ast.LockTable{$1.(ast.LockTable)}
	}
	| lock_table_list ',' lock_table
	{
		$$ = append(asLockTables($1), $3.(ast.LockTable))
	}

lock_table:
	table_name opt_lock_alias lock_type
	{
		lt := ast.LockTable{Table: asIdent($1), LockType: asStr($3)}
		if $2 != nil {
			lt.Alias = asStr($2)
			lt.HasAlias = true
		}
		$$ = lt
	}

opt_lock_alias:
	/*empty*/
	{
		$$ = nil
	}
	| AS NAME_OB
	{
		$$ = asIdentName($2)
	}
	| NAME_OB
	{
		$$ = asIdentName($1)
	}

lock_type:
	READ
	{
		$$ = "READ"
	}
	| READ LOCAL
	{
		$$ = "READ LOCAL"
	}
	| WRITE
	{
		$$ = "WRITE"
	}
	| LOW_PRIORITY WRITE
	{
		$$ = "LOW_PRIORITY WRITE"
	}

unlock_tables_stmt:
	UNLOCK TABLES
	{
		$$ = &ast.UnlockTablesStmt{}
	}

describe_stmt:
	DESC table_name
	{
		$$ = &ast.ShowStmt{Kind: "describe table", Target: asIdent($2)}
	}
	| DESCRIBE table_name
	{
		$$ = &ast.ShowStmt{Kind: "describe table", Target: asIdent($2)}
	}
