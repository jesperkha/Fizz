package stmt

import (
	"errors"

	"github.com/jesperkha/Fizz/expr"
)

var (
	ErrNoSemicolon = errors.New("expected ; to end expression, line %d")
	ErrInvalidStmtType = errors.New("internal: invalid statement type") // Will never happen
)

const (
	// Statement types
	ExpressionStmt = iota
	Print
	Variable
)

type Statement struct {
	Type 	   int
	Expression *expr.Expression
}