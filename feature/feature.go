// Package feature implements feature-tag based version gating for OceanBase
// 4.x MySQL-mode grammar features. Each feature is identified by a stable
// FeatureID and is available within a [Since, Until] version interval.
package feature

import (
	"fmt"
	"sync"

	"github.com/zetatez/ob-sql-parser/version"
)

// FeatureID is a stable, human-readable identifier for a syntax feature.
type FeatureID string

// Feature tags relevant to the current statement coverage. The set grows as
// coverage expands toward full OB 4.x MySQL-mode syntax.
const (
	FNone             FeatureID = "none"
	FWithWindows      FeatureID = "with_windows"       // named/versioned windows (4.2+)
	FAutoPartition    FeatureID = "auto_partition"     // automatic partition syntax (4.3+)
	FExplainExtended  FeatureID = "explain_extended"   // EXPLAIN EXTENDED (added 4.1, removed 4.3)
	FRowWindowFrame   FeatureID = "row_window_frame"   // ROWS/RANGE window frame (4.0+)
	FWithCte          FeatureID = "with_cte"           // WITH ... AS common table expressions (4.0+)
	FJsonTableFunc    FeatureID = "json_table"         // JSON_TABLE (4.2+)
	FColumnStoreIndex FeatureID = "column_store_index" // column store secondary index (4.3+)

	// Features introduced after 4.2.x. Since values calibrated against the
	// statement inventories extracted from the official v4.0.0 / v4.1.0 / v4.2.x /
	// v4.3.2_CE_BETA / v4.4.2_CE grammar files (docs/version-coverage.md).
	FMaterializedView    FeatureID = "materialized_view"     // 物化视图 + MVIEW LOG (4.3+)
	FExternalCatalog     FeatureID = "external_catalog"      // EXTERNAL CATALOG (4.4+)
	FExternalLocation    FeatureID = "external_location"     // LOCATION / REMOVE FILES (4.4+)
	FRoutineLoad         FeatureID = "routine_load"          // ROUTINE LOAD (5.0+, 不在 4.4 文法)
	FCCLRule             FeatureID = "ccl_rule"              // CONCURRENT_LIMITING_RULE (4.4+)
	FSensitiveRule       FeatureID = "sensitive_rule"        // SENSITIVE RULE (4.4+)
	FTenantSnapshot      FeatureID = "tenant_snapshot"       // CREATE/DROP SNAPSHOT ... FOR TENANT (4.3+)
	FCloneTenant         FeatureID = "clone_tenant"          // CLONE TENANT (4.3+)
	FAlterLS             FeatureID = "alter_ls"              // ALTER SYSTEM CREATE/MODIFY LS (4.4+)
	FFlashbackStandbyLog FeatureID = "flashback_standby_log" // ALTER SYSTEM FLASHBACK STANDBY LOG (4.4+)
	FMerge               FeatureID = "merge"                 // MERGE INTO (4.4+)
	FGroupingSets        FeatureID = "grouping_sets"         // GROUPING SETS/CUBE/ROLLUP (4.4+)
	FExchangeValidation  FeatureID = "exchange_validation"   // EXCHANGE PARTITION WITH[OUT] VALIDATION (4.4+)

	// 4.2.x minor-series refinements (calibrated: 4.2.2 has transfer_partition,
	// 4.2.4 has role/service, 4.2.5 has XA RECOVER; 4.2.0/4.2.1 have none).
	FTransferPartition FeatureID = "transfer_partition" // TRANSFER PARTITION (4.2.2+)
	FRoleManagement    FeatureID = "role_management"    // CREATE/DROP/SET ROLE (4.2.4+)
	FServiceName       FeatureID = "service_name"       // ALTER SYSTEM … SERVICE (4.2.4+)
	FXaRecover         FeatureID = "xa_recover"         // XA RECOVER (4.2.5+)
)

// Rule describes the version interval in which a feature is available.
type Rule struct {
	Feature FeatureID
	Since   version.Version // inclusive lower bound
	Until   version.Version // exclusive upper bound; zero means "no upper bound"
	Desc    string
}

func (r Rule) ok(v version.Version) bool {
	if v.IsZero() {
		return true // unversioned target => accept everything
	}
	if v.Before(r.Since) {
		return false
	}
	if !r.Until.IsZero() && !v.Before(r.Until) {
		return false
	}
	return true
}

// Registry holds feature rules. It is safe for concurrent use.
type Registry struct {
	mu    sync.RWMutex
	rules map[FeatureID]Rule
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{rules: make(map[FeatureID]Rule)}
}

// Register adds or replaces a rule for a feature.
func (f *Registry) Register(r Rule) error {
	if r.Since.IsZero() {
		return fmt.Errorf("feature: %q: Since version is required", r.Feature)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rules[r.Feature] = r
	return nil
}

// Enabled reports whether the feature is enabled at version v.
func (f *Registry) Enabled(feat FeatureID, v version.Version) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	r, ok := f.rules[feat]
	if !ok {
		// Features without a registered rule are always enabled.
		return true
	}
	return r.ok(v)
}

