package ast_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/token"
)

// isNilNode reports whether n is nil or a typed-nil pointer.
func isNilNode(n ast.Node) bool {
	if n == nil {
		return true
	}
	rv := reflect.ValueOf(n)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}

// factories covers every concrete node type in the package: each instance is
// exercised through Finalize / Children / NodeKind / JSON marshal so that
// constructor-less "shape" nodes are smoke-tested too.
var factories = map[string]func() ast.Node{
	"AlterDatabaseStmt":       func() ast.Node { return &ast.AlterDatabaseStmt{} },
	"AlterLSStmt":             func() ast.Node { return &ast.AlterLSStmt{} },
	"AlterMViewStmt":          func() ast.Node { return &ast.AlterMViewStmt{} },
	"AlterOutlineStmt":        func() ast.Node { return &ast.AlterOutlineStmt{} },
	"AlterResourcePoolStmt":   func() ast.Node { return &ast.AlterResourcePoolStmt{} },
	"AlterResourceUnitStmt":   func() ast.Node { return &ast.AlterResourceUnitStmt{} },
	"AlterSequenceStmt":       func() ast.Node { return &ast.AlterSequenceStmt{} },
	"AlterSystemStmt":         func() ast.Node { return &ast.AlterSystemStmt{} },
	"AlterTablegroupStmt":     func() ast.Node { return &ast.AlterTablegroupStmt{} },
	"AlterTablespaceStmt":     func() ast.Node { return &ast.AlterTablespaceStmt{} },
	"AlterTableStmt":          func() ast.Node { return &ast.AlterTableStmt{} },
	"AlterTenantStmt":         func() ast.Node { return &ast.AlterTenantStmt{} },
	"AlterViewStmt":           func() ast.Node { return &ast.AlterViewStmt{} },
	"AuditStmt":               func() ast.Node { return &ast.AuditStmt{} },
	"BeginStmt":               func() ast.Node { return &ast.BeginStmt{} },
	"BetweenExpr":             func() ast.Node { return &ast.BetweenExpr{} },
	"BinaryExpr":              func() ast.Node { return &ast.BinaryExpr{} },
	"CallStmt":                func() ast.Node { return &ast.CallStmt{} },
	"CaseExpr":                func() ast.Node { return &ast.CaseExpr{} },
	"CastExpr":                func() ast.Node { return &ast.CastExpr{} },
	"CatalogStmt":             func() ast.Node { return &ast.CatalogStmt{} },
	"CCLLimitRuleStmt":        func() ast.Node { return &ast.CCLLimitRuleStmt{} },
	"CloneTenantStmt":         func() ast.Node { return &ast.CloneTenantStmt{} },
	"CommitStmt":              func() ast.Node { return &ast.CommitStmt{} },
	"CreateDatabaseStmt":      func() ast.Node { return &ast.CreateDatabaseStmt{} },
	"CreateIndexStmt":         func() ast.Node { return &ast.CreateIndexStmt{} },
	"CreateLocationStmt":      func() ast.Node { return &ast.CreateLocationStmt{} },
	"CreateMViewStmt":         func() ast.Node { return &ast.CreateMViewStmt{} },
	"CreateOutlineStmt":       func() ast.Node { return &ast.CreateOutlineStmt{} },
	"CreateResourcePoolStmt":  func() ast.Node { return &ast.CreateResourcePoolStmt{} },
	"CreateResourceUnitStmt":  func() ast.Node { return &ast.CreateResourceUnitStmt{} },
	"CreateRoleStmt":          func() ast.Node { return &ast.CreateRoleStmt{} },
	"CreateSequenceStmt":      func() ast.Node { return &ast.CreateSequenceStmt{} },
	"CreateStandbyTenantStmt": func() ast.Node { return &ast.CreateStandbyTenantStmt{} },
	"CreateTablegroupStmt":    func() ast.Node { return &ast.CreateTablegroupStmt{} },
	"CreateTableLikeStmt":     func() ast.Node { return &ast.CreateTableLikeStmt{} },
	"CreateTablespaceStmt":    func() ast.Node { return &ast.CreateTablespaceStmt{} },
	"CreateTableStmt":         func() ast.Node { return &ast.CreateTableStmt{} },
	"CreateTenantStmt":        func() ast.Node { return &ast.CreateTenantStmt{} },
	"CreateUserStmt":          func() ast.Node { return &ast.CreateUserStmt{} },
	"CreateViewStmt":          func() ast.Node { return &ast.CreateViewStmt{} },
	"DblinkStmt":              func() ast.Node { return &ast.DblinkStmt{} },
	"DeallocateStmt":          func() ast.Node { return &ast.DeallocateStmt{} },
	"DeleteStmt":              func() ast.Node { return &ast.DeleteStmt{} },
	"DoStmt":                  func() ast.Node { return &ast.DoStmt{} },
	"DropDatabaseStmt":        func() ast.Node { return &ast.DropDatabaseStmt{} },
	"DropIndexStmt":           func() ast.Node { return &ast.DropIndexStmt{} },
	"DropMViewStmt":           func() ast.Node { return &ast.DropMViewStmt{} },
	"DropOutlineStmt":         func() ast.Node { return &ast.DropOutlineStmt{} },
	"DropResourcePoolStmt":    func() ast.Node { return &ast.DropResourcePoolStmt{} },
	"DropResourceUnitStmt":    func() ast.Node { return &ast.DropResourceUnitStmt{} },
	"DropRoleStmt":            func() ast.Node { return &ast.DropRoleStmt{} },
	"DropRoutineStmt":         func() ast.Node { return &ast.DropRoutineStmt{} },
	"DropSequenceStmt":        func() ast.Node { return &ast.DropSequenceStmt{} },
	"DropTablegroupStmt":      func() ast.Node { return &ast.DropTablegroupStmt{} },
	"DropTablespaceStmt":      func() ast.Node { return &ast.DropTablespaceStmt{} },
	"DropTableStmt":           func() ast.Node { return &ast.DropTableStmt{} },
	"DropTenantStmt":          func() ast.Node { return &ast.DropTenantStmt{} },
	"DropUserStmt":            func() ast.Node { return &ast.DropUserStmt{} },
	"DropViewStmt":            func() ast.Node { return &ast.DropViewStmt{} },
	"DumpMemoryStmt":          func() ast.Node { return &ast.DumpMemoryStmt{} },
	"ExecuteStmt":             func() ast.Node { return &ast.ExecuteStmt{} },
	"ExistsExpr":              func() ast.Node { return &ast.ExistsExpr{} },
	"ExplainStmt":             func() ast.Node { return &ast.ExplainStmt{} },
	"FlashbackStandbyLogStmt": func() ast.Node { return &ast.FlashbackStandbyLogStmt{} },
	"FlashbackTableStmt":      func() ast.Node { return &ast.FlashbackTableStmt{} },
	"FlushStmt":               func() ast.Node { return &ast.FlushStmt{} },
	"FuncCall":                func() ast.Node { return &ast.FuncCall{} },
	"GetDiagnosticsStmt":      func() ast.Node { return &ast.GetDiagnosticsStmt{} },
	"GrantStmt":               func() ast.Node { return &ast.GrantStmt{} },
	"GroupingSetExpr":         func() ast.Node { return &ast.GroupingSetExpr{} },
	"HandlerStmt":             func() ast.Node { return &ast.HandlerStmt{} },
	"HelpStmt":                func() ast.Node { return &ast.HelpStmt{} },
	"Identifier":              func() ast.Node { return &ast.Identifier{} },
	"InExpr":                  func() ast.Node { return &ast.InExpr{} },
	"InsertStmt":              func() ast.Node { return &ast.InsertStmt{} },
	"InstallPluginStmt":       func() ast.Node { return &ast.InstallPluginStmt{} },
	"IntervalExpr":            func() ast.Node { return &ast.IntervalExpr{} },
	"IsExpr":                  func() ast.Node { return &ast.IsExpr{} },
	"KillStmt":                func() ast.Node { return &ast.KillStmt{} },
	"LikeExpr":                func() ast.Node { return &ast.LikeExpr{} },
	"Literal":                 func() ast.Node { return &ast.Literal{} },
	"LoadDataStmt":            func() ast.Node { return &ast.LoadDataStmt{} },
	"LocationUtilsStmt":       func() ast.Node { return &ast.LocationUtilsStmt{} },
	"LockTablesStmt":          func() ast.Node { return &ast.LockTablesStmt{} },
	"LockUserStmt":            func() ast.Node { return &ast.LockUserStmt{} },
	"MatchAgainstExpr":        func() ast.Node { return &ast.MatchAgainstExpr{} },
	"MemberOfExpr":            func() ast.Node { return &ast.MemberOfExpr{} },
	"MergeStmt":               func() ast.Node { return &ast.MergeStmt{} },
	"MLogStmt":                func() ast.Node { return &ast.MLogStmt{} },
	"ParamMarker":             func() ast.Node { return &ast.ParamMarker{} },
	"ParenExpr":               func() ast.Node { return &ast.ParenExpr{} },
	"PLBlock":                 func() ast.Node { return &ast.PLBlock{} },
	"PLCaseStmt":              func() ast.Node { return &ast.PLCaseStmt{} },
	"PLConditionDeclStmt":     func() ast.Node { return &ast.PLConditionDeclStmt{} },
	"PLConditionHandlerStmt":  func() ast.Node { return &ast.PLConditionHandlerStmt{} },
	"PLCursorDeclStmt":        func() ast.Node { return &ast.PLCursorDeclStmt{} },
	"PLCursorStmt":            func() ast.Node { return &ast.PLCursorStmt{} },
	"PLDeclareStmt":           func() ast.Node { return &ast.PLDeclareStmt{} },
	"PLIfStmt":                func() ast.Node { return &ast.PLIfStmt{} },
	"PLIterateStmt":           func() ast.Node { return &ast.PLIterateStmt{} },
	"PLLeaveStmt":             func() ast.Node { return &ast.PLLeaveStmt{} },
	"PLLoopStmt":              func() ast.Node { return &ast.PLLoopStmt{} },
	"PLRepeatStmt":            func() ast.Node { return &ast.PLRepeatStmt{} },
	"PLReturnStmt":            func() ast.Node { return &ast.PLReturnStmt{} },
	"PLSetStmt":               func() ast.Node { return &ast.PLSetStmt{} },
	"PLWhileStmt":             func() ast.Node { return &ast.PLWhileStmt{} },
	"PrepareStmt":             func() ast.Node { return &ast.PrepareStmt{} },
	"PurgeStmt":               func() ast.Node { return &ast.PurgeStmt{} },
	"RecoverTenantStmt":       func() ast.Node { return &ast.RecoverTenantStmt{} },
	"ReleaseSavepointStmt":    func() ast.Node { return &ast.ReleaseSavepointStmt{} },
	"RenameTableStmt":         func() ast.Node { return &ast.RenameTableStmt{} },
	"RenameUserStmt":          func() ast.Node { return &ast.RenameUserStmt{} },
	"RestorePointStmt":        func() ast.Node { return &ast.RestorePointStmt{} },
	"RevokeStmt":              func() ast.Node { return &ast.RevokeStmt{} },
	"RollbackStmt":            func() ast.Node { return &ast.RollbackStmt{} },
	"RollbackToSavepointStmt": func() ast.Node { return &ast.RollbackToSavepointStmt{} },
	"RotateMasterKeyStmt":     func() ast.Node { return &ast.RotateMasterKeyStmt{} },
	"RoutineLoadActionStmt":   func() ast.Node { return &ast.RoutineLoadActionStmt{} },
	"RoutineLoadStmt":         func() ast.Node { return &ast.RoutineLoadStmt{} },
	"SavepointStmt":           func() ast.Node { return &ast.SavepointStmt{} },
	"SelectBody":              func() ast.Node { return &ast.SelectBody{} },
	"SelectStmt":              func() ast.Node { return &ast.SelectStmt{} },
	"SensitiveRuleStmt":       func() ast.Node { return &ast.SensitiveRuleStmt{} },
	"ServiceStmt":             func() ast.Node { return &ast.ServiceStmt{} },
	"SetPasswordStmt":         func() ast.Node { return &ast.SetPasswordStmt{} },
	"SetRoleStmt":             func() ast.Node { return &ast.SetRoleStmt{} },
	"SetStmt":                 func() ast.Node { return &ast.SetStmt{} },
	"SetTransactionStmt":      func() ast.Node { return &ast.SetTransactionStmt{} },
	"ShowStmt":                func() ast.Node { return &ast.ShowStmt{} },
	"SignalStmt":              func() ast.Node { return &ast.SignalStmt{} },
	"StarLiteral":             func() ast.Node { return &ast.StarLiteral{} },
	"StoredProcStmt":          func() ast.Node { return &ast.StoredProcStmt{} },
	"SubqueryExpr":            func() ast.Node { return &ast.SubqueryExpr{} },
	"SwitchoverTenantStmt":    func() ast.Node { return &ast.SwitchoverTenantStmt{} },
	"TableMaintStmt":          func() ast.Node { return &ast.TableMaintStmt{} },
	"TableRef":                func() ast.Node { return &ast.TableRef{} },
	"TenantSnapshotStmt":      func() ast.Node { return &ast.TenantSnapshotStmt{} },
	"TransferPartitionStmt":   func() ast.Node { return &ast.TransferPartitionStmt{} },
	"TruncateStmt":            func() ast.Node { return &ast.TruncateStmt{} },
	"UnaryExpr":               func() ast.Node { return &ast.UnaryExpr{} },
	"UninstallPluginStmt":     func() ast.Node { return &ast.UninstallPluginStmt{} },
	"UnlockTablesStmt":        func() ast.Node { return &ast.UnlockTablesStmt{} },
	"UpdateStmt":              func() ast.Node { return &ast.UpdateStmt{} },
	"UseDbStmt":               func() ast.Node { return &ast.UseDbStmt{} },
	"Variable":                func() ast.Node { return &ast.Variable{} },
	"XAStmt":                  func() ast.Node { return &ast.XAStmt{} },
}

