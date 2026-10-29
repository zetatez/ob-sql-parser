package formatter

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
)

func (f *formatter) createTable(v *ast.CreateTableStmt) {
	f.raw("CREATE TABLE ")
	if v.IfNotExists {
		f.raw("IF NOT EXISTS ")
	}
	f.raw(f.ident(v.Table))
	if len(v.Columns) > 0 || len(v.Constraints) > 0 {
		f.raw(" (")
		f.depth++
		n := 0
		for _, c := range v.Columns {
			if n > 0 {
				f.raw(",")
			}
			f.nl()
			f.raw(f.columnDef(c))
			n++
		}
		for _, c := range v.Constraints {
			if n > 0 {
				f.raw(",")
			}
			f.nl()
			f.raw(f.constraintDef(c))
			n++
		}
		f.depth--
		if f.opts.Pretty {
			f.raw("\n")
			f.raw(")")
		} else {
			f.raw(")")
		}
	}
	for _, o := range v.Options {
		f.printf(" %s=%s", stringsUpper(o.Name), o.Value)
	}
	if v.Partition != nil {
		f.raw(" ")
		f.partitionSpec(v.Partition)
	}
	if v.Select != nil {
		f.raw(" AS ")
		f.selectStmt(v.Select)
	}
}
func (f *formatter) partitionSpec(ps *ast.PartitionSpec) {
	f.raw("PARTITION BY ")
	f.raw(string(ps.Kind))
	if len(ps.Columns) > 0 {
		f.raw(" ")
		f.colsParen(ps.Columns)
	} else if ps.Expr != nil {
		f.raw(" (")
		f.expr(ps.Expr)
		f.raw(")")
	}
	if ps.Auto {
		if ps.Size != nil {
			f.raw(" PARTITION SIZE ")
			f.expr(ps.Size)
		}
		f.raw(" PARTITIONS AUTO")
	} else if ps.Count != nil {
		f.raw(" PARTITIONS ")
		f.expr(ps.Count)
	}
	if ps.SubKind != "" {
		f.raw(" SUBPARTITION BY ")
		f.raw(string(ps.SubKind))
		if len(ps.SubColumns) > 0 {
			f.raw(" ")
			f.colsParen(ps.SubColumns)
		} else if ps.SubExpr != nil {
			f.raw(" (")
			f.expr(ps.SubExpr)
			f.raw(")")
		}
		if ps.SubCount != nil {
			f.raw(" SUBPARTITIONS ")
			f.expr(ps.SubCount)
		}
	}
	if len(ps.Partitions) > 0 {
		f.raw(" (")
		for i := range ps.Partitions {
			if i > 0 {
				f.raw(", ")
			}
			f.partitionDef(&ps.Partitions[i])
		}
		f.raw(")")
	}
}
func (f *formatter) partitionDef(pd *ast.PartitionDef) {
	f.raw("PARTITION ")
	f.raw(f.quoteIfNeeded(pd.Name))
	if pd.ValueKind != "" {
		f.raw(" VALUES ")
		f.raw(pd.ValueKind)
		f.raw(" (")
		for i, v := range pd.Values {
			if i > 0 {
				f.raw(", ")
			}
			f.expr(v)
		}
		f.raw(")")
	}
	for _, o := range pd.Options {
		f.printf(" %s=%s", stringsUpper(o.Name), o.Value)
	}
}
func (f *formatter) columnDef(c *ast.ColumnDef) string {
	var sb strings.Builder
	sb.WriteString(f.quoteIfNeeded(c.Name))
	sb.WriteString(" ")
	sb.WriteString(c.Type)
	if len(c.TypeParams) > 0 {
		sb.WriteString("(")
		for i, a := range c.TypeParams {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.exprString(a))
		}
		sb.WriteString(")")
	}
	if c.Unsigned {
		sb.WriteString(" UNSIGNED")
	}
	if c.Generated != nil {
		sb.WriteString(" GENERATED ALWAYS AS (")
		sb.WriteString(f.exprString(c.Generated))
		sb.WriteString(")")
		if c.GeneratedStored {
			sb.WriteString(" STORED")
		} else {
			sb.WriteString(" VIRTUAL")
		}
	}
	if c.Nullable == ast.NullableNo {
		sb.WriteString(" NOT NULL")
	}
	if c.Nullable == ast.NullableYes {
		sb.WriteString(" NULL")
	}
	if c.Default != nil {
		sb.WriteString(" DEFAULT ")
		sb.WriteString(f.exprString(c.Default))
	}
	if c.OnUpdate != nil {
		sb.WriteString(" ON UPDATE ")
		sb.WriteString(f.exprString(c.OnUpdate))
	}
	if c.AutoIncrement {
		sb.WriteString(" AUTO_INCREMENT")
	}
	if c.Comment != "" {
		sb.WriteString(" COMMENT ")
		sb.WriteString(f.stringLiteral(c.Comment))
	}
	if c.Charset != "" {
		sb.WriteString(" CHARSET ")
		sb.WriteString(c.Charset)
	}
	if c.Collate != "" {
		sb.WriteString(" COLLATE ")
		sb.WriteString(c.Collate)
	}
	return sb.String()
}
func (f *formatter) constraintDef(c *ast.ConstraintDef) string {
	var sb strings.Builder
	if c.Name != "" {
		sb.WriteString("CONSTRAINT ")
		sb.WriteString(f.quoteIfNeeded(c.Name))
		sb.WriteString(" ")
	}
	switch c.Kind {
	case "PRIMARY KEY":
		sb.WriteString("PRIMARY KEY ")
	case "UNIQUE":
		sb.WriteString("UNIQUE ")
	case "FOREIGN KEY":
		sb.WriteString("FOREIGN KEY ")
	case "CHECK":
		sb.WriteString("CHECK (")
		if c.Check != nil {
			sb.WriteString(f.exprString(c.Check))
		}
		sb.WriteString(")")
		return sb.String()
	default:
		sb.WriteString("KEY ")
	}
	sb.WriteString("(")
	for i, col := range c.Columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(f.ident(col))
	}
	sb.WriteString(")")
	return sb.String()
}
func (f *formatter) alterTable(v *ast.AlterTableStmt) {
	f.printf("ALTER TABLE %s ", f.ident(v.Table))
	for i := range v.Actions {
		if i > 0 {
			f.raw(", ")
		}
		f.alterAction(&v.Actions[i])
	}
}
func (f *formatter) alterAction(a *ast.AlterAction) {
	switch a.Kind {
	case ast.AlterAddColumn:
		f.raw("ADD COLUMN ")
		if a.Column != nil {
			f.raw(f.columnDef(a.Column))
		}
		f.colPosition(a)
	case ast.AlterDropColumn:
		f.printf("DROP COLUMN %s", f.quoteIfNeeded(a.Name))
	case ast.AlterModifyColumn:
		f.raw("MODIFY COLUMN ")
		if a.Column != nil {
			f.raw(f.columnDef(a.Column))
		}
		f.colPosition(a)
	case ast.AlterChangeColumn:
		f.printf("CHANGE COLUMN %s ", f.quoteIfNeeded(a.OldName))
		if a.Column != nil {
			f.raw(f.columnDef(a.Column))
		}
	case ast.AlterRenameColumn:
		f.printf("RENAME COLUMN %s TO %s", f.quoteIfNeeded(a.Name), f.quoteIfNeeded(a.NewName))
	case ast.AlterRenameTable:
		if a.Target != nil {
			f.printf("RENAME TO %s", f.ident(a.Target))
		} else {
			f.printf("RENAME TO %s", f.quoteIfNeeded(a.NewName))
		}
	case ast.AlterAddIndex:
		f.raw("ADD ")
		if a.Unique {
			f.raw("UNIQUE ")
		}
		f.raw("INDEX ")
		if a.Name != "" {
			f.printf("%s ", f.quoteIfNeeded(a.Name))
		}
		f.colsParen(a.Columns)
	case ast.AlterDropIndex:
		f.printf("DROP INDEX %s", f.quoteIfNeeded(a.Name))
	case ast.AlterRenameIndex:
		f.printf("RENAME INDEX %s TO %s", f.quoteIfNeeded(a.Name), f.quoteIfNeeded(a.NewName))
	case ast.AlterDropPrimaryKey:
		f.raw("DROP PRIMARY KEY")
	case ast.AlterAddConstraint:
		f.raw("ADD ")
		if a.Constraint != nil {
			f.raw(f.constraintDef(a.Constraint))
		}
	case ast.AlterSetDefault:
		f.printf("ALTER COLUMN %s SET DEFAULT ", f.quoteIfNeeded(a.Name))
		f.expr(a.Default)
	case ast.AlterDropDefault:
		f.printf("ALTER COLUMN %s DROP DEFAULT", f.quoteIfNeeded(a.Name))
	case ast.AlterAddPartition:
		f.raw("ADD PARTITION (")
		for i := range a.PartitionDefs {
			if i > 0 {
				f.raw(", ")
			}
			f.partitionDef(&a.PartitionDefs[i])
		}
		f.raw(")")
	case ast.AlterDropPartition, ast.AlterTruncatePartition:
		if a.All {
			f.printf("%s ALL", a.Kind)
		} else {
			f.printf("%s %s", a.Kind, strings.Join(a.Names, ", "))
		}
	case ast.AlterCoalescePartition:
		f.raw("COALESCE PARTITION ")
		if a.Count != nil {
			f.expr(a.Count)
		}
	case ast.AlterExchangePartition:
		f.printf("EXCHANGE PARTITION %s WITH TABLE %s", strings.Join(a.Names, ", "), f.ident(a.With))
		if a.Validation != "" {
			f.printf(" %s", a.Validation)
		}
	case ast.AlterReorganizePartition:
		if a.All {
			f.raw("REORGANIZE PARTITION ALL INTO (")
		} else {
			f.printf("REORGANIZE PARTITION %s INTO (", strings.Join(a.Names, ", "))
		}
		for i := range a.PartitionDefs {
			if i > 0 {
				f.raw(", ")
			}
			f.partitionDef(&a.PartitionDefs[i])
		}
		f.raw(")")
	case ast.AlterRemovePartitioning:
		f.raw("REMOVE PARTITIONING")
	case ast.AlterAddColumns:
		f.raw("ADD (")
		for i, cd := range a.ColumnDefs {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.columnDef(cd))
		}
		f.raw(")")
	case ast.AlterConvertCharset:
		f.printf("CONVERT TO CHARACTER SET %s", a.Charset)
		if a.Collate != "" {
			f.printf(" COLLATE %s", a.Collate)
		}
	case ast.AlterAutoIncrement:
		f.raw("AUTO_INCREMENT = ")
		f.expr(a.Count)
	}
}
func (f *formatter) propKV(o *ast.Option) {
	f.printf("%s", stringsUpper(o.Name))
	if len(o.Values) > 0 {
		f.raw(" (")
		for i, v := range o.Values {
			if i > 0 {
				f.raw(", ")
			}
			f.expr(v)
		}
		f.raw(")")
	} else {
		f.raw(" = ")
		f.expr(o.Value)
	}
}
func (f *formatter) propList(opts []ast.Option) {
	for i := range opts {
		if i > 0 {
			f.raw(", ")
		}
		f.propKV(&opts[i])
	}
}
func (f *formatter) optionKV(o *ast.Option) {
	f.printf("%s = ", stringsUpper(o.Name))
	if len(o.Values) > 0 {
		f.raw("(")
		for i, v := range o.Values {
			if i > 0 {
				f.raw(", ")
			}
			f.expr(v)
		}
		f.raw(")")
	} else {
		f.expr(o.Value)
	}
}
func (f *formatter) colPosition(a *ast.AlterAction) {
	if a.Position == "" {
		return
	}
	f.printf(" %s", a.Position)
	if a.Position == "AFTER" && a.RefCol != "" {
		f.printf(" %s", f.quoteIfNeeded(a.RefCol))
	}
}
func (f *formatter) colsParen(cols []*ast.Identifier) {
	f.raw("(")
	for i, c := range cols {
		if i > 0 {
			f.raw(", ")
		}
		f.raw(f.ident(c))
	}
	f.raw(")")
}
