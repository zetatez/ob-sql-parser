package formatter

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
)

func (f *formatter) expr(e ast.Expression) { f.raw(f.exprString(e)) }
func (f *formatter) exprString(e ast.Expression) string {
	switch v := e.(type) {
	case *ast.Literal:
		return f.literalString(v)
	case *ast.Identifier:
		return f.ident(v)
	case *ast.StarLiteral:
		if v.Qualifier != "" {
			return v.Qualifier + ".*"
		}
		return "*"
	case *ast.ParamMarker:
		return "?"
	case *ast.Variable:
		return f.variableString(v)
	case *ast.UnaryExpr:
		return f.unaryString(v)
	case *ast.BinaryExpr:
		return f.binaryString(v)
	case *ast.FuncCall:
		return f.funcString(v)
	case *ast.CastExpr:
		return f.castString(v)
	case *ast.CaseExpr:
		return f.caseString(v)
	case *ast.BetweenExpr:
		return f.betweenString(v)
	case *ast.InExpr:
		return f.inString(v)
	case *ast.LikeExpr:
		return f.likeString(v)
	case *ast.IsExpr:
		return f.isString(v)
	case *ast.IntervalExpr:
		return f.intervalString(v)
	case *ast.ExistsExpr:
		n := "NOT "
		if !v.Not {
			n = ""
		}
		return n + "EXISTS (" + Format(v.Query) + ")"
	case *ast.SubqueryExpr:
		return "(" + Format(v.Query) + ")"
	case *ast.ParenExpr:
		return "(" + f.exprString(v.Expr) + ")"
	case *ast.MemberOfExpr:
		var sb strings.Builder
		sb.WriteString(f.exprString(v.Expr))
		sb.WriteString(" MEMBER OF (")
		for i, x := range v.List {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.exprString(x))
		}
		sb.WriteString(")")
		return sb.String()
	case *ast.MatchAgainstExpr:
		var sb strings.Builder
		sb.WriteString("MATCH (")
		for i, c := range v.Columns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.ident(c))
		}
		sb.WriteString(") AGAINST (")
		sb.WriteString(f.exprString(v.Against))
		if v.Modifier != "" {
			sb.WriteString(" ")
			sb.WriteString(v.Modifier)
		}
		sb.WriteString(")")
		return sb.String()
	case *ast.GroupingSetExpr:
		var sb strings.Builder
		sb.WriteString(v.Kind)
		sb.WriteString(" (")
		if len(v.Sets) > 0 {
			for i, set := range v.Sets {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString("(")
				for j, x := range set {
					if j > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(f.exprString(x))
				}
				sb.WriteString(")")
			}
		} else {
			for i, x := range v.Exprs {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(f.exprString(x))
			}
		}
		sb.WriteString(")")
		return sb.String()
	}
	return sprintf("%v", e)
}
func (f *formatter) literalString(v *ast.Literal) string {
	switch v.LitType {
	case ast.LitNull:
		return "NULL"
	case ast.LitBool:
		if v.Value == "true" {
			return "TRUE"
		}
		return "FALSE"
	case ast.LitString:
		val := v.Value
		if i := strings.IndexByte(val, '\x00'); i >= 0 {

		}
		prefix := ""
		if v.Charset != "" {
			prefix = "_" + v.Charset
		}
		return prefix + f.stringLiteral(val)
	case ast.LitHex, ast.LitBit:
		if v.Text != "" {
			return v.Text
		}
		return v.Value
	default:
		return v.Value
	}
}
func (f *formatter) variableString(v *ast.Variable) string {
	if v.IsSystem {
		if v.Scope != "" {
			return "@@" + v.Scope + "." + v.Name
		}
		return "@@" + v.Name
	}
	return "@" + v.Name
}
func (f *formatter) unaryString(v *ast.UnaryExpr) string {
	if isAlpha(v.Op) {
		return v.Op + " " + f.exprString(v.Operand)
	}
	return v.Op + f.exprString(v.Operand)
}
func (f *formatter) binaryString(v *ast.BinaryExpr) string {
	return f.exprString(v.Left) + " " + v.Op + " " + f.exprString(v.Right)
}
func (f *formatter) funcString(v *ast.FuncCall) string {
	var sb strings.Builder
	sb.WriteString(stringsLower(v.Name))
	sb.WriteString("(")
	if v.Distinct {
		sb.WriteString("DISTINCT ")
	}
	if v.Star {
		sb.WriteString("*")
	} else {
		for i, a := range v.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.exprString(a))
		}
	}
	sb.WriteString(")")
	if v.Over != nil {
		sb.WriteString(f.windowSpecString(v.Over))
	}
	return sb.String()
}
func (f *formatter) windowSpecString(w *ast.WindowSpec) string {
	var sb strings.Builder
	sb.WriteString(" OVER(")
	if w.Name != "" {
		sb.WriteString(w.Name)
	}
	if len(w.PartitionBy) > 0 {
		sb.WriteString(" PARTITION BY ")
		for i, e := range w.PartitionBy {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.exprString(e))
		}
	}
	if len(w.OrderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		for i, o := range w.OrderBy {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.exprString(o.Expr))
			if o.Desc {
				sb.WriteString(" DESC")
			}
		}
	}
	if w.Frame != "" {
		sb.WriteString(" ")
		sb.WriteString(w.Frame)
		sb.WriteString(" ")
		if w.FrameEnd != "" {
			sb.WriteString("BETWEEN ")
			sb.WriteString(w.FrameStart)
			sb.WriteString(" AND ")
			sb.WriteString(w.FrameEnd)
		} else {
			sb.WriteString(w.FrameStart)
		}
	}
	sb.WriteString(")")
	return sb.String()
}
func (f *formatter) castString(v *ast.CastExpr) string {
	if v.IsConvert {
		return "CONVERT(" + f.exprString(v.Expr) + ", " + f.dataTypeString(v) + ")"
	}
	return "CAST(" + f.exprString(v.Expr) + " AS " + f.dataTypeString(v) + ")"
}
func (f *formatter) dataTypeString(v *ast.CastExpr) string {
	var sb strings.Builder
	sb.WriteString(v.TypeName)
	if len(v.TypeArgs) > 0 {
		sb.WriteString("(")
		for i, a := range v.TypeArgs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.exprString(a))
		}
		sb.WriteString(")")
	}
	return sb.String()
}
func (f *formatter) caseString(v *ast.CaseExpr) string {
	var sb strings.Builder
	sb.WriteString("CASE")
	if v.Operand != nil {
		sb.WriteString(" ")
		sb.WriteString(f.exprString(v.Operand))
	}
	for _, c := range v.Whens {
		sb.WriteString(" WHEN ")
		sb.WriteString(f.exprString(c.When))
		sb.WriteString(" THEN ")
		sb.WriteString(f.exprString(c.Then))
	}
	if v.Else != nil {
		sb.WriteString(" ELSE ")
		sb.WriteString(f.exprString(v.Else))
	}
	sb.WriteString(" END")
	return sb.String()
}
func (f *formatter) betweenString(v *ast.BetweenExpr) string {
	n := " "
	if v.Not {
		n = " NOT "
	}
	return f.exprString(v.Expr) + n + "BETWEEN " + f.exprString(v.Low) + " AND " + f.exprString(v.High)
}
func (f *formatter) inString(v *ast.InExpr) string {
	n := " "
	if v.Not {
		n = " NOT "
	}
	var sb strings.Builder
	sb.WriteString(f.exprString(v.Expr))
	sb.WriteString(n)
	sb.WriteString("IN (")
	if v.Subquery != nil {
		sb.WriteString(Format(v.Subquery))
	} else {
		for i, e := range v.List {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.exprString(e))
		}
	}
	sb.WriteString(")")
	return sb.String()
}
func (f *formatter) likeString(v *ast.LikeExpr) string {
	n := " "
	if v.Not {
		n = " NOT "
	}
	var sb strings.Builder
	sb.WriteString(f.exprString(v.Expr))
	sb.WriteString(n)
	sb.WriteString("LIKE ")
	sb.WriteString(f.exprString(v.Pattern))
	if v.Escape != nil {
		sb.WriteString(" ESCAPE ")
		sb.WriteString(f.exprString(v.Escape))
	}
	return sb.String()
}
func (f *formatter) isString(v *ast.IsExpr) string {
	n := " IS "
	if v.Not {
		n = " IS NOT "
	}
	return f.exprString(v.Expr) + n + stringsUpper(v.Truth)
}
func (f *formatter) intervalString(v *ast.IntervalExpr) string {
	return "INTERVAL " + f.exprString(v.Value) + " " + stringsUpper(v.Unit)
}
