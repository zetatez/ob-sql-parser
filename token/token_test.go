package token_test

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/token"
)

func TestKeywordLookup(t *testing.T) {
	cases := []struct {
		word     string
		isKw     bool
		reserved bool
	}{
		{"SELECT", true, true},
		{"select", true, true},
		{"FROM", true, true},
		{"indexes", true, false}, // non-reserved keyword
		{"keys", true, false},
		{"nope_not_a_kw", false, false},
		{"", false, false},
		{"USER", true, false}, // context keyword examples should exist in table
	}
	for _, tc := range cases {
		kw, res := token.LookupKeyword(tc.word)
		if kw != tc.isKw || res != tc.reserved {
			t.Errorf("LookupKeyword(%q) = (%v,%v), want (%v,%v)", tc.word, kw, res, tc.isKw, tc.reserved)
		}
	}
}

func TestReservedWord(t *testing.T) {
	if !token.IsReservedWord("SELECT") {
		t.Fatal("SELECT should be reserved")
	}
	if token.IsReservedWord("identifier") {
		t.Fatal("identifier should not be reserved")
	}
}

func TestIsKeyword(t *testing.T) {
	if !token.IsKeyword("FROM") {
		t.Fatal("FROM is a keyword")
	}
	if token.IsKeyword("xyz_unlikely") {
		t.Fatal("xyz_unlikely should not be a keyword")
	}
}

func TestTypeString(t *testing.T) {
	for _, tt := range []struct {
		ty   token.Type
		want string
	}{
		{token.IDENT, "IDENT"},
		{token.STRING, "STRING"},
		{token.EQEQ, "EQEQ"},
		{token.OTHER, "OTHER"},
	} {
		if got := tt.ty.String(); got != tt.want {
			t.Errorf("Type(%d).String() = %q, want %q", tt.ty, got, tt.want)
		}
	}
	// Unknown type falls back to numeric formatting and must not panic.
	_ = token.Type(9999).String()
}

func TestPositionIsValid(t *testing.T) {
	if (token.Position{Line: 1, Col: 1}).IsValid() != true {
		t.Fatal("1-based position should be valid")
	}
	if (token.Position{}).IsValid() {
		t.Fatal("zero position should be invalid")
	}
}

func TestTokenSpanText(t *testing.T) {
	tok := token.Token{
		Type:    token.IDENT,
		Literal: "abc",
		Pos:     token.Position{Offset: 0, Line: 1, Col: 1},
		EndPos:  token.Position{Offset: 3, Line: 1, Col: 4},
	}
	if got := tok.Span(); got.Start.Col != 1 || got.End.Col != 4 {
		t.Fatalf("span = %+v", got)
	}
	if got := tok.Text(); got != "abc" {
		t.Fatalf("text = %q", got)
	}
}
