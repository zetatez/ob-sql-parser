package ast

// OrderItem is an ORDER BY item.
type OrderItem struct {
	Expr       Expression
	Desc       bool
	HasNulls   bool
	NullsFirst bool
}

// WindowSpec describes an OVER (...) clause.
type WindowSpec struct {
	Name        string
	PartitionBy []Expression
	OrderBy     []OrderItem
	Frame       string // "", "ROWS", "RANGE"
	FrameStart  string
	FrameEnd    string
}

// TableRefKind discriminates table references.
type TableRefKind string

const (
	RefTable    TableRefKind = "table"
	RefSubquery TableRefKind = "subquery"
	RefJoin     TableRefKind = "join"
)

// TableRef is a FROM/UPDATE/DELETE source: a table, a derived subquery, or a join.
type TableRef struct {
	Base
	RefKind  TableRefKind
	Name     *Identifier // for RefTable
	Alias    string
	HasAlias bool

	// RefTable: MySQL index hints (USE/FORCE/IGNORE INDEX ...) and an optional
	// partition selection PARTITION (p0, p1).
	IndexHints []IndexHint
	Partitions []string

	// RefSubquery
	Subquery *SelectStmt
	Lateral  bool // LATERAL (subquery)

	// RefJoin
	JoinType string // "JOIN","INNER JOIN","LEFT [OUTER] JOIN","RIGHT","CROSS","STRAIGHT_JOIN","NATURAL ..."
	Left     *TableRef
	Right    *TableRef
	On       Expression
	Using    []*Identifier
}

// IndexHintType is the kind of an index hint.
type IndexHintType string

const (
	IndexHintUse    IndexHintType = "USE"
	IndexHintForce  IndexHintType = "FORCE"
	IndexHintIgnore IndexHintType = "IGNORE"
)

// IndexHint is a MySQL index hint: USE|FORCE|IGNORE INDEX [FOR ...] (...).
type IndexHint struct {
	Type    IndexHintType
	For     string   // "", "JOIN", "ORDER BY", "GROUP BY"
	Indexes []string // index names, or "PRIMARY"
}

// ColumnItem is a SELECT projection item.
type ColumnItem struct {
	Expr     Expression
	Alias    string
	HasAlias bool
}

// LimitClause is LIMIT [offset,] count or LIMIT count OFFSET offset.
type LimitClause struct {
	Count     Expression
	Offset    Expression
	HasOffset bool
}

// SelectBody is a single SELECT (before any set operator).
type SelectBody struct {
	Base
	Distinct    bool
	Items       []ColumnItem
	From        []*TableRef
	Where       Expression
	GroupBy     []Expression
	Rollup      bool // GROUP BY ... WITH ROLLUP
	Having      Expression
	Qualify     Expression // QUALIFY expr (window filter)
	OrderBy     []OrderItem
	Limit       *LimitClause
	ForUpdate   bool
	LockWait    string // "","NOWAIT","SKIP LOCKED"
	InShareMode bool   // LOCK IN SHARE MODE
	NamedWindow []*NamedWindow

	// PROCEDURE ANALYSE(...) (legacy MySQL, kept for completeness).
	HasProcAnalyse bool
	ProcAnalyse    []Expression
}

// NamedWindow is WINDOW name AS (spec) clause.
type NamedWindow struct {
	Name string
	Spec *WindowSpec
}

// SetOp is a UNION/EXCEPT/INTERSECT joining two SelectBodies.
type SetOp struct {
	Op    string // "UNION","EXCEPT","INTERSECT"
	All   bool
	Right *SelectBody
}

// SelectStmt is a complete SELECT, possibly composed of several bodies via set
// operators, with optional trailing ORDER BY / LIMIT / INTO.
type SelectStmt struct {
	Base
	Bodies  []*SelectBody
	SetOps  []*SetOp
	OrderBy []OrderItem
	Limit   *LimitClause

	// INTO OUTFILE/DUMPFILE.
	IntoKind   string // "","OUTFILE","DUMPFILE","@vars"
	IntoFile   string
	IntoVars   []*Variable // SELECT ... INTO @a, @b
	IntoFields string      // FIELDS TERMINATED BY
	IntoLines  string
}
