package ast

// This file implements Children() for every concrete node. Children() powers
// the generic visitor walker. It reconstructs children from typed fields on
// demand so there is no duplicated [children] linkage.

func (n *Literal) Children() []Node     { return nil }
func (n *Identifier) Children() []Node  { return nil }
func (n *StarLiteral) Children() []Node { return nil }
func (n *ParamMarker) Children() []Node { return nil }
func (n *Variable) Children() []Node    { return nil }

func (n *UnaryExpr) Children() []Node { return childrenOf(n.Operand) }

func (n *BinaryExpr) Children() []Node { return childrenOf(n.Left, n.Right) }

func (n *FuncCall) Children() []Node {
	out := enode(n.Args)
	if n.Over != nil {
		out = append(out, windowSpecChildren(n.Over)...)
	}
	return out
}

func (n *CastExpr) Children() []Node { return childrenOf(n.Expr) }

func (n *CaseExpr) Children() []Node {
	out := childrenOf(n.Operand)
	for _, w := range n.Whens {
		out = append(out, w.When, w.Then)
	}
	out = append(out, n.Else)
	return out
}

func (n *BetweenExpr) Children() []Node { return childrenOf(n.Expr, n.Low, n.High) }

func (n *InExpr) Children() []Node {
	out := childrenOf(n.Expr)
	out = append(out, enode(n.List)...)
	if n.Subquery != nil {
		out = append(out, n.Subquery)
	}
	return out
}

func (n *LikeExpr) Children() []Node { return childrenOf(n.Expr, n.Pattern, n.Escape) }

func (n *IsExpr) Children() []Node { return childrenOf(n.Expr) }

func (n *IntervalExpr) Children() []Node { return childrenOf(n.Value) }

func (n *ExistsExpr) Children() []Node { return childrenOf(n.Query) }

func (n *MemberOfExpr) Children() []Node {
	out := childrenOf(n.Expr)
	return append(out, enode(n.List)...)
}

func (n *MatchAgainstExpr) Children() []Node { return childrenOf(n.Against) }

func (n *GroupingSetExpr) Children() []Node {
	var out []Node
	for _, s := range n.Sets {
		out = append(out, enode(s)...)
	}
	out = append(out, enode(n.Exprs)...)
	return out
}

func (n *SubqueryExpr) Children() []Node { return childrenOf(n.Query) }
func (n *ParenExpr) Children() []Node    { return childrenOf(n.Expr) }

func (n *CreateDatabaseStmt) Children() []Node      { return nil }
func (n *DropDatabaseStmt) Children() []Node        { return nil }
func (n *AlterDatabaseStmt) Children() []Node       { return nil }
func (n *AlterSystemStmt) Children() []Node         { return childrenOf(n.Value) }
func (n *TransferPartitionStmt) Children() []Node   { return childrenOf(n.TableID, n.ObjectID, n.LS) }
func (n *CreateViewStmt) Children() []Node          { return childrenOf(n.Select) }
func (n *DropViewStmt) Children() []Node            { return nil }
func (n *GrantStmt) Children() []Node               { return nil }
func (n *RevokeStmt) Children() []Node              { return nil }
func (n *CreateUserStmt) Children() []Node          { return nil }
func (n *DropUserStmt) Children() []Node            { return nil }
func (n *PrepareStmt) Children() []Node             { return childrenOf(n.Source) }
func (n *DeallocateStmt) Children() []Node          { return nil }
func (n *ExecuteStmt) Children() []Node             { return nil }
func (n *CreateRoleStmt) Children() []Node          { return nil }
func (n *DropRoleStmt) Children() []Node            { return nil }
func (n *SetRoleStmt) Children() []Node             { return nil }
func (n *SetPasswordStmt) Children() []Node         { return nil }
func (n *KillStmt) Children() []Node                { return childrenOf(n.Target) }
func (n *HelpStmt) Children() []Node                { return childrenOf(n.Target) }
func (n *TableMaintStmt) Children() []Node          { return nil }
func (n *RenameTableStmt) Children() []Node         { return nil }
func (n *PurgeStmt) Children() []Node               { return nil }
func (n *FlashbackTableStmt) Children() []Node      { return nil }
func (n *CreateMViewStmt) Children() []Node         { return childrenOf(n.Select) }
func (n *AlterViewStmt) Children() []Node           { return childrenOf(n.Select) }
func (n *AlterMViewStmt) Children() []Node          { return nil }
func (n *SwitchoverTenantStmt) Children() []Node    { return nil }
func (n *AuditStmt) Children() []Node               { return nil }
func (n *RoutineLoadStmt) Children() []Node         { return optionExprs(n.LoadProps) }
func (n *RoutineLoadActionStmt) Children() []Node   { return nil }
func (n *CCLLimitRuleStmt) Children() []Node        { return optionExprs(n.Options) }
func (n *SensitiveRuleStmt) Children() []Node       { return nil }
func (n *DblinkStmt) Children() []Node              { return childrenOf(n.Port) }
func (n *ServiceStmt) Children() []Node             { return nil }
func (n *InstallPluginStmt) Children() []Node       { return nil }
func (n *UninstallPluginStmt) Children() []Node     { return nil }
func (n *HandlerStmt) Children() []Node             { return childrenOf(n.Where) }
func (n *CreateLocationStmt) Children() []Node      { return nil }
func (n *CatalogStmt) Children() []Node             { return optionExprs(n.Options) }
func (n *LocationUtilsStmt) Children() []Node       { return nil }
func (n *DumpMemoryStmt) Children() []Node          { return childrenOf(n.FirstVal, n.SecondVal) }
func (n *FlashbackStandbyLogStmt) Children() []Node { return childrenOf(n.SCN) }
func (n *AlterLSStmt) Children() []Node             { return optionExprs(n.Options) }
func (n *CreateStandbyTenantStmt) Children() []Node { return optionExprs(n.Options) }
func (n *CloneTenantStmt) Children() []Node         { return optionExprs(n.Options) }
func (n *RestorePointStmt) Children() []Node        { return nil }
func (n *TenantSnapshotStmt) Children() []Node      { return nil }
func (n *RecoverTenantStmt) Children() []Node       { return childrenOf(n.UntilValue) }
func (n *MLogStmt) Children() []Node                { return nil }
func (n *CreateOutlineStmt) Children() []Node       { return childrenOf(n.Select) }
func (n *AlterOutlineStmt) Children() []Node        { return childrenOf(n.Select) }
func (n *DropOutlineStmt) Children() []Node         { return nil }
func (n *AlterTablespaceStmt) Children() []Node     { return optionExprs(n.Options) }
func (n *RotateMasterKeyStmt) Children() []Node     { return nil }
func (n *RenameUserStmt) Children() []Node          { return nil }
func (n *LockUserStmt) Children() []Node            { return nil }

