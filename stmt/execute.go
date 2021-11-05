package stmt

import (
	"fmt"
	"reflect"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
)

// Goes through list of statements and executes them. Error is returned from statements exec method.
func ExecuteStatements(stmts []Statement) (err error) {
	currentIdx := 0 // For changing dynamically

	for currentIdx < len(stmts) {
		statement := stmts[currentIdx]
		line := statement.Line
		currentIdx++

		// Ignore expression statements
		if statement.Type == ExpressionStmt {
			continue
		}

		// Check conditional statements
		if execFunc, ok := execConTable[statement.Type]; ok {
			if err = execFunc(statement, &currentIdx); err != nil {
				return formatError(err, line)
			}

			continue
		}
		
		// Check block sepeartly because go maps are gay
		if statement.Type == Block {
			if err = execBlock(statement); err != nil {
				return formatError(err, line)
			}

			continue
		}

		// Check ramining statement types
		if execFunc, ok := execStatementTable[statement.Type]; ok {
			if err = execFunc(statement); err != nil {
				return formatError(err, line)
			}

			continue
		}
		
		// Will never be returned since all types are pre-defined.
		// However it is nice to have in case reword is done and types
		// get mixed up or new types are only partially added.
		return ErrInvalidStmtType
	}

	return err
}

type execTable map[int]func(stmt Statement) error 

var execConTable = map[int]func(stmt Statement, idx *int) error {}

var execStatementTable = execTable {
	Print: 	    execPrint,
	Variable:   execVariable,
	Assignment: execAssignment,
	// Break:		execBreak,
}

// Do nothing. Handled in loop exec functions
// func execBreak(stmt Statement) (err error) {
// 	return ErrBreakStatement{}
// }

// Evaluates statement expression and prints out to terminal
func execPrint(stmt Statement) (err error) {
	value, err := expr.EvaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

// Adds variable init value to current environment table
func execVariable(stmt Statement) (err error) {
	if stmt.InitExpression != nil {
		val, err := expr.EvaluateExpression(stmt.InitExpression)
		if err != nil {
			return err
		}

		return env.Declare(stmt.Name, val)
	}

	return env.Declare(stmt.Name, nil)
}

// Assigns right side expression value to variable
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
	if reflect.TypeOf(oldVal) != reflect.TypeOf(val) {
		return ErrInvalidStatement
	}

	// String addition
	if reflect.TypeOf(val) == reflect.TypeOf("") {
		if stmt.Operator == lexer.MINUS_EQUAL {
			return ErrInvalidOperator
		} 

		return env.Assign(stmt.Name, oldVal.(string) + val.(string))
	}

	// Float addition / subtraction
	if reflect.TypeOf(val) == reflect.TypeOf(float64(1)) {
		a := val.(float64)
		b := oldVal.(float64)
		if stmt.Operator == lexer.MINUS_EQUAL {
			a *= -1
		}
	
		return env.Assign(stmt.Name, a + b)
	}

	return ErrDifferentTypes
}

// Executes all statements within block scope
func execBlock(stmt Statement) (err error) {
	env.PushScope()
	err = ExecuteStatements(stmt.Statements)
	env.PopScope()
	return err
}

// Skips trailing block statement if expression is false
func execIf(stmt Statement, idx *int) (err error) {
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

// Runs block if expression is true or no expression
func execWhile(stmt Statement, idx *int) (err error) {
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
		// Todo Handle break here
		if err != nil {
			return err
		}
	}

	return err
}