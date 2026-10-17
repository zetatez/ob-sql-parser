package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/zetatez/ob-sql-parser/ast"
	"github.com/zetatez/ob-sql-parser/token"
)

// Unmarshal reconstructs an AST node (or a statement list produced by Dump/S)
// from JSON produced by this package. Every node must carry its kind tag.
func Unmarshal(data []byte) (ast.Node, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("json: expected an object, got %T", v)
	}
	if rb, ok := m["results"]; ok {
		arr, _ := rb.([]interface{})
		nodes := make([]ast.Node, 0, len(arr))
		for _, it := range arr {
			nm, ok := it.(map[string]interface{})
			if !ok {
				continue
			}
			n, err := decodeNode(nm)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, n)
		}
		return &statementList{nodes: nodes}, nil
	}
	return decodeNode(m)
}

// statementList wraps the results array of Dump.
type statementList struct{ nodes []ast.Node }

// Nodes returns the decoded statements.
func (l *statementList) Nodes() []ast.Node { return l.nodes }

func (l *statementList) NodeKind() ast.NodeKind { return "StatementList" }
func (l *statementList) Span() token.Span       { return token.Span{} }
func (l *statementList) Children() []ast.Node   { return l.nodes }

// Dump encodes several nodes as one JSON document {"results": [...]}.
func Dump(nodes []ast.Node) ([]byte, error) {
	arr := make([]json.RawMessage, 0, len(nodes))
	for _, n := range nodes {
		b, err := Marshal(n)
		if err != nil {
			return nil, err
		}
		arr = append(arr, b)
	}
	body, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return bytes.Join([][]byte{[]byte(`{"results":`), body, []byte(`}`)}, nil), nil
}

var (
	nodeType     = reflect.TypeOf((*ast.Node)(nil)).Elem()
	positionType = reflect.TypeOf(token.Position{})
)

// implementsNode reports whether t has the ast.Node method set.
func implementsNode(t reflect.Type) bool { return t != nil && t.Implements(nodeType) }

// isNodeInterface reports whether t is an interface node slot (Statement /
// Expression / Node).
func isNodeInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface && t.Implements(nodeType)
}

func kindOf(m map[string]interface{}) (string, bool) {
	if k, ok := m["kind"].(string); ok {
		return k, true
	}
	k, ok := m["Kind"].(string)
	return k, ok
}

func decodeNode(m map[string]interface{}) (ast.Node, error) {
	kind, ok := kindOf(m)
	if !ok || kind == "" {
		return nil, fmt.Errorf("json: missing node kind")
	}
	n := newNode(kind)
	if n == nil {
		return nil, fmt.Errorf("json: unknown node kind %q", kind)
	}
	if err := fillStruct(reflect.ValueOf(n).Elem(), m); err != nil {
		return nil, err
	}
	return n, nil
}

// decodeInto sets rv (must be addressable) from the decoded JSON value raw.
func decodeInto(rv reflect.Value, raw interface{}) error {
	if !rv.CanSet() || raw == nil {
		return nil
	}
	t := rv.Type()
	switch {
	case isNodeInterface(t) || (t.Kind() == reflect.Ptr && implementsNode(t)):
		return setNodeSlot(rv, raw, t)
	case t.Kind() == reflect.Ptr:
		// pointer to a shape struct (e.g. *WindowSpec, *LimitClause,
		// *PartitionSpec) or plain value; recurse into the element.
		elm := reflect.New(t.Elem())
		if err := decodeInto(elm.Elem(), raw); err != nil {
			return err
		}
		rv.Set(elm)
		return nil
	case t.Kind() == reflect.Slice:
		return setSlice(rv, raw)
	case t.Kind() == reflect.Struct:
		m, ok := raw.(map[string]interface{})
		if !ok {
			return jsonUnmarshalBack(raw, rv)
		}
		if t == positionType {
			return jsonUnmarshalBack(raw, rv)
		}
		return fillStruct(rv, m)
	default:
		return jsonUnmarshalBack(raw, rv)
	}
}

func setNodeSlot(rv reflect.Value, raw interface{}, t reflect.Type) error {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return jsonUnmarshalBack(raw, rv)
	}
	n, err := decodeNode(m)
	if err != nil {
		return err
	}
	nv := reflect.ValueOf(n)
	if isNodeInterface(t) {
		if nv.Type().Implements(t) {
			rv.Set(nv)
		}
		return nil
	}
	// pointer slot
	rv.Set(nv)
	return nil
}

func setSlice(rv reflect.Value, raw interface{}) error {
	items, ok := raw.([]interface{})
	if !ok {
		return jsonUnmarshalBack(raw, rv)
	}
	et := rv.Type().Elem()
	sl := reflect.MakeSlice(rv.Type(), 0, len(items))
	for _, it := range items {
		if it == nil {
			sl = reflect.Append(sl, reflect.Zero(et))
			continue
		}
		switch {
		case isNodeInterface(et) || (et.Kind() == reflect.Ptr && implementsNode(et)):
			m, ok := it.(map[string]interface{})
			if !ok {
				sl = reflect.Append(sl, reflect.Zero(et))
				continue
			}
			n, err := decodeNode(m)
			if err != nil {
				return err
			}
			if et.Kind() == reflect.Interface {
				sl = reflect.Append(sl, reflect.ValueOf(n))
			} else {
				sl = reflect.Append(sl, reflect.ValueOf(n))
			}
		default:
			elm := reflect.New(et)
			if err := decodeInto(elm.Elem(), it); err != nil {
				return err
			}
			sl = reflect.Append(sl, elm.Elem())
		}
	}
	rv.Set(sl)
	return nil
}

func fillStruct(v reflect.Value, m map[string]interface{}) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" {
			continue // unexported
		}
		raw, ok := m[sf.Name]
		if !ok {
			continue
		}
		if err := decodeInto(v.Field(i), raw); err != nil {
			return fmt.Errorf("json: %s.%s: %w", t.Name(), sf.Name, err)
		}
	}
	return nil
}

func jsonUnmarshalBack(raw interface{}, dst reflect.Value) error {
	b, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	nv := reflect.New(dst.Type())
	if err := json.Unmarshal(b, nv.Interface()); err != nil {
		return err
	}
	dst.Set(nv.Elem())
	return nil
}
