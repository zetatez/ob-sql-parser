// Package ast defines the typed abstract syntax tree for OceanBase MySQL-mode
// SQL, plus node Kind tags and source spans.

package ast

import (
	"reflect"

	"github.com/zetatez/ob-sql-parser/token"
)

// NodeKind is a stable string tag used for visitor dispatch and JSON output.
type NodeKind string

// Node is any AST node.
type Node interface {
	NodeKind() NodeKind
	Span() token.Span
	Children() []Node
}

// Statement is a top-level SQL statement.
type Statement interface {
	Node
	statementNode()
}

// Expression is an expression node.
type Expression interface {
	Node
	expressionNode()
}

// Base is embedded by all nodes. It holds the kind tag and source span.
// Fields are exported so encoding/json can serialize them.
type Base struct {
	Kind  NodeKind
	Start token.Position
	End   token.Position
}

func (b *Base) NodeKind() NodeKind {

	if b.Kind != "" {
		return b.Kind
	}
	t := reflect.TypeOf(b)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return NodeKind(t.Name())
}
func (b *Base) Span() token.Span { return token.Span{Start: b.Start, End: b.End} }

// SetSpan sets the node's source range. It is used by the goyacc driver to
// attach statement-level positions (child positions are not tracked yet).
func (b *Base) SetSpan(start, end token.Position) {
	b.Start = start
	b.End = end
}

// childrenOf is a helper returning a slice; used by concrete Children() methods.
func childrenOf(nodes ...Node) []Node {
	out := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		if n != nil {
			out = append(out, n)
		}
	}
	return out
}
func inode(s []*TableRef) []Node {
	out := make([]Node, 0, len(s))
	for _, n := range s {
		if n != nil {
			out = append(out, n)
		}
	}
	return out
}
func enode(s []Expression) []Node {
	out := make([]Node, 0, len(s))
	for _, n := range s {
		if n != nil {
			out = append(out, n)
		}
	}
	return out
}
