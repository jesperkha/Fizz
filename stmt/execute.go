package stmt

import (
	"fmt"

	"github.com/jesperkha/Fizz/expr"
)

// Goes through list of statements and executes them. Error is returned from statements exec method.
func ExecuteStatements(stmts []Statement) (err error) {
	currentIdx := 0 // For changing dynamically

	for currentIdx < len(stmts) {
		statement := stmts[currentIdx]

		if execMethod, ok := statementTable[statement.Type]; ok {
			err = execMethod(statement)
			if err != nil {
				return err
			}
			
			currentIdx++
			continue
		}

		// Will never be returned since all types are pre-defined.
		// However it is nice to have in case reword is done and types
		// get mixed up or new types are only partially added.
		return ErrInvalidStmtType
	}

	return err
}

var statementTable = map[int]func(stmt Statement) error {
	Print: execPrint,
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

