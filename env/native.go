package env

import (
	"bufio"
	"fmt"
	"os"
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)

	Declare("input", Callable{
		NumArgs: 1,
		Call: func(i ...interface{}) (interface{}, error) {
			// Todo: raise error if stdin is not an interactive terminal
			// https://github.com/mattn/go-isatty
			fmt.Print(i...)
			scanner.Scan()
			input := scanner.Text()
			return input, nil
		},
	})
}