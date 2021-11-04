package stmt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
)

var (
	ErrNoSemicolon 		  = errors.New("expected ; after statement, line %d")
	ErrInvalidStmtType    = errors.New("execute: invalid statement type")
	ErrExpectedExpression = errors.New("expected expression in statement, line %d")
	ErrNoStatement		  = errors.New("expected statement before semicolon, line %d")
	ErrExpectedIdentifier = errors.New("expected identifier at variable declaration, line %d")
	ErrInvalidStatement	  = errors.New("invalid statement, line %d")
	ErrNoBrace			  = errors.New("expected } after block statement, line %d")
	ErrExpectedBlock	  = errors.New("expected block after statememt, line %d")
	ErrExpectedIf		  = errors.New("expected if statement before else, line %d")
)

const (
	ExpressionStmt = iota
	Print
	Variable
	Assignment
	Block
	If
	Else
)

type Statement struct {
	Type 	  	   int
	Line	   	   int
	Name	  	   string
	Expression 	   *expr.Expression
	InitExpression *expr.Expression
	Statements	   []Statement
	Then		   *Statement
	Else		   *Statement
}

// Format error with line numbers for local errors, but ignore for errors passed from
// expression parsing as they are already formatted with line numbers.
func formatError(err error, line int) error {
	if err == nil {
		return err
	}

	if strings.Contains(err.Error(), "%d") {
		return fmt.Errorf(err.Error(), line)
	}

	return err
}

func init() {
	execConTable[If] = execIf

	pconTable[lexer.IF] = parseIf
	pconTable[lexer.LEFT_BRACE] = parseBlock
}
