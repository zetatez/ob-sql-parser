package ast

// ColumnDef is a CREATE TABLE column definition.
type ColumnDef struct {
	Name            string
	Type            string
	TypeParams      []Expression
	Unsigned        bool
	Zerofill        bool
	Nullable        Nullable
	Default         Expression
	OnUpdate        Expression // ON UPDATE CURRENT_TIMESTAMP[(n)] / NOW() / ...
	Comment         string
	AutoIncrement   bool
	Charset         string
	Collate         string
	Invisible       bool
	Generated       Expression // GENERATED ALWAYS AS (expr)
	GeneratedStored bool       // true = STORED, false = VIRTUAL

	// 列级约束
	PrimaryKey bool
	UniqueKey  bool
	RefTable   *Identifier // 列级 REFERENCES
	RefColumns []*Identifier
	OnDelete   string
	FkOnUpdate string
}

// Nullable is a tri-state column nullability.
type Nullable int

const (
	NullableDefault Nullable = iota
	NullableYes
	NullableNo
)

func (n Nullable) String() string {
	switch n {
	case NullableYes:
		return "NULL"
	case NullableNo:
		return "NOT NULL"
	}
	return ""
}

// ConstraintDef is a table-level constraint.
type ConstraintDef struct {
	Name       string
	Kind       string        // "PRIMARY KEY","UNIQUE","KEY","INDEX","FOREIGN KEY","CHECK"
	Columns    []*Identifier // index/pr key columns (as names)
	IndexCols  []IndexCol    // KEY/INDEX/UNIQUE 带顺序/长度的列
	Check      Expression    // CHECK (expr)
	IndexType  string
	RefTable   *Identifier   // FOREIGN KEY ... REFERENCES table
	RefColumns []*Identifier // REFERENCES table (cols)
	OnDelete   string        // ON DELETE RESTRICT|CASCADE|SET NULL|NO ACTION|SET DEFAULT
	OnUpdate   string        // ON UPDATE ...
}

// TableOption is a generic CREATE TABLE option (ENGINE=.., COMMENT=.., ...).
type TableOption struct {
	Name  string
	Value string
}

// CreateTableStmt is CREATE [TEMPORARY] TABLE.
type CreateTableStmt struct {
	Base
	Temp        bool
	IfNotExists bool
	Table       *Identifier
	Columns     []*ColumnDef
	Constraints []*ConstraintDef
	Options     []TableOption
	Partition   *PartitionSpec
	Select      *SelectStmt // CREATE TABLE ... AS SELECT
}

// PartitionByKind is the partitioning method.
type PartitionByKind string

const (
	PartitionRange        PartitionByKind = "RANGE"
	PartitionRangeColumns PartitionByKind = "RANGE COLUMNS"
	PartitionList         PartitionByKind = "LIST"
	PartitionListColumns  PartitionByKind = "LIST COLUMNS"
	PartitionHash         PartitionByKind = "HASH"
	PartitionKey          PartitionByKind = "KEY"
)

// PartitionSpec is a PARTITION BY clause on CREATE TABLE.
type PartitionSpec struct {
	Kind       PartitionByKind
	Expr       Expression    // RANGE/LIST/HASH key expression
	Columns    []*Identifier // RANGE/LIST COLUMNS, KEY columns
	Count      Expression    // PARTITIONS n
	Auto       bool          // PARTITIONS AUTO
	Size       Expression    // PARTITION SIZE n (auto partition)
	Partitions []PartitionDef

	// SUBPARTITION BY ...
	SubKind    PartitionByKind
	SubExpr    Expression
	SubColumns []*Identifier
	SubCount   Expression
}

// PartitionDef is one PARTITION p... definition.
type PartitionDef struct {
	Name      string
	ValueKind string // "LESS THAN" or "IN"
	Values    []Expression
	MaxValue  bool // VALUES LESS THAN MAXVALUE (无括号)
	Options   []TableOption
}

// CreateDatabaseStmt is CREATE DATABASE/SCHEMA.
type CreateDatabaseStmt struct {
	Base
	IfNotExists bool
	Name        string
	Charset     string
	Collate     string
}

// DropDatabaseStmt is DROP DATABASE/SCHEMA.
type DropDatabaseStmt struct {
	Base
	IfExists bool
	Name     string
}

// AlterDatabaseStmt is ALTER DATABASE/SCHEMA.
type AlterDatabaseStmt struct {
	Base
	Name    string
	Charset string
	Collate string
}

// TransferPartitionStmt is ALTER SYSTEM [CANCEL] TRANSFER PARTITION.
type TransferPartitionStmt struct {
	Base
	Cancel   bool
	All      bool
	TableID  Expression
	ObjectID Expression
	LS       Expression
	Tenant   string
}

