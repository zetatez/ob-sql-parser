package formatter

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
)

func (f *formatter) options(opts []ast.Option) {
	for i := range opts {
		if i == 0 {
			f.raw(" ")
		} else {
			f.raw(", ")
		}
		f.optionKV(&opts[i])
	}
}
func (f *formatter) userSpecs(users []ast.UserSpec) {
	for i := range users {
		if i > 0 {
			f.raw(", ")
		}
		f.userSpec(&users[i])
	}
}
func (f *formatter) userSpec(u *ast.UserSpec) {
	f.raw("'" + strings.ReplaceAll(u.Name, "'", "''") + "'")
	if u.HasHost {
		f.raw("@'" + strings.ReplaceAll(u.Host, "'", "''") + "'")
	}
	if u.Plugin != "" {
		f.raw(" IDENTIFIED WITH " + u.Plugin)
		if u.HasPassword {
			f.raw(" BY '" + strings.ReplaceAll(u.Password, "'", "''") + "'")
		}
	} else if u.HasPassword {
		f.raw(" IDENTIFIED BY '" + strings.ReplaceAll(u.Password, "'", "''") + "'")
	}
}
func (f *formatter) showStmt(v *ast.ShowStmt) {
	if v.Kind == "describe table" {
		f.raw("DESCRIBE ")
		if v.Target != nil {
			f.raw(f.ident(v.Target))
		}
		return
	}
	f.raw("SHOW ")
	if v.Full {
		f.raw("FULL ")
	}
	if v.Scope != "" {
		f.printf("%s ", v.Scope)
	}
	f.raw(stringsUpper(v.Kind))
	if v.Target != nil {
		if showKindNeedsFrom(v.Kind) {
			f.printf(" FROM %s", f.ident(v.Target))
		} else {
			f.printf(" %s", f.ident(v.Target))
		}
	}
	if v.From != nil {
		f.printf(" FROM %s", f.ident(v.From))
	}
	if v.User != nil {
		if v.Kind == "grants" {
			f.raw(" FOR ")
		} else {
			f.raw(" ")
		}
		f.userSpec(v.User)
	}
	if v.Mode != "" {
		f.printf(" %s", v.Mode)
	}
	if v.Like != nil {
		f.raw(" LIKE ")
		f.expr(v.Like)
	}
	if v.Where != nil {
		f.raw(" WHERE ")
		f.expr(v.Where)
	}
}
func (f *formatter) setStmt(v *ast.SetStmt) {
	f.raw("SET ")
	for i, it := range v.Items {
		if i > 0 {
			f.raw(", ")
		}
		if it.Scope != "" {
			f.printf("%s ", stringsUpper(it.Scope))
		}
		f.raw(it.Name)
		if it.NoValue {
			continue
		}
		f.raw(" = ")
		f.expr(it.Value)
	}
}
