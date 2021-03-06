package stmt

import (
	"strings"

	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

// Parses lexer tokens into list of statements
func ParseStatements(tokens []lexer.Token) (statements []Statement, err error) {
	currentIdx := 0

	for currentIdx < len(tokens) {
		startIndex := currentIdx
		firstToken := tokens[currentIdx]

		var currentStmt Statement
		line := firstToken.Line

		// Check conditional statements seperatly because the parse funcs need
		// a currentIndex pointer. Note: Full list of tokens is given
		currentStmt, err = parseComplexStatement(firstToken.Type, tokens, &currentIdx)
		if err != nil {
			return statements, util.FormatError(err, line)
		}

		// Parse any other type of statement.
		if currentStmt.Type == NotStatement {
			// Seeks a semicolon since all other statements end with a semicolon
			endIdx, eof := seekToken(tokens, startIndex, lexer.SEMICOLON)
			if eof {
				return statements, util.FormatError(ErrNoSemicolon, line)
			}

			currentIdx = endIdx // Skip to end of statement to section off token list

			// Get tokens in interval between last semicolon and current one
			tokenInterval := tokens[startIndex:currentIdx]
			if len(tokenInterval) == 0 {
				return statements, util.FormatError(ErrNoStatement, line)
			}

			// Parse statement
			currentStmt, err = parseStatement(firstToken.Type, tokenInterval)
			if err != nil {
				return statements, util.FormatError(err, line)
			}
		}

		currentStmt.Line = line
		statements = append(statements, currentStmt)
		currentIdx++
	}

	return statements, err
}

// Defaults to ExpressionStatement
func parseStatement(typ int, tokens []lexer.Token) (stmt Statement, err error) {
	switch typ {
	case lexer.IDENTIFIER:
		return parseAssignment(tokens)
	case lexer.PRINT:
		return parsePrint(tokens)
	case lexer.ELSE:
		return parseElse(tokens)
	case lexer.BREAK:
		return parseBreak(tokens)
	case lexer.SKIP:
		return parseSkip(tokens)
	case lexer.RETURN:
		return parseReturn(tokens)
	case lexer.EXIT:
		return parseExit(tokens)
	case lexer.ERROR:
		return parseError(tokens)
	case lexer.IMPORT:
		return parseImport(tokens)
	case lexer.INCLUDE:
		return parseInclude(tokens)
	}

	return parseExpression(tokens)
}

// Parses statements where current index would be modified or the length of the statement is unknown
func parseComplexStatement(typ int, tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	switch typ {
	case lexer.IF:
		return parseIf(tokens, idx)
	case lexer.WHILE:
		return parseWhile(tokens, idx)
	case lexer.LEFT_BRACE:
		return parseBlock(tokens, idx)
	case lexer.REPEAT:
		return parseRepeat(tokens, idx)
	case lexer.FUNC:
		return parseFunc(tokens, idx)
	case lexer.DEFINE:
		return parseObject(tokens, idx)
	case lexer.ENUM:
		return parseEnum(tokens, idx)
	case lexer.RANGE:
		return parseRange(tokens, idx)
	}

	return stmt, err
}

// Returns index of target
func seekToken(tokens []lexer.Token, start int, target int) (endIdx int, eof bool) {
	for i := start; i < len(tokens); i++ {
		// Missing semicolon if there are multiple statement identifiers
		if i > start && tokens[i].Type >= lexer.FUNC {
			return 0, true
		}

		if tokens[i].Type == target {
			return i, false
		}
	}

	return 0, true
}

func parseRange(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	i := *idx
	if len(tokens[i:]) < 6 {
		return stmt, ErrInvalidStatement
	}

	// Dont care about expression because of manual parsing
	_, block, err := getExpressionAndBlock(tokens, idx, true)
	if err != nil {
		return stmt, err
	}

	// Get left and right side of 'in'
	endIdx, _ := seekToken(tokens, i, lexer.LEFT_BRACE)
	interval := tokens[i+1 : endIdx]
	split := util.SplitByToken(interval, lexer.IN)
	if len(split) != 2 {
		return stmt, ErrInvalidStatement
	}

	left, err := expr.ParseExpression(split[0])
	if err != nil {
		return stmt, err
	}

	// Only variables are identifiers
	if left.Type != expr.Variable {
		return stmt, ErrExpectedIdentifier
	}

	// Right side is just expression passed into private function
	right, err := expr.ParseExpression(split[1])
	return Statement{Type: Range, Expression: &right, Name: left.Name, Then: &block}, err
}

func parseEnum(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	i := *idx
	if len(tokens[i:]) < 3 {
		return stmt, ErrInvalidStatement
	}

	if tokens[i+1].Type != lexer.LEFT_BRACE {
		return stmt, ErrExpectedBlock
	}

	endIdx, eof := seekToken(tokens, i, lexer.RIGHT_BRACE)
	if eof {
		return stmt, ErrNoBrace
	}

	names := []string{}
	for _, t := range tokens[i+2 : endIdx] {
		if t.Type != lexer.IDENTIFIER {
			return stmt, ErrExpectedIdentifier
		}

		names = append(names, t.Lexeme)
	}

	*idx = endIdx
	return Statement{Type: Enum, Params: names}, err
}

func parseImport(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) != 2 {
		return stmt, ErrInvalidStatement
	}

	if str, ok := tokens[1].Literal.(string); ok {
		name := strings.TrimSuffix(str, ".fizz")
		return Statement{Type: Import, Name: name}, err
	}

	return stmt, ErrExpectedName
}

