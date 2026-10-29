// Semantic helpers for the goyacc grammar: AST builders (asX), type
// conversions, and node constructors invoked by sql.y actions.

package gparser

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/token"
)

// dtype is the parsed form of a SQL data type (name + optional type arguments).
type dtype struct {
	name     string
	args     []ast.Expression
	unsigned bool
}

// colAttrs accumulates the attributes of a CREATE TABLE column definition.
type colAttrs struct {
	nullable      ast.Nullable
	def           ast.Expression
	onUpdate      ast.Expression
	comment       string
	autoIncrement bool
	charset       string
	collate       string
	unsigned      bool
	invisible     bool
	generated     ast.Expression
	genStored     bool
	hasGenerated  bool
}

// nullsOpt is the parsed NULLS FIRST/LAST ordering option.
type nullsOpt struct {
	has   bool
	first bool
}

// setOp is a parsed set operator.
type setOp struct {
	op  string
	all bool
}

// inList distinguishes an IN list from an IN subquery.
type inList struct {
	list []ast.Expression
	sub  *ast.SelectStmt
}

// callArgs is the parsed argument list of a function call.
type callArgs struct {
	distinct bool
	star     bool
	args     []ast.Expression
}

// insertSource is the source of an INSERT/REPLACE.
type insertSource struct {
	rows [][]ast.Expression
	set  []ast.Assignment
	sel  *ast.SelectStmt
}

// joinCond is an ON/USING condition on a join.
type joinCond struct {
	on    ast.Expression
	using []*ast.Identifier
}

// createDefs accumulates CREATE TABLE column definitions and table constraints.
type createDefs struct {
	cols []*ast.ColumnDef
	cons []*ast.ConstraintDef
}

// tableBody is the parsed tail of a CREATE TABLE (either a column list with
// options, or a LIKE source).
type tableBody struct {
	defs   createDefs
	opts   []ast.TableOption
	part   *ast.PartitionSpec
	sel    *ast.SelectStmt
	isLike bool
	like   *ast.Identifier
}

func asNodes(v interface{}) []ast.Node {
	if s, ok := v.([]ast.Node); ok {
		return s
	}
	if v == nil {
		return nil
	}
	return []ast.Node{asNode(v)}
}
func asNode(v interface{}) ast.Node {
	if n, ok := v.(ast.Node); ok {
		return n
	}
	return nil
}
func asProcParams(v interface{}) []ast.ProcParam {
	if s, ok := v.([]ast.ProcParam); ok {
		return s
	}
	return nil
}
func asPLCaseWhens(v interface{}) []ast.PLCaseWhen {
	if s, ok := v.([]ast.PLCaseWhen); ok {
		return s
	}
	return nil
}
func asPLBranches(v interface{}) []ast.PLIfBranch {
	if s, ok := v.([]ast.PLIfBranch); ok {
		return s
	}
	return nil
}
func asPLBlock(v interface{}) *ast.PLBlock {
	if b, ok := v.(*ast.PLBlock); ok {
		return b
	}
	return nil
}
func asHandlerConditions(v interface{}) []ast.HandlerCondition {
	if s, ok := v.([]ast.HandlerCondition); ok {
		return s
	}
	return nil
}

// reduceSpan computes the committed-token index at which the just-reduced rule
// began, and stamps the produced node (if any) with its source span. It is
// called from the generated parser's reduce epilogue (see genpos.go) which
// maintains a parallel start-index stack for every LALR stack frame.
func reduceSpan(yylex yyLexer, n interface{}, yyPos []int, yypt, yyR, yyTok int) int {
	start := yyTok
	if yyR > 0 && yypt-yyR+1 >= 0 && yypt-yyR+1 < len(yyPos) {
		start = yyPos[yypt-yyR+1]
	}
	if s, ok := n.(ast.Node); ok {

		if !s.Span().Start.IsValid() {
			if sp, ok2 := s.(spannable); ok2 {
				if l, ok3 := yylex.(*obLexer); ok3 {
					if start >= 0 && start < len(l.toks) && yyTok >= 0 && yyTok < len(l.toks) {
						sp.SetSpan(l.toks[start].Pos, l.toks[yyTok].EndPos)
					}
				}
			}
		}
	}
	return start
}

