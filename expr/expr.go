package expr

import (
	"errors"

	"github.com/jesperkha/Fizz/lexer"
)

var (
	ErrParenError           = errors.New("unmatched parenthesies, line %d")
	ErrBracketError         = errors.New("unmatched brackets, line %d")
	ErrInvalidUnaryOperator = errors.New("invalid unary operator '%s', line %d")
	ErrInvalidOperatorType  = errors.New("invalid operator '%s' for type %s, line %d")
	ErrInvalidOperatorTypes = errors.New("invalid operator '%s' for types %s and %s, line %d")
	ErrDivideByZero         = errors.New("division by 0, line %d")
	ErrNoExpression         = errors.New("empty expression, line %d")
	ErrInvalidExpression    = errors.New("invalid expression, line %d")
	ErrExpectedExpression   = errors.New("expected expression in group, line %d")
	ErrNotInteger           = errors.New("index must be integer, line %d")
	ErrCommaError           = errors.New("comma error, line %d")
	ErrIncorrectArgs        = errors.New("%s() expected %d args, got %d, line %d")
	ErrNotFunction          = errors.New("type %s is not a function, line %d")
	ErrNilValueError        = errors.New("unexpected nil value in expression, line %d")
	ErrNotObject            = errors.New("type %s has no attributes, line %d")
	ErrInvalidType          = errors.New("expr: unknown expression type, line %d")
	ErrExpectedName         = errors.New("expected name after dot, line %d")
	ErrIllegalType          = errors.New("unknown type '%s'")
)

const (
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
