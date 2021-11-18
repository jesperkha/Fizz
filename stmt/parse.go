package stmt

import (
	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

// Parses lexer tokens into list of statements
func ParseStatements(tokens []lexer.Token) (statements []Statement, err error) {
	currentIdx := 0

	for currentIdx < len(tokens) {
		startIndex  := currentIdx
		firstToken  := tokens[currentIdx]
		
		var currentStmt Statement
		currentLine := firstToken.Line

		// Check conditional statements seperatly because the parse funcs need
		// a currentIndex pointer
		if parseFunc, ok := pconTable[firstToken.Type]; ok {
			currentStmt, err = parseFunc(tokens, &currentIdx)
			if err != nil {
				return statements, util.FormatError(err, currentLine)
			}
			
			currentStmt.Line = currentLine
			statements = append(statements, currentStmt)
			currentIdx++
			continue
		}

		// Finally parse remainig statements.
		// Seeks a semicolon since all other statements end with a semicolon
		endIdx, eof := seekToken(tokens, startIndex, lexer.SEMICOLON)
		if eof {
			return statements, util.FormatError(ErrNoSemicolon, currentLine)
		}

		currentIdx = endIdx

		// Get tokens in interval between last semicolon and current one
		tokenInterval := tokens[startIndex:currentIdx]
		if len(tokenInterval) == 0 {
			return statements, util.FormatError(err, currentLine)
		}
		
		if parseFunc, ok := parseStatementTable[firstToken.Type]; ok {
			currentStmt, err = parseFunc(tokenInterval)
		} else {
			// Default to parsing expression statement
			currentStmt, err = parseExpression(tokenInterval)
		}

		if err != nil {
			return statements, util.FormatError(err, currentLine)
		}

		currentStmt.Line = currentLine
		statements = append(statements, currentStmt)
		currentIdx++
	}

	return statements, err
}

// Seeks target token type and returns the index of target and eof
func seekToken(tokens []lexer.Token, start int, target int) (endIdx int, eof bool) {
	for i := start; i < len(tokens); i++ {
		if tokens[i].Type == target {
			return i, false
		}
	}

	return 0, true
}

type parseTable map[int]func(tokens []lexer.Token) (stmt Statement, err error)

// For the functions in this map the tokens are passed as the complete list
// This is to ensure that the result of seekToken mathes up with the current
// index value.
var pconTable = map[int]func([]lexer.Token, *int)(Statement, error){}

var parseStatementTable = parseTable {
	lexer.IDENTIFIER: parseAssignment,
	lexer.PRINT: 	  parsePrint,
	lexer.VAR: 	 	  parseVariable,
	lexer.ELSE:		  parseElse,
	lexer.BREAK: 	  parseBreak,
	lexer.SKIP: 	  parseSkip,
}

// Just returns the statement with a type. Implementation is handled in exec.
func parseBreak(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) > 1 {
		return stmt, ErrInvalidStatement
	}

	return Statement{Type: Break}, err
}

// Same as break
func parseSkip(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) > 1 {
		return stmt, ErrInvalidStatement
	}

	return Statement{Type: Skip}, err
}

// Does literally nothing lol
func parseExpression(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) == 0 {
		return Statement{Type: ExpressionStmt}, err
	}

	expr, err := expr.ParseExpression(tokens)
	return Statement{Type: ExpressionStmt, Expression: &expr}, err
}

// Else statements are consumed by the if parser so if one is found its an error
func parseElse(tokens []lexer.Token) (stmt Statement, err error) {
	return stmt, ErrExpectedIf
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
	if len(tokens) < 3 {
		return parseExpression(tokens)
	}

	t := tokens[1].Type
	if t != lexer.EQUAL && t != lexer.PLUS_EQUAL && t != lexer.MINUS_EQUAL {
		return parseExpression(tokens)
	}

	name := tokens[0].Lexeme
	rightExpr := tokens[2:]
	expr, err := expr.ParseExpression(rightExpr)
	return Statement{Type: Assignment, Name: name, Expression: &expr, Operator: t}, err
}

// Gets a trailing block statement
func getBlockStatement(tokens []lexer.Token, idx *int) (block Statement, err error) {
	start := *idx
	if tokens[start].Type != lexer.LEFT_BRACE {
		return block, ErrExpectedBlock
	}

	numEndBraces := 0
	foundEndBrace := false

	// Loop over until finds brace ending a nested block
	for *idx < len(tokens) {
		switch tokens[*idx].Type {
			case lexer.LEFT_BRACE: numEndBraces++
			case lexer.RIGHT_BRACE: numEndBraces--
		}
		
		if numEndBraces == 0 {
			foundEndBrace = true
			break
		}

		*idx++
	}

	if !foundEndBrace {
		return block, ErrNoBrace
	}

	blockTokens := tokens[start + 1:*idx]
	statements, err := ParseStatements(blockTokens)
	return Statement{Type: Block, Statements: statements}, err
}

// Parses all statements within block
func parseBlock(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	return getBlockStatement(tokens, idx)
}

// Gets statements between keyword and block, also gets block
func getStatementAndBlock(tokens []lexer.Token, idx *int) (stmt Statement, block Statement, err error) {
	startBlock, eof := seekToken(tokens, *idx, lexer.LEFT_BRACE)
	if eof {
		return stmt, block, ErrExpectedBlock
	}

	stmt, err = parseExpression(tokens[*idx + 1:startBlock])
	if err != nil {
		return stmt, block, err
	}

	*idx = startBlock
	block, err = getBlockStatement(tokens, idx)
	return stmt, block, err
}

// Finds trailing block and parses expression between block and if token
// as well as the block. Adds else statement if found
func parseIf(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	stmt, block, err := getStatementAndBlock(tokens, idx)
	if err != nil {
		return stmt, err
	}
	
	if stmt.Expression == nil {
		return stmt, ErrExpectedExpression
	}
	
	// Check for else statement
	if *idx + 1 < len(tokens) && tokens[*idx + 1].Type == lexer.ELSE {
		if *idx + 2 >= len(tokens) {
			return stmt, ErrExpectedBlock
		}

		*idx += 2 // Skip to block
		elseBlock, err := getBlockStatement(tokens, idx)
		return Statement{Type: If, Expression: stmt.Expression, Then: &block, Else: &elseBlock}, err
	}

	return Statement{Type: If, Expression: stmt.Expression, Then: &block}, err
}

// Parses while with expression and block. No expression means always true
func parseWhile(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	stmt, block, err := getStatementAndBlock(tokens, idx)
	if err != nil {
		return stmt, err
	}

	if stmt.Expression == nil {
		return Statement{Type: While, Then: &block}, err
	}

	if stmt.Type == ExpressionStmt {
		return Statement{Type: While, Expression: stmt.Expression, Then: &block}, err
	}

	return stmt, ErrExpectedExpression
}

// Parses repeat loop expression and checks if it is correct
func parseRepeat(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	stmt, block, err := getStatementAndBlock(tokens, idx)
	if stmt.Expression == nil {
		return stmt, ErrExpectedExpression
	}

	if stmt.Expression.Type != expr.Binary {
		return stmt, ErrInvalidStatement
	}

	notLess := stmt.Expression.Operand.Type != lexer.LESS
	notIdentifier := stmt.Expression.Left.Type != expr.Variable
	if notLess || notIdentifier {
		return stmt, ErrInvalidStatement
	}

	return Statement{Type: Repeat, Expression: stmt.Expression, Then: &block}, err
}