// lockOpt is a parsed FOR UPDATE [NOWAIT|SKIP LOCKED] / LOCK IN SHARE MODE clause.
type lockOpt struct {
	update bool
	wait   string
	share  bool
}

func asLockOpt(v interface{}) lockOpt {
	if o, ok := v.(lockOpt); ok {
		return o
	}
	return lockOpt{}
}

// partInfo is a parsed TABLE_ID = n, OBJECT_ID = m pair.
type partInfo struct {
	table ast.Expression
	obj   ast.Expression
}

// identAuth is a parsed IDENTIFIED BY/WITH clause for CREATE USER.
type identAuth struct {
	pwd    string
	plugin string
}

// colPos is a parsed FIRST | AFTER col column position.
type colPos struct {
	position string
	ref      string
}

// nameOrAll is a partition name list or ALL.
type nameOrAll struct {
	names []string
	all   bool
}

func asColumnDefs(v interface{}) []*ast.ColumnDef {
	if s, ok := v.([]*ast.ColumnDef); ok {
		return s
	}
	return nil
}

// dblinkTail is the optional HOST/PORT/CLUSTER tail of CREATE DATABASE LINK.
type dblinkTail struct {
	host    string
	port    ast.Expression
	cluster string
}

func idNamesList(ids []*ast.Identifier) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, strings.Join(id.Parts, "."))
	}
	return out
}

// recUntil is a parsed RECOVER ... UNTIL/CANCEL clause.
type recUntil struct {
	until  string
	value  ast.Expression
	cancel bool
}

func asOptPtr(v interface{}) *ast.Option {
	if o, ok := v.(*ast.Option); ok {
		return o
	}
	return nil
}

// optionValue is a parsed NAME = value | NAME = (values) option value.
type optionValue struct {
	value  ast.Expression
	values []ast.Expression
}

func asOptions(v interface{}) []ast.Option {
	if s, ok := v.([]ast.Option); ok {
		return s
	}
	return nil
}

// intoFile is a parsed SELECT ... INTO OUTFILE/DUMPFILE clause.
type intoFile struct {
	kind   string
	file   string
	vars   []*ast.Variable
	fields string
	lines  string
}
type intoFields struct {
	fields string
	lines  string
}

// procAnalyse is a parsed PROCEDURE ANALYSE(...) clause.
type procAnalyse struct {
	has  bool
	args []ast.Expression
}

// autoPart is a parsed PARTITIONS AUTO / PARTITION SIZE n clause.
type autoPart struct {
	auto bool
	size ast.Expression
}

func asAutoPart(v interface{}) autoPart {
	if a, ok := v.(autoPart); ok {
		return a
	}
	return autoPart{}
}

// subSpec is a parsed SUBPARTITION BY clause.
type subSpec struct {
	kind  ast.PartitionByKind
	expr  ast.Expression
	cols  []*ast.Identifier
	count ast.Expression
}

func asSubSpec(v interface{}) subSpec {
	if s, ok := v.(subSpec); ok {
		return s
	}
	return subSpec{}
}

// partitionValues is a parsed PARTITION VALUES clause.
type partitionValues struct {
	kind   string
	values []ast.Expression
}

// dbOpts accumulates CREATE DATABASE options.
type dbOpts struct {
	charset string
	collate string
}

// frameSpec is a parsed window frame (ROWS/RANGE ...).
type frameSpec struct {
	frame string
	start string
	end   string
}

// groupOpt is a parsed GROUP BY clause result.
type groupOpt struct {
	exprs  []ast.Expression
	rollup bool
}

// mlogOpts accumulates materialized-view-log options.
type mlogOpts struct {
	withPrimary  bool
	includingNew bool
}

func asMLogOpts(v interface{}) mlogOpts {
	if o, ok := v.(mlogOpts); ok {
		return o
	}
	return mlogOpts{}
}

// mviewOpts accumulates materialized-view options.
type mviewOpts struct {
	refresh string
	mode    string
}

func asMViewOpts(v interface{}) mviewOpts {
	if o, ok := v.(mviewOpts); ok {
		return o
	}
	return mviewOpts{}
}

