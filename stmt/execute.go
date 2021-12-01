package stmt

import (
	"errors"
	"fmt"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

var currentReturnValue interface{}

// Goes through list of statements and executes them. Error is returned from statements exec method.
func ExecuteStatements(stmts []Statement) (err error) {
	for _, statement := range stmts {
		line := statement.Line
		if err = executeStatement(statement); err != nil {
			return util.FormatError(err, line)
		}
	}

	return err
}

func executeStatement(stmt Statement) error {
	switch stmt.Type {
	case ExpressionStmt:
		return execExpression(stmt)
	case Block:
		return execBlock(stmt)
	case Print:
		return execPrint(stmt)
	case Variable:
		return execVariable(stmt)
	case Assignment:
		return execAssignment(stmt)
	case Break:
		return execBreak(stmt)
	case Skip:
		return execSkip(stmt)
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
	case Object:
		return execObject(stmt)
	}

	// Will never be returned since all types are pre-defined.
	// However it is nice to have in case rework is done and types
	// get mixed up or new types are only partially added.
	return ErrInvalidStmtType
}

func execExit(stmt Statement) (err error) {
	return ErrProgramExit
}

// Raises error and assigns expr value to global currentReturnValue
func execReturn(stmt Statement) (err error) {
	if stmt.Expression == nil {
		currentReturnValue = nil
		return ErrReturnOutsideFunc
	}

	value, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	currentReturnValue = value
	return ErrReturnOutsideFunc
}

func execExpression(stmt Statement) (err error) {
	_, err = expr.EvaluateExpression(stmt.Expression)
	return err
}

func execBreak(stmt Statement) (err error) {
	return ErrBeakOutsideLoop
}

func execSkip(stmt Statement) (err error) {
	return ErrSkipOutsideLoop
}

func execPrint(stmt Statement) (err error) {
	value, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func execVariable(stmt Statement) (err error) {
	if stmt.Expression != nil {
		val, err := expr.EvaluateExpression(stmt.Expression)
		if err != nil {
			return err
		}

		return env.Declare(stmt.Name, val)
	}

	return env.Declare(stmt.Name, nil)
}

func execAssignment(stmt Statement) (err error) {
	val, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	// Plain assignment
	if stmt.Operator == lexer.EQUAL {
		return env.Assign(stmt.Name, val)
	}

	// Plus equals
	oldVal, err := env.Get(stmt.Name)
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

		return env.Assign(stmt.Name, oldVal.(string)+val.(string))
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

		return env.Assign(stmt.Name, newVal)
	}

	return ErrInvalidStatement
}

func execBlock(stmt Statement) (err error) {
	env.PushScope()
	err = ExecuteStatements(stmt.Statements)
	env.PopScope()
	return err
}

// Todo: add max recursion limit
func execFunction(stmt Statement) (err error) {
	err = env.Declare(stmt.Name, env.Callable{
		NumArgs: len(stmt.Params),

		// Call function and set param variables to scope
		Call: func(args ...interface{}) (interface{}, error) {
			env.PushScope()

			// Declare args
			for idx, arg := range args {
				// Cannot raise error because block is in own scope
				env.Declare(stmt.Params[idx], arg)
			}

			err = ExecuteStatements(stmt.Then.Statements)
			env.PopScope()
			if errors.Is(err, ErrReturnOutsideFunc) {
				return currentReturnValue, nil
			}

			return nil, err
		},
	})

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

// Runs block if expression is nil too
func execWhile(stmt Statement) (err error) {
	for {
		if stmt.Expression != nil {
			val, err := expr.EvaluateExpression(stmt.Expression)
			if err != nil {
				return err
			}

			if val == nil || val == false {
				break
			}
		}

		err = ExecuteStatements(stmt.Then.Statements)
		if errors.Is(err, ErrBeakOutsideLoop) {
			return nil
		}

		if errors.Is(err, ErrSkipOutsideLoop) {
			continue
		}

		if err != nil {
			return err
		}
	}

	return err
}

func execRepeat(stmt Statement) (err error) {
	name := stmt.Expression.Left.Name
	// Push new scope to avoid clashing when defining new variable
	// Block is in child scope anyway
	env.PushScope()
	env.Declare(name, float64(0))

	for {
		val, err := expr.EvaluateExpression(stmt.Expression)
		if err != nil {
			return err
		}

		if val == false { // Explicit check for false because interface, stfu
			break
		}

		// Break and skip return errors that are handled here
		// Todo: add line number for break, skip, and return 'errors'
		err = ExecuteStatements(stmt.Then.Statements)
		if errors.Is(err, ErrBeakOutsideLoop) {
			return nil
		}

		if err == nil || errors.Is(err, ErrSkipOutsideLoop) {
			if oldVal, err := env.Get(name); err == nil {
				env.Assign(name, oldVal.(float64)+1)
			}

			continue
		}

		return err
	}

	env.PopScope()
	return err
}

func execObject(stmt Statement) (err error) {
	err = env.Declare(stmt.Name, env.Callable{
		NumArgs: len(stmt.Params),

		Call: func(args ...interface{}) (interface{}, error) {
			obj := env.Object{Fields: map[string]interface{}{}}
			for i, field := range stmt.Params {
				obj.Fields[field] = args[i]
			}

			return obj, err
		},
	})

	return err
}