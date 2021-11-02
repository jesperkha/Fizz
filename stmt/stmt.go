package stmt

import (
	"errors"

	"github.com/jesperkha/Fizz/expr"
)

var (
	ErrNoSemicolon 		  = errors.New("expected ; after statement, line %d")
	ErrInvalidStmtType    = errors.New("execute: invalid statement type")
	ErrExpectedExpression = errors.New("expected expression in statement, line %d")
	ErrNoStatement		  = errors.New("expected statement before semicolon, line %d")
	ErrExpectedIdentifier = errors.New("expected identifier at variable declaration, line %d")
	ErrInvalidStatement	  = errors.New("invalid statement, line %d")
	ErrNoBrace			  = errors.New("expected } after block statement, line %d")
)

const (
	// Statement types
	ExpressionStmt = iota
	Print
	Variable
	Assignment
	Block
)

type Statement struct {
	Type 	  	   int
	Line	   	   int
	Name	  	   string
	Expression 	   *expr.Expression
	InitExpression *expr.Expression
	Statements	   []Statement
}