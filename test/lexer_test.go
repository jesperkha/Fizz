package test

import (
	"testing"

	"github.com/jesperkha/Fizz/lexer"
)

func TestTokenLexer(t *testing.T) {
	input := "var name = \"John\";"
	tokens, err := lexer.GetTokens(input)
	if err != nil {
		t.Error(err)
	}

	// Check string lexing
	strToken := tokens[3]
	if strToken.Type != lexer.STRING || strToken.Literal != "John" {
		t.Error("failed to tokenize string")
	}
}