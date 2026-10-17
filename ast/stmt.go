package ast

// UserSpec is a user@host account reference.
type UserSpec struct {
	Name        string
	Host        string
	HasHost     bool
	Password    string
	HasPassword bool
	Plugin      string // IDENTIFIED WITH <plugin>
}

// CreateUserStmt is CREATE USER.
type CreateUserStmt struct {
	Base
	IfNotExists bool
	Users       []UserSpec
}

// DropUserStmt is DROP USER.
type DropUserStmt struct {
	Base
	IfExists bool
	Users    []UserSpec
}

// Option is a generic NAME = value | NAME = (v1, v2) option (cluster admin).
type Option struct {
	Name   string
	Value  Expression
	Values []Expression
}

// CreateTenantStmt is CREATE TENANT.
type CreateTenantStmt struct {
	Base
	IfNotExists bool
	Name        string
	Options     []Option
}

// DropTenantStmt is DROP TENANT.
type DropTenantStmt struct {
	Base
	IfExists bool
	Name     string
	Purge    bool
}

// AlterTenantStmt is ALTER TENANT.
type AlterTenantStmt struct {
	Base
	Name    string
	Options []Option
}

// CreateResourceUnitStmt is CREATE RESOURCE UNIT.
type CreateResourceUnitStmt struct {
	Base
	IfNotExists bool
	Name        string
	Options     []Option
}

// DropResourceUnitStmt is DROP RESOURCE UNIT.
type DropResourceUnitStmt struct {
	Base
	IfExists bool
	Name     string
}

// CreateResourcePoolStmt is CREATE RESOURCE POOL.
type CreateResourcePoolStmt struct {
	Base
	IfNotExists bool
	Name        string
	Options     []Option
}

// DropResourcePoolStmt is DROP RESOURCE POOL.
type DropResourcePoolStmt struct {
	Base
	IfExists bool
	Name     string
}

// SignalItem is one NAME = value in SIGNAL/RESIGNAL SET.
type SignalItem struct {
	Name  string
	Value Expression
}

// SignalStmt is SIGNAL SQLSTATE ... [SET ...] / RESIGNAL ....
type SignalStmt struct {
	Base
	IsResignal  bool
	SqlState    string
	HasSqlState bool
	Items       []SignalItem
}

// DiagItem is one target = source in GET DIAGNOSTICS.
type DiagItem struct {
	Target string
	Source string
}

// GetDiagnosticsStmt is GET [CURRENT|STACKED] DIAGNOSTICS [CONDITION n] ....
type GetDiagnosticsStmt struct {
	Base
	Scope        string
	Condition    Expression
	HasCondition bool
	Items        []DiagItem
}

// CCLLimitRuleStmt is CREATE|ALTER|DROP CONCURRENT_LIMITING_RULE.
type CCLLimitRuleStmt struct {
	Base
	Action   string // "CREATE","ALTER","DROP"
	IfExists bool
	Name     string
	Table    *Identifier
	User     UserSpec
	HasUser  bool
	ForName  string
	Options  []Option
}

// SensitiveRuleStmt is CREATE|ALTER|DROP SENSITIVE RULE.
type SensitiveRuleStmt struct {
	Base
	Action string
	Name   string
	Fields []string
	Using  string
}

// DblinkStmt is CREATE|DROP DATABASE LINK.
type DblinkStmt struct {
	Base
	Action   string
	Name     string
	User     string
	Database string
	Password string
	Host     string
	Port     Expression
	Cluster  string
}

// RoutineLoadStmt is CREATE ROUTINE LOAD job ON table … FROM STREAM (props).
type RoutineLoadStmt struct {
	Base
	Job       string
	Table     *Identifier
	LoadProps []Option
	JobProps  []Option
}

// RoutineLoadActionStmt is PAUSE|RESUME|STOP ROUTINE LOAD FOR job.
type RoutineLoadActionStmt struct {
	Base
	Action string // "PAUSE","RESUME","STOP"
	Job    string
}

