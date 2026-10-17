// Package version models an OceanBase server version, primarily the 4.x line.
package version

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is a parsed OceanBase version such as "4.2.1" or "4.2.4.2_CE".
// The comparison semantics follow OB's dotted numeric scheme; a missing
// component compares as zero.
type Version struct {
	Major  int
	Minor  int
	Patch  int
	Build  int // optional 4th numeric component (e.g. 4.2.4.2)
	Suffix string
}

// String reconstructs the canonical form "M.m.p" or "M.m.p.b".
func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Build != 0 {
		s += fmt.Sprintf(".%d", v.Build)
	}
	if v.Suffix != "" {
		s += "_" + v.Suffix
	}
	return s
}

// Compare returns -1, 0 or +1.
func (v Version) Compare(o Version) int {
	for _, pair := range [][2]int{{v.Major, o.Major}, {v.Minor, o.Minor}, {v.Patch, o.Patch}, {v.Build, o.Build}} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	return 0
}

// After reports whether v > o.
func (v Version) After(o Version) bool { return v.Compare(o) > 0 }

// Before reports whether v < o.
func (v Version) Before(o Version) bool { return v.Compare(o) < 0 }

// AtLeast reports whether v >= o.
func (v Version) AtLeast(o Version) bool { return v.Compare(o) >= 0 }

// IsZero reports whether no numeric component was set.
func (v Version) IsZero() bool { return v.Major == 0 && v.Minor == 0 && v.Patch == 0 && v.Build == 0 }

// Parse parses a version string. Accepted forms: "4", "4.2", "4.2.4",
// "4.2.4.2", and any of those with a trailing "_SUFFIX". It is tolerant of
// leading/trailing whitespace.
func Parse(s string) (Version, error) {
	var v Version
	s = strings.TrimSpace(s)
	if s == "" {
		return v, fmt.Errorf("version: empty string")
	}
	// optional _suffix, e.g. _CE or _BP1
	if i := strings.IndexByte(s, '_'); i >= 0 {
		v.Suffix = s[i+1:]
		s = s[:i]
	}
	parts := strings.Split(s, ".")
	if len(parts) > 4 {
		return v, fmt.Errorf("version: too many components in %q", s)
	}
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			return v, fmt.Errorf("version: empty component in %q", s)
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return v, fmt.Errorf("version: invalid component %q in %q", p, s)
		}
		nums = append(nums, n)
	}
	assign := []*int{&v.Major, &v.Minor, &v.Patch, &v.Build}
	for i, n := range nums {
		*assign[i] = n
	}
	return v, nil
}

// MustParse is Parse that panics on error; useful for constants and tests.
func MustParse(s string) Version {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Convenience constants for the OceanBase 4.x minor versions.
var (
	V4_0 = MustParse("4.0.0")
	V4_1 = MustParse("4.1.0")
	V4_2 = MustParse("4.2.0")
	V4_3 = MustParse("4.3.0")
	V4_4 = MustParse("4.4.0")
)
