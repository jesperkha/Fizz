package stmt

import (
	"fmt"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

var (
	CurrentOrigin string
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
		return execVariable(stmt)
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
	case Object:
		return execObject(stmt)
	case Import:
		return nil // Handled in interp
	}

	// Will never be returned since all types are pre-defined.
	// However it is nice to have in case rework is done and types
	// get mixed up or new types are only partially added.
	return ErrInvalidStmtType
}

func execExit(stmt Statement) (err error) {
	if stmt.Expression != nil {
		if err = execPrint(stmt); err != nil {
			return err
		}
	}

	return ErrProgramExit
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

func formatPrintValue(val interface{}) interface{} {
	switch val.(type) {
	case float64, string, bool:
		return val
	case nil:
		return "nil"
	}

	if o, ok := val.(env.Object); ok {
		str := "{\n"
		for key, value := range o.Fields {
			str += fmt.Sprintf("    %s: %v\n", key, formatPrintValue(value))
		}
	
		return str + "}"
	}

	if o, ok := val.(env.Callable); ok {
		return o.Name + "()"
	}

	if o, ok := val.(env.Array); ok {
		str := "["
		for i, v := range o.Values {
			if i != 0 {
				str += ", "
			}

			str += fmt.Sprintf("%v", formatPrintValue(v))
		}

		return str + "]"
	}

	return ""
}

func execPrint(stmt Statement) (err error) {
	value, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	fmt.Println(formatPrintValue(value))
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

// Helper for execAssignment
func assignToObject(objTokens []lexer.Token, name string, value interface{}) (err error) {
	// Not object value assignment
	if len(objTokens) == 0 {
		return env.Assign(name, value)
	}

	// Get object to assign to
	v, err := expr.ParseAndEval(objTokens[:len(objTokens)-2]) // exclude name
	if err != nil {
		return err
	}

	if obj, ok := v.(env.Object); ok {
		return obj.Set(name, value)
	}
	
	return ErrNonAssignable
}

func getObjectValue(stmt Statement) (value interface{}, err error) {
	if len(stmt.ObjTokens) == 0 {
		return env.Get(stmt.Name)
	}

	return expr.ParseAndEval(stmt.ObjTokens)
}

func execAssignment(stmt Statement) (err error) {
	val, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	// Plain assignment
	if stmt.Operator == lexer.EQUAL {
		return assignToObject(stmt.ObjTokens, stmt.Name, val)
	}

	// Declare variable with special := operator
	if stmt.Operator == lexer.DEF_EQUAL {
		return execVariable(stmt)
	}

	oldVal, err := getObjectValue(stmt)
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

		return assignToObject(stmt.ObjTokens, stmt.Name, oldVal.(string)+val.(string))
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

		return assignToObject(stmt.ObjTokens, stmt.Name, newVal)
	}

	return ErrInvalidStatement
}

func execBlock(stmt Statement) (err error) {
	env.PushScope()
	err = ExecuteStatements(stmt.Statements)
	env.PopScope()
	return err
}

// Todo: implement callstack (add recursion limit when doing so)
// when the interp model is redone to fit a more object/instance
// approach the callstack should be implemented
func execFunction(stmt Statement) (err error) {
	// Store origin at point of function declaration as well as scope around it
	originCache := CurrentOrigin
	envCache := env.GetCurrentEnv()

	function := env.Callable{
		Name: stmt.Name,
		NumArgs: len(stmt.Params),
		Origin: CurrentOrigin,
		// Call function and set param variables to scope
		Call: func(args ...interface{}) (interface{}, error) {
			env.PushTempEnv(envCache)
			env.PushScope()

			// fmt.Println(envCache)

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
			
			return nil, util.WrapFilename(originCache, err)
		},
	}

	return env.Declare(stmt.Name, function)
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
		if e, ok := err.(ConditionalError); ok {
			switch e.Type {
			case BREAK:
				return nil
			case SKIP:
				continue
			}
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
		err = ExecuteStatements(stmt.Then.Statements)
		if e, ok := err.(ConditionalError); ok {
			switch e.Type {
			case BREAK:
				return nil
			case SKIP:
				if oldVal, err := env.Get(name); err == nil {
					env.Assign(name, oldVal.(float64)+1)
				}
				continue
			}
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
			obj := env.Object{Fields: map[string]interface{}{}, Name: stmt.Name}
			for i, field := range stmt.Params {
				obj.Fields[field] = args[i]
			}

			return obj, err
		},
	})

	return err
}
