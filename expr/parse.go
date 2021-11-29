package expr

import (
	"log"

	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

// Generates ptokens and parses them into an expression.
func ParseExpression(tokens []lexer.Token) (expr Expression, err error) {
	ptokens, err := generateParseTokens(tokens)
	if err != nil {
		return expr, err
	}
	
	return *parsePTokens(ptokens), err
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

	// Todo: make sure getters dont cause an error here
	for currentIdx < len(tokens) {
		token := tokens[currentIdx]
		line := token.Line

		// Find end paren and call recursive for inner part
		if token.Type == lexer.LEFT_PAREN {
			endIdx, eof := util.SeekClosingBracket(tokens, currentIdx, lexer.LEFT_PAREN, lexer.RIGHT_PAREN)
			if eof {
				return ptokens, util.FormatError(ErrParenError, line)
			}

			if endIdx-currentIdx == 1 {
				return ptokens, util.FormatError(ErrExpectedExpression, line)
			}

			// Generate ptokens between start and end paren
			tokenGroup, err := generateParseTokens(tokens[currentIdx+1 : endIdx])
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
			return ptokens, util.FormatError(ErrParenError, line)
		}

		// Check and handle function call expression
		if token.Type == lexer.IDENTIFIER {
			if currentIdx+1 >= len(tokens) || tokens[currentIdx+1].Type != lexer.LEFT_PAREN {
				ptokens = append(ptokens, ParseToken{Type: Single, Token: token})
				currentIdx++
				continue
			}

			// Increment to skip identifier
			currentIdx++
			endIdx, eof := util.SeekClosingBracket(tokens, currentIdx, lexer.LEFT_PAREN, lexer.RIGHT_PAREN)
			if eof {
				return ptokens, util.FormatError(ErrParenError, line)
			}

			// +1 to skip start paren
			interval := tokens[currentIdx+1 : endIdx]
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
					return ptokens, util.FormatError(ErrCommaError, line)
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

		// Not expression type symbol
		if token.Type > lexer.IDENTIFIER && token.Type != lexer.DOT {
			return ptokens, ErrInvalidExpression
		}

		ptokens = append(ptokens, ParseToken{Type: Single, Token: token})
		currentIdx++
	}

	return ptokens, err
}

// Parses slice of ParseTokens into final AST
func parsePTokens(tokens []ParseToken) *Expression {
	if len(tokens) == 0 {
		return &Expression{}
	}

	line := tokens[0].Token.Line

	// Literal, Variable, or Group expression
	if len(tokens) == 1 {
		token := tokens[0]
		if token.Type == TokenGroup {
			return &Expression{Type: Group, Line: line, Inner: parsePTokens(token.Inner)}
		}

		// Parse call expression
		if token.Type == CallGroup {
			argExpressions := []Expression{}
			for _, arg := range token.Args {
				argExpressions = append(argExpressions, *parsePTokens(arg))
			}

			return &Expression{Type: Call, Line: line, Name: token.Token.Lexeme, Exprs: argExpressions}
		}

		// Variable
		if token.Token.Type == lexer.IDENTIFIER {
			return &Expression{Type: Variable, Line: line, Name: token.Token.Lexeme}
		}

		// Defualts to literal
		return &Expression{Type: Literal, Line: line, Value: token.Token}
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

	// Invalid expression might have lowest as end token
	// Move to middle and let evaluation handle error
	if lowestIdx == len(tokens)-1 {
		newIdx := int(len(tokens) / 2)
		lowestIdx = newIdx
		lowest = tokens[lowestIdx].Token
	}

	// Todo: parse chained dots similar to a chained plus expression but order
	if tokens[1].Type == lexer.DOT {
		log.Fatal("dot")
	}
	
	// Unary expression
	unaryTokens := []int{lexer.MINUS, lexer.TYPE, lexer.NOT}
	if len(tokens) == 2 || (util.Contains(unaryTokens, lowest.Type) && lowest.Type == tokens[0].Type)  {
		return &Expression{Type: Unary, Line: line, Operand: tokens[0].Token, Right: parsePTokens(tokens[1:])}
	}

	right, left := parsePTokens(tokens[lowestIdx+1:]), parsePTokens(tokens[:lowestIdx])
	return &Expression{Type: Binary, Line: line, Operand: lowest, Left: left, Right: right}
}
