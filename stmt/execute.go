package stmt

import (
	"fmt"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/expr"
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
}

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

	return env.Assign(stmt.Name, val)
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
	}
	
	return err
}
