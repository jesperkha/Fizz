package expr

import (
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

func ParseExpression(tokens []lexer.Token) (expr Expression, err error) {
	if len(tokens) == 0 {
		return Expression{Type: EmptyExpression}, err
	}

	// Todo: parse dot getters
	line := tokens[0].Line

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

	// Unary expressions
	// Check if first token is a valid unary token type, if not skip to binary
	unaryOperators := []int{lexer.MINUS, lexer.TYPE, lexer.NOT}
	if util.Contains(unaryOperators, tokens[0].Type) {
		right, err := ParseExpression(tokens[1:])
		return Expression{Type: Unary, Right: &right, Operand: tokens[0], Line: line}, err
	}

	// Binary expressions
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

	// Group expression
	// Binary is already parsed so there cannot be another expression after the group. This means the last
	// token MUST be a closing paren. Parses inner expression.
	if tokens[0].Type == lexer.LEFT_PAREN {
		if tokens[len(tokens)-1].Type != lexer.RIGHT_PAREN {
			return expr, ErrParenError
		}

		inner, err := ParseExpression(tokens[1:len(tokens)-1])
		return Expression{Type: Group, Inner: &inner, Line: line}, err
	}

	// Call expression
	// Search for left paren, then parse the left part of the expression. Also parse the args of the caller.
	// Start cannot be set to 0 in loop because that would be a group expression
	targetIdx, eof := util.SeekBreakPoint(tokens, func(i int, t lexer.Token) bool {
		return t.Type == lexer.LEFT_PAREN
	})

	if !eof {
		callee, err := ParseExpression(tokens[:targetIdx])
		if err != nil {
			return expr, err
		}

		args, err := ParseExpression(tokens[targetIdx:])
		return Expression{Type: Call, Left: &callee, Inner: &args, Line: line}, err
	}

	// All expression types below only require one token. If more are found its and error.
	if len(tokens) != 1 {
		return expr, ErrInvalidExpression
	}

	// Variable expression
	if tokens[0].Type == lexer.IDENTIFIER {
		return Expression{Type: Variable, Name: tokens[0].Lexeme, Line: line}, err
	}

	// Literal expression
	return Expression{Type: Literal, Value: tokens[0], Line: line}, err
}