// AlterSystemStmt is ALTER SYSTEM SET name = value.
type AlterSystemStmt struct {
	Base
	Name  string
	Value Expression
}

// CreateSequenceStmt is CREATE SEQUENCE.
type CreateSequenceStmt struct {
	Base
	IfNotExists bool
	Name        string
	StartWith   Expression
	IncrementBy Expression
	MinValue    Expression
	MaxValue    Expression
	Cache       Expression
	Cycle       string // "", "CYCLE", "NOCYCLE"
	Order       string // "", "ORDER", "NOORDER"
}

// DropSequenceStmt is DROP SEQUENCE.
type DropSequenceStmt struct {
	Base
	IfExists bool
	Name     string
}

// CreateViewStmt is CREATE [OR REPLACE] VIEW.
type CreateViewStmt struct {
	Base
	OrReplace bool
	Name      *Identifier
	Columns   []*Identifier
	Select    *SelectStmt
}

// DropViewStmt is DROP VIEW.
type DropViewStmt struct {
	Base
	IfExists bool
	Names    []*Identifier
}

// CreateTablegroupStmt is CREATE TABLEGROUP.
type CreateTablegroupStmt struct {
	Base
	IfNotExists bool
	Name        string
	Options     []Option
}

// DropTablegroupStmt is DROP TABLEGROUP.
type DropTablegroupStmt struct {
	Base
	IfExists bool
	Name     string
}

// CreateTablespaceStmt is CREATE TABLESPACE.
type CreateTablespaceStmt struct {
	Base
	Name    string
	Options []Option
}

// DropTablespaceStmt is DROP TABLESPACE.
type DropTablespaceStmt struct {
	Base
	IfExists bool
	Name     string
}

// AlterMViewStmt is ALTER MATERIALIZED VIEW.
type AlterMViewStmt struct {
	Base
	Name    *Identifier
	Refresh string // "FAST","COMPLETE","FORCE",""
	Mode    string // "ON DEMAND","ON COMMIT",""
}

// MLogStmt is CREATE|ALTER|DROP MATERIALIZED VIEW LOG ON table.
type MLogStmt struct {
	Base
	Action       string // "CREATE","ALTER","DROP"
	IfExists     bool
	Table        *Identifier
	WithPrimary  bool
	IncludingNew bool
}

// CreateOutlineStmt is CREATE OUTLINE.
type CreateOutlineStmt struct {
	Base
	Name   *Identifier
	Text   string // ON 'sql text'
	Select *SelectStmt
}

// AlterOutlineStmt is ALTER OUTLINE.
type AlterOutlineStmt struct {
	Base
	Name   *Identifier
	Text   string
	Select *SelectStmt
}

// DropOutlineStmt is DROP OUTLINE.
type DropOutlineStmt struct {
	Base
	IfExists bool
	Name     *Identifier
}

// AlterTablespaceStmt is ALTER TABLESPACE.
type AlterTablespaceStmt struct {
	Base
	Name    string
	Options []Option
}

// AlterTablegroupStmt is ALTER TABLEGROUP.
type AlterTablegroupStmt struct {
	Base
	Name     string
	Options  []Option
	AddTable *Identifier // ALTER TABLEGROUP ... ADD TABLE t
}

// CreateMViewStmt is CREATE MATERIALIZED VIEW.
type CreateMViewStmt struct {
	Base
	Name    *Identifier
	Columns []*Identifier
	Refresh string // "FAST","COMPLETE","FORCE",""
	Mode    string // "ON DEMAND","ON COMMIT",""
	Select  *SelectStmt
}

// DropMViewStmt is DROP MATERIALIZED VIEW.
type DropMViewStmt struct {
	Base
	IfExists bool
	Names    []*Identifier
}

// FlashbackTableStmt is FLASHBACK TABLE ... TO BEFORE DROP | TO SNAPSHOT ... .
type FlashbackTableStmt struct {
	Base
	Tables       []*Identifier
	ToBeforeDrop bool
	Snapshot     string
	RenameTo     *Identifier
}

// RenamePair is one "a TO b" mapping in RENAME TABLE.
type RenamePair struct {
	From *Identifier
	To   *Identifier
}

// RenameTableStmt is RENAME TABLE a TO b[, c TO d, …].
type RenameTableStmt struct {
	Base
	Pairs []RenamePair
}

// CreateTableLikeStmt is CREATE TABLE t LIKE src.
type CreateTableLikeStmt struct {
	Base
	Table *Identifier
	Like  *Identifier
}

