package std

import (
	"bufio"
	"fmt"
	"os"
)

type i interface{}

// Bad pracice, but ok here. Bigger values should be stored/created upon calling an
// init function or something similar from your library.
var scanner = bufio.NewScanner(os.Stdin)

// Gets input from standard input. Returns said input. Assumes terminal is interactive.
func GetStdinInput(prompt string) (input i, err error) {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text(), nil
}