package lexer

import (
	"fmt"
	"strings"

	"github.com/zetatez/ob-sql-parser/token"
)

func (l *Lexer) scanQuotedIdent(startPos int) (token.Token, error) {
	// startPos is the opening backtick/quote; the caller has consumed it.
	var sb strings.Builder
	for {
		c, ok := l.peekByte()
		if !ok {
			return token.Token{}, fmt.Errorf("lexer: unterminated quoted identifier at %d:%d", l.line, l.col)
		}
		if c == '`' || (c == '"') {

			if l.peekAt(1) == c {
				sb.WriteByte(c)
				l.advance()
				l.advance()
				continue
			}
			l.advance()
			break
		}
		if c == '\\' && l.peekAt(1) == '`' {
			sb.WriteByte('`')
			l.advance()
			l.advance()
			continue
		}
		sb.WriteByte(c)
		l.advance()
	}
	tok := l.mkToken(token.IDENT, l.text(startPos, l.pos), startPos)
	tok.Value = sb.String()

	tok.Keyword = ""
	tok.Reserved = false
	return tok, nil
}
func (l *Lexer) scanHexBit(startPos int, kind token.Type) (token.Token, error) {

	l.advance()
	l.advance()
	var sb strings.Builder
	for {
		c, ok := l.peekByte()
		if !ok {
			return token.Token{}, fmt.Errorf("lexer: unterminated %s literal", kind)
		}
		if c == '\'' {

			if l.peekAt(1) == '\'' {
				sb.WriteString("''")
				l.advance()
				l.advance()
				continue
			}
			l.advance()
			break
		}
		if kind == token.HEXSTRING && !isHexDigit(c) {
			return token.Token{}, fmt.Errorf("lexer: invalid hex char %q", string(c))
		}
		if kind == token.BITSTRING && c != '0' && c != '1' {
			return token.Token{}, fmt.Errorf("lexer: invalid bit char %q", string(c))
		}
		sb.WriteByte(c)
		l.advance()
	}
	tok := l.mkToken(kind, l.text(startPos, l.pos), startPos)
	tok.Value = sb.String()
	return tok, nil
}
func (l *Lexer) scanNumber() (token.Token, error) {
	start := l.pos

	if c, _ := l.peekByte(); c == '0' {
		nc, _ := l.peekAtByte(1)
		if nc == 'x' || nc == 'X' {
			p := l.pos
			l.advance()
			l.advance()
			digitsStart := l.pos
			for {
				c, ok := l.peekByte()
				if !ok || !isHexDigit(c) {
					break
				}
				l.advance()
			}
			if l.pos == digitsStart {
				return token.Token{}, fmt.Errorf("lexer: hex literal with no digits at %d:%d", l.line, l.col)
			}
			return l.mkToken(token.HEXNUM, l.text(p, l.pos), start), nil
		}
		if nc == 'b' || nc == 'B' {
			p := l.pos
			l.advance()
			l.advance()
			digitsStart := l.pos
			for {
				c, ok := l.peekByte()
				if !ok || (c != '0' && c != '1') {
					break
				}
				l.advance()
			}
			if l.pos == digitsStart {
				return token.Token{}, fmt.Errorf("lexer: bit literal with no digits at %d:%d", l.line, l.col)
			}
			return l.mkToken(token.BITNUM, l.text(p, l.pos), start), nil
		}
	}

	kind := token.INT
	for {
		c, ok := l.peekByte()
		if !ok {
			break
		}
		if isDigit(c) {
			l.advance()
			continue
		}
		if c == '.' {

			kind = token.DECIMAL
			l.advance()
			continue
		}
		if c == 'e' || c == 'E' {
			kind = token.FLOAT
			l.advance()
			if n, _ := l.peekByte(); n == '+' || n == '-' {
				l.advance()
			}
			continue
		}
		break
	}
	tok := l.mkToken(kind, l.text(start, l.pos), start)
	tok.Value = tok.Literal
	return tok, nil
}
func (l *Lexer) scanString() (token.Token, error) {
	start := l.pos
	quote, _ := l.peekByte()
	l.advance()
	var sb strings.Builder
	for {
		c, ok := l.peekByte()
		if !ok {
			return token.Token{}, fmt.Errorf("lexer: unterminated string literal at %d:%d", l.line, l.col)
		}
		if c == quote {
			if l.peekAt(1) == quote {
				sb.WriteByte(quote)
				l.advance()
				l.advance()
				continue
			}
			l.advance()
			break
		}
		if c == '\\' {
			esc, err := l.scanEscape()
			if err != nil {
				return token.Token{}, err
			}
			sb.WriteString(esc)
			continue
		}
		sb.WriteByte(c)
		l.advance()
	}
	tok := l.mkToken(token.STRING, l.text(start, l.pos), start)
	tok.Value = sb.String()
	return tok, nil
}
func (l *Lexer) scanEscape() (string, error) {
	l.advance()
	c, ok := l.peekByte()
	if !ok {
		return "", fmt.Errorf("lexer: dangling escape at %d:%d", l.line, l.col)
	}
	l.advance()
	switch c {
	case '0':
		return "\x00", nil
	case 'b':
		return "\b", nil
	case 'n':
		return "\n", nil
	case 'r':
		return "\r", nil
	case 't':
		return "\t", nil
	case 'Z':
		return "\x1a", nil
	case '\\':
		return "\\", nil
	case '\'':
		return "'", nil
	case '"':
		return "\"", nil
	case '`':
		return "`", nil
	case '%':
		return "%", nil
	case '_':
		return "_", nil
	default:
		return string(c), nil
	}
}

// charset introducer: _<identname>'...'
func (l *Lexer) tryScanIntroducer(start int) (token.Token, bool) {
	pos := l.pos
	l.advance()
	nameStart := l.pos
	for {
		c, ok := l.peekByte()
		if !ok || !isIdentChar(c) {
			break
		}
		l.advance()
	}
	name := l.text(nameStart, l.pos)
	c, _ := l.peekByte()
	if c != '\'' {

		l.pos = pos
		return token.Token{}, false
	}

	strTok, err := l.scanString()
	if err != nil {
		return token.Token{}, false
	}
	tok := l.mkToken(token.STRING, l.text(start, l.pos), start)
	tok.Value = strTok.Value

	tok.Value = name + "\x00" + strTok.Value
	return tok, true
}
func (l *Lexer) scanVariable() (token.Token, error) {
	start := l.pos
	l.advance()
	kind := token.VARIABLE
	if c, _ := l.peekByte(); c == '@' {
		l.advance()

		if g := l.tryMatchWord("global"); g {
			if c, _ := l.peekByte(); c == '.' {
				l.advance()
				kind = token.SYSVARIABLE
			}
		} else if s := l.tryMatchWord("session"); s {
			if c, _ := l.peekByte(); c == '.' {
				l.advance()
				kind = token.SYSVARIABLE
			}
		} else {
			kind = token.SYSVARIABLE
		}
	}

	for {
		c, ok := l.peekByte()
		if !ok || !isVarChar(c) {
			break
		}
		l.advance()
	}
	return l.mkToken(kind, l.text(start, l.pos), start), nil
}
func (l *Lexer) tryMatchWord(w string) bool {
	if l.pos+len(w) > len(l.src) {
		return false
	}
	sub := l.src[l.pos : l.pos+len(w)]
	if !strings.EqualFold(string(sub), w) {
		return false
	}
	l.advanceN(len(w))
	return true
}
func (l *Lexer) identFollows(off int) bool {
	c, ok := l.peekAtByte(off)
	if !ok {
		return false
	}
	return isIdentStart(c)
}