// InstallPluginStmt is INSTALL PLUGIN name SONAME 'file'.
type InstallPluginStmt struct {
	Base
	Name   string
	Soname string
}

// UninstallPluginStmt is UNINSTALL PLUGIN name.
type UninstallPluginStmt struct {
	Base
	Name string
}

// HandlerStmt is HANDLER table OPEN|CLOSE|READ ….
type HandlerStmt struct {
	Base
	Op        string // "OPEN","CLOSE","READ"
	Table     *Identifier
	Alias     string
	Direction string // "FIRST","NEXT","PREV","LAST"
	Columns   []*Identifier
	Where     Expression
	Limit     *LimitClause
}

// AuditStmt is AUDIT|NOAUDIT priv[,…] ON level.
type AuditStmt struct {
	Base
	Audit      bool // true = AUDIT, false = NOAUDIT
	Privileges []string
	Level      string
}

// DumpMemoryStmt is DUMP ENTITY|CHUNK ALL | DUMP ENTITY|CHUNK name = value, ….
type DumpMemoryStmt struct {
	Base
	Kind       string // "ENTITY","CHUNK"
	All        bool
	FirstName  string
	FirstVal   Expression
	SecondName string
	SecondVal  Expression
}

// ServiceStmt is ALTER SYSTEM CREATE|DELETE SERVICE name [TENANT t].
type ServiceStmt struct {
	Base
	Create bool
	Name   string
	Tenant string
}

// CatalogStmt is CREATE|DROP|SET CATALOG.
type CatalogStmt struct {
	Base
	Action      string // "CREATE","DROP","SET"
	IfNotExists bool
	IfExists    bool
	Name        string
	Options     []Option
}

// CreateLocationStmt is CREATE|ALTER|DROP LOCATION.
type CreateLocationStmt struct {
	Base
	Action    string // "CREATE","ALTER","DROP"
	OrReplace bool
	Name      string
	URL       string
}

// LocationUtilsStmt is REMOVE FILES IN LOCATION @path.
type LocationUtilsStmt struct {
	Base
	Path string
}

// FlashbackStandbyLogStmt is ALTER SYSTEM FLASHBACK STANDBY LOG TO SCN = n.
type FlashbackStandbyLogStmt struct {
	Base
	SCN    Expression
	Tenant string
}

// SwitchoverTenantStmt is ALTER SYSTEM {ACTIVATE STANDBY|SWITCHOVER TO …} TENANT.
type SwitchoverTenantStmt struct {
	Base
	Clause string // "ACTIVATE STANDBY","SWITCHOVER TO PRIMARY","SWITCHOVER TO STANDBY"
	Name   string
}

// RecoverTenantStmt is ALTER SYSTEM RECOVER STANDBY [TENANT t] [UNTIL …|CANCEL].
type RecoverTenantStmt struct {
	Base
	Tenant     string
	Until      string // "UNLIMITED","TIME","SCN"
	UntilValue Expression
	Cancel     bool
}

// AlterLSStmt is ALTER SYSTEM CREATE|MODIFY LS ….
type AlterLSStmt struct {
	Base
	Op      string // "CREATE","MODIFY"
	Options []Option
	Tenant  string
}

// CreateStandbyTenantStmt is CREATE STANDBY TENANT.
type CreateStandbyTenantStmt struct {
	Base
	IfNotExists   bool
	Name          string
	RestoreSource *Option
	Options       []Option
}

// CloneTenantStmt is CREATE TENANT t FROM src [USING SNAPSHOT s] WITH opts.
type CloneTenantStmt struct {
	Base
	IfNotExists bool
	Name        string
	From        string
	Snapshot    string
	Options     []Option
}

// RestorePointStmt is CREATE|DROP RESTORE POINT.
type RestorePointStmt struct {
	Base
	Action string
	Name   string
}

// TenantSnapshotStmt is CREATE|DROP SNAPSHOT [FOR TENANT t].
type TenantSnapshotStmt struct {
	Base
	Action    string
	Name      string
	Tenant    string
	HasTenant bool
}

