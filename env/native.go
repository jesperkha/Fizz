package env

import (
	"bufio"
	"os"
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)

	// Todo: raise error if stdin is not an interactive terminal
	// https://github.com/mattn/go-isatty
	Declare("input", Callable{
		NumArgs: 0,
		Call: func(i ...interface{}) (interface{}, error) {
			scanner.Scan()
			input := scanner.Text()
			return input, nil
		},
	})
}
