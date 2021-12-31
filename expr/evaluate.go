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
	case Getter:
		return evalGetter(expr)
	case Array:
		return evalArray(expr)
	case Index:
		return evalIndex(expr)
	}

	// Wont be reached
	return expr, ErrInvalidExpression
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

	// Todo: add in operator for array

	// If non of the previous checks worked the expression is invalid
	typeLeft, typeRight := util.GetType(left), util.GetType(right)
	op, line := binary.Operand.Lexeme, binary.Line
	return nil, fmt.Errorf(ErrInvalidOperatorTypes.Error(), op, typeLeft, typeRight, line)
}

// Todo: function caching if flag is active
func evalCall(call *Expression) (value interface{}, err error) {
	callee, err := EvaluateExpression(call.Left)
	if err != nil {
		return value, err
	}

	// Function should be of type env.Callable
	if f, ok := callee.(*env.Callable); ok {
		argToken := call.Inner.Inner
		args := []interface{}{}

		// Single argument
		if argToken.Type != Args && argToken.Type != EmptyExpression {
			arg, err := EvaluateExpression(call.Inner)
			if err != nil {
				return value, err
			}

			args = append(args, arg)
		}

		// Argument list
		if argToken.Type == Args {
			for _, arg := range argToken.Exprs {
				val, err := EvaluateExpression(&arg)
				if err != nil {
					return value, err
				}

				args = append(args, val)
			}
		}

		// -1 is set from /lib and should be ignored as it is handled there
		if len(args) != f.NumArgs && f.NumArgs != -1 {
			return value, fmt.Errorf(ErrIncorrectArgs.Error(), f.Name, f.NumArgs, len(args), call.Line)
		}

		// Errors from lib need line format
		value, err = f.Call(args...)
		return value, util.FormatError(err, call.Line)
	}

	return value, fmt.Errorf(ErrNotFunction.Error(), util.GetType(callee), call.Line)
}

func evalGetter(getter *Expression) (value interface{}, err error) {
	line := getter.Line
	name := getter.Right.Name
	// No name before dot raises error here. No name after dot raises error in lexer.
	if getter.Left.Type == EmptyExpression {
		return value, fmt.Errorf(ErrInvalidExpression.Error(), line)
	}

	// Recursively get parent expression, must be object
	parent, err := EvaluateExpression(getter.Left)
	if err != nil {
		return value, err
	}

	// Only objects allow getter expressions
	if obj, ok := parent.(*env.Object); ok {
		value, err = obj.Get(name)
		if err != nil {
			return value, fmt.Errorf(err.Error(), obj.Name, name, line)
		}

		return value, err
	}

	return value, fmt.Errorf(ErrNotObject.Error(), util.GetType(parent), line)
}

func evalArray(array *Expression) (value interface{}, err error) {
	inner := array.Inner
	if inner.Type == EmptyExpression {
		return &env.Array{}, err
	}

	values := []interface{}{}
	// Single argument
	if inner.Type != Args && inner.Type != EmptyExpression {
		v, err := EvaluateExpression(inner)
		if err != nil {
			return value, err
		}

		values = append(values, v)
	}

	// Multiple arguments
	if inner.Type == Args {
		for _, expr := range inner.Exprs {
			v, err := EvaluateExpression(&expr)
			if err != nil {
				return value, err
			}

			values = append(values, v)
		}
	}

	return &env.Array{Values: values, Length: len(values)}, err
}

func evalIndex(array *Expression) (value interface{}, err error) {
	line := array.Line
	arr, err := EvaluateExpression(array.Left)
	if err != nil {
		return value, err
	}

	index, err := EvaluateExpression(array.Right)
	if err != nil {
		return value, err
	}

	// Get index as integer. If not return error
	indexInt, ok := util.IsInt(index)
	if !ok {
		return value, fmt.Errorf(ErrNotInteger.Error(), line)
	}

	if a, ok := arr.(*env.Array); ok {
		// Env handles getting index and errors for out of range etc
		value, err = a.Get(indexInt)
		if err != nil {
			return value, fmt.Errorf(err.Error(), line)
		}

		return value, err
	}

	// Get string index
	if s, ok := arr.(string); ok {
		if indexInt > len(s) {
			return value, env.ErrIndexOutOfRange
		}

		return string(s[indexInt]), err
	}

	// arr is not array (or string)
	return value, fmt.Errorf(env.ErrNotArray.Error(), util.GetType(arr), line)
}

func isTruthy(value interface{}) bool {
	return value != false && value != nil
}

func isNumber(value interface{}) bool {
	return util.GetType(value) == "number"
}
