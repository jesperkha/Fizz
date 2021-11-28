package test

import (
	"testing"

	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/stmt"
)

var (
	// Valid statements, should not return error
	ExampleStatements1 = []string{
		"print 1;",
		"var age = 20; age += 1;",
		"if true { 1+1; }",
		"while { break; }",
		"repeat n < 10 { skip; }",
		"func main(a, b) {return 1;} print main(1, 2);",
	}

	// Invalid statements, should return error, but not panic
	ExampleStatements2 = []string{
		"printf 20;",
		"var a a = 0;",
		"func () {}",
		"func main(a) {} main();",
		"if {}",
		"while a + 2 {}",
		"repeat 20 {}",
	}
)

func parseStatement(input string) (err error) {
	tokens, err := lexer.GetTokens(input)
	if err != nil {
		return err
	}

	stmts, err := stmt.ParseStatements(tokens)
	if err != nil {
		return err
	}

	return stmt.ExecuteStatements(stmts)
}

func TestValidStatements(t *testing.T) {
	for idx, s := range ExampleStatements1 {
		if err := parseStatement(s); err != nil {
			t.Errorf("expected no error, valid case %d, got err: %s", idx, err)
		}
	}
}

func TestInalidStatements(t *testing.T) {
	for idx, s := range ExampleStatements2 {
		if parseStatement(s) == nil {
			t.Errorf("expected error, invalid case %d", idx)
		}
	}
}
