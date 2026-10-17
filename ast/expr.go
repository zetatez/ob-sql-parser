package ast

// LiteralKind describes the literal category.
type LiteralKind string

const (
	LitString  LiteralKind = "string"
	LitInt     LiteralKind = "int"
	LitDecimal LiteralKind = "decimal"
	LitFloat   LiteralKind = "float"
	LitHex     LiteralKind = "hex"
	LitBit     LiteralKind = "bit"
	LitNull    LiteralKind = "null"
	LitBool    LiteralKind = "bool"
)

// Literal is a constant literal.
type Literal struct {
	Base
	LitType LiteralKind
	Value   string // unescaped value (for hex/bit keeps textual digits)
	Text    string // exact source text including quote/introducer
	Charset string // charset introducer, if any (e.g. "utf8mb4")
}

// Identifier is a (possibly qualified) name such as a, a.b, a.b.c.
type Identifier struct {
	Base
	Parts []string
}

// StarLiteral is `*` or `tbl.*`.
type StarLiteral struct {
	Base
	Qualifier string // empty, or a table alias prefix like "t"
}

// ParamMarker is a prepared-statement placeholder `?`.
type ParamMarker struct {
	Base
	Index int // 1-based position, assigned during parse
}

// Variable is `@x`, `@@x`, `@@session.x`, `@@global.x`.
type Variable struct {
	Base
	Name     string
	IsSystem bool
	Scope    string // "global", "session", or ""
}

// UnaryExpr applies a unary operator.
type UnaryExpr struct {
	Base
	Op      string
	Operand Expression
}

// BinaryExpr applies a binary operator.
type BinaryExpr struct {
	Base
	Op    string
	Left  Expression
	Right Expression
}

// FuncCall is a function invocation (ordinary, aggregate, or windowed).
type FuncCall struct {
	Base
	Name         string
	Distinct     bool
	Args         []Expression
	Star         bool        // count(*) style
	Over         *WindowSpec // window function if non-nil
	DistinctFlag bool
}

// CastExpr is CAST(x AS T) or CONVERT(x, T).
type CastExpr struct {
	Base
	Expr      Expression
	TypeName  string
	TypeArgs  []Expression
	IsConvert bool
}

// CaseClause is one WHEN ... THEN ... branch.
type CaseClause struct {
	When Expression
	Then Expression
}

// CaseExpr is a CASE expression.
type CaseExpr struct {
	Base
	Operand Expression
	Whens   []CaseClause
	Else    Expression
}

// BetweenExpr is expr [NOT] BETWEEN low AND high.
type BetweenExpr struct {
	Base
	Expr Expression
	Low  Expression
	High Expression
	Not  bool
}

// InExpr is expr [NOT] IN (list | subquery).
type InExpr struct {
	Base
	Expr     Expression
	List     []Expression
	Subquery *SelectStmt
	Not      bool
}

// LikeExpr is expr [NOT] LIKE pattern [ESCAPE e].
type LikeExpr struct {
	Base
	Expr    Expression
	Pattern Expression
	Escape  Expression
	Not     bool
}

// IsExpr is expr IS [NOT] (NULL | TRUE | FALSE).
type IsExpr struct {
	Base
	Expr  Expression
	Not   bool
	Truth string // "null", "true", "false"
}

// IntervalExpr is INTERVAL n UNIT.
type IntervalExpr struct {
	Base
	Value Expression
	Unit  string
}

// MemberOfExpr is expr MEMBER OF (json_array).
type MemberOfExpr struct {
	Base
	Expr Expression
	List []Expression
}

// MatchAgainstExpr is MATCH (cols) AGAINST (expr [modifier]).
type MatchAgainstExpr struct {
	Base
	Columns  []*Identifier
	Against  Expression
	Modifier string // "", "IN BOOLEAN MODE", "WITH QUERY EXPANSION"
}

// GroupingSetExpr is GROUPING SETS ((...),(...)) | CUBE(...) | ROLLUP(...).
type GroupingSetExpr struct {
	Base
	Kind  string // "GROUPING SETS", "CUBE", "ROLLUP"
	Sets  [][]Expression
	Exprs []Expression
}

// ExistsExpr is [NOT] EXISTS (subquery).
type ExistsExpr struct {
	Base
	Query *SelectStmt
	Not   bool
}

// SubqueryExpr is a scalar subquery used as an expression.
type SubqueryExpr struct {
	Base
	Query *SelectStmt
}

// ParenExpr wraps a parenthesized expression (kept to preserve grouping).
type ParenExpr struct {
	Base
	Expr Expression
}
