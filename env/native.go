package env

import (
	"bufio"
	"fmt"
	"os"
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)

	// Todo: raise error if stdin is not an interactive terminal
	// https://github.com/mattn/go-isatty
	Declare("input", Callable{
		NumArgs: 1,
		Call: func(i ...interface{}) (interface{}, error) {
			fmt.Print(i...)
			scanner.Scan()
			input := scanner.Text()
			return input, nil
		},
	})
}