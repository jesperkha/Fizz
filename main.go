package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/parser"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func check(expected string, value interface{}) {
	fmt.Printf("Expected value: %s. Got: %s\n", expected, value)
}

func main() {
	file, err := os.Open("./main.fz")
	handleError(err)

	var buf bytes.Buffer
	bufio.NewReader(file).WriteTo(&buf)

	tokens, err := lexer.GetTokens(buf.String())
	handleError(err)

	ptokens, err := parser.GenerateParseTokens(tokens)
	handleError(err)

	expr := parser.ParseExpression(ptokens)
	check("1", expr.Left.Value.Lexeme)
	check("+", expr.Operand.Lexeme)
	check("2", expr.Right.Inner.Left.Value.Lexeme)
	check("-", expr.Right.Inner.Operand.Lexeme)
	check("3", expr.Right.Inner.Right.Value.Lexeme)
}