// seqOpts accumulates CREATE SEQUENCE options.
type seqOpts struct {
	start ast.Expression
	incr  ast.Expression
	min   ast.Expression
	max   ast.Expression
	cache ast.Expression
	cycle string
	order string
}

func asSeqOpts(v interface{}) seqOpts {
	if s, ok := v.(seqOpts); ok {
		return s
	}
	return seqOpts{}
}
func mergeSeqOpts(a, b interface{}) seqOpts {
	var out seqOpts
	if x, ok := a.(seqOpts); ok {
		out = x
	}
	if y, ok := b.(seqOpts); ok {
		if y.start != nil {
			out.start = y.start
		}
		if y.incr != nil {
			out.incr = y.incr
		}
		if y.min != nil {
			out.min = y.min
		}
		if y.max != nil {
			out.max = y.max
		}
		if y.cache != nil {
			out.cache = y.cache
		}
		if y.cycle != "" {
			out.cycle = y.cycle
		}
		if y.order != "" {
			out.order = y.order
		}
	}
	return out
}
func newIdent(name string) *ast.Identifier    { return &ast.Identifier{Parts: []string{name}} }
func idParts(parts ...string) *ast.Identifier { return &ast.Identifier{Parts: parts} }
func asExpr(v interface{}) ast.Expression {
	if e, ok := v.(ast.Expression); ok {
		return e
	}
	return nil
}
func asIdent(v interface{}) *ast.Identifier {
	if i, ok := v.(*ast.Identifier); ok {
		return i
	}
	return nil
}
func asIdentName(v interface{}) string {
	if i := asIdent(v); i != nil && len(i.Parts) > 0 {
		return i.Parts[len(i.Parts)-1]
	}
	return ""
}
func asIdentList(v interface{}) []*ast.Identifier {
	if s, ok := v.([]*ast.Identifier); ok {
		return s
	}
	return nil
}
func asExprList(v interface{}) []ast.Expression {
	if s, ok := v.([]ast.Expression); ok {
		return s
	}
	return nil
}
func asItems(v interface{}) []ast.ColumnItem {
	if s, ok := v.([]ast.ColumnItem); ok {
		return s
	}
	return nil
}
func asItem(v interface{}) ast.ColumnItem {
	if it, ok := v.(ast.ColumnItem); ok {
		return it
	}
	return ast.ColumnItem{}
}
func asBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}
func asStr(v interface{}) string {
	s, _ := v.(string)
	return s
}
func asStmt(v interface{}) ast.Statement {
	if s, ok := v.(ast.Statement); ok {
		return s
	}
	return nil
}
func asTableRef(v interface{}) *ast.TableRef {
	if r, ok := v.(*ast.TableRef); ok {
		return r
	}
	return nil
}
func asSelect(v interface{}) *ast.SelectStmt {
	if s, ok := v.(*ast.SelectStmt); ok {
		return s
	}
	return nil
}
func asBody(v interface{}) *ast.SelectBody {
	if s, ok := v.(*ast.SelectBody); ok {
		return s
	}
	return nil
}
func asDType(v interface{}) dtype {
	if d, ok := v.(dtype); ok {
		return d
	}
	return dtype{}
}
func asLitValue(v interface{}) string {
	if l, ok := v.(*ast.Literal); ok {
		return l.Value
	}
	return ""
}

