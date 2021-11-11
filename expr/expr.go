package expr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesperkha/Fizz/lexer"
)

var (
	ErrParenError			= errors.New("unmatched parenthesies, line %d")
	ErrInvalidUnaryOperator = errors.New("invalid unary operator '%s', line %d")
	ErrInvalidOperatorType  = errors.New("invalid operator '%s' for type '%s', line %d")
	ErrInvalidOperatorTypes = errors.New("invalid operator '%s' for types '%s' and '%s', line %d")
	ErrDivideByZero			= errors.New("division by 0, line %d")
	ErrInvalidExpression    = errors.New("invalid expression, line %d")
	ErrExpectedExpression   = errors.New("expected expression in group, line %d")
)

const (
	// Expression types
	Literal = iota
	Unary
	Binary
	Group
	Variable
	Call

	// ParseToken types
	Single
	TokenGroup
	CallGroup
)

type Expression struct {
	Type     int
	Name	 string
	Operand  lexer.Token
	Callee   lexer.Token
	Value    lexer.Token
	Left	 *Expression
	Right 	 *Expression
	Inner 	 *Expression
	Exprs	 []Expression
}

type ParseToken struct {
	Type   int
	Token  lexer.Token
	Inner  []ParseToken
}

// Format error with line numbers for local errors, but ignore for errors passed from
// expression parsing as they are already formatted with line numbers.
func formatError(err error, line int) error {
	if err == nil {
		return err
	}

	if strings.Contains(err.Error(), "%d") {
		return fmt.Errorf(err.Error(), line)
	}

	return err
}