package expr

import (
	"fmt"

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

// Parse comma separated vaules for function arguments and array values
func parseCSV(tokens []lexer.Token, start int, beginT int, endT int) (ptokens []ParseToken, endIdx int, err error) {
	endIdx, eof := util.SeekClosingBracket(tokens, start, beginT, endT)
	if eof {
		return ptokens, 0, ErrBracketError
	}
	
	interval := tokens[start+1:endIdx]
	group := []ParseToken{}
	split := util.SplitByToken(interval, lexer.COMMA)
	for _, expr := range split {
		t, err := generateParseTokens(expr)
		if err != nil {
			return ptokens, 0, err
		}

		group = append(group, ParseToken{Type: TokenGroup, Inner: t})
	}

	return group, endIdx, err
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

		// Parse array expression
		if token.Type == lexer.LEFT_SQUARE {
			group, endIdx, err := parseCSV(tokens, currentIdx, lexer.LEFT_SQUARE, lexer.RIGHT_SQUARE)
			if err != nil {
				return ptokens, err
			}

			ptokens = append(ptokens, ParseToken{Type: ArrayGroup, Inner: group})
			currentIdx = endIdx + 1
			continue
		}

		// Check and handle function call expression
		// Todo: make function for parsing comma separated values
		// connect to array parsing too
		if token.Type == lexer.IDENTIFIER {
			// Check if just identifier
			if currentIdx+1 >= len(tokens) || tokens[currentIdx+1].Type != lexer.LEFT_PAREN {
				ptokens = append(ptokens, ParseToken{Type: Single, Token: token})
				currentIdx++
				continue
			}

			currentIdx++ // Skip function name
			group, endIdx, err := parseCSV(tokens, currentIdx, lexer.LEFT_PAREN, lexer.RIGHT_PAREN)
			if err != nil {
				return ptokens, err
			}

			ptokens = append(ptokens, ParseToken{Type: CallGroup, Token: token, Inner: group})
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

		// Parse array group or call group, both are comma separated
		if token.Type == ArrayGroup || token.Type == CallGroup {
			exprs := []Expression{}
			for _, grp := range token.Inner {
				exprs = append(exprs, *parsePTokens(grp.Inner))
			}
			
			if token.Type ==  ArrayGroup {
				return &Expression{Type: Array, Line: line, Exprs: exprs}
			}
			
			return &Expression{Type: Call, Line: line, Name: token.Token.Lexeme, Exprs: exprs}
		}

		// Parse single variable
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
	
	// Is getter expression (last parsed)
	if tokens[1].Token.Type == lexer.DOT && lowest.Type == lexer.IDENTIFIER {
		if len(tokens)%2 == 0 {
			return &Expression{}
		}
		
		exprs := []Expression{}
		for i := 0; i < len(tokens)-1; i += 2 {
			idn := tokens[i]
			if tokens[i+1].Token.Type != lexer.DOT {
				return &Expression{}
			}
			
			exprs = append(exprs, *parsePTokens([]ParseToken{idn}))
		}
		
		// Get last token
		exprs = append(exprs, *parsePTokens([]ParseToken{tokens[len(tokens)-1]}))
		fullName := ""
		for _, e := range exprs {
			fullName += fmt.Sprintf(".%s", e.Name)
		}

		return &Expression{Type: Getter, Exprs: exprs, Line: line, Name: fullName[1:]}
	}
	
	// Invalid expression might have lowest as end token
	// Move to middle and let evaluation handle error
	if lowestIdx == len(tokens)-1 {
		newIdx := int(len(tokens) / 2)
		lowestIdx = newIdx
		lowest = tokens[lowestIdx].Token
	}
	
	// Unary expression
	unaryTokens := []int{lexer.MINUS, lexer.TYPE, lexer.NOT}
	if len(tokens) == 2 || (util.Contains(unaryTokens, lowest.Type) && lowest.Type == tokens[0].Type) {
		return &Expression{Type: Unary, Line: line, Operand: tokens[0].Token, Right: parsePTokens(tokens[1:])}
	}

	right, left := parsePTokens(tokens[lowestIdx+1:]), parsePTokens(tokens[:lowestIdx])
	return &Expression{Type: Binary, Line: line, Operand: lowest, Left: left, Right: right}
}
