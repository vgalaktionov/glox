package parser_test

import (
	"testing"

	"github.com/vgalaktionov/glox/parser"
	"github.com/vgalaktionov/glox/scanner"
)

func TestAstPrinter(t *testing.T) {
	input := parser.Binary{
		parser.Unary{
			scanner.Token{Type: scanner.Minus, Lexeme: "-", Literal: nil, Line: 1},
			parser.Literal{123}},
		scanner.Token{Type: scanner.Star, Lexeme: "*", Literal: nil, Line: 1},
		parser.Grouping{
			parser.Literal{45.67}}}

	expected := "(* (- 123) (group 45.67))"

	if parser.AstPrinter(input) != expected {
		t.Errorf("input %v not equal to %s", input, expected)
	}
}