// CreateIndexStmt is CREATE [UNIQUE] INDEX.
// CreateIndexStmt is CREATE [UNIQUE] INDEX.
type CreateIndexStmt struct {
	Base
	Unique      bool
	IfNotExists bool
	Name        *Identifier
	Table       *Identifier
	Columns     []IndexCol
	Method      string   // USING BTREE/HASH
	Options     []string // GLOBAL/LOCAL/KEY_BLOCK_SIZE/COMMENT/...
	Partition   *PartitionSpec
}

// IndexCol 是索引列：可带前缀长度与排序方向。
type IndexCol struct {
	Name     *Identifier
	Length   Expression // a(10)
	Desc     bool
	HasOrder bool // 显式 ASC/DESC
}

// DropTableStmt is DROP TABLE.
type DropTableStmt struct {
	Base
	IfExists bool
	Tables   []*Identifier
}

// DropIndexStmt is DROP INDEX.
type DropIndexStmt struct {
	Base
	Name  *Identifier
	Table *Identifier
}

// AlterActionKind enumerates ALTER TABLE sub-actions.
type AlterActionKind string

const (
	AlterAddColumn      AlterActionKind = "ADD COLUMN"
	AlterDropColumn     AlterActionKind = "DROP COLUMN"
	AlterModifyColumn   AlterActionKind = "MODIFY COLUMN"
	AlterChangeColumn   AlterActionKind = "CHANGE COLUMN"
	AlterRenameColumn   AlterActionKind = "RENAME COLUMN"
	AlterRenameTable    AlterActionKind = "RENAME TABLE"
	AlterAddIndex       AlterActionKind = "ADD INDEX"
	AlterDropIndex      AlterActionKind = "DROP INDEX"
	AlterRenameIndex    AlterActionKind = "RENAME INDEX"
	AlterDropPrimaryKey AlterActionKind = "DROP PRIMARY KEY"
	AlterAddConstraint  AlterActionKind = "ADD CONSTRAINT"
	AlterSetDefault     AlterActionKind = "ALTER COLUMN SET DEFAULT"
	AlterDropDefault    AlterActionKind = "ALTER COLUMN DROP DEFAULT"

	AlterAddPartition        AlterActionKind = "ADD PARTITION"
	AlterDropPartition       AlterActionKind = "DROP PARTITION"
	AlterTruncatePartition   AlterActionKind = "TRUNCATE PARTITION"
	AlterCoalescePartition   AlterActionKind = "COALESCE PARTITION"
	AlterExchangePartition   AlterActionKind = "EXCHANGE PARTITION"
	AlterReorganizePartition AlterActionKind = "REORGANIZE PARTITION"
	AlterRemovePartitioning  AlterActionKind = "REMOVE PARTITIONING"

	AlterAddColumns     AlterActionKind = "ADD COLUMNS"
	AlterConvertCharset AlterActionKind = "CONVERT TO CHARACTER SET"
	AlterAutoIncrement  AlterActionKind = "AUTO_INCREMENT"
	AlterVisibleIndex   AlterActionKind = "ALTER INDEX VISIBLE"
	AlterInvisibleIndex AlterActionKind = "ALTER INDEX INVISIBLE"
	AlterEnableKeys     AlterActionKind = "ENABLE KEYS"
	AlterDisableKeys    AlterActionKind = "DISABLE KEYS"
)

// AlterAction is one comma-separated action in an ALTER TABLE statement.
type AlterAction struct {
	Kind       AlterActionKind
	Column     *ColumnDef  // ADD/MODIFY/CHANGE COLUMN definition
	OldName    string      // CHANGE COLUMN old name
	Name       string      // index/constraint/column name; RENAME INDEX old name
	NewName    string      // RENAME INDEX new name
	Target     *Identifier // RENAME TO table
	Columns    []*Identifier
	Unique     bool
	Constraint *ConstraintDef
	Default    Expression

	// column position / convert charset
	Position     string // "","FIRST","AFTER"
	RefCol       string
	ColumnDefs   []*ColumnDef
	Charset      string
	Collate      string
	IndexOptions []string // ADD INDEX ... LOCAL/GLOBAL/COMMENT '..'/KEY_BLOCK_SIZE=n
	Visible      bool     // ALTER INDEX ... VISIBLE/INVISIBLE
	IndexCols    []IndexCol

	// partition operations
	PartitionDefs []PartitionDef
	Names         []string
	All           bool
	Count         Expression
	With          *Identifier
	Validation    string // EXCHANGE PARTITION ... [WITH|WITHOUT] VALIDATION
}

// AlterTableStmt is ALTER TABLE with one or more actions.
type AlterTableStmt struct {
	Base
	Table   *Identifier
	Actions []AlterAction
}

// TruncateStmt is TRUNCATE TABLE.
type TruncateStmt struct {
	Base
	Table *Identifier
}
