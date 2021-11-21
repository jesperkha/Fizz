package stmt

import (
	"errors"

	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
)

var (
	ErrNoSemicolon 		  = errors.New("expected ; after statement, line %d")
	ErrInvalidStmtType    = errors.New("invalid statement type, check statement parsing")
	ErrExpectedExpression = errors.New("expected expression in statement, line %d")
	ErrNoStatement		  = errors.New("expected statement before semicolon, line %d")
	ErrExpectedIdentifier = errors.New("expected identifier, line %d")
	ErrInvalidStatement	  = errors.New("invalid statement, line %d")
	ErrNoBrace			  = errors.New("expected } after block statement, line %d")
	ErrExpectedBlock	  = errors.New("expected block after statememt, line %d")
	ErrExpectedIf		  = errors.New("expected if statement before else, line %d")
	ErrInvalidOperator	  = errors.New("invalid statement operator, line %d")
	ErrDifferentTypes	  = errors.New("different types in statement, line %d")
	ErrBeakOutsideLoop    = errors.New("cannot use break outside of a loop")
	ErrSkipOutsideLoop    = errors.New("cannot use skip outside of a loop")
	ErrNonCallable		  = errors.New("cannot call non-callable type, line %d")
)

const (
	ExpressionStmt = iota
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
)

type Statement struct {
	Type 	  	   int
	Line	   	   int
	Operator	   int
	Name	  	   string
	Expression 	   *expr.Expression
	InitExpression *expr.Expression
	Statements	   []Statement
	Then		   *Statement
	Else		   *Statement
	Params	   	   []string
}

func init() {
	// Recursive functions set at init because go is weird
	// Todo fix this shit pls
	execConTable[If] = execIf
	execConTable[While] = execWhile
	execConTable[Repeat] = execRepeat
	execConTable[Function] = execFunction

	pconTable[lexer.IF] = parseIf
	pconTable[lexer.WHILE] = parseWhile
	pconTable[lexer.LEFT_BRACE] = parseBlock
	pconTable[lexer.REPEAT] = parseRepeat
	pconTable[lexer.FUNC] = parseFunc
}
