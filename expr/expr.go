package expr

import (
	"errors"

	"github.com/jesperkha/Fizz/lexer"
)

var (
	ErrParenError			= errors.New("unmatched parenthesies, line %d")
	ErrInvalidUnaryOperator = errors.New("invalid unary operator '%s', line %d")
	ErrInvalidOperatorType  = errors.New("invalid operator '%s' for type '%s', line %d")
	ErrInvalidOperatorTypes = errors.New("invalid operator '%s' for types '%s' and '%s', line %d")
	ErrDivideByZero			= errors.New("division by 0, line %d")
	ErrInvalidExpression   = errors.New("invalid expression, line %d")
)

const (
	// Expression types
	Literal = iota
	Unary
	Binary
	Group
	Variable

	// ParseToken types
	Single
	TokenGroup
)

type Expression struct {
	Type     int
	Name	 string
	Operand  lexer.Token
	Terminal lexer.Token
	Value    lexer.Token
	Left	 *Expression
	Right 	 *Expression
	Inner 	 *Expression
}

type ParseToken struct {
	Type   int
	Token  lexer.Token
	Inner  []ParseToken
}