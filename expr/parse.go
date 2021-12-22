package expr

import (
	"fmt"

	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/util"
)

func ParseExpression(tokens []lexer.Token) (expr Expression, err error) {
	if len(tokens) == 0 {
		return Expression{Type: EmptyExpression}, err
	}

	line := tokens[0].Line

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

		fmt.Println(args)
		return Expression{Type: Args, Exprs: args, Line: line}, err
	}

	// UNARY
	// Check if first token is a valid unary token type, if not skip to binary
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
	// Binary is already parsed so there cannot be another expression after the group. This means the last
	// token MUST be a closing paren. Parses inner expression.
	if tokens[0].Type == lexer.LEFT_PAREN {
		if tokens[len(tokens)-1].Type != lexer.RIGHT_PAREN {
			return expr, ErrParenError
		}

		inner, err := ParseExpression(tokens[1:len(tokens)-1])
		return Expression{Type: Group, Inner: &inner, Line: line}, err
	}

	// ARRAY
	// Same as group parsing but with square bracket. Handled differently in eval.
	// Todo: implement array parsing

	// CALL
	// Search for left paren, then parse the left part of the expression. Also parse the args of the caller.
	// Start cannot be set to 0 in loop because that would be a group expression
	targetCall, eofCall := util.SeekBreakPoint(tokens, func(i int, t lexer.Token) bool {
		return t.Type == lexer.LEFT_PAREN
	})

	// Also check where the closest dot is. If a dot comes after the last call, it should be parsed as a getter.
	// The only other option is for the call to be last (or error).
	targetDot, eofDot := util.SeekBreakPoint(tokens, func(i int, t lexer.Token) bool {
		return t.Type == lexer.DOT
	})

	if !eofCall && targetCall > targetDot {
		callee, err := ParseExpression(tokens[:targetCall])
		if err != nil {
			return expr, err
		}

		args, err := ParseExpression(tokens[targetCall:])
		return Expression{Type: Call, Left: &callee, Inner: &args, Line: line}, err
	}

	// GETTER
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

	// All expression types below only require one token. If more are found its and error.
	if len(tokens) != 1 {
		return expr, ErrInvalidExpression
	}

	// VARIABLE
	// Variables have a different expression type
	if tokens[0].Type == lexer.IDENTIFIER {
		return Expression{Type: Variable, Name: tokens[0].Lexeme, Line: line}, err
	}

	// LITERAL
	// Only other option is a literal, the only error this can cause is an undefined variable
	return Expression{Type: Literal, Value: tokens[0], Line: line}, err
}