// RotateMasterKeyStmt is ALTER SYSTEM ROTATE [INNODB] MASTER KEY.
type RotateMasterKeyStmt struct {
	Base
	Innodb bool
}

// AlterViewStmt is ALTER VIEW.
type AlterViewStmt struct {
	Base
	Name    *Identifier
	Columns []*Identifier
	Select  *SelectStmt
}

// AlterSequenceStmt is ALTER SEQUENCE.
type AlterSequenceStmt struct {
	Base
	Name        string
	StartWith   Expression
	IncrementBy Expression
	MinValue    Expression
	MaxValue    Expression
	Cache       Expression
	Cycle       string // "", "CYCLE", "NOCYCLE"
	Order       string // "", "ORDER", "NOORDER"
}

// AlterResourceUnitStmt is ALTER RESOURCE UNIT.
type AlterResourceUnitStmt struct {
	Base
	Name    string
	Options []Option
}

// AlterResourcePoolStmt is ALTER RESOURCE POOL.
type AlterResourcePoolStmt struct {
	Base
	Name    string
	Options []Option
}

// RenameUserStmt is RENAME USER old TO new.
type RenameUserStmt struct {
	Base
	Old UserSpec
	New UserSpec
}

// LockUserStmt is LOCK USER / UNLOCK USER.
type LockUserStmt struct {
	Base
	Unlock bool
	Users  []UserSpec
}

// XAStmt is an XA transaction control statement.
type XAStmt struct {
	Base
	Action string       // "START","END","PREPARE","COMMIT","ROLLBACK","RECOVER"
	XID    []Expression // up to 3 parts (gtrid[, bqual[, formatID]])
	Option string       // "JOIN","RESUME","SUSPEND","SUSPEND FOR MIGRATE","ONE PHASE","CONVERT XID"
}

// PurgeStmt is PURGE TABLE|INDEX|RECYCLEBIN.
type PurgeStmt struct {
	Base
	Kind  string // "TABLE","INDEX","RECYCLEBIN"
	Table *Identifier
	Index *Identifier
}

// LoadDataStmt is LOAD DATA [LOCAL] INFILE … .
type LoadDataStmt struct {
	Base
	Hint               string // /*+ DIRECT ... */ direct-load hint, if any
	Local              bool
	Infile             string
	Replace            bool
	Ignore             bool
	Table              *Identifier
	FieldsTerminatedBy string
	LinesTerminatedBy  string
	IgnoreLines        Expression
	Columns            []*Identifier
}

// FlushStmt is FLUSH [LOCAL] {TABLES …|PRIVILEGES|STATUS}.
type FlushStmt struct {
	Base
	Local  bool
	Kind   string // "TABLES", "PRIVILEGES", "STATUS"
	Tables []*Identifier
	Option string // "", "WITH READ LOCK", "FOR EXPORT"
}

// CreateRoleStmt is CREATE ROLE.
type CreateRoleStmt struct {
	Base
	IfNotExists bool
	Roles       []string
}

// DropRoleStmt is DROP ROLE.
type DropRoleStmt struct {
	Base
	IfExists bool
	Roles    []string
}

// SetRoleStmt is SET ROLE {role[,…] | NONE | DEFAULT}.
type SetRoleStmt struct {
	Base
	Roles []string
}

// SetPasswordStmt is SET PASSWORD [FOR 'u'@'h'] = 'pwd'.
type SetPasswordStmt struct {
	Base
	User     UserSpec
	HasUser  bool
	Password string
}

// KillStmt is KILL [CONNECTION|QUERY] n.
type KillStmt struct {
	Base
	Kind   string
	Target Expression
}

// HelpStmt is HELP expr.
type HelpStmt struct {
	Base
	Target Expression
}

// TableMaintStmt is ANALYZE|OPTIMIZE|REPAIR TABLE t[, …].
type TableMaintStmt struct {
	Base
	Kind      string // "ANALYZE", "OPTIMIZE", "REPAIR"
	Tables    []*Identifier
	MethodOpt string // ANALYZE TABLE … METHOD_OPT '…'
}