// normalizeFuncStar rewrites f(*) parsed as a single bare-star argument into
// the dedicated Star form, so both spellings produce the same AST.
func normalizeFuncStar(fc *ast.FuncCall) {
	if fc.Star || len(fc.Args) != 1 {
		return
	}
	if s, ok := fc.Args[0].(*ast.StarLiteral); ok && s.Qualifier == "" {
		fc.Star = true
		fc.Args = nil
	}
}
func bin(op string, l, r ast.Expression) *ast.BinaryExpr {
	return &ast.BinaryExpr{Op: op, Left: l, Right: r}
}
func un(op string, x ast.Expression) *ast.UnaryExpr {
	return &ast.UnaryExpr{Op: op, Operand: x}
}
func mkIn(e ast.Expression, in interface{}, not bool) *ast.InExpr {
	il, _ := in.(inList)
	return &ast.InExpr{Expr: e, List: il.list, Subquery: il.sub, Not: not}
}
func mkJoin(left, right *ast.TableRef, joinType string, cond interface{}) *ast.TableRef {
	ref := &ast.TableRef{RefKind: ast.RefJoin, JoinType: joinType, Left: left, Right: right}
	if c, ok := cond.(joinCond); ok {
		ref.On = c.on
		ref.Using = c.using
	}
	return ref
}
func appendDef(d interface{}, v interface{}) createDefs {
	var out createDefs
	if prev, ok := d.(createDefs); ok {
		out = prev
	}
	switch x := v.(type) {
	case *ast.ColumnDef:
		out.cols = append(out.cols, x)
	case *ast.ConstraintDef:
		out.cons = append(out.cons, x)
	}
	return out
}
func buildColumnDef(name interface{}, dt interface{}, attrs interface{}) *ast.ColumnDef {
	d := asDType(dt)
	cd := &ast.ColumnDef{Name: asIdentName(name), Type: d.name, TypeParams: d.args, Unsigned: d.unsigned}
	if a, ok := attrs.(*colAttrs); ok && a != nil {
		cd.Nullable = a.nullable
		cd.Default = a.def
		cd.OnUpdate = a.onUpdate
		cd.Comment = a.comment
		cd.AutoIncrement = a.autoIncrement
		cd.Charset = a.charset
		cd.Collate = a.collate
		cd.Unsigned = cd.Unsigned || a.unsigned
		cd.Invisible = a.invisible
		if a.hasGenerated {
			cd.Generated = a.generated
			cd.GeneratedStored = a.genStored
		}
	}
	return cd
}
func asOrderList(v interface{}) []ast.OrderItem {
	if s, ok := v.([]ast.OrderItem); ok {
		return s
	}
	return nil
}
func asLimit(v interface{}) *ast.LimitClause {
	if l, ok := v.(*ast.LimitClause); ok {
		return l
	}
	return nil
}
func asAssignList(v interface{}) []ast.Assignment {
	if s, ok := v.([]ast.Assignment); ok {
		return s
	}
	return nil
}
func asTableRefList(v interface{}) []*ast.TableRef {
	if s, ok := v.([]*ast.TableRef); ok {
		return s
	}
	return nil
}
func asCaseWhens(v interface{}) []ast.CaseClause {
	if s, ok := v.([]ast.CaseClause); ok {
		return s
	}
	return nil
}
func asTableOptions(v interface{}) []ast.TableOption {
	if s, ok := v.([]ast.TableOption); ok {
		return s
	}
	return nil
}
func asDefs(v interface{}) createDefs {
	if d, ok := v.(createDefs); ok {
		return d
	}
	return createDefs{}
}
func asTableName(v interface{}) *ast.Identifier {
	return asIdent(v)
}
func asAlterActions(v interface{}) []ast.AlterAction {
	if s, ok := v.([]ast.AlterAction); ok {
		return s
	}
	return nil
}
func asIndexHints(v interface{}) []ast.IndexHint {
	if s, ok := v.([]ast.IndexHint); ok {
		return s
	}
	return nil
}
func asStringList(v interface{}) []string {
	if s, ok := v.([]string); ok {
		return s
	}
	return nil
}
func asPartition(v interface{}) *ast.PartitionSpec {
	if p, ok := v.(*ast.PartitionSpec); ok {
		return p
	}
	return nil
}
func asSignalItems(v interface{}) []ast.SignalItem {
	if s, ok := v.([]ast.SignalItem); ok {
		return s
	}
	return nil
}
func asDiagItems(v interface{}) []ast.DiagItem {
	if s, ok := v.([]ast.DiagItem); ok {
		return s
	}
	return nil
}
func asMergeWhens(v interface{}) []ast.MergeWhen {
	if s, ok := v.([]ast.MergeWhen); ok {
		return s
	}
	return nil
}
func asExprSets(v interface{}) [][]ast.Expression {
	if s, ok := v.([][]ast.Expression); ok {
		return s
	}
	return nil
}
func asRenamePairs(v interface{}) []ast.RenamePair {
	if s, ok := v.([]ast.RenamePair); ok {
		return s
	}
	return nil
}
func asVariableList(v interface{}) []*ast.Variable {
	if s, ok := v.([]*ast.Variable); ok {
		return s
	}
	return nil
}
func asUserSpecs(v interface{}) []ast.UserSpec {
	if s, ok := v.([]ast.UserSpec); ok {
		return s
	}
	return nil
}
func joinIdentNames(ids []*ast.Identifier) string {
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		names = append(names, asIdentName(id))
	}
	return strings.Join(names, ", ")
}
func asLockTables(v interface{}) []ast.LockTable {
	if s, ok := v.([]ast.LockTable); ok {
		return s
	}
	return nil
}
func asPartitionDefs(v interface{}) []ast.PartitionDef {
	if s, ok := v.([]ast.PartitionDef); ok {
		return s
	}
	return nil
}
func asPartitionValues(v interface{}) partitionValues {
	if p, ok := v.(partitionValues); ok {
		return p
	}
	return partitionValues{}
}
func asGroupOpt(v interface{}) groupOpt {
	if g, ok := v.(groupOpt); ok {
		return g
	}
	return groupOpt{}
}

