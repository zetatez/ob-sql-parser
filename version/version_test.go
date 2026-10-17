package version

import "testing"

func TestParseAndString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"4", "4.0.0"},
		{"4.2", "4.2.0"},
		{"4.2.4", "4.2.4"},
		{"4.2.4.2", "4.2.4.2"},
		{"4.2.4_CE", "4.2.4_CE"},
		{" 4.3.0 ", "4.3.0"},
	}
	for _, c := range cases {
		v, err := Parse(c.in)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.in, err)
		}
		if v.String() != c.want {
			t.Fatalf("Parse(%q).String() = %q, want %q", c.in, v.String(), c.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, bad := range []string{"", "4..2", "abc", "4.2.1.2.3"} {
		if _, err := Parse(bad); err == nil {
			t.Fatalf("expected error for %q", bad)
		}
	}
}

func TestCompare(t *testing.T) {
	if !V4_2.After(V4_1) {
		t.Fatal("4.2 should be after 4.1")
	}
	if !V4_2.Before(V4_3) {
		t.Fatal("4.2 should be before 4.3")
	}
	if !V4_2.AtLeast(V4_2) {
		t.Fatal("4.2 should be at least 4.2")
	}
	if V4_0.Compare(V4_0) != 0 {
		t.Fatal("equal versions should compare 0")
	}
}

func TestParsePatchOrdering(t *testing.T) {
	a := MustParse("4.2.4")
	b := MustParse("4.2.10")
	if b.Before(a) {
		t.Fatal("4.2.10 should be after 4.2.4")
	}
}
