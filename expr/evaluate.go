package expr

import (
	"fmt"
	"math"
	"strings"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

// Evaluates expression tree. Hands off to helper methods which can also recursively call to
// resolve nested expressions. Returned value is result of expression and is Go literal.
func EvaluateExpression(expr *Expression) (value interface{}, err error) {
	switch expr.Type {
	case Literal:
		return evalLiteral(expr)
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
	case Array:
		return evalArray(expr)
	case Index:
		return evalIndex(expr)
	}

	// Wont be reached
	return expr, ErrInvalidType
}

// Token types >= string are valid literal types
func evalLiteral(literal *Expression) (value interface{}, err error) {
	if literal.Value.Type >= lexer.STRING {
		return literal.Value.Literal, err
	}

	return value, ErrInvalidExpression
}

func evalUnary(unary *Expression) (value interface{}, err error) {
	right, err := EvaluateExpression(unary.Right)
	if err != nil {
		return value, err
	}

	// Matches to operator
	switch unary.Operand.Type {
	case lexer.MINUS:
		if isNumber(right) {
			return -right.(float64), err
		}
		op, typ, line := unary.Operand.Lexeme, util.GetType(right), unary.Line
		return nil, fmt.Errorf(ErrInvalidOperatorType.Error(), op, typ, line)
	case lexer.NOT:
		return !isTruthy(right), err
	case lexer.TYPE:
		return util.GetType(right), err
	}

	// If none of the mentioned operators are present its an invalid one
	op, line := unary.Operand.Lexeme, unary.Line
	return value, fmt.Errorf(ErrInvalidUnaryOperator.Error(), op, line)
}

func evalBinary(binary *Expression) (value interface{}, err error) {
	// Recursivly evaluates left and right expressions
	left, err := EvaluateExpression(binary.Left)
	if err != nil {
		return nil, err
	}

	right, err := EvaluateExpression(binary.Right)
	if err != nil {
		return nil, err
	}

	// Operations if both are number types
	if isNumber(right) && isNumber(left) {
		vl, vr := left.(float64), right.(float64)
		switch binary.Operand.Type {
		case lexer.PLUS:
			return vl + vr, err
		case lexer.MINUS:
			return vl - vr, err
		case lexer.STAR:
			return vl * vr, err
		case lexer.HAT:
			return math.Pow(vl, vr), err
		case lexer.GREATER:
			return vl > vr, err
		case lexer.LESS:
			return vl < vr, err
		case lexer.LESS_EQUAL:
			return vl <= vr, err
		case lexer.GREATER_EQUAL:
			return vl >= vr, err
		case lexer.MODULO:
			return float64(int(vl) % int(vr)), err
		case lexer.SLASH:
			if vr == 0 {
				return nil, util.FormatError(ErrDivideByZero, binary.Line)
			}
			return vl / vr, err
		}
	}

	// Types do not need to match for comparisons
	switch binary.Operand.Type {
	case lexer.EQUAL_EQUAL:
		return left == right, err
	case lexer.NOT_EQUAL:
		return left != right, err
	case lexer.AND:
		return isTruthy(left) && isTruthy(right), err
	case lexer.OR:
		return isTruthy(left) || isTruthy(right), err
	}

	// Support string addition
	if util.GetType(left) == "string" && util.GetType(right) == "string" && binary.Operand.Type == lexer.PLUS {
		return strings.Join([]string{left.(string), right.(string)}, ""), err
	}

	// If non of the previous checks worked the expression is invalid
	typeLeft, typeRight := util.GetType(left), util.GetType(right)
	op, line := binary.Operand.Lexeme, binary.Line
	return nil, fmt.Errorf(ErrInvalidOperatorTypes.Error(), op, typeLeft, typeRight, line)
}

func evalCall(call *Expression) (value interface{}, err error) {
	callee, err := EvaluateExpression(call.Left)
	if err != nil {
		return value, err
	}
	
	// Function should be of type env.Callable
	if f, ok := callee.(env.Callable); ok {
		args := []interface{}{}
		// Single argument
		if call.Inner.Type != Args && call.Inner.Inner.Type != EmptyExpression {
			arg, err := EvaluateExpression(call.Inner)
			if err != nil {
				return value, err
			}

			args = append(args, arg)
		}
		
		// Argument list
		if call.Inner.Type == Args {
			for _, arg := range call.Inner.Exprs {
				val, err := EvaluateExpression(&arg)
				if err != nil {
					return value, err
				}
				
				args = append(args, val)
			}
		}

		if len(args) != f.NumArgs {
			return value, fmt.Errorf(ErrIncorrectArgs.Error(), f.Name, f.NumArgs, len(args), call.Line)
		}

		return f.Call(args...)
	}

	// Todo: fix error message for not function error
	return value, fmt.Errorf(ErrNotFunction.Error(), call.Name, call.Line)
}

func evalArray(array *Expression) (value interface{}, err error) {
	values := []interface{}{}
	for _, expr := range array.Exprs {
		if expr.Type == EmptyExpression {
			return value, fmt.Errorf(ErrNoExpression.Error(), array.Line)
		}
		
		val, err := EvaluateExpression(&expr)
		if err != nil {
			return value, err
		}

		values = append(values, val)
	}

	return env.Array{Values: values, Length: len(values)}, err
}

func evalIndex(array *Expression) (value interface{}, err error) {
	arr, err := env.Get(array.Name)
	if err != nil {
		return value, err
	}

	a, ok := arr.(env.Array)
	if !ok {
		return value, fmt.Errorf(ErrNotArray.Error(), array.Name, array.Line)
	}

	if array.Inner.Type == EmptyExpression {
		return value, fmt.Errorf(ErrNoExpression.Error(), array.Line)
	}

	index, err := EvaluateExpression(array.Inner)
	if err != nil {
		return value, err
	}

	if index, ok := index.(float64); ok {
		if index == float64(int(index)) {
			value, err = a.Get(int(index))
			if err != nil {
				return value, fmt.Errorf(err.Error(), array.Line)
			}
			
			return value, err
		}
	}
	
	// Not integer
	return value, fmt.Errorf(ErrNotInteger.Error(), array.Line)
}

func isTruthy(value interface{}) bool {
	return value != false && value != nil
}

func isNumber(value interface{}) bool {
	return util.GetType(value) == "number"
}