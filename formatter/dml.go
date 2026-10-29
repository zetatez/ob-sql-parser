package formatter

import "github.com/zetatez/ob-sql-parser/ast"

func (f *formatter) insertStmt(v *ast.InsertStmt) {
	if v.Replace {
		f.raw("REPLACE ")
	} else {
		f.raw("INSERT ")
	}
	if v.Ignore {
		f.raw("IGNORE ")
	}
	f.printf("INTO %s", f.ident(v.Table))
	if len(v.Columns) > 0 {
		f.raw(" (")
		for i, c := range v.Columns {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(c))
		}
		f.raw(")")
	}
	if len(v.Rows) > 0 {
		f.raw(" VALUES")
		f.depth++
		for i, row := range v.Rows {
			f.nl()
			f.raw("(")
			for j, e := range row {
				if j > 0 {
					f.raw(", ")
				}
				f.expr(e)
			}
			f.raw(")")
			if i+1 < len(v.Rows) {
				f.raw(",")
			}
		}
		f.depth--
	} else if v.Source != nil {
		f.raw(" ")
		f.selectStmt(v.Source)
	} else if len(v.Set) > 0 {
		f.raw(" SET ")
		f.assignments(v.Set)
	}
	if len(v.OnDuplicate) > 0 {
		f.raw(" ON DUPLICATE KEY UPDATE ")
		f.assignments(v.OnDuplicate)
	}
	if len(v.Returning) > 0 {
		f.raw(" RETURNING ")
		for i, c := range v.Returning {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(c))
		}
	}
}
func (f *formatter) updateStmt(v *ast.UpdateStmt) {
	f.raw("UPDATE ")
	if v.LowPriority {
		f.raw("LOW_PRIORITY ")
	}
	if v.Ignore {
		f.raw("IGNORE ")
	}
	for i, t := range v.Tables {
		if i > 0 {
			f.raw(", ")
		}
		f.tableRef(t)
	}
	f.raw(" SET ")
	f.assignments(v.Set)
	if v.Where != nil {
		f.raw(" WHERE ")
		f.expr(v.Where)
	}
	if len(v.OrderBy) > 0 {
		f.raw(" ORDER BY ")
		f.orderItems(v.OrderBy)
	}
	if v.Limit != nil {
		f.raw(" ")
		f.limit(v.Limit)
	}
	if len(v.Returning) > 0 {
		f.raw(" RETURNING ")
		f.idents(v.Returning)
	}
}
func (f *formatter) deleteStmt(v *ast.DeleteStmt) {
	f.raw("DELETE ")
	if v.Mods != "" {
		f.raw(v.Mods + " ")
	}
	if len(v.Tables) > 0 {
		// 多表删除：DELETE t1, t2 FROM ...
		for i, t := range v.Tables {
			if i > 0 {
				f.raw(", ")
			}
			f.raw(f.ident(t))
		}
		f.raw(" FROM ")
	} else {
		f.raw("FROM ")
	}
	for i, t := range v.From {
		if i > 0 {
			f.raw(", ")
		}
		f.tableRef(t)
	}
	if len(v.Using) > 0 {
		f.raw(" USING ")
		for i, t := range v.Using {
			if i > 0 {
				f.raw(", ")
			}
			f.tableRef(t)
		}
	}
	if v.Where != nil {
		f.raw(" WHERE ")
		f.expr(v.Where)
	}
	if len(v.OrderBy) > 0 {
		f.raw(" ORDER BY ")
		f.orderItems(v.OrderBy)
	}
	if v.Limit != nil {
		f.raw(" ")
		f.limit(v.Limit)
	}
	if len(v.Returning) > 0 {
		f.raw(" RETURNING ")
		f.idents(v.Returning)
	}
}
func (f *formatter) assignments(assigns []ast.Assignment) {
	for i, a := range assigns {
		if i > 0 {
			f.raw(", ")
		}
		f.printf("%s = ", f.ident(a.Column))
		f.expr(a.Expr)
	}
}
func (f *formatter) mergeWhen(w *ast.MergeWhen) {
	f.raw(" WHEN ")
	if w.Matched {
		f.raw("MATCHED ")
	} else {
		f.raw("NOT MATCHED ")
	}
	f.raw("THEN ")
	switch w.Action {
	case "UPDATE":
		f.raw("UPDATE SET ")
		f.assignments(w.Set)
		if w.DeleteCond != nil {
			f.raw(" DELETE WHERE ")
			f.expr(w.DeleteCond)
		}
	case "DELETE":
		f.raw("DELETE")
		if w.DeleteCond != nil {
			f.raw(" WHERE ")
			f.expr(w.DeleteCond)
		}
	case "INSERT":
		f.raw("INSERT")
		if len(w.Columns) > 0 {
			f.raw(" ")
			f.colsParen(w.Columns)
		}
		f.raw(" VALUES (")
		for i, e := range w.Values {
			if i > 0 {
				f.raw(", ")
			}
			f.expr(e)
		}
		f.raw(")")
	}
}
