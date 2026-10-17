package lexer

import (
	"testing"

	"github.com/zetatez/ob-sql-parser/token"
)

func tokTypes(src string) []token.Type {
	toks, err := LexAll(src, Options{})
	if err != nil {
		panic(err)
	}
	var out []token.Type
	for _, t := range toks {
		if t.Type == token.EOF {
			break
		}
		out = append(out, t.Type)
	}
	return out
}

func TestOperators(t *testing.T) {
	cases := map[string][]token.Type{
		"a<=b":  {token.IDENT, token.LE, token.IDENT},
		"a<>b":  {token.IDENT, token.NE, token.IDENT},
		"a<<b":  {token.IDENT, token.SHL, token.IDENT},
		"a>>b":  {token.IDENT, token.SHR, token.IDENT},
		"a<=>b": {token.IDENT, token.NULLSAFE_EQ, token.IDENT},
		"a:=b":  {token.IDENT, token.ASSIGN, token.IDENT},
		"a!=b":  {token.IDENT, token.NE, token.IDENT},
		"a=b":   {token.IDENT, token.EQ, token.IDENT},
	}
	for src, want := range cases {
		got := tokTypes(src)
		if len(got) != len(want) {
			t.Fatalf("%q: got %v want %v", src, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("%q: token %d = %v, want %v", src, i, got[i], want[i])
			}
		}
	}
}

func TestLiterals(t *testing.T) {
	cases := []struct {
		sql  string
		kind token.Type
		val  string
	}{
		{"123", token.INT, "123"},
		{"1.5", token.DECIMAL, "1.5"},
		{"1e3", token.FLOAT, "1e3"},
		{"0xFF", token.HEXNUM, "0xFF"},
		{"0b101", token.BITNUM, "0b101"},
		{"x'41'", token.HEXSTRING, "41"},
		{"'hi'", token.STRING, "hi"},
		{"'a''b'", token.STRING, "a'b"},
		{"n'x'", token.STRING, "x"},
		{"_utf8mb4'hi'", token.STRING, "utf8mb4\x00hi"},
	}
	for _, c := range cases {
		toks, err := LexAll(c.sql, Options{})
		if err != nil {
			t.Fatalf("%q: %v", c.sql, err)
		}
		if len(toks) < 2 {
			t.Fatalf("%q: expected 1 token + EOF", c.sql)
		}
		tk := toks[0]
		if tk.Type != c.kind {
			t.Fatalf("%q: type=%v want %v (lit=%q)", c.sql, tk.Type, c.kind, tk.Literal)
		}
		if tk.Value != c.val {
			t.Fatalf("%q: value=%q want %q", c.sql, tk.Value, c.val)
		}
	}
}

func TestQuotedIdent(t *testing.T) {
	toks, _ := LexAll("`from` `a b`", Options{})
	// comments none
	var idents []string
	for _, tk := range toks {
		if tk.Type == token.IDENT {
			idents = append(idents, tk.Value)
		}
	}
	if len(idents) != 2 || idents[0] != "from" || idents[1] != "a b" {
		t.Fatalf("got %v", idents)
	}
}

func TestPositions(t *testing.T) {
	toks, _ := LexAll("SELECT\n  a", Options{})
	// find 'a'
	for _, tk := range toks {
		if tk.Literal == "a" {
			if tk.Pos.Line != 2 || tk.Pos.Col != 3 {
				t.Fatalf("a position = %d:%d, want 2:3", tk.Pos.Line, tk.Pos.Col)
			}
		}
	}
}

func TestCommentsSkipped(t *testing.T) {
	toks, _ := LexAll("a /* hi */ b -- trailing\nc", Options{})
	var lits []string
	for _, tk := range toks {
		if tk.Type == token.IDENT {
			lits = append(lits, tk.Literal)
		}
	}
	if len(lits) != 3 || lits[0] != "a" || lits[1] != "b" || lits[2] != "c" {
		t.Fatalf("idents = %v", lits)
	}
}

func TestLexerOptions(t *testing.T) {
	// Double-quoted identifiers under ANSI_QUOTES option.
	toks, err := LexAll(`"col"=1`, Options{AllowDoubleQuotedIdentifiers: true})
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Type != token.IDENT || toks[0].Value != "col" {
		t.Fatalf("double-quoted ident: %v %q (%q)", toks[0].Type, toks[0].Literal, toks[0].Value)
	}

	// Fullwidth option must not panic and still tokenize plain text.
	if toks, err := LexAll("SELECT 1", Options{Fullwidth: true}); err != nil || len(toks) < 3 {
		t.Fatalf("fullwidth mode: %v (%d tokens)", err, len(toks))
	}
}

func TestLexerErrors(t *testing.T) {
	bad := []string{
		"'unterminated",
		"`unterminated",
		"/* unterminated",
		"0x",
	}
	for _, src := range bad {
		if _, err := LexAll(src, Options{}); err == nil {
			t.Errorf("LexAll(%q) should error", src)
		}
	}
}

func TestLexerEdgeLiterals(t *testing.T) {
	cases := []struct {
		sql  string
		kind token.Type
		val  string
	}{
		{".5", token.DECIMAL, ".5"},
		{"5.", token.DECIMAL, "5."},
		{"1E+3", token.FLOAT, "1E+3"},
		{"1e-3", token.FLOAT, "1e-3"},
		{"0b0001", token.BITNUM, "0b0001"},
		{"'it''s'", token.STRING, "it's"},
		{`'a\'b'`, token.STRING, "a'b"},
		{"'\\n'", token.STRING, "\n"},
		{"x'0fF'", token.HEXSTRING, "0fF"},
		{"b'101'", token.BITSTRING, "101"},
	}
	for _, tc := range cases {
		toks, err := LexAll(tc.sql, Options{})
		if err != nil {
			t.Errorf("%q: %v", tc.sql, err)
			continue
		}
		if len(toks) < 2 || toks[0].Type != tc.kind {
			t.Errorf("%q: first token = %v (%q), want %v", tc.sql, toks[0].Type, toks[0].Literal, tc.kind)
			continue
		}
		if toks[0].Value != tc.val {
			t.Errorf("%q: value = %q, want %q", tc.sql, toks[0].Value, tc.val)
		}
	}
}

func TestVariablesAndParams(t *testing.T) {
	cases := []struct {
		sql  string
		kind token.Type
		lit  string
	}{
		{"@x", token.VARIABLE, "@x"},
		{"@@sql_mode", token.SYSVARIABLE, "@@sql_mode"},
		{"@@session.x", token.SYSVARIABLE, "@@session.x"},
		{"?", token.QUESTION, "?"},
	}
	for _, tc := range cases {
		toks, err := LexAll(tc.sql, Options{})
		if err != nil {
			t.Fatalf("%q: %v", tc.sql, err)
		}
		if toks[0].Type != tc.kind {
			t.Errorf("%q: got %v want %v", tc.sql, toks[0].Type, tc.kind)
		}
	}
}

func TestCharsetIntroducer(t *testing.T) {
	toks, err := LexAll("_utf8mb4'hi'", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Type != token.STRING {
		t.Fatalf("introducer token = %v", toks[0].Type)
	}
	if toks[0].Value != "utf8mb4\x00hi" {
		t.Fatalf("introducer value = %q", toks[0].Value)
	}
}
