package expr

import (
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

func ParseExpression(tokens []lexer.Token) (expr Expression, err error) {
	if len(tokens) == 0 {
		return Expression{Type: EmptyExpression}, err
	}

	line := tokens[0].Line

	// Safeguard against unmatched parens/brackets later on. This means the bracket seek functions dont
	// need to check for eof (unless they are seek in a sliced off token list)
	if eofParens, eofBrackets := util.HasMatchedParens(tokens); eofParens || eofBrackets {
		if eofParens {
			return expr, ErrParenError
		} else {
			return expr, ErrBracketError
		}
	}

	if len(tokens) == 1 {
		// VARIABLE
		// Variables have a different expression type
		if tokens[0].Type == lexer.IDENTIFIER {
			return Expression{Type: Variable, Name: tokens[0].Lexeme, Line: line}, err
		}

		// LITERAL
		// Only other option is a literal, the only error this can cause is an undefined variable
		return Expression{Type: Literal, Value: tokens[0], Line: line}, err
	}

	// ARGUMENTS
	// Argument list for function calls. If the arg list is empty it is handled as an empty group.
	if splits := util.SplitByToken(tokens, lexer.COMMA); len(splits) != 1 {
		args := []Expression{}
		for _, list := range splits {
			arg, err := ParseExpression(list)
			if err != nil {
				return expr, err
			}

			args = append(args, arg)
		}

		return Expression{Type: Args, Exprs: args, Line: line}, err
	}

	// Todo: fix bug where unary expressions are parsed instead of binary ones
	// example: print type 1 == type "hello" (prints bool)

	// UNARY
	// Check if first token is a valid unary token type
	unaryOperators := []int{lexer.MINUS, lexer.TYPE, lexer.NOT}
	if util.Contains(unaryOperators, tokens[0].Type) {
		right, err := ParseExpression(tokens[1:])
		return Expression{Type: Unary, Right: &right, Operand: tokens[0], Line: line}, err
	}

	// BINARY
	// Find lowest precedense operator in token list. Split at that token and put into expression tree.
	// Only check if not in a group. Then check if valid operator, skip if not.
	lowest, lowestIdx := lexer.Token{Type: 999}, 0
	util.SeekBreakPoint(tokens, func(i int, t lexer.Token) bool {
		if t.Type < lowest.Type {
			lowest, lowestIdx = t, i
		}

		return false
	})

	t := lowest.Type
	if (t >= lexer.AND && t <= lexer.MINUS) || (t >= lexer.STAR && t <= lexer.HAT) {
		left, err := ParseExpression(tokens[:lowestIdx])
		if err != nil {
			return expr, err
		}

		right, err := ParseExpression(tokens[lowestIdx+1:])
		return Expression{Type: Binary, Left: &left, Right: &right, Operand: tokens[lowestIdx], Line: line}, err
	}

	// GROUP
	// Binary is already parsed so there can either be a single group or a chained call / getter.
	// First gets the endIdx of the first group. Its a single group expression is the closing paren
	// is the last token in the list. Paren error if eof, doesnt matter how many calls come after.
	endIdx, _ := util.SeekClosingBracket(tokens, 0, lexer.LEFT_PAREN, lexer.RIGHT_PAREN)
	if tokens[0].Type == lexer.LEFT_PAREN && endIdx == len(tokens)-1 {
		inner, err := ParseExpression(tokens[1 : len(tokens)-1])
		return Expression{Type: Group, Inner: &inner, Line: line}, err
	}

	// ARRAY LITERAL
	// Same as group parsing but with square bracket.
	endIdx, _ = util.SeekClosingBracket(tokens, 0, lexer.LEFT_SQUARE, lexer.RIGHT_SQUARE)
	if tokens[0].Type == lexer.LEFT_SQUARE && endIdx == len(tokens)-1 {
		inner, err := ParseExpression(tokens[1 : len(tokens)-1])
		return Expression{Type: Array, Inner: &inner, Line: line}, err
	}

	// ARRAY GETTER
	// Array index getter has same priority as call
	targetIndex, eofIndex := util.SeekBreakPoint(tokens, func(i int, t lexer.Token) bool {
		return t.Type == lexer.LEFT_SQUARE
	})

	// Also check where the closest dot is. If a dot comes after the last call, it should be parsed as a getter.
	// The only other option is for the call to be last (or error).
	targetDot, eofDot := util.SeekBreakPoint(tokens, func(i int, t lexer.Token) bool {
		return t.Type == lexer.DOT
	})

	if !eofIndex && targetIndex > targetDot {
		array, err := ParseExpression(tokens[:targetIndex])
		if err != nil {
			return expr, err
		}

		endIdx, _ := util.SeekClosingBracket(tokens, targetIndex, lexer.LEFT_SQUARE, lexer.RIGHT_SQUARE)
		arg, err := ParseExpression(tokens[targetIndex+1 : endIdx])
		return Expression{Type: Index, Left: &array, Right: &arg, Line: line}, err
	}

	// FUNCTION CALL
	// Search for left paren, then parse the left part of the expression. Also parse the args of the caller.
	// Start cannot be set to 0 in loop because that would be a group expression
	targetCall, eofCall := util.SeekBreakPoint(tokens, func(i int, t lexer.Token) bool {
		return t.Type == lexer.LEFT_PAREN
	})

	if !eofCall && targetCall > targetDot {
		callee, err := ParseExpression(tokens[:targetCall])
		if err != nil {
			return expr, err
		}

		args, err := ParseExpression(tokens[targetCall:])
		return Expression{Type: Call, Left: &callee, Inner: &args, Line: line}, err
	}

	// OBJECT GETTER
	// Splits by dot and parses the left side recursively
	if !eofDot {
		left, err := ParseExpression(tokens[:targetDot])
		if err != nil {
			return expr, err
		}

		right, err := ParseExpression(tokens[targetDot+1:])
		if right.Type != Variable {
			return expr, ErrExpectedName
		}

		return Expression{Type: Getter, Left: &left, Right: &right, Line: line}, err
	}

	return expr, ErrInvalidExpression
}
