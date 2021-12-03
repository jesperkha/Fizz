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
	}

	// Wont be reached
	return expr, ErrInvalidExpression
}

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

	op, line := unary.Operand.Lexeme, unary.Line
	return value, fmt.Errorf(ErrInvalidUnaryOperator.Error(), op, line)
}

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

	typeLeft, typeRight := util.GetType(left), util.GetType(right)
	op, line := binary.Operand.Lexeme, binary.Line
	return nil, fmt.Errorf(ErrInvalidOperatorTypes.Error(), op, typeLeft, typeRight, line)
}

// Error is returned for function fail, unmatched arg number, or undefined name
func evalCall(call *Expression) (value interface{}, err error) {
	notFuncErr := fmt.Errorf(ErrNotFunction.Error(), call.Name, call.Line)
	f, err := env.Get(call.Name)
	if err != nil {
		return value, notFuncErr
	}

	if function, ok := f.(env.Callable); ok {
		return evalCallableObject(call, function)
	}

	return value, notFuncErr
}

// Seperate function because evalGetter needs to evaluate functions that are not in the current scope but rather value passed in nested objects
func evalCallableObject(call *Expression, f env.Callable) (value interface{}, err error) {
	numArgs := len(call.Exprs)
	if numArgs != f.NumArgs {
		name, expected, line := call.Name, f.NumArgs, call.Line
		return value, fmt.Errorf(ErrIncorrectArgs.Error(), name, expected, numArgs, line)
	}

	// Evaluated args to be passed to function
	args := []interface{}{}
	for _, arg := range call.Exprs {
		v, err := EvaluateExpression(&arg)
		if err != nil {
			return value, err
		}

		args = append(args, v)
	}

	val, err := f.Call(args...)
	if !validFizzType(val) {
		return val, fmt.Errorf(ErrIllegalType.Error(), util.GetType(val))
	}

	return val, err
}

func evalGetter(getter *Expression) (value interface{}, err error) {
	var last interface{}
	line := getter.Line
	for idx, gtr := range getter.Exprs {
		// For first token
		if util.GetType(last) == "nil" {
			value, err = EvaluateExpression(&gtr)
			if err != nil {
				return value, err
			}

			last = value
			continue
		}

		// Last has to be object to continue. If gtr.Name is empty its not a identifier
		if util.GetType(last) != "object" || gtr.Name == "" {
			value, err = EvaluateExpression(&gtr)
			if err != nil {
				return value, err
			}

			return value, fmt.Errorf(ErrNotObject.Error(), util.GetType(value), line)
		}

		// Will not raise error because checked before
		lastObj := last.(env.Object)

		// Get current value
		current, err := lastObj.Get(gtr.Name)
		if err != nil {
			return value, fmt.Errorf(err.Error(), lastObj.Name, gtr.Name, line)
		}

		// If function, set current as the return value of the function + args from gtr
		if gtr.Type == Call {
			if util.GetType(current) != "function" {
				return value, fmt.Errorf(ErrNotFunction.Error(), getter.Name, line)
			}

			current, err = evalCallableObject(&gtr, current.(env.Callable))
			if err != nil {
				return value, err
			}
		}

		// For last token return the value
		if idx == len(getter.Exprs)-1 {
			return current, err
		}

		last = current
	}

	return value, err
}

func isTruthy(value interface{}) bool {
	return value != false && value != nil
}

func isNumber(value interface{}) bool {
	switch value.(type) {
	case int:
		return true
	case float32:
		return true
	case float64:
		return true
	}

	return false
}

// Checks if type is valid for fizz (used for native functions)
func validFizzType(val interface{}) bool {
	return util.SContains(LegalTypes, util.GetType(val))
}
