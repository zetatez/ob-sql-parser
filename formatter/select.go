package formatter

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
)

func (f *formatter) selectStmt(v *ast.SelectStmt) {
	if v == nil {
		return
	}
	for i, body := range v.Bodies {
		if i > 0 {
			op := v.SetOps[i-1]
			f.printf(" %s", op.Op)
			if op.All {
				f.raw(" ALL")
			}
			f.raw(" ")
		}
		f.selectBody(body)
	}
	if len(v.OrderBy) > 0 {
		f.nl()
		f.raw("ORDER BY ")
		f.orderItems(v.OrderBy)
	}
	if v.Limit != nil {
		f.nl()
		f.limit(v.Limit)
	}
	if v.IntoKind == "OUTFILE" {
		f.printf(" INTO OUTFILE %s", f.stringLiteral(v.IntoFile))
		if v.IntoFields != "" {
			f.printf(" FIELDS TERMINATED BY %s", f.stringLiteral(v.IntoFields))
		}
		if v.IntoLines != "" {
			f.printf(" LINES TERMINATED BY %s", f.stringLiteral(v.IntoLines))
		}
	} else if v.IntoKind == "DUMPFILE" {
		f.printf(" INTO DUMPFILE %s", f.stringLiteral(v.IntoFile))
	} else if v.IntoKind == "@vars" {
		f.raw(" INTO ")
		for i, vr := range v.IntoVars {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.variableString(vr))
		}
	}
}
func (f *formatter) selectBody(b *ast.SelectBody) {
	f.raw("SELECT ")
	if b.Distinct {
		f.raw("DISTINCT ")
	}
	for i, it := range b.Items {
		if i > 0 {
			f.raw(", ")
			if f.opts.Pretty {
				f.nl()
			}
		}
		f.projectItem(&it)
	}
	f.depth++
	if len(b.From) > 0 {
		f.nl()
		f.raw("FROM ")
		for i, r := range b.From {
			if i > 0 {
				f.raw(", ")
			}
			f.tableRef(r)
		}
	}
	if b.Where != nil {
		f.nl()
		f.raw("WHERE ")
		f.expr(b.Where)
	}
	if len(b.GroupBy) > 0 {
		f.nl()
		f.raw("GROUP BY ")
		for i, e := range b.GroupBy {
			if i > 0 {
				f.raw(", ")
			}
			f.expr(e)
		}
		if b.Rollup {
			f.raw(" WITH ROLLUP")
		}
	}
	if b.Having != nil {
		f.nl()
		f.raw("HAVING ")
		f.expr(b.Having)
	}
	if b.Qualify != nil {
		f.nl()
		f.raw("QUALIFY ")
		f.expr(b.Qualify)
	}
	if len(b.OrderBy) > 0 {
		f.nl()
		f.raw("ORDER BY ")
		f.orderItems(b.OrderBy)
	}
	if b.Limit != nil {
		f.nl()
		f.limit(b.Limit)
	}
	if len(b.NamedWindow) > 0 {
		f.nl()
		f.raw("WINDOW ")
		for i, w := range b.NamedWindow {
			if i > 0 {
				f.raw(", ")
			}
			f.printf("%s AS (%s)", w.Name, f.windowSpecInner(w.Spec))
		}
	}
	if b.InShareMode {
		f.nl()
		if b.LockWait == "SHARE" {
			f.raw("FOR SHARE")
		} else {
			f.raw("LOCK IN SHARE MODE")
		}
	}
	if b.ForUpdate {
		f.nl()
		f.raw("FOR UPDATE")
		if b.LockWait != "" {
			f.printf(" %s", b.LockWait)
		}
	}
	f.depth--
	if b.HasProcAnalyse {
		f.raw(" PROCEDURE ANALYSE(")
		for i, e := range b.ProcAnalyse {
			if i > 0 {
				f.raw(", ")
			}
			f.expr(e)
		}
		f.raw(")")
	}
}
func (f *formatter) projectItem(it *ast.ColumnItem) {
	f.expr(it.Expr)
	if it.HasAlias {
		f.printf(" AS %s", f.quoteIfNeeded(it.Alias))
	}
}
func (f *formatter) orderItems(items []ast.OrderItem) {
	for i, o := range items {
		if i > 0 {
			f.raw(", ")
		}
		f.expr(o.Expr)
		if o.Desc {
			f.raw(" DESC")
		}
		if o.HasNulls {
			if o.NullsFirst {
				f.raw(" NULLS FIRST")
			} else {
				f.raw(" NULLS LAST")
			}
		}
	}
}
func (f *formatter) limit(l *ast.LimitClause) {
	f.raw("LIMIT ")
	if l.HasOffset && l.Offset != nil {
		f.expr(l.Offset)
		f.raw(", ")
		f.expr(l.Count)
		return
	}
	if l.Count != nil {
		f.expr(l.Count)
	}
}
func (f *formatter) tableRef(r *ast.TableRef) {
	if r == nil {
		return
	}
	switch r.RefKind {
	case ast.RefTable:
		f.raw(f.ident(r.Name))
		if len(r.Partitions) > 0 {
			f.raw(" PARTITION (")
			for i, p := range r.Partitions {
				if i > 0 {
					f.raw(", ")
				}
				f.raw(f.quoteIfNeeded(p))
			}
			f.raw(")")
		}
		if r.HasAlias {
			f.printf(" AS %s", f.quoteIfNeeded(r.Alias))
		}
		if len(r.IndexHints) > 0 {
			f.indexHints(r.IndexHints)
		}
	case ast.RefSubquery:
		if r.Lateral {
			f.raw("LATERAL ")
		}
		f.raw("(")
		f.selectStmt(r.Subquery)
		f.raw(")")
		if r.HasAlias {
			f.printf(" AS %s", f.quoteIfNeeded(r.Alias))
		}
	case ast.RefJoin:
		f.tableRef(r.Left)
		f.printf(" %s ", r.JoinType)
		f.tableRef(r.Right)
		if r.On != nil {
			f.raw(" ON ")
			f.expr(r.On)
		} else if len(r.Using) > 0 {
			f.raw(" USING (")
			for i, u := range r.Using {
				if i > 0 {
					f.raw(", ")
				}
				f.raw(f.ident(u))
			}
			f.raw(")")
		}
	}
}
func (f *formatter) indexHints(hints []ast.IndexHint) {
	for _, h := range hints {
		f.raw(" ")
		f.raw(string(h.Type))
		f.raw(" INDEX")
		if h.For != "" {
			f.printf(" FOR %s", h.For)
		}
		f.raw(" (")
		for i, idx := range h.Indexes {
			if i > 0 {
				f.raw(", ")
			}
			if strings.EqualFold(idx, "PRIMARY") {
				f.raw("PRIMARY")
			} else {
				f.raw(f.quoteIfNeeded(idx))
			}
		}
		f.raw(")")
	}
}
