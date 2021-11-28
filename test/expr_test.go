package test

import (
	"testing"

	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
)

var (
	// Valid expressions and should not return an error
	ExampleExpressions1 = []string{
		"1 + 1",
		"(5 + 1) * 2",
		"-1 + -1",
		"8.2 % 2",
		"type true == \"bool\"",
		"5 - ((1^1) + 2)",
	}

	// Invalid expressions and should return an error, but not panic
	ExampleExpressions2 = []string{
		"+1 - 8",
		"() * 5",
		"x + 1",
		"1.0.2 - 1",
		"1 --+ 1",
		"2 + +",
	}
)

func getValue(input string) (value interface{}, err error) {
	tokens, err := lexer.GetTokens(input)
	if err != nil {
		return
	}

	e, err := expr.ParseExpression(tokens)
	if err != nil {
		return
	}

	value, err = expr.EvaluateExpression(&e)
	if err != nil {
		return
	}

	return value, err
}

func TestValidSyntax(t *testing.T) {
	for idx, input := range ExampleExpressions1 {
		_, err := getValue(input)
		if err != nil {
			t.Errorf("expected no error, valid case %d, got err: %s", idx, err)
		}
	}
}

func TestInvalidSyntax(t *testing.T) {
	for idx, input := range ExampleExpressions2 {
		if _, err := getValue(input); err == nil {
			t.Errorf("expected error, invalid case %d", idx)
		}
	}
}

func TestTypes(t *testing.T) {
	input := "1 + true"
	_, err := getValue(input)
	if err == nil {
		t.Error("expected type error")
	}
}

func TestUnary(t *testing.T) {
	input, expect := "(type true) == (type 1)", false
	value, err := getValue(input)
	if err != nil {
		t.Error(err)
	}

	if value != expect {
		t.Errorf("expected false, got %v", value)
	}
}