func (n *AlterSequenceStmt) Children() []Node {
	return childrenOf(n.StartWith, n.IncrementBy, n.MinValue, n.MaxValue, n.Cache)
}
func (n *AlterTablegroupStmt) Children() []Node   { return optionExprs(n.Options) }
func (n *AlterResourceUnitStmt) Children() []Node { return optionExprs(n.Options) }
func (n *AlterResourcePoolStmt) Children() []Node { return optionExprs(n.Options) }
func (n *GetDiagnosticsStmt) Children() []Node    { return childrenOf(n.Condition) }

func (n *CreateTenantStmt) Children() []Node       { return optionExprs(n.Options) }
func (n *AlterTenantStmt) Children() []Node        { return optionExprs(n.Options) }
func (n *CreateResourceUnitStmt) Children() []Node { return optionExprs(n.Options) }
func (n *CreateResourcePoolStmt) Children() []Node { return optionExprs(n.Options) }
func (n *CreateTablegroupStmt) Children() []Node   { return optionExprs(n.Options) }
func (n *CreateTablespaceStmt) Children() []Node   { return optionExprs(n.Options) }
func (n *DropTenantStmt) Children() []Node         { return nil }
func (n *DropResourceUnitStmt) Children() []Node   { return nil }
func (n *DropResourcePoolStmt) Children() []Node   { return nil }
func (n *DropTablegroupStmt) Children() []Node     { return nil }
func (n *DropTablespaceStmt) Children() []Node     { return nil }

func optionExprs(opts []Option) []Node {
	var out []Node
	for i := range opts {
		out = append(out, opts[i].Value)
		out = append(out, enode(opts[i].Values)...)
	}
	return out
}

func (n *SignalStmt) Children() []Node {
	var out []Node
	for i := range n.Items {
		out = append(out, n.Items[i].Value)
	}
	return out
}
func (n *DropMViewStmt) Children() []Node { return nil }

func (n *MergeStmt) Children() []Node {
	var out []Node
	out = append(out, n.Target, n.Source, n.On)
	for i := range n.Whens {
		w := &n.Whens[i]
		for j := range w.Set {
			out = append(out, w.Set[j].Expr)
		}
		out = append(out, w.DeleteCond)
		out = append(out, enode(w.Values)...)
	}
	return out
}

func (n *XAStmt) Children() []Node {
	var out []Node
	for _, e := range n.XID {
		out = append(out, e)
	}
	return out
}
func (n *LoadDataStmt) Children() []Node     { return childrenOf(n.IgnoreLines) }
func (n *FlushStmt) Children() []Node        { return nil }
func (n *DropSequenceStmt) Children() []Node { return nil }

func (n *CreateSequenceStmt) Children() []Node {
	return childrenOf(n.StartWith, n.IncrementBy, n.MinValue, n.MaxValue, n.Cache)
}

func (n *TableRef) Children() []Node {
	switch n.RefKind {
	case RefSubquery:
		return childrenOf(n.Subquery)
	case RefJoin:
		return childrenOf(n.Left, n.Right, n.On)
	default:
		return nil
	}
}

