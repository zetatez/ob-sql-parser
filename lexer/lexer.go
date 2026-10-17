// Package lexer implements a hand-written tokenizer for OceanBase MySQL-mode
// SQL. It is the only lexical layer feeding the gparser goyacc engine.

package lexer

import (
	"fmt"

	"github.com/zetatez/ob-sql-parser/token"
)

// Options tweak scanner behavior.
type Options struct {
	// Fullwidth converts full-width (ideographic) punctuation to ASCII when
	// true, mirroring OceanBase's full-width converter.
	Fullwidth bool
	// AllowDoubleQuotedIdentifiers treats "..." as a quoted identifier
	// (ANSI_QUOTES). When false, double quotes are treated like string
	// literals (MySQL default).
	AllowDoubleQuotedIdentifiers bool
}

// Lexer scans a SQL source string into tokens.
type Lexer struct {
	src  []byte
	pos  int
	line int
	col  int
	opts Options
}

// New returns a Lexer over src.
func New(src string, opts Options) *Lexer {
	return &Lexer{src: []byte(src), line: 1, col: 1, opts: opts}
}

// LexAll scans the entire source into a token slice (including an EOF token).
func LexAll(src string, opts Options) ([]token.Token, error) {
	l := New(src, opts)
	var toks []token.Token
	for {
		t, err := l.Next()
		if err != nil {
			return toks, err
		}
		toks = append(toks, t)
		if t.Type == token.EOF {
			return toks, nil
		}
	}
}

