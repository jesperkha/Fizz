package main

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/stmt"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func runFile(filename string) {
	file, err := os.Open(filename)
	handleError(err)
	var buf bytes.Buffer
	bufio.NewReader(file).WriteTo(&buf)

	tokens, err := lexer.GetTokens(buf.String())
	handleError(err)

	stmts, err := stmt.ParseStatements(tokens)
	handleError(err)

	err = stmt.ExecuteStatements(stmts)
	handleError(err)
}

func main() {
	runFile("./main.fz")
}
