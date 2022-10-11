package scanner

import (
	"strconv"
	"unicode"

	"github.com/vgalaktionov/glox/util"
)

type Scanner struct {
	source   []rune
	tokens   []Token
	start    int
	current  int
	line     int
	reporter util.ErrorReporter
}

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func NewScanner(source string, reporter util.ErrorReporter) *Scanner {
	return &Scanner{[]rune(source), nil, 0, 0, 1, reporter}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	case '/':
		if s.match('/') {
			// comments go on until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if unicode.IsDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.reporter.Error(s.line, "Unexpected character: %s", string(c))
		}
	}
}

func isAlpha(c rune) bool {
	return c == '_' || unicode.IsLetter(c)
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || unicode.IsDigit(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) advance() rune {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t TokenType, literal interface{}) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, NewToken(t, text, literal, s.line))
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.reporter.Error(s.line, "Unterminated string.")
		return
	}

	// closing "
	s.advance()

	// trim quotes
	value := string(s.source[s.start+1 : s.current-1])
	s.addTokenLiteral(String, value)
}

func (s *Scanner) number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	// look for fractional part
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// consume the .
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	num := string(s.source[s.start:s.current])
	value, err := strconv.ParseFloat(num, 64)
	if err != nil {
		s.reporter.Error(s.line, "Failed to parse float %s: %s", num, err)
		return
	}
	s.addTokenLiteral(Number, value)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, found := keywords[text]
	if !found {
		tokenType = Identifier
	}

	s.addToken(tokenType)
}