// Next returns the next token.
func (l *Lexer) Next() (token.Token, error) {
	for {
		c, ok := l.peekByte()
		if !ok {
			return l.mkToken(token.EOF, "", l.pos), nil
		}
		if isSpace(c) {
			l.advance()
			continue
		}

		if c == '#' {
			return l.scanLineComment()
		}
		if c == '-' && l.peekAt(1) == '-' {

			third, ok := l.peekAtByte(2)
			if !ok || isSpace(third) {
				return l.scanLineComment()
			}
		}
		if c == '/' && l.peekAt(1) == '*' {
			return l.scanBlockComment()
		}
		break
	}

	start := l.pos
	c, _ := l.peekByte()

	if c == '`' {
		l.advance()
		return l.scanQuotedIdent(start)
	}
	if c == '"' && l.opts.AllowDoubleQuotedIdentifiers {
		l.advance()
		return l.scanQuotedIdent(start)
	}

	if (c == 'x' || c == 'X') && l.peekAt(1) == '\'' {
		return l.scanHexBit(start, token.HEXSTRING)
	}
	if (c == 'b' || c == 'B') && l.peekAt(1) == '\'' {
		return l.scanHexBit(start, token.BITSTRING)
	}
	if (c == 'n' || c == 'N') && l.peekAt(1) == '\'' {

		l.advance()
		return l.scanString()
	}

	if c == '_' && l.peekAt(1) != '\'' && l.identFollows(1) {
		if t, ok := l.tryScanIntroducer(start); ok {
			return t, nil
		}
	}

	if isIdentStart(c) {
		return l.scanIdentifier()
	}

	if isDigit(c) {
		return l.scanNumber()
	}

	switch c {
	case 'x', 'X':
		if l.peekAt(1) == '\'' {
			return l.scanHexBit(l.advancePos(), token.HEXSTRING)
		}
	case 'b', 'B':
		if l.peekAt(1) == '\'' {
			return l.scanHexBit(l.advancePos(), token.BITSTRING)
		}
	case 'n', 'N':

		if l.peekAt(1) == '\'' {
			l.advance()
			return l.scanString()
		}
	}

	if c == '_' && l.peekAt(1) != '\'' && l.identFollows(1) {
		if t, ok := l.tryScanIntroducer(start); ok {
			return t, nil
		}
	}

	if c == '\'' || c == '"' {
		return l.scanString()
	}

	switch c {
	case '?':
		l.advance()
		return l.mkToken(token.QUESTION, "?", start), nil
	case '@':
		return l.scanVariable()
	case ',':
		l.advance()
		return l.mkToken(token.COMMA, ",", start), nil
	case '.':

		if isDigit(l.peekAt(1)) {
			return l.scanNumber()
		}
		l.advance()
		return l.mkToken(token.DOT, ".", start), nil
	case ';':
		l.advance()
		return l.mkToken(token.SEMI, ";", start), nil
	case '(':
		l.advance()
		return l.mkToken(token.LPAREN, "(", start), nil
	case ')':
		l.advance()
		return l.mkToken(token.RPAREN, ")", start), nil
	case '+':
		l.advance()
		return l.mkToken(token.PLUS, "+", start), nil
	case '-':
		l.advance()
		return l.mkToken(token.MINUS, "-", start), nil
	case '*':
		l.advance()
		return l.mkToken(token.STAR, "*", start), nil
	case '/':
		l.advance()
		return l.mkToken(token.SLASH, "/", start), nil
	case '%':
		l.advance()
		return l.mkToken(token.PERCENT, "%", start), nil
	case '^':
		l.advance()
		return l.mkToken(token.POW, "^", start), nil
	case '&':
		if l.peekAt(1) == '&' {
			l.advance()
			l.advance()
			return l.mkToken(token.OTHER, "&&", start), nil
		}
		l.advance()
		return l.mkToken(token.AMP, "&", start), nil
	case '|':
		if l.peekAt(1) == '|' {
			l.advance()
			l.advance()
			return l.mkToken(token.OTHER, "||", start), nil
		}
		l.advance()
		return l.mkToken(token.BITOR, "|", start), nil
	case '~':
		l.advance()
		return l.mkToken(token.TILDE, "~", start), nil
	case '<':
		switch l.peekAt(1) {
		case '=':
			if l.peekAt(2) == '>' {
				l.advanceN(3)
				return l.mkToken(token.NULLSAFE_EQ, "<=>", start), nil
			}
			l.advanceN(2)
			return l.mkToken(token.LE, "<=", start), nil
		case '>':
			l.advanceN(2)
			return l.mkToken(token.NE, "<>", start), nil
		case '<':
			if l.peekAt(2) == '=' {
				l.advanceN(3)
				return l.mkToken(token.OTHER, "<<=", start), nil
			}
			l.advanceN(2)
			return l.mkToken(token.SHL, "<<", start), nil
		}
		l.advance()
		return l.mkToken(token.LT, "<", start), nil
	case '>':
		switch l.peekAt(1) {
		case '=':
			l.advanceN(2)
			return l.mkToken(token.GE, ">=", start), nil
		case '>':
			if l.peekAt(2) == '=' {
				l.advanceN(3)
				return l.mkToken(token.OTHER, ">>=", start), nil
			}
			l.advanceN(2)
			return l.mkToken(token.SHR, ">>", start), nil
		}
		l.advance()
		return l.mkToken(token.GT, ">", start), nil
	case '=':
		if l.peekAt(1) == '=' {
			l.advanceN(2)
			return l.mkToken(token.EQEQ, "==", start), nil
		}
		l.advance()
		return l.mkToken(token.EQ, "=", start), nil
	case '!':
		if l.peekAt(1) == '=' {
			l.advanceN(2)
			return l.mkToken(token.NE, "!=", start), nil
		}
		l.advance()
		return l.mkToken(token.OTHER, "!", start), nil
	case ':':
		if l.peekAt(1) == '=' {
			l.advanceN(2)
			return l.mkToken(token.ASSIGN, ":=", start), nil
		}

		l.advance()
		return l.mkToken(token.OTHER, ":", start), nil
	}

	return token.Token{}, fmt.Errorf("lexer: illegal character %q at %d:%d", string(c), l.line, l.col)
}
func (l *Lexer) scanLineComment() (token.Token, error) {
	start := l.pos
	for {
		c, ok := l.peekByte()
		if !ok || c == '\n' {
			break
		}
		l.advance()
	}
	return l.mkToken(token.COMMENT, l.text(start, l.pos), start), nil
}
func (l *Lexer) scanBlockComment() (token.Token, error) {
	start := l.pos
	l.advance()
	l.advance()
	for {
		c, ok := l.peekByte()
		if !ok {
			return token.Token{}, fmt.Errorf("lexer: unterminated block comment at %d:%d", l.line, l.col)
		}
		if c == '*' && l.peekAt(1) == '/' {
			l.advance()
			l.advance()
			break
		}
		l.advance()
	}
	return l.mkToken(token.COMMENT, l.text(start, l.pos), start), nil
}
func (l *Lexer) scanIdentifier() (token.Token, error) {
	start := l.pos
	for {
		c, ok := l.peekByte()
		if !ok || !isIdentChar(c) {
			break
		}
		l.advance()
	}
	word := l.text(start, l.pos)
	tok := l.mkToken(token.IDENT, word, start)
	tok.Value = word
	if isKeyword, reserved := token.LookupKeyword(word); isKeyword {
		tok.Keyword = lower(word)
		tok.Reserved = reserved
	}
	return tok, nil
}
func lower(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}
func (l *Lexer) mkToken(typ token.Type, literal string, start int) token.Token {
	return token.Token{
		Type:    typ,
		Literal: literal,
		Value:   literal,
		Pos:     l.posTo(start),
		EndPos:  l.posTo(l.pos),
	}
}
func (l *Lexer) posTo(offset int) token.Position {

	line, col := 1, 1
	for i := 0; i < offset && i < len(l.src); i++ {
		if l.src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return token.Position{Offset: offset, Line: line, Col: col}
}
func (l *Lexer) text(from, to int) string {
	return string(l.src[from:to])
}
func (l *Lexer) advancePos() int {
	p := l.pos
	l.advance()
	return p
}
func (l *Lexer) advance() {
	if l.pos < len(l.src) {
		if l.opts.Fullwidth {
			l.pos += l.applyFullwidth(l.pos)
		} else {
			l.pos++
		}
	}
}
func (l *Lexer) applyFullwidth(i int) int {
	w := fullwidthWidth(l.src[i])
	l.updateLineCol(l.src[i])
	return w
}
func (l *Lexer) advanceN(n int) {
	for k := 0; k < n && l.pos < len(l.src); k++ {
		l.advance()
	}
}
func (l *Lexer) updateLineCol(c byte) {
	if c == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
}
func (l *Lexer) peekByte() (byte, bool) {
	if l.pos >= len(l.src) {
		return 0, false
	}
	return l.src[l.pos], true
}
func (l *Lexer) peekAtByte(i int) (byte, bool) {
	if l.pos+i >= len(l.src) {
		return 0, false
	}
	return l.src[l.pos+i], true
}
func (l *Lexer) peekAt(i int) byte {
	c, _ := l.peekAtByte(i)
	return c
}
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}
func isDigit(c byte) bool { return c >= '0' && c <= '9' }
func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
func isIdentStart(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c >= 0x80
}
func isIdentChar(c byte) bool {
	return isIdentStart(c) || isDigit(c) || c == '$'
}
func isVarChar(c byte) bool {
	return isIdentChar(c) || c == '.' || c == '$'
}
func fullwidthWidth(c byte) int { return 1 }
