package formatter

import (
	"strings"

	"github.com/zetatez/ob-sql-parser/ast"
)

func (f *formatter) plBody(b *ast.PLBlock) {
	for _, s := range b.Statements {
		f.plStmt(s)
		f.raw("; ")
	}
}
func (f *formatter) plBlock(b *ast.PLBlock) {
	if b == nil {
		return
	}
	if b.Label != "" {
		f.printf("%s: BEGIN ", b.Label)
		f.plBody(b)
		f.raw("END")
		return
	}
	f.raw("BEGIN ")
	f.plBody(b)
	f.raw("END")
}
func (f *formatter) plStmt(n ast.Node) {
	switch v := n.(type) {
	case *ast.PLBlock:
		f.plBlock(v)
	case *ast.PLSetStmt:
		f.printf("SET %s = ", v.Name)
		f.expr(v.Expr)
	case *ast.PLDeclareStmt:
		f.printf("DECLARE %s %s", strings.Join(v.Names, ", "), v.Type)
		if len(v.TypeArgs) > 0 {
			f.raw(" (")
			for i, a := range v.TypeArgs {
				if i > 0 {
					f.raw(", ")
				}
				f.raw(f.exprString(a))
			}
			f.raw(")")
		}
		if v.Default != nil {
			f.raw(" DEFAULT ")
			f.expr(v.Default)
		}
	case *ast.PLIfStmt:
		for i, br := range v.Branches {
			if i == 0 {
				f.raw("IF ")
			} else {
				f.raw("ELSEIF ")
			}
			f.expr(br.Cond)
			f.raw(" THEN ")
			f.plBody(br.Body)
		}
		if v.Else != nil {
			f.raw("ELSE ")
			f.plBody(v.Else)
		}
		f.raw("END IF")
	case *ast.PLWhileStmt:
		if v.Label != "" {
			f.printf("%s: ", v.Label)
		}
		f.raw("WHILE ")
		f.expr(v.Cond)
		f.raw(" DO ")
		f.plBody(v.Body)
		f.raw("END WHILE")
	case *ast.PLRepeatStmt:
		if v.Label != "" {
			f.printf("%s: ", v.Label)
		}
		f.raw("REPEAT ")
		f.plBody(v.Body)
		f.raw("UNTIL ")
		f.expr(v.Until)
		f.raw(" END REPEAT")
	case *ast.PLLoopStmt:
		if v.Label != "" {
			f.printf("%s: ", v.Label)
		}
		f.raw("LOOP ")
		f.plBody(v.Body)
		f.raw("END LOOP")
	case *ast.PLReturnStmt:
		f.raw("RETURN ")
		f.expr(v.Expr)
	case *ast.PLCursorDeclStmt:
		f.printf("DECLARE %s CURSOR FOR ", v.Name)
		f.selectStmt(v.Select)
	case *ast.PLCursorStmt:
		f.printf("%s %s", v.Op, v.Name)
		if v.Op == "FETCH" && len(v.Into) > 0 {
			f.printf(" INTO %s", strings.Join(v.Into, ", "))
		}
	case *ast.PLLeaveStmt:
		f.printf("LEAVE %s", v.Name)
	case *ast.PLIterateStmt:
		f.printf("ITERATE %s", v.Name)
	case *ast.PLConditionDeclStmt:
		f.printf("DECLARE %s CONDITION FOR ", v.Name)
		if v.Kind == "sqlstate" {
			f.printf("SQLSTATE %s", f.stringLiteral(v.Value))
		} else {
			f.raw(v.Value)
		}
	case *ast.PLCaseStmt:
		f.raw("CASE")
		if v.Operand != nil {
			f.raw(" ")
			f.expr(v.Operand)
		}
		for _, br := range v.Branches {
			f.raw(" WHEN ")
			f.expr(br.When)
			f.raw(" THEN ")
			f.plBody(br.Body)
		}
		if v.Else != nil {
			f.raw(" ELSE ")
			f.plBody(v.Else)
		}
		f.raw("END CASE")
	case *ast.PLConditionHandlerStmt:
		f.printf("DECLARE %s HANDLER FOR ", v.Action)
		for i, c := range v.Conditions {
			if i > 0 {
				f.raw(", ")
			}
			switch c.Kind {
			case "sqlstate":
				f.printf("SQLSTATE %s", f.stringLiteral(c.Value))
			case "notfound":
				f.raw("NOT FOUND")
			default:
				f.raw(c.Value)
			}
		}
		f.raw(" ")
		f.plStmt(v.Body)
	case ast.Statement:
		f.stmt(v)
	default:
		f.nodeAny(n)
	}
}
