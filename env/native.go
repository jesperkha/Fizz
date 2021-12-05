package env

import (
	"bufio"
	"os"
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)

	Declare("input", Callable{
		NumArgs: 0,
		Call: func(i ...interface{}) (interface{}, error) {
			scanner.Scan()
			input := scanner.Text()
			return input, nil
		},
	})
}