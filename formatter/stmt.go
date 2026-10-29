package formatter

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
)

func (f *formatter) stmt(s ast.Statement) {
	switch v := s.(type) {
	case *ast.SelectStmt:
		f.selectStmt(v)
	case *ast.InsertStmt:
		f.insertStmt(v)
	case *ast.UpdateStmt:
		f.updateStmt(v)
	case *ast.DeleteStmt:
		f.deleteStmt(v)
	case *ast.CreateTableStmt:
		f.createTable(v)
	case *ast.CreateTableLikeStmt:
		f.printf("CREATE TABLE %s LIKE %s", f.ident(v.Table), f.ident(v.Like))
	case *ast.CreateDatabaseStmt:
		f.printf("CREATE DATABASE ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		if v.Charset != "" {
			f.printf(" DEFAULT CHARACTER SET %s", v.Charset)
		}
		if v.Collate != "" {
			f.printf(" COLLATE %s", v.Collate)
		}
	case *ast.DropDatabaseStmt:
		f.printf("DROP DATABASE ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
	case *ast.AlterDatabaseStmt:
		f.printf("ALTER DATABASE %s", f.quoteIfNeeded(v.Name))
		if v.Charset != "" {
			f.printf(" DEFAULT CHARACTER SET %s", v.Charset)
		}
		if v.Collate != "" {
			f.printf(" COLLATE %s", v.Collate)
		}
	case *ast.CreateSequenceStmt:
		f.raw("CREATE SEQUENCE ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		writeSequenceOpts(f, v.StartWith, v.IncrementBy, v.MinValue, v.MaxValue, v.Cache, v.Cycle, v.Order)
	case *ast.DropSequenceStmt:
		f.raw("DROP SEQUENCE ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
	case *ast.TransferPartitionStmt:
		f.raw("ALTER SYSTEM ")
		if v.Cancel {
			f.raw("CANCEL TRANSFER PARTITION ")
		} else {
			f.raw("TRANSFER PARTITION ")
		}
		if v.All {
			f.raw("ALL")
		} else {
			f.raw("TABLE_ID = ")
			f.expr(v.TableID)
			f.raw(", OBJECT_ID = ")
			f.expr(v.ObjectID)
			if !v.Cancel {
				f.raw(" TO LS ")
				f.expr(v.LS)
			}
		}
		if v.Tenant != "" {
			f.printf(" TENANT %s", f.quoteIfNeeded(v.Tenant))
		}
	case *ast.AlterSystemStmt:
		f.printf("ALTER SYSTEM SET %s = ", f.quoteIfNeeded(v.Name))
		f.expr(v.Value)
	case *ast.CCLLimitRuleStmt:
		f.printf("%s CONCURRENT_LIMITING_RULE ", v.Action)
		if v.Action == "DROP" && v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		if v.HasUser {
			f.printf(" ON %s TO ", f.ident(v.Table))
			f.userSpec(&v.User)
		}
		if v.ForName != "" {
			f.printf(" FOR %s", v.ForName)
		}
		f.options(v.Options)
	case *ast.SensitiveRuleStmt:
		f.printf("%s SENSITIVE RULE %s", v.Action, f.quoteIfNeeded(v.Name))
		if v.Action != "DROP" {
			f.printf(" ON %s", strings.Join(v.Fields, ", "))
			if v.Using != "" {
				f.printf(" USING %s", v.Using)
			}
		}
	case *ast.DblinkStmt:
		f.printf("%s DATABASE LINK %s", v.Action, f.quoteIfNeeded(v.Name))
		if v.Action == "CREATE" {
			f.printf(" CONNECT TO %s DATABASE %s IDENTIFIED BY %s",
				f.stringLiteral(v.User), f.stringLiteral(v.Database), f.stringLiteral(v.Password))
			if v.Host != "" {
				f.printf(" HOST %s", f.stringLiteral(v.Host))
				if v.Port != nil {
					f.raw(" PORT ")
					f.expr(v.Port)
				}
			}
			if v.Cluster != "" {
				f.printf(" CLUSTER %s", f.quoteIfNeeded(v.Cluster))
			}
		}
	case *ast.RoutineLoadStmt:
		f.printf("CREATE ROUTINE LOAD %s ON %s", f.quoteIfNeeded(v.Job), f.ident(v.Table))
		if len(v.LoadProps) > 0 {
			f.raw(" ")
			f.propList(v.LoadProps)
		}
		f.raw(" FROM STREAM (")
		f.propList(v.JobProps)
		f.raw(")")
	case *ast.RoutineLoadActionStmt:
		f.printf("%s ROUTINE LOAD FOR %s", v.Action, f.quoteIfNeeded(v.Job))
	case *ast.AuditStmt:
		if v.Audit {
			f.raw("AUDIT ")
		} else {
			f.raw("NOAUDIT ")
		}
		f.raw(strings.Join(v.Privileges, ", "))
		f.printf(" ON %s", v.Level)
	case *ast.DumpMemoryStmt:
		f.printf("DUMP %s ", v.Kind)
		if v.All {
			f.raw("ALL")
		} else {
			f.printf("%s = ", v.FirstName)
			f.expr(v.FirstVal)
			f.printf(", %s = ", v.SecondName)
			f.expr(v.SecondVal)
		}
	case *ast.ServiceStmt:
		f.raw("ALTER SYSTEM ")
		if v.Create {
			f.raw("CREATE ")
		} else {
			f.raw("DELETE ")
		}
		f.printf("SERVICE %s", f.quoteIfNeeded(v.Name))
		if v.Tenant != "" {
			f.printf(" TENANT %s", f.quoteIfNeeded(v.Tenant))
		}
	case *ast.StoredProcStmt:
		f.raw("CREATE ")
		if v.OrReplace {
			f.raw("OR REPLACE ")
		}
		if v.Function {
			f.printf("FUNCTION %s (", v.Name)
		} else {
			f.printf("PROCEDURE %s (", v.Name)
		}
		for i, p := range v.Parameters {
			if i > 0 {
				f.raw(", ")
			}
			if p.Mode == "OUT" || p.Mode == "INOUT" {
				f.printf("%s ", p.Mode)
			}
			f.printf("%s %s", f.quoteIfNeeded(p.Name), p.Type)
			if len(p.TypeArgs) > 0 {
				f.raw(" (")
				for j, a := range p.TypeArgs {
					if j > 0 {
						f.raw(", ")
					}
					f.raw(f.exprString(a))
				}
				f.raw(")")
			}
		}
		f.raw(")")
		if v.Function {
			f.printf(" RETURNS %s", v.Returns)
			if len(v.ReturnsArgs) > 0 {
				f.raw(" (")
				for i, a := range v.ReturnsArgs {
					if i > 0 {
						f.raw(", ")
					}
					f.raw(f.exprString(a))
				}
				f.raw(")")
			}
		}
		f.raw(" ")
		f.plBlock(v.Body)
	case *ast.DropRoutineStmt:
		if v.Function {
			f.raw("DROP FUNCTION ")
		} else {
			f.raw("DROP PROCEDURE ")
		}
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
	case *ast.DoStmt:
		f.raw("DO ")
		f.expr(v.Expr)
	case *ast.CallStmt:
		f.printf("CALL %s(", v.Name)
		for i, a := range v.Args {
			if i > 0 {
				f.raw(", ")
			}
			f.expr(a)
		}
		f.raw(")")
	case *ast.CatalogStmt:
		switch v.Action {
		case "CREATE":
			f.raw("CREATE EXTERNAL CATALOG ")
			if v.IfNotExists {
				f.raw("IF NOT EXISTS ")
			}
			f.raw(f.quoteIfNeeded(v.Name))
			f.raw(" PROPERTIES (")
			f.propList(v.Options)
			f.raw(")")
		case "DROP":
			f.raw("DROP CATALOG ")
			if v.IfExists {
				f.raw("IF EXISTS ")
			}
			f.raw(f.quoteIfNeeded(v.Name))
		case "SET":
			f.printf("SET CATALOG %s", f.quoteIfNeeded(v.Name))
		}
	case *ast.CreateLocationStmt:
		switch v.Action {
		case "DROP":
			f.printf("DROP LOCATION %s", f.quoteIfNeeded(v.Name))
		case "ALTER":
			f.printf("ALTER LOCATION %s SET URL %s", f.quoteIfNeeded(v.Name), f.stringLiteral(v.URL))
		default:
			f.raw("CREATE ")
			if v.OrReplace {
				f.raw("OR REPLACE ")
			}
			f.printf("LOCATION %s URL %s", f.quoteIfNeeded(v.Name), f.stringLiteral(v.URL))
		}
	case *ast.LocationUtilsStmt:
		f.printf("REMOVE FILES IN LOCATION %s", v.Path)
	case *ast.FlashbackStandbyLogStmt:
		f.raw("ALTER SYSTEM FLASHBACK STANDBY LOG TO SCN = ")
		f.expr(v.SCN)
		if v.Tenant != "" {
			f.printf(" TENANT %s", f.quoteIfNeeded(v.Tenant))
		}
	case *ast.SwitchoverTenantStmt:
		f.printf("ALTER SYSTEM %s", v.Clause)
		if v.Name != "" {
			f.printf(" TENANT %s", f.quoteIfNeeded(v.Name))
		}
	case *ast.RecoverTenantStmt:
		f.raw("ALTER SYSTEM RECOVER STANDBY")
		if v.Tenant != "" {
			f.printf(" TENANT %s", f.quoteIfNeeded(v.Tenant))
		}
		if v.Cancel {
			f.raw(" CANCEL")
		} else if v.Until == "UNLIMITED" {
			f.raw(" UNTIL UNLIMITED")
		} else if v.Until != "" {
			f.printf(" UNTIL %s ", v.Until)
			f.expr(v.UntilValue)
		}
	case *ast.AlterLSStmt:
		f.printf("ALTER SYSTEM %s LS", v.Op)
		f.options(v.Options)
		if v.Tenant != "" {
			f.printf(" TENANT %s", f.quoteIfNeeded(v.Tenant))
		}
	case *ast.CreateStandbyTenantStmt:
		f.raw("CREATE STANDBY TENANT ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		if v.RestoreSource != nil {
			f.raw(" ")
			f.optionKV(v.RestoreSource)
		}
		f.options(v.Options)
	case *ast.CloneTenantStmt:
		f.raw("CREATE TENANT ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.printf("%s FROM %s", f.quoteIfNeeded(v.Name), f.quoteIfNeeded(v.From))
		if v.Snapshot != "" {
			f.printf(" USING SNAPSHOT %s", f.quoteIfNeeded(v.Snapshot))
		}
		f.raw(" WITH")
		for i := range v.Options {
			f.raw(" ")
			f.optionKV(&v.Options[i])
			if i < len(v.Options)-1 {
				f.raw(", ")
			}
		}
	case *ast.RestorePointStmt:
		f.printf("%s RESTORE POINT %s", v.Action, f.quoteIfNeeded(v.Name))
	case *ast.TenantSnapshotStmt:
		f.printf("%s SNAPSHOT %s", v.Action, f.quoteIfNeeded(v.Name))
		if v.HasTenant {
			f.printf(" FOR TENANT %s", f.quoteIfNeeded(v.Tenant))
		}
	case *ast.AlterMViewStmt:
		f.printf("ALTER MATERIALIZED VIEW %s", f.ident(v.Name))
		if v.Refresh != "" {
			f.printf(" REFRESH %s", v.Refresh)
		}
		if v.Mode != "" {
			f.printf(" %s", v.Mode)
		}
	case *ast.MLogStmt:
		f.printf("%s MATERIALIZED VIEW LOG ", v.Action)
		if v.Action == "DROP" && v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.printf("ON %s", f.ident(v.Table))
		if v.WithPrimary {
			f.raw(" WITH PRIMARY KEY")
		}
		if v.IncludingNew {
			f.raw(" INCLUDING NEW VALUES")
		}
	case *ast.CreateOutlineStmt:
		f.printf("CREATE OUTLINE %s ON ", f.ident(v.Name))
		if v.Select != nil {
			f.selectStmt(v.Select)
		} else {
			f.raw(f.stringLiteral(v.Text))
		}
	case *ast.AlterOutlineStmt:
		f.printf("ALTER OUTLINE %s ON ", f.ident(v.Name))
		if v.Select != nil {
			f.selectStmt(v.Select)
		} else {
			f.raw(f.stringLiteral(v.Text))
		}
	case *ast.DropOutlineStmt:
		f.raw("DROP OUTLINE ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.ident(v.Name))
	case *ast.AlterTablespaceStmt:
		f.printf("ALTER TABLESPACE %s", f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.RotateMasterKeyStmt:
		f.raw("ALTER SYSTEM ROTATE")
		if v.Innodb {
			f.raw(" INNODB")
		}
		f.raw(" MASTER KEY")
	case *ast.AlterViewStmt:
		f.printf("ALTER VIEW %s", f.ident(v.Name))
		if len(v.Columns) > 0 {
			f.raw(" ")
			f.colsParen(v.Columns)
		}
		f.raw(" AS ")
		f.selectStmt(v.Select)
	case *ast.AlterSequenceStmt:
		f.printf("ALTER SEQUENCE %s", f.quoteIfNeeded(v.Name))
		writeSequenceOpts(f, v.StartWith, v.IncrementBy, v.MinValue, v.MaxValue, v.Cache, v.Cycle, v.Order)
	case *ast.AlterTablegroupStmt:
		f.printf("ALTER TABLEGROUP %s", f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.AlterResourceUnitStmt:
		f.printf("ALTER RESOURCE UNIT %s", f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.AlterResourcePoolStmt:
		f.printf("ALTER RESOURCE POOL %s", f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.RenameUserStmt:
		f.raw("RENAME USER ")
		f.userSpec(&v.Old)
		f.raw(" TO ")
		f.userSpec(&v.New)
	case *ast.LockUserStmt:
		if v.Unlock {
			f.raw("UNLOCK USER ")
		} else {
			f.raw("LOCK USER ")
		}
		f.userSpecs(v.Users)
		if v.AccountLock != "" {
			f.printf(" ACCOUNT %s", stringsUpper(v.AccountLock))
		}
	case *ast.CreateViewStmt:
		f.raw("CREATE ")
		if v.OrReplace {
			f.raw("OR REPLACE ")
		}
		f.printf("VIEW %s", f.ident(v.Name))
		if len(v.Columns) > 0 {
			f.raw(" ")
			f.colsParen(v.Columns)
		}
		f.raw(" AS ")
		f.selectStmt(v.Select)
	case *ast.DropViewStmt:
		f.raw("DROP VIEW ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		for i, n := range v.Names {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(n))
		}
	case *ast.GrantStmt:
		f.raw("GRANT ")
		f.raw(strings.Join(v.Privileges, ", "))
		f.printf(" ON %s TO ", v.Level)
		f.userSpecs(v.Users)
		if v.WithGrant {
			f.raw(" WITH GRANT OPTION")
		}
	case *ast.RevokeStmt:
		f.raw("REVOKE ")
		f.raw(strings.Join(v.Privileges, ", "))
		f.printf(" ON %s FROM ", v.Level)
		f.userSpecs(v.Users)
	case *ast.CreateUserStmt:
		f.raw("CREATE USER ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.userSpecs(v.Users)
	case *ast.DropUserStmt:
		f.raw("DROP USER ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.userSpecs(v.Users)
	case *ast.PrepareStmt:
		f.printf("PREPARE %s FROM ", f.quoteIfNeeded(v.Name))
		f.expr(v.Source)
	case *ast.ExecuteStmt:
		f.printf("EXECUTE %s", f.quoteIfNeeded(v.Name))
		if len(v.Using) > 0 {
			f.raw(" USING ")
			for i, u := range v.Using {
				if i > 0 {
					f.raw(", ")
				}
				f.raw(f.variableString(u))
			}
		}
	case *ast.DeallocateStmt:
		f.printf("DEALLOCATE PREPARE %s", f.quoteIfNeeded(v.Name))
	case *ast.CreateRoleStmt:
		f.raw("CREATE ROLE ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(strings.Join(v.Roles, ", "))
	case *ast.DropRoleStmt:
		f.raw("DROP ROLE ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(strings.Join(v.Roles, ", "))
	case *ast.SetRoleStmt:
		f.raw("SET ROLE ")
		f.raw(strings.Join(v.Roles, ", "))
	case *ast.SetPasswordStmt:
		f.raw("SET PASSWORD")
		if v.HasUser {
			f.raw(" FOR ")
			f.userSpec(&v.User)
		}
		f.raw(" = '" + strings.ReplaceAll(v.Password, "'", "''") + "'")
	case *ast.KillStmt:
		f.raw("KILL ")
		if v.Kind != "" {
			f.printf("%s ", v.Kind)
		}
		f.expr(v.Target)
	case *ast.HelpStmt:
		f.raw("HELP ")
		f.expr(v.Target)
	case *ast.TableMaintStmt:
		f.printf("%s TABLE ", v.Kind)
		for i, t := range v.Tables {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(t))
		}
		if v.MethodOpt != "" {
			f.printf(" METHOD_OPT %s", f.stringLiteral(v.MethodOpt))
		}
	case *ast.InstallPluginStmt:
		f.printf("INSTALL PLUGIN %s SONAME %s", f.quoteIfNeeded(v.Name), f.stringLiteral(v.Soname))
	case *ast.UninstallPluginStmt:
		f.printf("UNINSTALL PLUGIN %s", f.quoteIfNeeded(v.Name))
	case *ast.HandlerStmt:
		f.printf("HANDLER %s %s", f.ident(v.Table), v.Op)
		if v.Op == "OPEN" && v.Alias != "" {
			f.printf(" AS %s", f.quoteIfNeeded(v.Alias))
		}
		if v.Op == "READ" && v.Direction != "" {
			f.printf(" %s", v.Direction)
			if len(v.Columns) > 0 {
				f.raw(" ")
				f.colsParen(v.Columns)
			}
			if v.Where != nil {
				f.raw(" WHERE ")
				f.expr(v.Where)
			}
			if v.Limit != nil {
				f.raw(" ")
				f.limit(v.Limit)
			}
		}
	case *ast.RenameTableStmt:
		f.raw("RENAME TABLE ")
		for i, p := range v.Pairs {
			if i > 0 {
				f.raw(", ")
			}
			f.printf("%s TO %s", f.ident(p.From), f.ident(p.To))
		}
	case *ast.LoadDataStmt:
		f.raw("LOAD DATA")
		if v.Hint != "" {
			f.printf(" %s", v.Hint)
		}
		if v.Local {
			f.raw(" LOCAL")
		}
		f.printf(" INFILE %s", f.stringLiteral(v.Infile))
		if v.Replace {
			f.raw(" REPLACE")
		} else if v.Ignore {
			f.raw(" IGNORE")
		}
		f.printf(" INTO TABLE %s", f.ident(v.Table))
		if v.FieldsTerminatedBy != "" {
			f.printf(" FIELDS TERMINATED BY %s", f.stringLiteral(v.FieldsTerminatedBy))
		}
		if v.LinesTerminatedBy != "" {
			f.printf(" LINES TERMINATED BY %s", f.stringLiteral(v.LinesTerminatedBy))
		}
		if v.IgnoreLines != nil {
			f.raw(" IGNORE ")
			f.expr(v.IgnoreLines)
			f.raw(" LINES")
		}
		if len(v.Columns) > 0 {
			f.raw(" ")
			f.colsParen(v.Columns)
		}
	case *ast.CreateTenantStmt:
		f.raw("CREATE TENANT ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.AlterTenantStmt:
		f.printf("ALTER TENANT %s", f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.DropTenantStmt:
		f.raw("DROP TENANT ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		if v.Purge {
			f.raw(" PURGE")
		}
	case *ast.CreateResourceUnitStmt:
		f.raw("CREATE RESOURCE UNIT ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.CreateResourcePoolStmt:
		f.raw("CREATE RESOURCE POOL ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.DropResourceUnitStmt:
		f.raw("DROP RESOURCE UNIT ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
	case *ast.DropResourcePoolStmt:
		f.raw("DROP RESOURCE POOL ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
	case *ast.CreateTablegroupStmt:
		f.raw("CREATE TABLEGROUP ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.DropTablegroupStmt:
		f.raw("DROP TABLEGROUP ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
	case *ast.CreateTablespaceStmt:
		f.printf("CREATE TABLESPACE %s", f.quoteIfNeeded(v.Name))
		f.options(v.Options)
	case *ast.DropTablespaceStmt:
		f.raw("DROP TABLESPACE ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(v.Name))
	case *ast.SignalStmt:
		if v.IsResignal {
			f.raw("RESIGNAL")
			if v.HasSqlState {
				f.printf(" SQLSTATE %s", f.stringLiteral(v.SqlState))
			}
		} else {
			f.printf("SIGNAL SQLSTATE %s", f.stringLiteral(v.SqlState))
		}
		if len(v.Items) > 0 {
			f.raw(" SET ")
			for i, it := range v.Items {
				if i > 0 {
					f.raw(", ")
				}
				f.printf("%s = ", it.Name)
				f.expr(it.Value)
			}
		}
	case *ast.GetDiagnosticsStmt:
		f.raw("GET ")
		if v.Scope != "" {
			f.printf("%s ", v.Scope)
		}
		f.raw("DIAGNOSTICS")
		if v.HasCondition {
			f.raw(" CONDITION ")
			f.expr(v.Condition)
		}
		f.raw(" ")
		for i, it := range v.Items {
			if i > 0 {
				f.raw(", ")
			}
			f.printf("%s = %s", it.Target, it.Source)
		}
	case *ast.CreateMViewStmt:
		f.printf("CREATE MATERIALIZED VIEW %s", f.ident(v.Name))
		if len(v.Columns) > 0 {
			f.raw(" ")
			f.colsParen(v.Columns)
		}
		if v.Refresh != "" {
			f.printf(" REFRESH %s", v.Refresh)
		}
		if v.Mode != "" {
			f.printf(" %s", v.Mode)
		}
		f.raw(" AS ")
		f.selectStmt(v.Select)
	case *ast.DropMViewStmt:
		f.raw("DROP MATERIALIZED VIEW ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		for i, n := range v.Names {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(n))
		}
	case *ast.MergeStmt:
		f.raw("MERGE INTO ")
		f.tableRef(v.Target)
		f.raw(" USING ")
		f.tableRef(v.Source)
		f.raw(" ON ")
		f.expr(v.On)
		for i := range v.Whens {
			f.mergeWhen(&v.Whens[i])
		}
	case *ast.FlashbackTableStmt:
		f.raw("FLASHBACK TABLE ")
		for i, t := range v.Tables {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(t))
		}
		if v.ToBeforeDrop {
			f.raw(" TO BEFORE DROP")
			if v.RenameTo != nil {
				f.printf(" RENAME TO %s", f.ident(v.RenameTo))
			}
		} else if v.Snapshot != "" {
			f.printf(" TO SNAPSHOT %s", f.quoteIfNeeded(v.Snapshot))
		}
	case *ast.XAStmt:
		f.printf("XA %s", v.Action)
		for i, e := range v.XID {
			if i == 0 {
				f.raw(" ")
			} else {
				f.raw(", ")
			}
			f.expr(e)
		}
		if v.Option != "" {
			f.printf(" %s", v.Option)
		}
	case *ast.PurgeStmt:
		switch v.Kind {
		case "INDEX":
			f.printf("PURGE INDEX %s ON %s", f.ident(v.Index), f.ident(v.Table))
		case "TABLE":
			f.printf("PURGE TABLE %s", f.ident(v.Table))
		default:
			f.printf("PURGE %s", v.Kind)
		}
	case *ast.FlushStmt:
		f.raw("FLUSH")
		if v.Local {
			f.raw(" LOCAL")
		}
		f.printf(" %s", v.Kind)
		if len(v.Tables) > 0 {
			f.raw(" ")
			for i, t := range v.Tables {
				if i > 0 {
					f.raw(", ")
				}
				f.raw(f.ident(t))
			}
		}
		if v.Option != "" {
			f.printf(" %s", v.Option)
		}
	case *ast.CreateIndexStmt:
		f.printf("CREATE ")
		if v.Unique {
			f.raw("UNIQUE ")
		}
		f.raw("INDEX ")
		if v.IfNotExists {
			f.raw("IF NOT EXISTS ")
		}
		f.printf("%s ", f.ident(v.Name))
		if v.Method != "" {
			f.printf("USING %s ", v.Method)
		}
		f.printf("ON %s ", f.ident(v.Table))
		f.indexColsParen(v.Columns)
		if len(v.Storing) > 0 {
			f.raw(" STORING (")
			f.idents(v.Storing)
			f.raw(")")
		}
		if len(v.WithOpts) > 0 {
			f.raw(" WITH (" + strings.Join(v.WithOpts, ", ") + ")")
		}
		for _, o := range v.Options {
			f.printf(" %s", o)
		}
		if v.Partition != nil {
			f.raw(" ")
			f.partitionSpec(v.Partition)
		}
		if v.ColumnGroup != "" {
			f.printf(" %s", v.ColumnGroup)
		}
	case *ast.DropTableStmt:
		f.printf("DROP TABLE ")
		if v.IfExists {
			f.raw("IF EXISTS ")
		}
		for i, t := range v.Tables {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(t))
		}
	case *ast.DropIndexStmt:
		f.printf("DROP INDEX %s", f.ident(v.Name))
		if v.Table != nil {
			f.printf(" ON %s", f.ident(v.Table))
		}
	case *ast.AlterTableStmt:
		f.alterTable(v)
	case *ast.TruncateStmt:
		f.printf("TRUNCATE TABLE %s", f.ident(v.Table))
	case *ast.ShowStmt:
		f.showStmt(v)
	case *ast.SetStmt:
		f.setStmt(v)
	case *ast.UseDbStmt:
		f.printf("USE %s", f.quoteIfNeeded(v.Name))
	case *ast.ExplainStmt:
		f.printf("EXPLAIN")
		if v.ExplainType != "" {
			f.printf(" %s", v.ExplainType)
		}
		f.raw(" ")
		if v.Stmt != nil {
			f.stmt(v.Stmt)
		}
	case *ast.BeginStmt:
		f.raw("BEGIN")
	case *ast.CommitStmt:
		f.raw("COMMIT")
		if v.Work {
			f.raw(" WORK")
		}
		if v.Chain {
			f.raw(" AND CHAIN")
		}
		if v.Release {
			f.raw(" AND RELEASE")
		}
	case *ast.RollbackStmt:
		f.raw("ROLLBACK")
		if v.Work {
			f.raw(" WORK")
		}
		if v.Chain {
			f.raw(" AND CHAIN")
		}
		if v.Release {
			f.raw(" AND RELEASE")
		}
	case *ast.AlterUserStmt:
		f.raw("ALTER USER ")
		f.userSpecs(v.Users)
		if v.Password != "" {
			if v.Plugin != "" {
				f.printf(" IDENTIFIED WITH %s BY ", stringsUpper(v.Plugin))
				f.raw(f.stringLiteral(v.Password))
			} else if v.Password == "RANDOM PASSWORD" {
				f.raw(" IDENTIFIED BY RANDOM PASSWORD")
			} else {
				f.raw(" IDENTIFIED BY ")
				f.raw(f.stringLiteral(v.Password))
			}
		} else if v.Plugin != "" {
			f.printf(" IDENTIFIED WITH %s", stringsUpper(v.Plugin))
		}
		if v.Account != "" {
			f.printf(" ACCOUNT %s", stringsUpper(v.Account))
		}
	case *ast.SetTransactionStmt:
		f.raw("SET ")
		if v.Scope != "" {
			f.printf("%s ", stringsUpper(v.Scope))
		}
		f.raw("TRANSACTION")
		for i, c := range v.Characteristics {
			if i == 0 {
				f.raw(" ")
			} else {
				f.raw(", ")
			}
			f.raw(c)
		}
	case *ast.SavepointStmt:
		f.printf("SAVEPOINT %s", f.quoteIfNeeded(v.Name))
	case *ast.ReleaseSavepointStmt:
		f.printf("RELEASE SAVEPOINT %s", f.quoteIfNeeded(v.Name))
	case *ast.RollbackToSavepointStmt:
		f.raw("ROLLBACK")
		if v.Work {
			f.raw(" WORK")
		}
		f.printf(" TO SAVEPOINT %s", f.quoteIfNeeded(v.Name))
	case *ast.LockTablesStmt:
		f.raw("LOCK TABLES ")
		for i, t := range v.Tables {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(t.Table))
			if t.HasAlias {
				f.printf(" AS %s", f.quoteIfNeeded(t.Alias))
			}
			f.printf(" %s", t.LockType)
		}
	case *ast.UnlockTablesStmt:
		f.raw("UNLOCK TABLES")
	default:
		f.nodeAny(s)
	}
}
func writeSequenceOpts(f *formatter, startWith, incrementBy, minValue, maxValue, cache ast.Expression, cycle, order string) {
	if startWith != nil {
		f.raw(" START WITH ")
		f.expr(startWith)
	}
	if incrementBy != nil {
		f.raw(" INCREMENT BY ")
		f.expr(incrementBy)
	}
	if isNoValue(minValue) {
		f.raw(" " + minValue.(*ast.Literal).Text)
	} else if minValue != nil {
		f.raw(" MINVALUE ")
		f.expr(minValue)
	}
	if isNoValue(maxValue) {
		f.raw(" " + maxValue.(*ast.Literal).Text)
	} else if maxValue != nil {
		f.raw(" MAXVALUE ")
		f.expr(maxValue)
	}
	if isNoValue(cache) {
		f.raw(" " + cache.(*ast.Literal).Text)
	} else if cache != nil {
		f.raw(" CACHE ")
		f.expr(cache)
	}
	if cycle != "" {
		f.printf(" %s", cycle)
	}
	if order != "" {
		f.printf(" %s", order)
	}
}

// isNoValue 判断是否为 NOCACHE/NOMINVALUE/NOMAXVALUE 等占位字面量。
func isNoValue(e ast.Expression) bool {
	if l, ok := e.(*ast.Literal); ok && l.LitType == ast.LitNull && l.Text != "" {
		return true
	}
	return false
}

// showKindNeedsFrom reports whether the SHOW kind takes a FROM <table> target
// (SHOW COLUMNS FROM t, SHOW INDEX FROM t).
func showKindNeedsFrom(kind string) bool {
	switch kind {
	case "columns", "index":
		return true
	}
	return false
}
