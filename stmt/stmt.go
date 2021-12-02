package stmt

import (
	"errors"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/expr"
)

var (
	ErrNoSemicolon        = errors.New("expected ; after statement, line %d")
	ErrInvalidStmtType    = errors.New("invalid statement type, check statement parsing")
	ErrExpectedExpression = errors.New("expected expression in statement, line %d")
	ErrNoStatement        = errors.New("expected statement before semicolon, line %d")
	ErrExpectedIdentifier = errors.New("expected identifier, line %d")
	ErrInvalidStatement   = errors.New("invalid statement, line %d")
	ErrNoBrace            = errors.New("expected } after block statement, line %d")
	ErrExpectedBlock      = errors.New("expected block after statememt, line %d")
	ErrExpectedIf         = errors.New("expected if statement before else, line %d")
	ErrInvalidOperator    = errors.New("invalid statement operator, line %d")
	ErrDifferentTypes     = errors.New("different types in statement, line %d")
	ErrNonCallable        = errors.New("cannot call non-callable type, line %d")
	ErrCommaError         = errors.New("comma error, line %d")

	ErrReturnOutsideFunc = errors.New("cannot use return outside of a function")
	ErrSkipOutsideLoop   = errors.New("cannot use skip outside of a loop")
	ErrBeakOutsideLoop   = errors.New("cannot use break outside of a loop")
	ErrProgramExit       = errors.New("")
)

const (
	NotStatement = iota
	ExpressionStmt
	Print
	Variable
	Assignment
	Block
	If
	While
	Repeat
	Break
	Skip
	Function
	Return
	Exit
	Object
)

type Statement struct {
	Type       int
	Line       int
	Operator   int
	Name       string
	Expression *expr.Expression
	Statements []Statement
	Then       *Statement
	Else       *Statement
	Params     []string
	Enviroment env.Environment
}
