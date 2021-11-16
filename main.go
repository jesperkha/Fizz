package main

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"github.com/jesperkha/Fizz/expr"
	"github.com/jesperkha/Fizz/lexer"
)

func main() {
	// run.RunInterpeter(os.Args)
	file, err := os.Open("./main.fizz")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	bufio.NewReader(file).WriteTo(&buf)
	input := buf.String()

	ltokens, err := lexer.GetTokens(input)
	if err != nil {
		log.Fatal(err)
	}

	expr, err := expr.ParseExpression(ltokens)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(expr.Exprs)
}
