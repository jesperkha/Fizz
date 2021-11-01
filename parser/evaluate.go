package parser

import (
	"errors"
	"fmt"
	"math"

	"github.com/jesperkha/Fizz/lexer"
)

var (
	ErrInvalidUnaryOperator = errors.New("invalid unary operator '%s', line %d")
	ErrInvalidOperatorType  = errors.New("invalid operator '%s' for type '%s', line %d")
	ErrInvalidOperatorTypes = errors.New("invalid operator '%s' for types '%s' and '%s', line %d")
	ErrDivideByZero			= errors.New("division by 0, line %d")
)

// Todo: handle floats (also in lexer/parser)

// Evaluates expression. Hands off to helper methods which can also recursively call to 
// resolve nested expressions. Returned value is result of expression and is Go literal.
func EvaluateExpression(expr *Expression) (value interface{}, err error) {
	switch expr.Type {
	case Literal:
		return expr.Value.Literal, err
	case Unary:
		return handleUnary(expr)
	case Binary:
		return handleBinary(expr)
	case Group:
		return EvaluateExpression(expr.Inner)
	}

	// Wont be reached
	return expr, nil
}

func handleUnary(unary *Expression) (value interface{}, err error) {
	right, err := EvaluateExpression(unary.Right)

	switch (unary.Operand.Type) {
		case lexer.MINUS: {
			if !isBool(right) {
				return -right.(int), err
			}

			op, typ, line := unary.Operand.Lexeme, getType(right), unary.Operand.Line
			return nil, fmt.Errorf(ErrInvalidOperatorType.Error(), op, typ, line)
		}
		case lexer.NOT: return !isTruthy(right), err
	}

	op, line := unary.Operand.Lexeme, unary.Operand.Line
	return value, fmt.Errorf(ErrInvalidUnaryOperator.Error(), op, line)
}

func handleBinary(binary *Expression) (value interface{}, err error) {
	left, err := EvaluateExpression(binary.Left)
	if err != nil {
		return nil, err
	}

	right, err := EvaluateExpression(binary.Right)
	if err != nil {
		return nil, err
	}

	if isNumber(right) && isNumber(left) {
		vl, vr := left.(int), right.(int)
		switch (binary.Operand.Type) {
		case lexer.PLUS: return vl + vr, err
		case lexer.MINUS: return vl - vr, err
		case lexer.STAR: return vl * vr, err
		case lexer.HAT: return math.Pow(float64(vl), float64(vr)), err
		case lexer.GREATER: return vl > vr, err
		case lexer.LESS: return vl < vr, err
		case lexer.LESS_EQUAL: return vl <= vr, err
		case lexer.GREATER_EQUAL: return vl >= vr, err
		case lexer.SLASH:
			if vr == 0 {
				return nil, fmt.Errorf(ErrDivideByZero.Error(), binary.Operand.Line)
			}
			return vl / vr, err
		}
	}

	switch (binary.Operand.Type) {
		case lexer.EQUAL_EQUAL: return left == right, err
		case lexer.NOT_EQUAL: return left != right, err
		case lexer.AND: return isTruthy(left) && isTruthy(right), err
		case lexer.OR: return isTruthy(left) || isTruthy(right), err
	}

	typeLeft, typeRight := getType(left), getType(right)
	op, line := binary.Operand.Lexeme, binary.Operand.Line
	return nil, fmt.Errorf(ErrInvalidOperatorTypes.Error(), op, typeLeft, typeRight, line)
}

func isTruthy(value interface{}) bool {
	return value != false && value != nil
}

func isBool(value interface{}) bool {
	return value == false || value == true
}

func isNumber(value interface{}) bool {
	switch value.(type) {
		case int: return true
		case float32: return true
		case float64: return true
	}

	return false
}

func getType(value interface{}) string {
	if isNumber(value) {
		return "number"
	}

	if isBool(value) {
		return fmt.Sprintf("%v", value)
	}

	return "identifier"
}