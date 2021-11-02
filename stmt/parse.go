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
		
		// Check and parse block statement first to not seek semicolon
		if firstToken.Type == lexer.LEFT_BRACE {
			numEndBraces := 0
			foundEndBrace := false

			// Loop over until finds brace ending a nested block
			for currentIdx < len(tokens) {
				switch tokens[currentIdx].Type {
					case lexer.LEFT_BRACE: numEndBraces++
					case lexer.RIGHT_BRACE: numEndBraces--
				}
				
				if numEndBraces == 0 {
					foundEndBrace = true
					break
				}

				currentIdx++
			}

			if !foundEndBrace {
				return statements, fmt.Errorf(ErrNoBrace.Error(), currentLine)
			}
			
			blockTokens := tokens[startIndex + 1:currentIdx]
			stmt, err := parseBlock(blockTokens)
			if err != nil {
				return statements, err // already formatted from this function
			}
			
			stmt.Line = currentLine
			statements = append(statements, stmt)
			currentIdx++
			continue
		}
		
		// Seeks a semicolon since all other statements end with a semicolon
		endIdx, eof := seekToken(tokens, startIndex, lexer.SEMICOLON)
		if eof {
			return statements, fmt.Errorf(ErrNoSemicolon.Error(), currentLine)
		}

		currentIdx = endIdx

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

// Seeks target token type and returns the end index and eof
func seekToken(tokens []lexer.Token, start int, target int) (endIdx int, eof bool) {
	for i := start; i < len(tokens); i++ {
		if tokens[i].Type == target {
			return i, false
		}
	}

	return 0, true
} 

func parseExpression(tokens []lexer.Token) (stmt Statement, err error) {
	expr, err := expr.ParseExpression(tokens)
	return Statement{Type: ExpressionStmt, Expression: &expr}, err
}

var parseTable = map[int]func(tokens []lexer.Token) (stmt Statement, err error) {
	lexer.PRINT: 	  parsePrint,
	lexer.VAR: 	 	  parseVariable,
	lexer.IDENTIFIER: parseAssignment,
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

// Assigns right hand expression value to variable name
func parseAssignment(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) < 3 || tokens[1].Type != lexer.EQUAL {
		return stmt, ErrInvalidStatement
	}

	name := tokens[0].Lexeme
	rightExpr := tokens[2:]
	expr, err := expr.ParseExpression(rightExpr)
	return Statement{Type: Assignment, Name: name, Expression: &expr}, err
}

// Parses all statements within block
func parseBlock(tokens []lexer.Token) (stmt Statement, err error) {
	statements, err := ParseStatements(tokens)
	if err != nil {
		return stmt, err
	}

	return Statement{Type: Block, Statements: statements}, err
}