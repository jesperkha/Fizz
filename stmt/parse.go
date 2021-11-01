package stmt

import (
	"fmt"

	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
)

// Parses lexer tokens into list of statements
func ParseStatements(tokens []lexer.Token) (statements []Statement, err error) {
	currentIdx := 0

	for currentIdx < len(tokens) {
		first := tokens[currentIdx]
		startIndex := currentIdx
		currentStmt := Statement{}

		// Todo swap to using map for funcs
		// Seeks a semicolon since all statements end with a semicolon
		foundSemicolon := false
		for i := startIndex; i < len(tokens); i++ {
			if tokens[i].Type == lexer.SEMICOLON {
				currentIdx = i
				foundSemicolon = true
				break
			}
		}

		if !foundSemicolon {
			return statements, fmt.Errorf(ErrNoSemicolon.Error(), first.Line)
		}

		tokenInterval := tokens[startIndex + 1:currentIdx]

		switch first.Type {
		case lexer.PRINT: // Print statement
			expr, err := expr.ParseExpression(tokenInterval)
			if err != nil {
				return statements, err
			}

			currentStmt.Type = Print
			currentStmt.Expression = &expr
		
		default: // Default to parsing expression statement
			expr, err := expr.ParseExpression(tokenInterval)
			if err != nil {
				return statements, err
			}
	
			currentStmt.Type = ExpressionStmt
			currentStmt.Expression = &expr
		}

		statements = append(statements, currentStmt)
		currentIdx++
	}

	return statements, err
}

var parseTable = map[int]func(tokens []lexer.Token, curIdx *int) (stmt Statement, err error) {
	Print: parsePrint,
	ExpressionStmt: parseExpression,
}

func parsePrint(tokens []lexer.Token, curIdx *int) (stmt Statement, err error) {

	return stmt, err
}

func parseExpression(tokens []lexer.Token, curIdx *int) (stmt Statement, err error) {

	return stmt, err
}