// RuleFor returns the rule for a feature (ok=false if unregistered).
func (f *Registry) RuleFor(feat FeatureID) (Rule, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	r, ok := f.rules[feat]
	return r, ok
}

// VersionInfo describes why a feature is not enabled at a version, for errors.
func (f *Registry) VersionInfo(feat FeatureID) string {
	r, ok := f.RuleFor(feat)
	if !ok {
		return ""
	}
	s := r.Since.String()
	if !r.Until.IsZero() {
		return fmt.Sprintf("%s supported from %s until (exclusive) %s", feat, s, r.Until.String())
	}
	return fmt.Sprintf("%s supported from %s", feat, s)
}

// DefaultRegistry returns a registry pre-populated for the OB 4.x line.
// Version gating current reflects known 4.x feature introductions.
func DefaultRegistry() *Registry {
	f := NewRegistry()
	_ = f.Register(Rule{Feature: FWithWindows, Since: version.MustParse("4.2.0"), Desc: "named/versioned window syntax"})
	// PARTITIONS AUTO is present in every 4.x grammar (4.0.0_CE onward), so it
	// is not gated for 4.x targets; the rule is kept for documentation only.
	_ = f.Register(Rule{Feature: FAutoPartition, Since: version.MustParse("4.0.0"), Desc: "automatic partition sizing (present since 4.0)"})
	_ = f.Register(Rule{Feature: FExplainExtended, Since: version.MustParse("4.1.0"), Until: version.MustParse("4.3.0"), Desc: "EXPLAIN EXTENDED"})
	_ = f.Register(Rule{Feature: FRowWindowFrame, Since: version.MustParse("4.0.0"), Desc: "ROWS/RANGE window frame"})
	_ = f.Register(Rule{Feature: FWithCte, Since: version.MustParse("4.0.0"), Desc: "common table expressions"})
	_ = f.Register(Rule{Feature: FJsonTableFunc, Since: version.MustParse("4.2.0"), Desc: "JSON_TABLE table function"})
	_ = f.Register(Rule{Feature: FColumnStoreIndex, Since: version.MustParse("4.3.0"), Desc: "column store secondary index"})
	// post-4.2.x features (versions calibrated from official grammar inventories)
	_ = f.Register(Rule{Feature: FMaterializedView, Since: version.MustParse("4.3.0"), Desc: "materialized view + MVIEW LOG"})
	_ = f.Register(Rule{Feature: FExternalCatalog, Since: version.MustParse("4.4.0"), Desc: "EXTERNAL CATALOG"})
	_ = f.Register(Rule{Feature: FExternalLocation, Since: version.MustParse("4.4.0"), Desc: "LOCATION / REMOVE FILES"})
	_ = f.Register(Rule{Feature: FRoutineLoad, Since: version.MustParse("5.0.0"), Desc: "ROUTINE LOAD"})
	_ = f.Register(Rule{Feature: FCCLRule, Since: version.MustParse("4.4.0"), Desc: "CONCURRENT_LIMITING_RULE"})
	_ = f.Register(Rule{Feature: FSensitiveRule, Since: version.MustParse("4.4.0"), Desc: "SENSITIVE RULE"})
	_ = f.Register(Rule{Feature: FTenantSnapshot, Since: version.MustParse("4.3.0"), Desc: "tenant snapshot"})
	_ = f.Register(Rule{Feature: FCloneTenant, Since: version.MustParse("4.3.0"), Desc: "CLONE TENANT"})
	_ = f.Register(Rule{Feature: FAlterLS, Since: version.MustParse("4.4.0"), Desc: "ALTER SYSTEM CREATE/MODIFY LS"})
	_ = f.Register(Rule{Feature: FFlashbackStandbyLog, Since: version.MustParse("4.4.0"), Desc: "FLASHBACK STANDBY LOG"})
	_ = f.Register(Rule{Feature: FMerge, Since: version.MustParse("4.4.0"), Desc: "MERGE INTO"})
	_ = f.Register(Rule{Feature: FGroupingSets, Since: version.MustParse("4.4.0"), Desc: "GROUPING SETS/CUBE/ROLLUP (4.4.2_CE)"})
	_ = f.Register(Rule{Feature: FExchangeValidation, Since: version.MustParse("4.4.0"), Desc: "EXCHANGE PARTITION WITH[OUT] VALIDATION (4.4.2_CE)"})
	// 4.2.x minor-series (earliest observed tag in parentheses)
	_ = f.Register(Rule{Feature: FTransferPartition, Since: version.MustParse("4.2.2"), Desc: "TRANSFER PARTITION (4.2.2_CE)"})
	_ = f.Register(Rule{Feature: FRoleManagement, Since: version.MustParse("4.2.4"), Desc: "CREATE/DROP/SET ROLE (4.2.4_CE)"})
	_ = f.Register(Rule{Feature: FServiceName, Since: version.MustParse("4.2.4"), Desc: "ALTER SYSTEM … SERVICE (4.2.4_CE)"})
	_ = f.Register(Rule{Feature: FXaRecover, Since: version.MustParse("4.2.5"), Desc: "XA RECOVER (4.2.5_CE)"})
	return f
}
