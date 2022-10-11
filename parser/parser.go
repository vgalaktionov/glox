package parser

import (
	"log"

	//lint:ignore ST1001 it's okay to avoid the line noise
	. "github.com/vgalaktionov/glox/scanner"
	"github.com/vgalaktionov/glox/util"
)

//go:generate go run ../tools/genast.go parser Expr
type Parser struct {
	tokens        []Token
	current       int
	errorReporter util.ErrorReporter
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}
func (p *Parser) Parse() Expr {
	var expr Expr
	defer func() {
		if r := recover(); r != nil {
			log.Println("encountered ParseError")
		}
	}()
	expr = p.expression()
	return expr
}

// Syntax productions here

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(Slash, Star) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right := p.unary()
		return Unary{operator, right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(False) {
		return Literal{false}
	}
	if p.match(True) {
		return Literal{true}
	}
	if p.match(Nil) {
		return Literal{nil}
	}

	if p.match(Number, String) {
		return Literal{p.previous().Literal}
	}

	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return Grouping{expr}
	}
	p.error(p.peek(), "Expect expression.")
	return nil
}

// Internal methods here

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}
func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}
	return p.error(p.peek(), message)
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token Token, message string) Token {
	p.errorReporter.Error(token, message)
	panic("ParseError")
}

//lint:ignore U1000 it will be needed soon
func (p *Parser) synchronize(token Token, message string) {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}

		switch p.peek().Type {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}

		p.advance()
	}
}
