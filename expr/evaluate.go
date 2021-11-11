package expr

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/lexer"
)

// Evaluates expression tree. Hands off to helper methods which can also recursively call to
// resolve nested expressions. Returned value is result of expression and is Go literal.
func EvaluateExpression(expr *Expression) (value interface{}, err error) {
	switch expr.Type {
	case Literal:
		return expr.Value.Literal, err
	case Unary:
		return evalUnary(expr)
	case Binary:
		return evalBinary(expr)
	case Group:
		return EvaluateExpression(expr.Inner)
	case Variable:
		return env.Get(expr.Name)
	case Call:
		return evalCall(expr)
	}

	// Wont be reached
	return expr, nil
}

// Helper for any unary expression with a valid operator. Returns an error
// if the type does not match the operator or the operator is invalid.
func evalUnary(unary *Expression) (value interface{}, err error) {
	right, err := EvaluateExpression(unary.Right)

	switch (unary.Operand.Type) {
		case lexer.MINUS: {
			if !isBool(right) {
				return -right.(float64), err
			}

			op, typ, line := unary.Operand.Lexeme, getType(right), unary.Operand.Line
			return nil, fmt.Errorf(ErrInvalidOperatorType.Error(), op, typ, line)
		}
		case lexer.NOT: return !isTruthy(right), err
		case lexer.TYPE: return getType(right), err
	}

	op, line := unary.Operand.Lexeme, unary.Operand.Line
	return value, fmt.Errorf(ErrInvalidUnaryOperator.Error(), op, line)
}

// Helper for evaluation for any binary expression with a valid operator.
// If the types do not match or the operator is invalid an error is returned.
func evalBinary(binary *Expression) (value interface{}, err error) {
	left, err := EvaluateExpression(binary.Left)
	if err != nil {
		return nil, err
	}

	right, err := EvaluateExpression(binary.Right)
	if err != nil {
		return nil, err
	}

	// Numbers are float64, set from lexer
	if isNumber(right) && isNumber(left) {
		vl, vr := left.(float64), right.(float64)
		switch (binary.Operand.Type) {
			case lexer.PLUS: return vl + vr, err
			case lexer.MINUS: return vl - vr, err
			case lexer.STAR: return vl * vr, err
			case lexer.HAT: return math.Pow(vl, vr), err
			case lexer.GREATER: return vl > vr, err
			case lexer.LESS: return vl < vr, err
			case lexer.LESS_EQUAL: return vl <= vr, err
			case lexer.GREATER_EQUAL: return vl >= vr, err
			case lexer.MODULO: return float64(int(vl) % int(vr)), err
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

	// Support string addition
	if getType(left) == "string" && getType(right) == "string" && binary.Operand.Type == lexer.PLUS {
		return strings.Join([]string{left.(string), right.(string)}, ""), err
	}
	
	typeLeft, typeRight := getType(left), getType(right)
	op, line := binary.Operand.Lexeme, binary.Operand.Line
	return nil, fmt.Errorf(ErrInvalidOperatorTypes.Error(), op, typeLeft, typeRight, line)
}

func evalCall(call *Expression) (value interface{}, err error) {

	return value, err
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
	if value == nil {
		return "nil"
	}

	return reflect.TypeOf(value).Name()
}