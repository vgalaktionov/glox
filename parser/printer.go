package parser

import (
	"fmt"
	"strings"
)

func AstPrinter(anyexpr Expr) string {
	switch expr := anyexpr.(type) {
	case Binary:
		return parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
	case Grouping:
		return parenthesize("group", expr.Expression)
	case Literal:
		if expr.Value == nil {
			return "nil"
		}
		return fmt.Sprint(expr.Value)
	case Unary:
		return parenthesize(expr.Operator.Lexeme, expr.Right)
	default:
		panic(fmt.Sprintf("unimplemented printing for expr type %T", expr))
	}

}

func parenthesize(name string, exprs ...Expr) string {
	sb := new(strings.Builder)

	sb.WriteByte('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteByte(' ')
		sb.WriteString(AstPrinter(expr))
	}
	sb.WriteByte(')')

	return sb.String()
}
