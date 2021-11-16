package expr

import (
	"github.com/jesperkha/Fizz/lexer"
)

// Generates ptokens and parses them into an expression.
func ParseExpression(tokens []lexer.Token) (expr Expression, err error) {
	ptokens, err := generateParseTokens(tokens)
	if err != nil {
		return expr, err
	}

	return *parsePTokens(ptokens), err
}

func seekEndParen(tokens []lexer.Token, start int) (endIdx int, eof bool) {
	numParen := 0
	for i := start; i < len(tokens); i++ {
		switch tokens[i].Type {
		case lexer.LEFT_PAREN: numParen++
		case lexer.RIGHT_PAREN: numParen--
		}

		if numParen == 0 {
			return i, false
		}
	}

	return endIdx, true
}

// Creates new ParseTokens from lexer tokens to simplify expression parsing. The ParseTokens can
// either be of type Single or TokenGroup. Symbols and identifiers are of the Single type while any
// expression within parens is a TokenGroup type. The single type has a .Token value which is the
// original lexer token. TokenGroup.Inner is a slice of ParseTokens which is retrieved from recursively
// calling this method. The result is that any parenthesized expression is put in a nested object to
// minimize operations and complexity when evaluating the actual expressions later.
// Example: 1, +, (, 2, +, 3, ) turns into 1, +, [2, +, 3]
func generateParseTokens(tokens []lexer.Token) (ptokens []ParseToken, err error) {
	currentIdx := 0

	for currentIdx < len(tokens) {
		token := tokens[currentIdx]
		line := token.Line
		
		// Find end paren and call recursive for inner part
		if token.Type == lexer.LEFT_PAREN {
			endIdx, eof := seekEndParen(tokens, currentIdx)
			if eof {
				return ptokens, formatError(ErrParenError, line)
			}

			if endIdx - currentIdx == 1 {
				return ptokens, formatError(ErrExpectedExpression, line)
			}

			// Generate ptokens between start and end paren
			tokenGroup, err := generateParseTokens(tokens[currentIdx + 1:endIdx])
			if err != nil {
				return ptokens, err
			}

			ptokens = append(ptokens, ParseToken{Type: TokenGroup, Inner: tokenGroup})
			// Skip past parenthesized section, +1 to skip closing paren
			currentIdx = endIdx + 1
			continue
		}

		// Unmatched open and closing parens
		if token.Type == lexer.RIGHT_PAREN {
			return ptokens, formatError(ErrParenError, line)
		}

		// Check and handle function call expression
		if token.Type == lexer.IDENTIFIER {
			if currentIdx + 1 >= len(tokens) || tokens[currentIdx + 1].Type != lexer.LEFT_PAREN {
				ptokens = append(ptokens, ParseToken{Type: Single, Token: token})
				currentIdx++
				continue
			}

			// Increment to skip identifier
			currentIdx++
			endIdx, eof := seekEndParen(tokens, currentIdx)
			if eof {
				return ptokens, formatError(ErrParenError, line)
			}

			// +1 to skip start paren
			interval := tokens[currentIdx + 1:endIdx]
			// Add end comma to parse last expression
			if len(interval) != 0 {
				interval = append(interval, lexer.Token{Type: lexer.COMMA})
			}

			args := [][]ParseToken{}
			exprStart := 0 // Start of arg expression (index)
			for idx, t := range interval {
				if t.Type != lexer.COMMA {
					continue
				}
				
				if exprStart == idx {
					return ptokens, formatError(ErrCommaError, line)
				}

				argToken, err := generateParseTokens(interval[exprStart:idx])
				if err != nil {
					return ptokens, err
				}

				args = append(args, argToken)
				exprStart = idx + 1
			}

			callToken := ParseToken{Type: CallGroup, Token: token, Args: args}
			ptokens = append(ptokens, callToken)
			currentIdx = endIdx + 1
			continue
		}

		ptokens = append(ptokens, ParseToken{Type: Single, Token: token})
		currentIdx++
	}

	return ptokens, err
}

// Parses slice of ParseTokens into final AST
func parsePTokens(tokens []ParseToken) *Expression {
	// Literal, Variable, or Group expression
	if len(tokens) == 1 {
		token := tokens[0]
		if token.Type == TokenGroup {
			return &Expression{Type: Group, Inner: parsePTokens(token.Inner)}
		}

		// Parse call expression
		if token.Type == CallGroup {
			argExpressions := []Expression{}
			for _, arg := range token.Args {
				argExpressions = append(argExpressions, *parsePTokens(arg))
			}

			return &Expression{Type: Call, Callee: token.Token, Exprs: argExpressions}
		}

		// Variable
		if token.Token.Type == lexer.IDENTIFIER {
			return &Expression{Type: Variable, Name: token.Token.Lexeme}
		}

		// Defualts to literal
		return &Expression{Type: Literal, Value: token.Token}
	}

	// Unary expression
	if len(tokens) == 2 {
		return &Expression{Type: Unary, Operand: tokens[0].Token, Right: parsePTokens(tokens[1:])}
	}

	// Binary expression
	lowest := lexer.Token{Type: 999}
	lowestIdx := 0
	for idx, token := range tokens {
		if token.Type == Single && token.Token.Type < lowest.Type {
			lowest = token.Token
			lowestIdx = idx
		}
	}

	right, left := parsePTokens(tokens[lowestIdx + 1:]), parsePTokens(tokens[:lowestIdx])
	return &Expression{Type: Binary, Operand: lowest, Left: left, Right: right}
}
