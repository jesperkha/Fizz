package std

import (
	"bufio"
	"fmt"
	"os"
)

type i interface{}

var Includes = map[string]interface{}{}

// Inits all the functions in this package to lib. Added to include.go when running
// script for adding new libraries.
func init() {
	Includes = map[string]interface{}{
		"input":    getStdin,
		"toString": toString,
	}
}

var scanner = bufio.NewScanner(os.Stdin)

// Gets input from standard input. Returns said input. Assumes terminal is interactive.
func getStdin(prompt string) (input i, err error) {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text(), nil
}

// Converts any value into string literal
func toString(val i) (str i, err error) {
	return fmt.Sprint(val), err
}
