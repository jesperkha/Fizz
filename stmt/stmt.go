package stmt

import (
	"errors"

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
	ErrNonAssignable      = errors.New("cannot assign value to non-subscriptable, line %d")
	ErrExpectedName       = errors.New("expected filename at import, line %d")
	ErrCannotImport       = errors.New("cannot import outside of global scope, line %d")
	ErrExpectedInteger    = errors.New("expected expression to be integer, line %d")
	ErrExpectedNumber     = errors.New("expected expression to be number, line %d")
	ErrInfiniteLoop       = errors.New("infinite loop in range statement not allowed, line %d")
	ErrProgramExit        = errors.New("")

	ErrReturnOutsideFunc = ConditionalError{Msg: "cannot use return outside of a function, line %d", Type: RETURN}
	ErrSkipOutsideLoop   = ConditionalError{Msg: "cannot use skip outside of a loop, line %d", Type: SKIP}
	ErrBeakOutsideLoop   = ConditionalError{Msg: "cannot use break outside of a loop, line %d", Type: BREAK}
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
	Import
	Include
	Error
	Enum
	Range
)

// Todo: make range statement and remove repeat
// Todo: enum statement

type Statement struct {
	Type       int
	Line       int
	Operator   int
	Name       string
	Params     []string
	Statements []Statement
	Then       *Statement
	Else       *Statement
	Expression *expr.Expression
	Left       *expr.Expression
}

const (
	SKIP = iota
	BREAK
	RETURN
)

type ConditionalError struct {
	Type  int
	Line  int
	Msg   string
	Value interface{}
}

func (c ConditionalError) Error() string {
	return c.Msg
}
