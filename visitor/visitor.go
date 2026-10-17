// Package visitor provides a generic AST traversal framework built on the
// Children() method of every ast.Node. Users implement Visitor and may return
// a (possibly stateful) sub-visitor, or nil to stop descending into a subtree.
package visitor

import (
	"reflect"

	"github.com/zetatez/ob-sql-parser/ast"
)

// Visitor dispatches on AST nodes during a walk.
// Returning nil from VisitNode prunes the node's subtree.
type Visitor interface {
	VisitNode(node ast.Node) Visitor
}

// Walk traverses node depth-first. VisitNode is invoked for every node. If it
// returns nil the node's subtree is pruned; otherwise traversal continues (the
// same visitor is reused for the children).
func Walk(node ast.Node, v Visitor) {
	if isNilNode(node) || v == nil {
		return
	}
	if v.VisitNode(node) == nil {
		return
	}
	for _, child := range node.Children() {
		Walk(child, v)
	}
}

// WalkStatements walks a list of statements.
func WalkStatements(stmts []ast.Statement, v Visitor) {
	for _, s := range stmts {
		Walk(s, v)
	}
}

// isNilNode reports whether n is nil or a typed-nil pointer. Children() of
// hand-built trees may contain typed-nil optionals (e.g. StoredProcStmt with a
// nil Body); treating those as absent keeps Walk panic-free.
func isNilNode(n ast.Node) bool {
	if n == nil {
		return true
	}
	rv := reflect.ValueOf(n)
	return rv.Kind() == reflect.Ptr && rv.IsNil()
}

// FuncVisitor adapts a function to the Visitor interface.
type FuncVisitor func(node ast.Node) Visitor

func (f FuncVisitor) VisitNode(node ast.Node) Visitor { return f(node) }

// NoOpVisitor always continues walking and does nothing.
type NoOpVisitor struct{}

func (NoOpVisitor) VisitNode(ast.Node) Visitor { return NoOpVisitor{} }

// Collect gathers all nodes of a specific interface type walking the tree.
func Collect[T ast.Node](node ast.Node, out *[]T) {
	Walk(node, FuncVisitor(func(n ast.Node) Visitor {
		if t, ok := n.(T); ok {
			*out = append(*out, t)
		}
		return fwdVisitor{}
	}))
}

type fwdVisitor struct{}

func (fwdVisitor) VisitNode(ast.Node) Visitor { return fwdVisitor{} }