func parseInclude(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) != 2 {
		return stmt, ErrInvalidStatement
	}

	if name, ok := tokens[1].Literal.(string); ok {
		return Statement{Type: Include, Name: name}, err
	}

	return stmt, ErrExpectedName
}

func parseExit(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) == 1 {
		return Statement{Type: Exit}, err
	}

	expr, err := expr.ParseExpression(tokens[1:])
	return Statement{Type: Exit, Expression: &expr}, err
}

func parseReturn(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) == 1 {
		return Statement{Type: Return}, err
	}

	expr, err := expr.ParseExpression(tokens[1:])
	return Statement{Type: Return, Expression: &expr}, err
}

func parseBreak(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) > 1 {
		return stmt, ErrInvalidStatement
	}

	return Statement{Type: Break}, err
}

func parseSkip(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) > 1 {
		return stmt, ErrInvalidStatement
	}

	return Statement{Type: Skip}, err
}

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

func parsePrint(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) == 1 {
		return stmt, ErrExpectedExpression
	}

	expr, err := expr.ParseExpression(tokens[1:])
	return Statement{Type: Print, Expression: &expr}, err
}

func parseError(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) == 1 {
		return stmt, ErrExpectedExpression
	}

	expr, err := expr.ParseExpression(tokens[1:])
	return Statement{Type: Error, Expression: &expr}, err
}

// Also parses variable declaration with := operator
func parseAssignment(tokens []lexer.Token) (stmt Statement, err error) {
	if len(tokens) < 3 {
		return parseExpression(tokens)
	}

	validOperands := []int{lexer.EQUAL, lexer.PLUS_EQUAL, lexer.MINUS_EQUAL, lexer.MULT_EQUAL, lexer.DIV_EQUAL, lexer.DEF_EQUAL}
	splits := util.SplitByTokens(tokens, validOperands)
	if len(splits) != 2 {
		return parseExpression(tokens)
	}

	left, err := expr.ParseExpression(splits[0])
	if err != nil {
		return stmt, err
	}

	operator := tokens[len(splits[0])].Type
	right, err := expr.ParseExpression(splits[1])
	return Statement{Type: Assignment, Expression: &right, Left: &left, Operator: operator}, err
}

func parseFunc(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	if len(tokens) < 6 {
		return stmt, ErrInvalidStatement
	}

	nameToken := tokens[*idx+1]
	if nameToken.Type != lexer.IDENTIFIER || tokens[*idx+2].Type != lexer.LEFT_PAREN {
		return stmt, ErrInvalidStatement // Missing identifier or block
	}

	*idx += 3 // Skip to start of param list
	endIdx, eof := seekToken(tokens, *idx, lexer.RIGHT_PAREN)
	if eof {
		return stmt, ErrInvalidStatement
	}

	// Get param names
	params := []string{}
	for _, p := range tokens[*idx:endIdx] {
		switch p.Type {
		case lexer.COMMA:
			continue
		case lexer.IDENTIFIER:
			params = append(params, p.Lexeme)
			continue
		}

		return stmt, ErrExpectedIdentifier
	}

	*idx = endIdx + 1 // Skip to start of block
	if tokens[*idx].Type != lexer.LEFT_BRACE {
		return stmt, ErrExpectedBlock
	}

	block, err := getBlockStatement(tokens, idx)
	return Statement{Type: Function, Name: nameToken.Lexeme, Params: params, Then: &block}, err
}