func (n *SelectBody) Children() []Node {
	var out []Node
	for i := range n.Items {
		out = append(out, expressionOfItem(&n.Items[i]))
	}
	out = append(out, inode(n.From)...)
	out = append(out, n.Where)
	out = append(out, enode(n.GroupBy)...)
	out = append(out, n.Having)
	out = append(out, n.Qualify)
	out = append(out, enode(n.ProcAnalyse)...)
	for i := range n.OrderBy {
		out = append(out, n.OrderBy[i].Expr)
	}
	if n.Limit != nil {
		out = append(out, n.Limit.Count, n.Limit.Offset)
	}
	return out
}

func (n *SelectStmt) Children() []Node {
	var out []Node
	for _, b := range n.Bodies {
		out = append(out, b)
	}
	for _, s := range n.SetOps {
		if s != nil {
			out = append(out, s.Right)
		}
	}
	for i := range n.OrderBy {
		out = append(out, n.OrderBy[i].Expr)
	}
	if n.Limit != nil {
		out = append(out, n.Limit.Count, n.Limit.Offset)
	}
	return out
}

func (n *InsertStmt) Children() []Node {
	var out []Node
	for _, r := range n.Rows {
		out = append(out, enode(r)...)
	}
	out = append(out, n.Source)
	for i := range n.Set {
		out = append(out, n.Set[i].Expr)
	}
	for i := range n.OnDuplicate {
		out = append(out, n.OnDuplicate[i].Expr)
	}
	return out
}

func (n *UpdateStmt) Children() []Node {
	var out []Node
	out = append(out, inode(n.Tables)...)
	for i := range n.Set {
		out = append(out, n.Set[i].Expr)
	}
	out = append(out, n.Where)
	out = append(out, idNodes(n.Returning)...)
	return out
}

func (n *DeleteStmt) Children() []Node {
	var out []Node
	out = append(out, idNodes(n.Tables)...)
	out = append(out, inode(n.From)...)
	out = append(out, inode(n.Using)...)
	out = append(out, n.Where)
	out = append(out, idNodes(n.Returning)...)
	return out
}

// idNodes converts []*Identifier to []Node (nil-safe).
func idNodes(ids []*Identifier) []Node {
	out := make([]Node, 0, len(ids))
	for _, id := range ids {
		if id != nil {
			out = append(out, id)
		}
	}
	return out
}

func (n *CreateTableStmt) Children() []Node {
	var out []Node
	for _, c := range n.Columns {
		out = append(out, childrenOf(c.Default, c.Generated, c.OnUpdate)...)
	}
	for _, k := range n.Constraints {
		out = append(out, childrenOf(k.Check)...)
		out = append(out, idNodes([]*Identifier{k.RefTable})...)
		out = append(out, idNodes(k.RefColumns)...)
	}
	if n.Select != nil {
		out = append(out, n.Select)
	}
	if n.Partition != nil {
		out = append(out, childrenOf(n.Partition.Expr, n.Partition.Count)...)
		for _, pd := range n.Partition.Partitions {
			out = append(out, enode(pd.Values)...)
		}
	}
	return out
}

func (n *CreateTableLikeStmt) Children() []Node { return nil }
func (n *CreateIndexStmt) Children() []Node     { return nil }
func (n *DropTableStmt) Children() []Node       { return nil }
func (n *DropIndexStmt) Children() []Node       { return nil }

func (n *AlterTableStmt) Children() []Node {
	var out []Node
	for i := range n.Actions {
		a := &n.Actions[i]
		if a.Column != nil {
			out = append(out, childrenOf(a.Column.Default)...)
		}
		out = append(out, childrenOf(a.Default)...)
	}
	return out
}
func (n *TruncateStmt) Children() []Node { return nil }

func (n *ShowStmt) Children() []Node { return childrenOf(n.Like, n.Where) }

func (n *SetStmt) Children() []Node {
	var out []Node
	for _, it := range n.Items {
		if it != nil {
			out = append(out, it.Value)
		}
	}
	return out
}

func (n *UseDbStmt) Children() []Node    { return nil }
func (n *ExplainStmt) Children() []Node  { return childrenOf(n.Stmt) }
func (n *BeginStmt) Children() []Node    { return nil }
func (n *CommitStmt) Children() []Node   { return nil }
func (n *RollbackStmt) Children() []Node { return nil }

func (n *SetTransactionStmt) Children() []Node      { return nil }
func (n *SavepointStmt) Children() []Node           { return nil }
func (n *ReleaseSavepointStmt) Children() []Node    { return nil }
func (n *RollbackToSavepointStmt) Children() []Node { return nil }
func (n *LockTablesStmt) Children() []Node          { return nil }
func (n *UnlockTablesStmt) Children() []Node        { return nil }

func expressionOfItem(it *ColumnItem) Node {
	if it == nil {
		return nil
	}
	return it.Expr
}

func windowSpecChildren(w *WindowSpec) []Node {
	var out []Node
	out = append(out, enode(w.PartitionBy)...)
	for i := range w.OrderBy {
		out = append(out, w.OrderBy[i].Expr)
	}
	return out
}