func TestAllNodeKindsShape(t *testing.T) {
	for name, mk := range factories {
		n := mk()
		if n == nil {
			t.Fatalf("%s: nil node", name)
		}
		ast.Finalize(n)
		// Finalize stamps Base.Kind for struct nodes that embed Base.
		if b, ok := n.(interface{ NodeKind() ast.NodeKind }); ok {
			if len(b.NodeKind()) == 0 {
				t.Errorf("%s: NodeKind empty after Finalize", name)
			}
		}
		// Children must be safe: no panics; nil / typed-nil entries are the
		// absent-field convention (visitor.Walk skips plain nils), so only
		// truly non-nil ones are checked.
		kids := n.Children()
		for _, k := range kids {
			if isNilNode(k) {
				continue
			}
			if k.Span().Start.Line < 0 {
				t.Fatalf("%s: invalid child span", name)
			}
		}
		// encoding/json smoke: default marshal must not panic.
		buf, err := json.Marshal(n)
		if err != nil {
			t.Errorf("%s: json.Marshal: %v", name, err)
			continue
		}
		if len(buf) == 0 {
			t.Errorf("%s: empty json output", name)
		}
	}
}

func TestNodeKindStable(t *testing.T) {
	sel := &ast.SelectStmt{}
	ast.Finalize(sel)
	if got := sel.NodeKind(); got != ast.NodeKind("SelectStmt") {
		t.Fatalf("kind = %q, want SelectStmt", got)
	}
}

func TestFinalizeSubtree(t *testing.T) {
	inner := &ast.Literal{Text: "1"}
	expr := &ast.BinaryExpr{Op: "+", Left: inner, Right: inner}
	ast.Finalize(expr)
	// Finalize stamps every reachable Base via reflection; the same literal
	// must be reachable and stamped (no infinite loop, no nil deref).
	if string(expr.NodeKind()) != "BinaryExpr" {
		t.Fatalf("kind = %q", expr.NodeKind())
	}
}

func TestSpanRoundTrip(t *testing.T) {
	expr := &ast.BinaryExpr{Op: "+", Left: &ast.Literal{Text: "1"}, Right: &ast.Literal{Text: "2"}}
	expr.SetSpan(token.Position{Offset: 0, Line: 1, Col: 1}, token.Position{Offset: 5, Line: 1, Col: 6})
	s := expr.Span()
	if !s.Start.IsValid() || s.Start.Line != 1 || s.Start.Col != 1 || s.End.Col != 6 {
		t.Fatalf("span not preserved: %+v", s)
	}
}
