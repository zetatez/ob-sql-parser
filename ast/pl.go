package ast

// PL block / stored routine AST. These nodes model the SQL/PSM-style
// compound-statement body of stored procedures and functions, plus CALL /
// DROP ROUTINE at the SQL statement level.

// ProcParam is one IN/OUT/INOUT parameter of a routine.
type ProcParam struct {
	Mode     string // "", "OUT", "INOUT"
	Name     string
	Type     string
	TypeArgs []Expression
}

// StoredProcStmt is CREATE [OR REPLACE] PROCEDURE|FUNCTION.
type StoredProcStmt struct {
	Base
	Function    bool
	OrReplace   bool
	Name        string
	Parameters  []ProcParam
	Returns     string // FUNCTION RETURNS type
	ReturnsArgs []Expression
	Body        *PLBlock
}

// DropRoutineStmt is DROP PROCEDURE|FUNCTION.
type DropRoutineStmt struct {
	Base
	Function bool
	IfExists bool
	Name     string
}

// CallStmt is CALL name(args).
type CallStmt struct {
	Base
	Name string
	Args []Expression
}

// DoStmt is DO expr (evaluate an expression, e.g. DO SLEEP(1)).
type DoStmt struct {
	Base
	Expr Expression
}

// PLBlock is a BEGIN ... END compound statement.
type PLBlock struct {
	Base
	Statements []Node
	Label      string // optional statement label: label: BEGIN ... END
}

// PLLeaveStmt is LEAVE label.
type PLLeaveStmt struct {
	Base
	Name string
}

// PLIterateStmt is ITERATE label.
type PLIterateStmt struct {
	Base
	Name string
}

// HandlerCondition is one condition of a DECLARE ... HANDLER.
type HandlerCondition struct {
	Kind  string // "sqlstate","name","notfound"
	Value string
}

// PLConditionHandlerStmt is DECLARE {CONTINUE|EXIT|UNDO} HANDLER FOR cond stmt.
type PLConditionHandlerStmt struct {
	Base
	Action     string // "CONTINUE","EXIT","UNDO"
	Conditions []HandlerCondition
	Body       Node
}

// PLConditionDeclStmt is DECLARE name CONDITION FOR SQLSTATE '...' | code.
type PLConditionDeclStmt struct {
	Base
	Name  string
	Kind  string // "sqlstate","code"
	Value string
}

// PLCaseWhen is one WHEN expr THEN stmts branch of a PL CASE statement.
type PLCaseWhen struct {
	When Expression
	Body *PLBlock
}

// PLCaseStmt is the statement-level CASE ... END CASE.
type PLCaseStmt struct {
	Base
	Operand  Expression
	Branches []PLCaseWhen
	Else     *PLBlock
}

// PLSetStmt is SET var = expr inside a routine.
type PLSetStmt struct {
	Base
	Name      string
	IsUserVar bool
	Expr      Expression
}

// PLDeclareStmt is DECLARE var[, var] type [DEFAULT expr].
type PLDeclareStmt struct {
	Base
	Names    []string
	Type     string
	TypeArgs []Expression
	Default  Expression
}

// PLIfBranch is one WHEN condition/THEN branch of a PL IF.
type PLIfBranch struct {
	Cond Expression
	Body *PLBlock
}

// PLIfStmt is IF ... THEN ... [ELSEIF ... THEN ...]* [ELSE ...] END IF.
type PLIfStmt struct {
	Base
	Branches []PLIfBranch
	Else     *PLBlock
}

// PLWhileStmt is [label:] WHILE cond DO ... END WHILE.
type PLWhileStmt struct {
	Base
	Label string
	Cond  Expression
	Body  *PLBlock
}

// PLRepeatStmt is [label:] REPEAT ... UNTIL cond END REPEAT.
type PLRepeatStmt struct {
	Base
	Label string
	Body  *PLBlock
	Until Expression
}

// PLLoopStmt is [label:] LOOP ... END LOOP.
type PLLoopStmt struct {
	Base
	Label string
	Body  *PLBlock
}

// PLReturnStmt is RETURN expr.
type PLReturnStmt struct {
	Base
	Expr Expression
}

// PLCursorDeclStmt is DECLARE c CURSOR FOR select.
type PLCursorDeclStmt struct {
	Base
	Name   string
	Select *SelectStmt
}

// PLCursorStmt is OPEN|CLOSE|FETCH ... INTO ... on a cursor.
type PLCursorStmt struct {
	Base
	Op   string // "OPEN","CLOSE","FETCH"
	Name string
	Into []string
}

// Children() for the PL nodes (see pl.go).

func (n *StoredProcStmt) Children() []Node {
	var out []Node
	out = append(out, n.Body)
	for i := range n.Parameters {
		out = append(out, enode(n.Parameters[i].TypeArgs)...)
	}
	out = append(out, enode(n.ReturnsArgs)...)
	return out
}

func (n *DropRoutineStmt) Children() []Node { return nil }
func (n *CallStmt) Children() []Node        { return enode(n.Args) }
func (n *DoStmt) Children() []Node          { return childrenOf(n.Expr) }

func (n *PLBlock) Children() []Node {
	return n.Statements
}

func (n *PLSetStmt) Children() []Node { return childrenOf(n.Expr) }

func (n *PLDeclareStmt) Children() []Node {
	return childrenOf(n.Default)
}

func (n *PLIfStmt) Children() []Node {
	var out []Node
	for i := range n.Branches {
		out = append(out, n.Branches[i].Cond, n.Branches[i].Body)
	}
	out = append(out, n.Else)
	return out
}

func (n *PLWhileStmt) Children() []Node      { return childrenOf(n.Cond, n.Body) }
func (n *PLRepeatStmt) Children() []Node     { return childrenOf(n.Body, n.Until) }
func (n *PLLoopStmt) Children() []Node       { return childrenOf(n.Body) }
func (n *PLReturnStmt) Children() []Node     { return childrenOf(n.Expr) }
func (n *PLCursorDeclStmt) Children() []Node { return childrenOf(n.Select) }
func (n *PLCursorStmt) Children() []Node     { return nil }
func (n *PLLeaveStmt) Children() []Node      { return nil }
func (n *PLIterateStmt) Children() []Node    { return nil }

func (n *PLConditionHandlerStmt) Children() []Node { return childrenOf(n.Body) }
func (n *PLConditionDeclStmt) Children() []Node    { return nil }

func (n *PLCaseStmt) Children() []Node {
	var out []Node
	out = append(out, n.Operand)
	for i := range n.Branches {
		out = append(out, n.Branches[i].When, n.Branches[i].Body)
	}
	out = append(out, n.Else)
	return out
}
