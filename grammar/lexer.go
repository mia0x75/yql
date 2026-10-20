package grammar

import (
	"fmt"
	"strings"
)

// Parse parses a YQL expression into an AST.
func Parse(input string) (*Query, error) {
	l := &lexer{input: input}
	if yyParse(l) != 0 {
		if l.err != nil {
			return nil, l.err
		}
		return nil, fmt.Errorf("syntax error")
	}
	if l.err != nil {
		return nil, l.err
	}
	if l.query == nil {
		return nil, fmt.Errorf("invalid yql")
	}
	return l.query, nil
}

type lexer struct {
	input string
	pos   int
	query *Query
	err   error
}

func (l *lexer) Error(s string) {
	if l.err == nil {
		l.err = fmt.Errorf("%s", s)
	}
}

func (l *lexer) errorf(format string, args ...any) int {
	if l.err == nil {
		l.err = fmt.Errorf(format, args...)
	}
	return 0
}

func (l *lexer) Lex(lval *yySymType) int {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\r', '\n':
			l.pos++
		default:
			goto scan
		}
	}
	return 0

scan:
	start := l.pos
	c := l.input[l.pos]
	switch {
	case c == '\'' || c == '"':
		return l.scanString(lval)
	case isLetter(c) || c == '_':
		return l.scanIdent(lval)
	case c == '+' || c == '-' || isDigit(c):
		return l.scanNumber(lval)
	case c == '(':
		l.pos++
		return LPAREN
	case c == ')':
		l.pos++
		return RPAREN
	case c == ',':
		l.pos++
		return COMMA
	case c == '.':
		l.pos++
		return DOT
	}

	rest := l.input[start:]
	switch {
	case strings.HasPrefix(rest, "!="):
		l.pos += len("!=")
		lval.str = "!="
		return CMPOP
	case strings.HasPrefix(rest, "!intersect"):
		l.pos += len("!intersect")
		lval.str = "!intersect"
		return SETOP
	case strings.HasPrefix(rest, "!in"):
		l.pos += len("!in")
		lval.str = "!in"
		return SETOP
	case strings.HasPrefix(rest, "<>"):
		l.pos += len("<>")
		lval.str = "<>"
		return CMPOP
	case strings.HasPrefix(rest, ">="):
		l.pos += len(">=")
		lval.str = ">="
		return CMPOP
	case strings.HasPrefix(rest, "<="):
		l.pos += len("<=")
		lval.str = "<="
		return CMPOP
	case strings.HasPrefix(rest, "="):
		l.pos += len("=")
		lval.str = "="
		return CMPOP
	case strings.HasPrefix(rest, ">"):
		l.pos += len(">")
		lval.str = ">"
		return CMPOP
	case strings.HasPrefix(rest, "<"):
		l.pos += len("<")
		lval.str = "<"
		return CMPOP
	}
	return l.errorf("unexpected character %q at position %d", c, start)
}

func (l *lexer) scanString(lval *yySymType) int {
	start := l.pos
	quote := l.input[l.pos]
	l.pos++
	for l.pos < len(l.input) {
		if l.input[l.pos] == quote {
			l.pos++
			lval.str = l.input[start:l.pos]
			return STRING
		}
		l.pos++
	}
	return l.errorf("unterminated string at position %d", start)
}

func (l *lexer) scanIdent(lval *yySymType) int {
	start := l.pos
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if !isLetter(c) && c != '_' {
			break
		}
		l.pos++
	}
	word := l.input[start:l.pos]
	switch strings.ToLower(word) {
	case "and":
		return AND
	case "or":
		return OR
	case "true":
		lval.str = word
		return TRUE
	case "false":
		lval.str = word
		return FALSE
	case "count", "sum", "avg", "max", "min":
		lval.str = word
		return FUNC
	}
	if word == "in" || word == "intersect" {
		lval.str = word
		return SETOP
	}
	lval.str = word
	return FIELDNAME
}

func (l *lexer) scanNumber(lval *yySymType) int {
	start := l.pos
	if l.input[l.pos] == '+' || l.input[l.pos] == '-' {
		l.pos++
	}
	if l.pos >= len(l.input) || !isDigit(l.input[l.pos]) {
		return l.errorf("invalid number at position %d", start)
	}
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		l.pos++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
		lval.str = l.input[start:l.pos]
		return FLOAT
	}
	lval.str = l.input[start:l.pos]
	return INT
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
