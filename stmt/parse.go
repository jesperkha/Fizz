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
		// a currentIndex pointer. Note: Full list of tokens is given
		currentStmt, err = parseComplexStatement(firstToken.Type, tokens, &currentIdx)
		if err != nil {
			return statements, util.FormatError(err, currentLine)
		}
		
		// Parse any other type of statement. Checks if not statement
		if currentStmt.Type == 0 {
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
	
			// Parse statement
			currentStmt, err = parseStatement(firstToken.Type, tokenInterval)
			if err != nil {
				return statements, util.FormatError(err, currentLine)
			}
		}

		currentStmt.Line = currentLine
		statements = append(statements, currentStmt)
		currentIdx++
	}

	return statements, err
}

// Helper. Returnes the value for the specific statement parse function. Defaults to ExpressionStatement
func parseStatement(typ int, tokens []lexer.Token) (stmt Statement, err error) {
	switch typ {
	case lexer.IDENTIFIER: return parseAssignment(tokens)
	case lexer.PRINT: return parsePrint(tokens)
	case lexer.VAR: return parseVariable(tokens)
	case lexer.ELSE: return parseElse(tokens)
	case lexer.BREAK: return parseBreak(tokens)
	case lexer.SKIP: return parseSkip(tokens)
	}
	
	return parseExpression(tokens)
}

// Parses statements where current index would be modified or the length of the statement is unknown
func parseComplexStatement(typ int, tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	switch typ {
	case lexer.IF: return parseIf(tokens, idx)
	case lexer.WHILE: return parseWhile(tokens, idx)
	case lexer.LEFT_BRACE: return parseBlock(tokens, idx)
	case lexer.REPEAT: return parseRepeat(tokens, idx)
	case lexer.FUNC: return parseFunc(tokens, idx)
	}

	return stmt, err
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

// Parse funcs

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

// Just returns the statement with a type. Implementation is handled in exec.
func parseFunc(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	if len(tokens) < 6 {
		return stmt, ErrInvalidStatement
	}

	nameToken := tokens[*idx + 1]

	hasName := nameToken.Type == lexer.IDENTIFIER
	hasParens := tokens[*idx + 2].Type == lexer.LEFT_PAREN
	if !hasName || !hasParens {
		return stmt, ErrInvalidStatement
	}

	// Get param names
	*idx += 3 // Skip to start of param list
	endIdx, eof := seekToken(tokens, *idx, lexer.RIGHT_PAREN)
	if eof {
		return stmt, ErrInvalidStatement
	}
	
	params := []string{}
	paramTokens := tokens[*idx:endIdx]
	for _, p := range paramTokens {
		if p.Type != lexer.IDENTIFIER {
			return stmt, ErrExpectedIdentifier
		}

		params = append(params, p.Lexeme)
	}

	*idx = endIdx + 1 // Skip to start of block
	if tokens[*idx].Type != lexer.LEFT_BRACE {
		return stmt, ErrExpectedBlock
	}

	block, err := getBlockStatement(tokens, idx)
	if err != nil {
		return stmt, err
	}

	return Statement{Type: Function, Name: nameToken.Lexeme, Params: params, Then: &block}, err
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
