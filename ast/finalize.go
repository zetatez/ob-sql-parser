package ast

import "reflect"

// Finalize walks the whole object graph (via reflection, not just Children)
// and stamps every reachable node's Base.Kind with its concrete type name.
// This produces a self-describing tree (each node carries a stable "Kind"
// discriminator) for JSON output/debugging. Reflection is used so nodes not
// exposed by Children — e.g. TableRef.Name/IndexHints — are stamped too.
func Finalize(n Node) {
	if n == nil {
		return
	}
	walkFinalize(reflect.ValueOf(n))
}

func walkFinalize(rv reflect.Value) {
	for rv.IsValid() {
		switch rv.Kind() {
		case reflect.Interface, reflect.Ptr:
			if rv.IsNil() {
				return
			}
			rv = rv.Elem()
			continue
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				walkFinalize(rv.Index(i))
			}
			return
		case reflect.Struct:
			stamp(rv)
			t := rv.Type()
			for i := 0; i < t.NumField(); i++ {
				if t.Field(i).PkgPath != "" {
					continue // unexported
				}
				walkFinalize(rv.Field(i))
			}
			return
		default:
			return
		}
	}
}

func stamp(rv reflect.Value) {
	bf := rv.FieldByName("Base")
	if !bf.IsValid() || !bf.CanSet() {
		return
	}
	kf := bf.FieldByName("Kind")
	if !kf.IsValid() || !kf.CanSet() {
		return
	}
	kf.SetString(rv.Type().Name())
}
