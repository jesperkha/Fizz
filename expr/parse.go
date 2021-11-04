package expr

import (
	"fmt"

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

		// Find end paren and call recursive for inner part
		if token.Type == lexer.LEFT_PAREN {
			endIdx := currentIdx
			for i := currentIdx; i < len(tokens); i++ {
				if tokens[i].Type == lexer.RIGHT_PAREN {
					endIdx = i
					break
				}
			}

			// Error for no closing paren
			if endIdx == currentIdx {
				return ptokens, fmt.Errorf(ErrParenError.Error(), token.Line)
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
			return ptokens, fmt.Errorf(ErrParenError.Error(), token.Line)
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

		if token.Token.Type == lexer.IDENTIFIER {
			return &Expression{Type: Variable, Name: token.Token.Lexeme}
		}

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
