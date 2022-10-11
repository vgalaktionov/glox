package scanner

import "fmt"

//go:generate stringer -type=TokenType
type TokenType int

const (
	// single char tokens
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	// one or two char tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	// literals
	Identifier
	String
	Number
	// keywords
	And
	Class
	Else
	False
	For
	Fun
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(t TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{t, lexeme, literal, line}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
}