// Modifies index to go to block end. First token must be left brace
func getBlockStatement(tokens []lexer.Token, idx *int) (block Statement, err error) {
	start := *idx
	if tokens[start].Type != lexer.LEFT_BRACE {
		return block, ErrExpectedBlock
	}

	if endIdx, eof := util.SeekClosingBracket(tokens, *idx, lexer.LEFT_BRACE, lexer.RIGHT_BRACE); !eof {
		*idx = endIdx
		blockTokens := tokens[start+1 : *idx]
		statements, err := ParseStatements(blockTokens)
		return Statement{Type: Block, Statements: statements}, err
	}

	return block, ErrNoBrace
}

func parseBlock(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	return getBlockStatement(tokens, idx)
}

// Parses expression and block after keyword
func getExpressionAndBlock(tokens []lexer.Token, idx *int, expectExpr bool) (expr Statement, block Statement, err error) {
	startBlock, eof := seekToken(tokens, *idx, lexer.LEFT_BRACE)
	if eof {
		return expr, block, ErrExpectedBlock
	}

	expr, err = parseExpression(tokens[*idx+1 : startBlock])
	if err != nil {
		return expr, block, err
	}

	if expectExpr && expr.Expression == nil {
		return expr, block, ErrExpectedExpression
	}

	*idx = startBlock
	block, err = getBlockStatement(tokens, idx)
	return expr, block, err
}

func parseIf(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	stmt, block, err := getExpressionAndBlock(tokens, idx, true)
	if err != nil {
		return stmt, err
	}

	// Check for else statement
	if *idx+1 < len(tokens) && tokens[*idx+1].Type == lexer.ELSE {
		if *idx+2 >= len(tokens) {
			return stmt, ErrExpectedBlock
		}

		*idx += 2 // Skip to block
		elseBlock, err := getBlockStatement(tokens, idx)
		return Statement{Type: If, Expression: stmt.Expression, Then: &block, Else: &elseBlock}, err
	}

	return Statement{Type: If, Expression: stmt.Expression, Then: &block}, err
}

func parseWhile(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	stmt, block, err := getExpressionAndBlock(tokens, idx, false)
	if err != nil {
		return stmt, err
	}

	if stmt.Expression == nil {
		return Statement{Type: While, Then: &block}, err
	}

	return Statement{Type: While, Expression: stmt.Expression, Then: &block}, err
}

func parseRepeat(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	stmt, block, err := getExpressionAndBlock(tokens, idx, true)
	if err != nil {
		return stmt, err
	}

	return Statement{Type: Repeat, Expression: stmt.Expression, Then: &block}, err
}

func parseObject(tokens []lexer.Token, idx *int) (stmt Statement, err error) {
	if len(tokens[*idx:]) < 4 {
		return stmt, ErrInvalidStatement
	}

	nameToken := tokens[*idx+1]
	if nameToken.Type != lexer.IDENTIFIER {
		return stmt, ErrExpectedIdentifier
	}

	*idx += 3 // Goto start of block
	if tokens[*idx-1].Type != lexer.LEFT_BRACE {
		return stmt, ErrExpectedBlock
	}

	endIdx, eof := seekToken(tokens, *idx, lexer.RIGHT_BRACE)
	if eof {
		return stmt, ErrNoBrace
	}

	fields := tokens[*idx:endIdx]
	fieldNames := []string{}
	for _, field := range fields {
		switch field.Type {
		case lexer.COMMA:
			continue
		case lexer.IDENTIFIER:
			fieldNames = append(fieldNames, field.Lexeme)
			continue
		}

		return stmt, ErrExpectedIdentifier
	}

	if len(fieldNames) == 0 {
		return stmt, ErrExpectedIdentifier
	}

	*idx = endIdx
	return Statement{Type: Object, Name: nameToken.Lexeme, Params: fieldNames}, err
}
