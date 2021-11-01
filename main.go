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

// func check(expected string, value interface{}) {
// 	fmt.Printf("Expected value: %s. Got: %s\n", expected, value)
// }

func Exec(input string) (value interface{}, err error) {
	tokens, err := lexer.GetTokens(input)
	if err != nil || len(tokens) == 0 {
		return nil, err
	}

	ptokens, err := parser.GenerateParseTokens(tokens)
	if err != nil {
		return nil, err
	}
	
	val, err := parser.EvaluateExpression(parser.ParseExpression(ptokens))
	if err != nil {
		return nil, err
	}

	return val, nil
}

func runFile(filename string) {
	file, err := os.Open(filename)
	handleError(err)
	var buf bytes.Buffer
	bufio.NewReader(file).WriteTo(&buf)
	fmt.Println(Exec(buf.String()))
}

func main() {
	runFile("./main.fz")
	// input := strings.Join(os.Args[1:], "")
	// val, err := Exec(input)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(val)
}