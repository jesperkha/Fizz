package test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/interp"
	"github.com/jesperkha/Fizz/lexer"
)

func runTestFile(path string) error {
	env.ThrowEnvironment = false
	byt, _ := os.ReadFile(path)
	split := bytes.Split(byt, []byte{'?'})

	for typ, content := range split {
		cases := bytes.Split(content, []byte{'>'})
		for idx, cas := range cases[1:] { // last case is empty because of split
			_, err := interp.Interperate("", string(cas)) // removed ; in split
			if err != nil && typ == 0 {
				return fmt.Errorf("valid case %d, got err: %s", idx+1, err.Error())
			}

			if err == nil && typ == 1 {
				return fmt.Errorf("invalid case %d, got no error: %s", idx+1, string(cas))
			}
		}
	}

	return nil
}

func TestExpressions(t *testing.T) {
	if err := runTestFile("./cases/expr.fizz"); err != nil {
		t.Error(err)
	}
}

func TestStatements(t *testing.T) {
	if err := runTestFile("./cases/stmt.fizz"); err != nil {
		t.Error(err)
	}
}

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