// exprText renders a simple expression to text (used for window frame bounds).
func exprText(e ast.Expression) string {
	switch v := e.(type) {
	case *ast.Literal:
		if v.Text != "" {
			return v.Text
		}
		return v.Value
	case *ast.Identifier:
		return strings.Join(v.Parts, ".")
	}
	if e == nil {
		return ""
	}
	return ""
}
func asColumnDef(v interface{}) *ast.ColumnDef {
	if c, ok := v.(*ast.ColumnDef); ok {
		return c
	}
	return nil
}
func asRows(v interface{}) [][]ast.Expression {
	if s, ok := v.([][]ast.Expression); ok {
		return s
	}
	return nil
}
func asNamedWindows(v interface{}) []*ast.NamedWindow {
	if s, ok := v.([]*ast.NamedWindow); ok {
		return s
	}
	return nil
}
func asTableBody(v interface{}) tableBody {
	if b, ok := v.(tableBody); ok {
		return b
	}
	return tableBody{}
}
func asDBOpts(v interface{}) dbOpts {
	if o, ok := v.(dbOpts); ok {
		return o
	}
	return dbOpts{}
}
func mergeAttrs(a, b *colAttrs) *colAttrs {
	if a == nil {
		a = &colAttrs{}
	}
	if b == nil {
		return a
	}
	if b.nullable != ast.NullableDefault {
		a.nullable = b.nullable
	}
	if b.def != nil {
		a.def = b.def
	}
	if b.onUpdate != nil {
		a.onUpdate = b.onUpdate
	}
	if b.comment != "" {
		a.comment = b.comment
	}
	if b.autoIncrement {
		a.autoIncrement = true
	}
	if b.charset != "" {
		a.charset = b.charset
	}
	if b.collate != "" {
		a.collate = b.collate
	}
	if b.unsigned {
		a.unsigned = true
	}
	if b.invisible {
		a.invisible = true
	}
	if b.hasGenerated {
		a.hasGenerated = true
		a.generated = b.generated
		a.genStored = b.genStored
	}
	return a
}
func appendOpts(v interface{}, o ast.TableOption) []ast.TableOption {
	var out []ast.TableOption
	if s, ok := v.([]ast.TableOption); ok {
		out = s
	}
	return append(out, o)
}
func variableFromToken(t token.Token) *ast.Variable {
	v := &ast.Variable{}
	if t.Type == token.SYSVARIABLE {
		v.IsSystem = true
		name := strings.TrimPrefix(t.Literal, "@@")
		low := strings.ToLower(name)
		for _, scope := range []string{"global", "session", "local"} {
			if strings.HasPrefix(low, scope+".") {
				v.Scope = scope
				name = name[len(scope)+1:]
				break
			}
		}
		v.Name = name
	} else {
		v.Name = strings.TrimPrefix(t.Literal, "@")
	}
	return v
}
func variableText(v *ast.Variable) string {
	if v.IsSystem {
		if v.Scope != "" {
			return "@@" + v.Scope + "." + v.Name
		}
		return "@@" + v.Name
	}
	return "@" + v.Name
}
