package stmt

import (
	"errors"

	"github.com/jesperkha/Fizz/expr"
)

var (
	ErrNoSemicolon 		  = errors.New("expected ; to end expression, line %d")
	ErrInvalidStmtType    = errors.New("execute: invalid statement type")
	ErrExpectedExpression = errors.New("expected expression in statement, line %d")
	ErrNoStatement		  = errors.New("expected statement before semicolon, line %d")
	ErrExpectedIdentifier = errors.New("expected identifier at variable declaration, line %d")
	ErrInvalidStatement	  = errors.New("invalid statement, line %d")
)

const (
	// Statement types
	ExpressionStmt = iota
	Print
	Variable
)

type Statement struct {
	Type 	  	   int
	Name	  	   string
	Line	   	   int
	Expression 	   *expr.Expression
	InitExpression *expr.Expression
}