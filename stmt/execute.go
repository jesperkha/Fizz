package stmt

import (
	"errors"
	"fmt"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

var (
	CurrentOrigin     string
	MaxRecursionDepth = 1000
)

// Goes through list of statements and executes them. Error is returned from statements exec method.
func ExecuteStatements(stmts []Statement) (err error) {
	for _, statement := range stmts {
		line := statement.Line
		if err = executeStatement(statement); err != nil {
			if cerr, ok := err.(ConditionalError); ok {
				cerr.Msg = fmt.Sprintf(cerr.Msg, line)
				return cerr
			}

			return util.FormatError(err, line)
		}
	}

	return err
}

func executeStatement(stmt Statement) error {
	switch stmt.Type {
	case ExpressionStmt:
		_, err := expr.EvaluateExpression(stmt.Expression)
		return err
	case Block:
		return execBlock(stmt)
	case Print:
		return execPrint(stmt)
	case Variable:
		return nil
	case Assignment:
		return execAssignment(stmt)
	case Break:
		return ErrBeakOutsideLoop
	case Skip:
		return ErrSkipOutsideLoop
	case Return:
		return execReturn(stmt)
	case If:
		return execIf(stmt)
	case While:
		return execWhile(stmt)
	case Repeat:
		return execRepeat(stmt)
	case Function:
		return execFunction(stmt)
	case Exit:
		return execExit(stmt)
	case Error:
		return execError(stmt)
	case Object:
		return execObject(stmt)
	case Enum:
		return execEnum(stmt)
	case Range:
		return execRange(stmt)
	case Import, Include:
		return nil // Handled in interp
	}

	// Will never be returned since all types are pre-defined.
	// However it is nice to have in case rework is done and types
	// get mixed up or new types are only partially added.
	return ErrInvalidStmtType
}

func execEnum(stmt Statement) (err error) {
	for curVal, name := range stmt.Params {
		err = env.Declare(name, float64(curVal))
		if err != nil {
			return err
		}
	}

	return err
}

func execExit(stmt Statement) (err error) {
	if stmt.Expression != nil {
		if err = execPrint(stmt); err != nil {
			return err
		}
	}

	return ErrProgramExit
}

func execError(stmt Statement) (err error) {
	value, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	return errors.New(util.FormatPrintValue(value))
}

// Raises error and assigns expr value to global currentReturnValue
func execReturn(stmt Statement) (err error) {
	e := ErrReturnOutsideFunc
	if stmt.Expression == nil {
		e.Value = nil
		return e
	}

	value, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	e.Value = value
	return e
}

func execPrint(stmt Statement) (err error) {
	value, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	fmt.Println(util.FormatPrintValue(value))
	return nil
}

func assignValue(left *expr.Expression, value interface{}) error {
	if left.Type == expr.Variable {
		return env.Assign(left.Name, value)
	}

	// First evaluate the entire expression to pluck out any
	// errors that are harder to check for later
	if _, err := expr.EvaluateExpression(left); err != nil {
		return err
	}

	// Get left expression of left expression (parent)
	val, err := expr.EvaluateExpression(left.Left)
	if err != nil {
		return err
	}

	// If object assign to name of parent expression
	if obj, ok := val.(*env.Object); ok {
		return obj.Set(left.Right.Name, value)
	}

	// If array assign value to index of parent expression
	if arr, ok := val.(*env.Array); ok {
		index, err := expr.EvaluateExpression(left.Right)
		if err != nil {
			return err
		}

		indexInt, ok := util.IsInt(index)
		if !ok {
			return expr.ErrNotInteger
		}

		return arr.Set(indexInt, value)
	}

	return ErrNonAssignable
}

func execAssignment(stmt Statement) (err error) {
	val, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	// Plain assignment
	if stmt.Operator == lexer.EQUAL {
		return assignValue(stmt.Left, val)
	}

	// Declare variable with special := operator
	if stmt.Operator == lexer.DEF_EQUAL {
		if stmt.Left == nil || stmt.Expression == nil {
			return ErrInvalidStatement
		}

		return env.Declare(stmt.Left.Name, val)
	}

	oldVal, err := expr.EvaluateExpression(stmt.Left)
	if err != nil {
		return err
	}

	// Not same type
	oldType, newType := util.GetType(oldVal), util.GetType(val)
	if oldType != newType {
		return ErrDifferentTypes
	}

	// String addition
	if newType == "string" {
		if stmt.Operator != lexer.PLUS_EQUAL {
			return ErrInvalidOperator
		}

		return assignValue(stmt.Left, oldVal.(string)+val.(string))
	}

	// Float addition / subtraction
	if newType == "number" {
		a := oldVal.(float64)
		b := val.(float64)
		var newVal float64

		switch stmt.Operator {
		case lexer.PLUS_EQUAL:
			newVal = a + b
		case lexer.MINUS_EQUAL:
			newVal = a - b
		case lexer.MULT_EQUAL:
			newVal = a * b
		case lexer.DIV_EQUAL:
			newVal = a / b
		}

		return assignValue(stmt.Left, newVal)
	}

	return ErrInvalidStatement
}

func execBlock(stmt Statement) (err error) {
	env.PushScope()
	err = ExecuteStatements(stmt.Statements)
	env.PopScope()
	return err
}

// Monitor recursion
var lastFunction = ""
var recursionDepth = 0

func execFunction(stmt Statement) (err error) {
	// Store origin at point of function declaration as well as scope around it
	originCache := CurrentOrigin
	var envCache env.Environment

	function := env.Callable{
		Name:    stmt.Name,
		NumArgs: len(stmt.Params),
		Origin:  CurrentOrigin,
		// Call function and set param variables to scope
		Call: func(args ...interface{}) (interface{}, error) {
			// Handle recursion errors
			name := stmt.Name
			if lastFunction == name {
				recursionDepth++
			} else {
				lastFunction = name
			}

			// Todo: better recursive checks for recursive limit
			if recursionDepth > MaxRecursionDepth {
				return nil, util.WrapFilename(originCache, ErrMaximumRecursion)
			}

			// Push closure scope into stack
			env.PushTempEnv(envCache)
			env.PushScope()

			// Declare args
			for idx, arg := range args {
				// Cannot raise error because block is in own scope
				env.Declare(stmt.Params[idx], arg)
			}

			err = ExecuteStatements(stmt.Then.Statements)
			env.PopScope()
			env.PopTempEnv()
			if e, ok := err.(ConditionalError); ok {
				return e.Value, nil
			}

			// Add to callstack
			if err != nil {
				env.FailCall(stmt.Name, originCache, stmt.Line)
			}

			return nil, util.WrapFilename(originCache, err)
		},
	}

	err = env.Declare(stmt.Name, &function)
	// Set after function is declared to allow using the function inside its body
	envCache = env.GetCurrentEnv()
	return err
}

func execIf(stmt Statement) (err error) {
	val, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	if val != nil && val != false {
		return ExecuteStatements(stmt.Then.Statements)
	} else if stmt.Else != nil {
		return ExecuteStatements(stmt.Else.Statements)
	}

	return err
}

func loopStatements(stmts []Statement) (brk bool, err error) {
	env.PushScope()
	err = ExecuteStatements(stmts)
	env.PopScope()
	if e, ok := err.(ConditionalError); ok {
		switch e.Type {
		case BREAK:
			return true, nil
		case SKIP:
			return false, nil
		}
	}

	return false, err
}

// Runs block if expression is nil too
func execWhile(stmt Statement) (err error) {
	for {
		if stmt.Expression != nil {
			val, err := expr.EvaluateExpression(stmt.Expression)
			if err != nil {
				return err
			}

			// Only falsy values for expression
			if val == nil || val == false {
				break
			}
		}

		brk, err := loopStatements(stmt.Then.Statements)
		if err != nil {
			return err
		}

		if brk {
			break
		}
	}

	return err
}

func execRepeat(stmt Statement) (err error) {
	v, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	r, ok := util.IsInt(v)
	if !ok {
		return ErrExpectedInteger
	}

	// Loop n times
	for i := 0; i < r; i++ {
		brk, err := loopStatements(stmt.Then.Statements)
		if err != nil {
			return err
		}

		if brk {
			break
		}
	}

	return err
}

func execObject(stmt Statement) (err error) {
	err = env.Declare(stmt.Name, &env.Callable{
		NumArgs: len(stmt.Params),
		Call: func(args ...interface{}) (interface{}, error) {
			obj := env.Object{Fields: map[string]interface{}{}, Name: stmt.Name}
			for i, field := range stmt.Params {
				obj.Fields[field] = args[i]
			}

			return &obj, err
		},
	})

	return err
}

func getRangeable(args ...expr.Expression) (v *env.Array, err error) {
	if len(args) > 3 {
		return v, ErrInvalidStatement
	}

	// If array just return the array as the list of values to loop over
	if len(args) == 1 {
		val, err := expr.EvaluateExpression(&args[0])
		if err != nil {
			return v, err
		}

		if arr, ok := val.(*env.Array); ok {
			return arr, err
		}
	}

	// Else create array of numbers in range
	a := []float64{}
	for _, e := range args {
		val, err := expr.EvaluateExpression(&e)
		if err != nil {
			return v, err
		}

		if num, ok := val.(float64); ok {
			a = append(a, num)
			continue
		}

		return v, ErrExpectedNumber
	}

	// Set each parameter based on how many there are
	nums := [3]float64{}
	switch len(a) {
	case 1:
		nums[0] = 0
		nums[1] = a[0]
		nums[2] = 1
	case 2:
		nums[0] = a[0]
		nums[1] = a[1]
		nums[2] = 1
	case 3:
		copy(nums[:], a)
	}

	// Check for infinite loop
	zero := nums[2] == 0
	negative := nums[1] > 0 && nums[2] <= 0
	if zero || negative {
		return v, ErrInfiniteLoop
	}

	// Create array
	arr := env.Array{Values: []interface{}{}}
	for i := nums[0]; i < nums[1]; i += nums[2] {
		arr.Values = append(arr.Values, i)
	}

	return &arr, err
}

func execRange(stmt Statement) (err error) {
	var rangeable *env.Array
	e := stmt.Expression
	name := stmt.Name

	// Set rangeable
	if e.Type == expr.Args {
		rangeable, err = getRangeable(e.Exprs...)
	} else {
		rangeable, err = getRangeable(*e)
	}

	if err != nil {
		return err
	}

	env.PushScope()
	env.Declare(name, 0.0)
	for _, val := range rangeable.Values {
		env.Assign(name, val)
		brk, err := loopStatements(stmt.Then.Statements)
		if err != nil {
			return err
		}

		if brk {
			break
		}
	}

	env.PopScope()
	return err
}