// PrepareStmt is PREPARE name FROM source.
type PrepareStmt struct {
	Base
	Name   string
	Source Expression
}

// ExecuteStmt is EXECUTE name [USING @a, @b, ...].
type ExecuteStmt struct {
	Base
	Name  string
	Using []*Variable
}

// DeallocateStmt is DEALLOCATE PREPARE name.
type DeallocateStmt struct {
	Base
	Name string
}

// GrantStmt is GRANT ... ON ... TO ....
type GrantStmt struct {
	Base
	Privileges []string // e.g. "SELECT", "ALL", "SELECT (a)"
	Level      string   // "*", "*.*", "db.*", "db.tbl"
	Users      []UserSpec
	WithGrant  bool
}

// RevokeStmt is REVOKE ... ON ... FROM ....
type RevokeStmt struct {
	Base
	Privileges []string
	Level      string
	Users      []UserSpec
}

// ShowStmt is a generic SHOW statement.
type ShowStmt struct {
	Base
	Scope  string // "GLOBAL"/"SESSION" for SHOW VARIABLES
	Kind   string // "tables","databases","create table","variables","columns","warnings","status","index","grants","processlist","collation","engines","schemas","table status",...
	Full   bool
	Target *Identifier // the object to describe (e.g. table for SHOW CREATE TABLE / SHOW COLUMNS FROM)
	From   *Identifier // FROM/IN qualifier (e.g. SHOW TABLES FROM db)
	Like   Expression  // LIKE 'pattern'
	Where  Expression  // WHERE clause
	User   *UserSpec   // SHOW GRANTS FOR 'u'@'h' / SHOW CREATE USER
	Mode   string      // SHOW ENGINE <n> STATUS|MUTEX, SHOW PROFILE FOR QUERY n
}

// SetStmt is SET ... = ... (variables, transaction, etc.).
type SetStmt struct {
	Base
	Items []*SetItem
}

// SetItem is one assignment in a SET statement.
type SetItem struct {
	Name       string
	IsVariable bool
	Scope      string // "global","session",""
	Value      Expression
	NoValue    bool // SET CHARACTER SET x style: no "= value" suffix
}

// UseDbStmt is USE db.
type UseDbStmt struct {
	Base
	Name string
}

// ExplainStmt is EXPLAIN.
type ExplainStmt struct {
	Base
	ExplainType string // "","EXTENDED","EXTENDED_NOADDR","PRETTY",...
	DisplayType string // "",
	Stmt        Statement
}

// BeginStmt is BEGIN / START TRANSACTION.
type BeginStmt struct {
	Base
	Work     bool
	ReadOnly bool
}

// CommitStmt is COMMIT.
type CommitStmt struct {
	Base
	Work bool
}

// RollbackStmt is ROLLBACK.
type RollbackStmt struct {
	Base
	Work bool
}

// SetTransactionStmt is SET [GLOBAL|SESSION] TRANSACTION <characteristics>.
type SetTransactionStmt struct {
	Base
	Scope           string // "GLOBAL", "SESSION", or ""
	Characteristics []string
}

// SavepointStmt is SAVEPOINT name.
type SavepointStmt struct {
	Base
	Name string
}

// ReleaseSavepointStmt is RELEASE SAVEPOINT name.
type ReleaseSavepointStmt struct {
	Base
	Name string
}

// RollbackToSavepointStmt is ROLLBACK [WORK] TO [SAVEPOINT] name.
type RollbackToSavepointStmt struct {
	Base
	Work bool
	Name string
}

// LockTable is one entry of LOCK TABLES: name [AS alias] READ|WRITE [LOCAL].
type LockTable struct {
	Table    *Identifier
	Alias    string
	HasAlias bool
	LockType string // "READ", "READ LOCAL", "WRITE", "LOW_PRIORITY WRITE"
}

// LockTablesStmt is LOCK TABLES t [AS a] READ, ... .
type LockTablesStmt struct {
	Base
	Tables []LockTable
}

// UnlockTablesStmt is UNLOCK TABLES.
type UnlockTablesStmt struct {
	Base
}
