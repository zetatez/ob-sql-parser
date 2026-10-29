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
		f.tableOption(o)
	}
	if v.ColumnGroup != "" {
		f.raw(" " + v.ColumnGroup)
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
		if pd.MaxValue {
			f.raw(" MAXVALUE")
		} else {
			f.raw(" (")
			for i, v := range pd.Values {
				if i > 0 {
					f.raw(", ")
				}
				f.expr(v)
			}
			f.raw(")")
		}
	}
	for _, o := range pd.Options {
		f.tableOption(o)
	}
	if len(pd.Subs) > 0 {
		f.raw(" (")
		for i := range pd.Subs {
			if i > 0 {
				f.raw(", ")
			}
			f.raw("SUBPARTITION " + f.quoteIfNeeded(pd.Subs[i].Name))
			for _, o := range pd.Subs[i].Options {
				f.tableOption(o)
			}
		}
		f.raw(")")
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
	if c.Zerofill {
		sb.WriteString(" ZEROFILL")
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
	if c.Check != nil {
		sb.WriteString(" CHECK (")
		sb.WriteString(f.exprString(c.Check))
		sb.WriteString(")")
	}
	if c.Srid != "" {
		sb.WriteString(" SRID " + c.Srid)
	}
	if c.ColumnFormat != "" {
		sb.WriteString(" COLUMN_FORMAT " + c.ColumnFormat)
	}
	if c.Storage != "" {
		sb.WriteString(" STORAGE " + c.Storage)
	}
	if c.BinaryAttr {
		sb.WriteString(" BINARY")
	}
	if c.PrimaryKey {
		sb.WriteString(" PRIMARY KEY")
	}
	if c.UniqueKey {
		sb.WriteString(" UNIQUE")
	}
	if c.RefTable != nil {
		sb.WriteString(" REFERENCES ")
		sb.WriteString(f.ident(c.RefTable))
		if len(c.RefColumns) > 0 {
			sb.WriteString(" (")
			for i, col := range c.RefColumns {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(f.ident(col))
			}
			sb.WriteString(")")
		}
		if c.OnDelete != "" {
			sb.WriteString(" ON DELETE " + c.OnDelete)
		}
		if c.FkOnUpdate != "" {
			sb.WriteString(" ON UPDATE " + c.FkOnUpdate)
		}
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
		if len(c.IndexCols) > 0 {
			sb.WriteString("UNIQUE KEY ")
		} else {
			sb.WriteString("UNIQUE ")
		}
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
	if len(c.IndexCols) > 0 {
		for i, col := range c.IndexCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.ident(col.Name))
			if col.Length != nil {
				sb.WriteString("(" + f.exprString(col.Length) + ")")
			}
			if col.Desc {
				sb.WriteString(" DESC")
			} else if col.HasOrder {
				sb.WriteString(" ASC")
			}
		}
	} else {
		for i, col := range c.Columns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.ident(col))
		}
	}
	sb.WriteString(")")
	if c.Kind == "FOREIGN KEY" && c.RefTable != nil {
		sb.WriteString(" REFERENCES ")
		sb.WriteString(f.ident(c.RefTable))
		if len(c.RefColumns) > 0 {
			sb.WriteString(" (")
			for i, col := range c.RefColumns {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(f.ident(col))
			}
			sb.WriteString(")")
		}
		if c.OnDelete != "" {
			sb.WriteString(" ON DELETE " + c.OnDelete)
		}
		if c.OnUpdate != "" {
			sb.WriteString(" ON UPDATE " + c.OnUpdate)
		}
	}
	if len(c.Storing) > 0 {
		sb.WriteString(" STORING (")
		for i, col := range c.Storing {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.ident(col))
		}
		sb.WriteString(")")
	}
	if len(c.WithOpts) > 0 {
		sb.WriteString(" WITH (" + strings.Join(c.WithOpts, ", ") + ")")
	}
	for _, o := range c.Options {
		sb.WriteString(" " + o)
	}
	if c.ColumnGroup != "" {
		sb.WriteString(" " + c.ColumnGroup)
	}
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
		f.raw("DROP COLUMN ")
		if a.IfExists {
			f.raw("IF EXISTS ")
		}
		f.raw(f.quoteIfNeeded(a.Name))
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
		f.colPosition(a)
	case ast.AlterRenameColumn:
		f.printf("RENAME COLUMN %s TO %s", f.quoteIfNeeded(a.Name), f.quoteIfNeeded(a.NewName))
	case ast.AlterRenameTable:
		if a.Target != nil {
			f.printf("RENAME TO %s", f.ident(a.Target))
		} else {
			f.printf("RENAME TO %s", f.quoteIfNeeded(a.NewName))
		}
	case ast.AlterAddFulltextIndex:
		f.raw("ADD FULLTEXT INDEX ")
		if a.Name != "" {
			f.printf("%s ", f.quoteIfNeeded(a.Name))
		}
		if len(a.IndexCols) > 0 {
			f.indexColsParen(a.IndexCols)
		} else {
			f.colsParen(a.Columns)
		}
		for _, o := range a.IndexOptions {
			f.printf(" %s", o)
		}
	case ast.AlterAddSpatialIndex:
		f.raw("ADD SPATIAL INDEX ")
		if a.Name != "" {
			f.printf("%s ", f.quoteIfNeeded(a.Name))
		}
		if len(a.IndexCols) > 0 {
			f.indexColsParen(a.IndexCols)
		} else {
			f.colsParen(a.Columns)
		}
		for _, o := range a.IndexOptions {
			f.printf(" %s", o)
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
		if len(a.IndexCols) > 0 {
			f.indexColsParen(a.IndexCols)
		} else {
			f.colsParen(a.Columns)
		}
		if len(a.Storing) > 0 {
			f.raw(" STORING (")
			f.idents(a.Storing)
			f.raw(")")
		}
		if len(a.WithOpts) > 0 {
			f.raw(" WITH (" + strings.Join(a.WithOpts, ", ") + ")")
		}
		for _, o := range a.IndexOptions {
			f.printf(" %s", o)
		}
		if a.ColumnGroup != "" {
			f.printf(" %s", a.ColumnGroup)
		}
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
		f.raw("ADD COLUMN (")
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
	case ast.AlterVisibleIndex:
		f.printf("ALTER INDEX %s VISIBLE", f.quoteIfNeeded(a.Name))
	case ast.AlterInvisibleIndex:
		f.printf("ALTER INDEX %s INVISIBLE", f.quoteIfNeeded(a.Name))
	case ast.AlterDropCheck:
		f.printf("DROP CHECK %s", f.quoteIfNeeded(a.Name))
	case ast.AlterDropConstraint:
		f.printf("DROP CONSTRAINT %s", f.quoteIfNeeded(a.Name))
	case ast.AlterDropForeignKey:
		f.printf("DROP FOREIGN KEY %s", f.quoteIfNeeded(a.Name))
	case ast.AlterSetColumnVisible:
		f.printf("ALTER COLUMN %s SET VISIBLE", f.quoteIfNeeded(a.Name))
	case ast.AlterSetColumnInvisible:
		f.printf("ALTER COLUMN %s SET INVISIBLE", f.quoteIfNeeded(a.Name))
	case ast.AlterForce:
		f.raw("FORCE")
	case ast.AlterRepartition:
		f.raw("PARTITION BY ")
		if a.Partition != nil {
			f.raw(string(a.Partition.Kind))
			f.raw(" (")
			f.expr(a.Partition.Expr)
			f.raw(")")
			if a.Partition.Count != nil {
				f.raw(" PARTITIONS ")
				f.expr(a.Partition.Count)
			}
		}
	case ast.AlterAlgorithm:
		f.printf("ALGORITHM = %s", a.Name)
	case ast.AlterLock:
		f.printf("LOCK = %s", a.Name)
	case ast.AlterAssignTablespace:
		f.printf("TABLESPACE %s", f.quoteIfNeeded(a.Name))
	case ast.AlterEnableKeys:
		f.raw("ENABLE KEYS")
	case ast.AlterDisableKeys:
		f.raw("DISABLE KEYS")
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

// indexColsParen 打印索引列列表（含前缀长度与排序方向）。
func (f *formatter) indexColsParen(cols []ast.IndexCol) {
	f.raw("(")
	for i, c := range cols {
		if i > 0 {
			f.raw(", ")
		}
		f.raw(f.ident(c.Name))
		if c.Length != nil {
			f.printf("(%s)", f.exprString(c.Length))
		}
		if c.Desc {
			f.raw(" DESC")
		} else if c.HasOrder {
			f.raw(" ASC")
		}
	}
	f.raw(")")
}

// tableOption 渲染 CREATE TABLE 选项；TABLESPACE 不带等号。
func (f *formatter) tableOption(o ast.TableOption) {
	if o.Name == "tablespace" {
		f.printf(" TABLESPACE %s", o.Value)
		return
	}
	f.printf(" %s=%s", stringsUpper(o.Name), o.Value)
}

// idents 以逗号分隔打印标识符列表（RETURNING 等）。
func (f *formatter) idents(cols []*ast.Identifier) {
	for i, c := range cols {
		if i > 0 {
			f.raw(", ")
		}
		f.raw(f.ident(c))
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
