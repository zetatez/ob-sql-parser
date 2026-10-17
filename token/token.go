// Package token defines lexical tokens, source positions, and the keyword
// classification table used by the lexer and identifier quoting.

package token

import "strconv"

// Type enumerates lexical token categories.
type Type int

const (
	ILLEGAL Type = iota
	EOF
	COMMENT

	// IDENT is emitted for any identifier-like word (a plain identifier, a
	// quoted identifier, or a keyword). Keyword status is carried on the Token.
	IDENT

	// Literals.
	STRING    // 'abc' (with optional _charset introducer)
	HEXSTRING // x'41' or X'41'
	BITSTRING // b'01' or B'01'
	INT       // 123
	DECIMAL   // 1.23
	FLOAT     // 1e5, 1.2e-3
	HEXNUM    // 0x41
	BITNUM    // 0b1010

	// Placeholders and variables.
	QUESTION    // ?
	VARIABLE    // @name
	SYSVARIABLE // @@global.x / @@session.x / @@x

	// Operators / punctuation.
	COMMA // ,
	DOT   // .
	SEMI  // ;
	LPAREN
	RPAREN
	LBRACKET // rarely used
	RBRACKET

	// Arithmetic.
	PLUS    // +
	MINUS   // -
	STAR    // *
	SLASH   // /
	PERCENT // %
	POW     // ^
	SHL     // <<
	SHR     // >>
	AMP     // &
	BITOR   // |
	TILDE   // ~

	// Comparison / logical assignment.
	EQ            // =
	NE            // <> or !=
	LT            // <
	GT            // >
	LE            // <=
	GE            // >=
	EQEQ          // ==   (non-standard, legacy)
	LEFTSHIFT_EQ  // <<= reserved
	RIGHTSHIFT_EQ // >>=
	ASSIGN        // :=
	NULLSAFE_EQ   // <=>

	LPAREN2 // sentinel, unused

	OTHER // falls through as illegal or custom
)

var typeNames = map[Type]string{
	ILLEGAL: "ILLEGAL", EOF: "EOF", COMMENT: "COMMENT",
	IDENT:  "IDENT",
	STRING: "STRING", HEXSTRING: "HEXSTRING", BITSTRING: "BITSTRING",
	INT: "INT", DECIMAL: "DECIMAL", FLOAT: "FLOAT", HEXNUM: "HEXNUM", BITNUM: "BITNUM",
	QUESTION: "QUESTION", VARIABLE: "VARIABLE", SYSVARIABLE: "SYSVARIABLE",
	COMMA: "COMMA", DOT: "DOT", SEMI: "SEMI", LPAREN: "LPAREN", RPAREN: "RPAREN",
	LBRACKET: "LBRACKET", RBRACKET: "RBRACKET",
	PLUS: "PLUS", MINUS: "MINUS", STAR: "STAR", SLASH: "SLASH", PERCENT: "PERCENT",
	POW: "POW", SHL: "SHL", SHR: "SHR", AMP: "AMP", BITOR: "BITOR", TILDE: "TILDE",
	EQ: "EQ", NE: "NE", LT: "LT", GT: "GT", LE: "LE", GE: "GE", EQEQ: "EQEQ",
	ASSIGN: "ASSIGN", NULLSAFE_EQ: "NULLSAFE_EQ", OTHER: "OTHER",
}

func (t Type) String() string {
	if s, ok := typeNames[t]; ok {
		return s
	}
	return "Type(" + strconv.Itoa(int(t)) + ")"
}

// Position is a 1-based line/column with a zero-based byte offset.
type Position struct {
	Offset int
	Line   int
	Col    int
}

// IsValid reports whether the position carries a line number.
func (p Position) IsValid() bool { return p.Line > 0 }

// Span is a half-open [Start, End) source range.
type Span struct {
	Start Position
	End   Position
}

// Token is a single lexical token.
type Token struct {
	Type     Type
	Literal  string // exact source text (for strings, unescaped value stored in Value)
	Value    string // semantic value for literals (unescaped string, numeric text, charset)
	Keyword  string // normalized (lowercase) keyword text, empty if not a keyword
	Reserved bool   // whether the keyword is reserved (only meaningful when Keyword != "")
	Pos      Position
	EndPos   Position
}

// Span returns the token's source range.
func (t Token) Span() Span { return Span{Start: t.Pos, End: t.EndPos} }

// Text returns the token's literal for diagnostic purposes.
func (t Token) Text() string { return t.Literal }
