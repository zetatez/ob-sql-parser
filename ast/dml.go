package ast

// Assignment is a SET / ON DUPLICATE KEY UPDATE column = value.
type Assignment struct {
	Column *Identifier
	Expr   Expression
}

// InsertStmt is INSERT/REPLACE.
type InsertStmt struct {
	Base
	Replace     bool
	Ignore      bool
	Table       *Identifier
	Columns     []*Identifier
	Rows        [][]Expression
	Source      *SelectStmt
	Set         []Assignment
	OnDuplicate []Assignment
	Returning   []*Identifier
}

// UpdateStmt is UPDATE.
type UpdateStmt struct {
	Base
	Ignore      bool
	LowPriority bool
	Tables      []*TableRef
	Set         []Assignment
	Where       Expression
	OrderBy     []OrderItem
	Limit       *LimitClause
	Returning   []*Identifier
}

// DeleteStmt is DELETE.
type DeleteStmt struct {
	Base
	Mods      string        // "","LOW_PRIORITY","QUICK","IGNORE"
	Tables    []*Identifier // 多表删除的目标表列表（DELETE t1, t2 FROM ...）
	From      []*TableRef
	Using     []*TableRef
	Where     Expression
	OrderBy   []OrderItem
	Limit     *LimitClause
	Returning []*Identifier
}

// MergeWhen is one WHEN [NOT] MATCHED branch of MERGE.
type MergeWhen struct {
	Matched    bool
	Action     string // "UPDATE", "DELETE", "INSERT"
	Set        []Assignment
	DeleteCond Expression
	Columns    []*Identifier
	Values     []Expression
}

// MergeStmt is MERGE INTO target USING source ON cond WHEN ... .
type MergeStmt struct {
	Base
	Target *TableRef
	Source *TableRef
	On     Expression
	Whens  []MergeWhen
}
