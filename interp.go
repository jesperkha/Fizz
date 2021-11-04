package main

import (
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/stmt"
)

// Interperates string of code. The string is tokenized in the lexer package
// where each token has a type, line, lexeme, and literal value. The parsed
// tokens are passed to the statement parser which looks for statements
// (identified by trailing semicolon or curly braces).

// The statements are parsed as statement tokens with a type, name, and
// expression values for the expression associated with the statement. For
// example: a variable declaration statement has an expression property for
// the initial value (if there is one) that is found after the equal sign.

// The statements are then executed and, when doing so, statement expressions
// are evaluated. Variable values are also assigned and manipulated in the
// variable environment found in the env package.

func Interperate(input string) (err error) {
	// Parses input characters into lexical tokens for single and double symbols,
	// identifiers, and keywords.
	lexicalTokens, err := lexer.GetTokens(input)
	if err != nil {
		return err
	}
	
	// Lexical tokens are analysed and put into statement tokens. These statements
	// contain all the information they need for execution and error handling.
	statements, err := stmt.ParseStatements(lexicalTokens)
	if err != nil {
		return err
	}

	// Finally executes statement tokens. This is the only step that has any effect
	// on the actual input program as the others were just braking it up into usable
	// pieces. While the interpreter is still running, the values of variables will be
	// remembered as the environments are never reset at runtime.
	return stmt.ExecuteStatements(statements)
}