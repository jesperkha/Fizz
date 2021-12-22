package expr

import (
	"errors"
	"fmt"

	"github.com/jesperkha/Fizz/lexer"
)

var (
	ErrParenError           = errors.New("unmatched parenthesies, line %d")
	ErrBracketError         = errors.New("unmatched brackets, line %d")
	ErrInvalidUnaryOperator = errors.New("invalid unary operator '%s', line %d")
	ErrInvalidOperatorType  = errors.New("invalid operator '%s' for type %s, line %d")
	ErrInvalidOperatorTypes = errors.New("invalid operator '%s' for types %s and %s, line %d")
	ErrDivideByZero         = errors.New("division by 0, line %d")
	ErrNoExpression			= errors.New("empty expression, line %d")
	ErrInvalidExpression    = errors.New("invalid expression, line %d")
	ErrExpectedExpression   = errors.New("expected expression in group, line %d")
	ErrNotInteger			= errors.New("index must be integer, line %d")
	ErrCommaError           = errors.New("comma error, line %d")
	ErrIncorrectArgs        = errors.New("%s() expected %d args, got %d, line %d")
	ErrNotFunction          = errors.New("type %s is not a function, line %d")
	ErrNilValueError        = errors.New("unexpected nil value in expression, line %d")
	ErrNotObject            = errors.New("type %s has no attributes, line %d")
	ErrNotArray				= errors.New("variable '%s' is not an array, line %d")
	ErrInvalidType			= errors.New("expr: unknown expression type, line %d")
	ErrExpectedName			= errors.New("expected name after dot, line %d")
	ErrIllegalType = errors.New("unknown type '%s'")
)

var LegalTypes = []string{
	"number",
	"nil",
	"string",
	"function",
	"bool",
	"object",
}

const (
	// Expression types
	EmptyExpression = iota
	Literal
	Unary
	Binary
	Group
	Variable
	Call
	Args

	Getter
	Array
	Index

	// ParseToken types
	Single
	TokenGroup
	CallGroup
	ArrayGroup
	ArrayGetter
)

type Expression struct {
	Type    int
	Line    int
	Name    string
	Operand lexer.Token
	Value   lexer.Token
	Left    *Expression
	Right   *Expression
	Inner   *Expression
	Exprs   []Expression
}

type ParseToken struct {
	Type  int
	Token lexer.Token
	Inner []ParseToken
}

func ParseAndEval(tokens []lexer.Token) (value interface{}, err error) {
	expr, err := ParseExpression(tokens)
	if err != nil {
		return value, err
	}

	return EvaluateExpression(&expr)
}

func PrintExpressionAST(expr Expression) {
	fmt.Println(printAST(expr))
}

func printAST(expr Expression) string {
	switch expr.Type {
	case Literal:
		return fmt.Sprintf("literal: %s", expr.Value.Lexeme)
	case Unary:
		return fmt.Sprintf("unary: %s [%s]", expr.Operand.Lexeme, printAST(*expr.Right))
	case Binary:
		left, right := printAST(*expr.Left), printAST(*expr.Right)
		return fmt.Sprintf("binary: %s [%s, %s]", expr.Operand.Lexeme, left, right)
	case Group:
		return fmt.Sprintf("group: [%s]", printAST(*expr.Inner))
	case Variable:
		return fmt.Sprintf("variable: %s", expr.Name)
	case Call:
		return fmt.Sprintf("call: %s()", expr.Name)
	}

	return ""
}
