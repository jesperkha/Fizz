package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/jesperkha/Fizz/lexer"
)

func main() {
	file, err := os.Open("./main.fz")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	bufio.NewReader(file).WriteTo(&buf)

	tokens, err := lexer.GetTokens(buf.String())
	if err != nil {
		log.Fatal(err)
	}

	for idx, token := range tokens {
		fmt.Println(idx, token.Lexeme)
	}
}