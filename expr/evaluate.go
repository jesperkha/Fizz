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
	// Todo: (doing) add eval for array indexing
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
	notFuncErr := fmt.Errorf(ErrNotFunction.Error(), call.Name, call.Line)
	f, err := env.Get(call.Name)
	if err != nil {
		// Use not function instead of not variable
		return value, notFuncErr
	}

	// Function should be of type env.Callable
	if function, ok := f.(env.Callable); ok {
		return evalCallableObject(call, function)
	}

	return value, notFuncErr
}

// Seperate function because evalGetter needs to be able to pass a reference to a function
// and not just a name to get from the current env
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

	// Call function, error is returned and handled to add support for custom errors in libraries
	val, err := f.Call(args...)
	// Perform type check to avoid runtime errors with illegal Fizz types
	if !validFizzType(val) {
		return val, fmt.Errorf(ErrIllegalType.Error(), util.GetType(val))
	}

	return val, err
}

func evalGetter(getter *Expression) (value interface{}, err error) {
	var last interface{}
	line := getter.Line
	for idx, gtr := range getter.Exprs {
		// First token has to be object as a single name without
		// a dot is just parsed as a normal variable
		if util.GetType(last) == "nil" {
			value, err = EvaluateExpression(&gtr)
			if err != nil {
				return value, err
			}

			last = value
			continue
		}

		// Last has to be object to continue. If gtr.Name is empty its not an identifier
		if util.GetType(last) != "object" || gtr.Name == "" {
			value, err = EvaluateExpression(&gtr)
			if err != nil {
				return value, err
			}

			return value, fmt.Errorf(ErrNotObject.Error(), util.GetType(value), line)
		}

		// Will not raise error because checked before
		lastObj := last.(env.Object)

		// Get the value of last.name
		current, err := lastObj.Get(gtr.Name)
		if err != nil {
			return value, fmt.Errorf(err.Error(), lastObj.Name, gtr.Name, line)
		}

		// If function call, set current as the return value of the function + args from gtr
		if gtr.Type == Call {
			// Also check if the value actually is a function
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

// Checks if type is valid for fizz (used for native functions)
func validFizzType(val interface{}) bool {
	return util.SContains(LegalTypes, util.GetType(val))
}
