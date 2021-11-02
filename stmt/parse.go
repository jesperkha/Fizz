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
		startIndex  := currentIdx
		firstToken  := tokens[currentIdx]
		currentLine := firstToken.Line

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
			return statements, fmt.Errorf(ErrNoSemicolon.Error(), currentLine)
		}

		// Get tokens in interval between last semicolon and current one
		var currentStmt Statement
		tokenInterval := tokens[startIndex:currentIdx]
		if len(tokenInterval) == 0 {
			return statements, fmt.Errorf(ErrNoStatement.Error(), currentLine)
		}
		
		if parseFunc, ok := parseTable[firstToken.Type]; ok {
			currentStmt, err = parseFunc(tokenInterval)
		} else {
			// Default to parsing expression statement
			currentStmt, err = parseExpression(tokenInterval)
		}

		if err != nil {
			return statements, fmt.Errorf(err.Error(), currentLine)
		}

		currentStmt.Line = currentLine
		statements = append(statements, currentStmt)
		currentIdx++
	}

	return statements, err
}

func parseExpression(tokens []lexer.Token) (stmt Statement, err error) {
	expr, err := expr.ParseExpression(tokens)
	return Statement{Type: ExpressionStmt, Expression: &expr}, err
}

var parseTable = map[int]func(tokens []lexer.Token) (stmt Statement, err error) {
	lexer.PRINT: parsePrint,
	lexer.VAR: 	 parseVariable,
}

// Parses print statement followed by expression
func parsePrint(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) == 1 {
		return stmt, ErrExpectedExpression
	}

	expr, err := expr.ParseExpression(tokens[1:])
	return Statement{Type: Print, Expression: &expr}, err
}

// Parses variable delcaration statement with either nil value init or expression
func parseVariable(tokens []lexer.Token) (stmt Statement, err error) {
	numTokens := len(tokens)
	if numTokens == 1 {
		return stmt, ErrExpectedExpression
	}

	if numTokens == 2 {
		name := tokens[1]
		if name.Type == lexer.IDENTIFIER {
			return Statement{Type: Variable}, err
		}

		return stmt, ErrExpectedIdentifier
	}

	if numTokens >= 4 {
		name := tokens[1]
		exprTokens := tokens[3:]
		equals := tokens[2].Type == lexer.EQUAL

		if name.Type == lexer.IDENTIFIER && equals {
			initExpr, err := expr.ParseExpression(exprTokens)
			return Statement{Type: Variable, Name: name.Lexeme, InitExpression: &initExpr}, err
		}
	}

	return stmt, ErrInvalidStatement
}