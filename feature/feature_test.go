package feature

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/version"
)

func TestDefaultGating(t *testing.T) {
	f := DefaultRegistry()

	// WithWindows: introduced at 4.2
	if f.Enabled(FWithWindows, version.V4_1) {
		t.Fatal("FWithWindows should be disabled at 4.1")
	}
	if !f.Enabled(FWithWindows, version.V4_2) {
		t.Fatal("FWithWindows should be enabled at 4.2")
	}

	// ExplainExtended: enabled [4.1, 4.3)
	if f.Enabled(FExplainExtended, version.V4_0) {
		t.Fatal("FExplainExtended should be disabled at 4.0")
	}
	if !f.Enabled(FExplainExtended, version.V4_1) {
		t.Fatal("FExplainExtended should be enabled at 4.1")
	}
	if f.Enabled(FExplainExtended, version.V4_3) {
		t.Fatal("FExplainExtended should be disabled at 4.3 (Until is exclusive)")
	}

	// Unregistered features are always enabled.
	if !f.Enabled("no_such_feature", version.V4_0) {
		t.Fatal("unregistered features should be enabled")
	}

	// Zero version means "everything enabled".
	if !f.Enabled(FWithWindows, version.Version{}) {
		t.Fatal("zero version should enable all features")
	}
}

func TestCustomRegistry(t *testing.T) {
	f := NewRegistry()
	if err := f.Register(Rule{Feature: "x", Since: version.MustParse("4.4.0")}); err != nil {
		t.Fatal(err)
	}
	if f.Enabled("x", version.V4_3) {
		t.Fatal("custom feature should be gated")
	}
	if !f.Enabled("x", version.V4_4) {
		t.Fatal("custom feature should be enabled at 4.4")
	}
	if _, ok := f.RuleFor("x"); !ok {
		t.Fatal("RuleFor should find registered rule")
	}
}
