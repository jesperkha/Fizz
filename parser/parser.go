package parser

import (
	"errors"
	"fmt"

	"github.com/jesperkha/Fizz/lexer"
)

const (
	Literal = iota
	Unary
	Binary
	Group
)

type Expression struct {
	Type     int
	Operand  lexer.Token
	Terminal lexer.Token
	Value    lexer.Token
	Left	 *Expression
	Right 	 *Expression
	Inner 	 *Expression
}

const (
	Single = iota
	TokenGroup
)

type ParseToken struct {
	Type   int
	Token  lexer.Token
	Inner  []ParseToken
}

var (
	ErrParenError = errors.New("unmatched parenthesies, line %d")
)

// Creates new ParseTokens from lexer tokens to simplify expression parsing. The ParseTokens can
// either be of type Single or TokenGroup. Symbols and identifiers are of the Single type while any
// expression within parens is a TokenGroup type. The single type has a .Token value which is the
// original lexer token. TokenGroup.Inner is a slice of ParseTokens which is retrieved from recursively
// calling this method. The result is that any parenthesized expression is put in a nested object to
// minimize operations and complexity when evaluating the actual expressions later.
// Example: 1, +, (, 2, +, 3, ) turns into 1, +, [2, +, 3]
func GenerateParseTokens(tokens []lexer.Token) (ptokens []ParseToken, err error) {
	currentIdx := 0

	for currentIdx < len(tokens) {
		token := tokens[currentIdx]

		// Find end paren and call recursive for inner part
		if token.Type == lexer.LEFT_PAREN {
			endIdx := currentIdx
			for i := currentIdx; i < len(tokens); i++ {
				if tokens[i].Type == lexer.RIGHT_PAREN && i > endIdx {
					endIdx = i
				}
			}

			// Error for no closing paren
			if endIdx == currentIdx {
				return ptokens, fmt.Errorf(ErrParenError.Error(), token.Line)
			}

			// Generate ptokens between start and end paren
			tokenGroup, err := GenerateParseTokens(tokens[currentIdx + 1:endIdx])
			if err != nil {
				return ptokens, err
			}

			ptokens = append(ptokens, ParseToken{Type: TokenGroup, Inner: tokenGroup})
			// Skip past parenthesized section, +1 to skip closing paren
			currentIdx += (endIdx - currentIdx) + 1
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
func ParseExpression(tokens []ParseToken) *Expression {
	// Literal or Group expression
	if len(tokens) == 1 {
		token := tokens[0]
		if token.Type == TokenGroup {
			return &Expression{Type: Group, Inner: ParseExpression(token.Inner)}
		}

		return &Expression{Type: Literal, Value: token.Token}
	}

	// Unary expression
	if len(tokens) == 2 {
		return &Expression{Type: Unary, Operand: tokens[0].Token, Value: tokens[1].Token}
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

	right, left := ParseExpression(tokens[lowestIdx + 1:]), ParseExpression(tokens[:lowestIdx])
	return &Expression{Type: Binary, Operand: lowest, Left: left, Right: right}
}
