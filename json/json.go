// Package json provides AST -> JSON serialization. Every node is rendered with
// a top-level "kind" key (derived from the concrete type) plus its exported
// fields, producing self-describing output suitable for debugging, storage, or
// cross-language exchange.
package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/zetatez/ob-sql-parser/ast"
)

// Marshal returns the JSON encoding of an AST node, including a "kind" key.
func Marshal(n ast.Node) ([]byte, error) {
	body, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		// not an object (should not happen for AST nodes); fall back to raw
		return body, nil
	}
	kb, _ := json.Marshal(string(n.NodeKind()))
	m["kind"] = kb
	return json.Marshal(m)
}

// MarshalIndent returns indented, human-readable JSON.
func MarshalIndent(n ast.Node) ([]byte, error) {
	b, err := Marshal(n)
	if err != nil {
		return nil, err
	}
	var out []byte
	if err := compactIndent(b, &out); err != nil {
		return b, nil
	}
	return out, nil
}

func compactIndent(in []byte, out *[]byte) error {
	var v interface{}
	if err := json.Unmarshal(in, &v); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	*out = b
	return nil
}

// Write writes the indented JSON form of node to w.
func Write(w io.Writer, n ast.Node) error {
	data, err := MarshalIndent(n)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// Decode is an alias of Unmarshal for symmetric Marshal/Decode naming.
func Decode(data []byte) (ast.Node, error) { return Unmarshal(data) }

// String returns the compact JSON string for n.
func String(n ast.Node) string {
	data, err := Marshal(n)
	if err != nil {
		return fmt.Sprintf("<json error: %v>", err)
	}
	return string(data)
}
