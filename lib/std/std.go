package std

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Todo: add math, string, memory, and io functions

type i interface{}

var Includes = map[string]interface{}{}

func init() {
	Includes = map[string]interface{}{
		"input":    input,
		"toString": toString,
		"toNumber": toNumber,
	}
}

var scanner = bufio.NewScanner(os.Stdin)

func input(prompt string) (input i, err error) {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text(), nil
}

func toString(val i) (str i, err error) {
	return fmt.Sprint(val), err
}

func toNumber(val string) (num i, err error) {
	return strconv.ParseFloat(val, 64)
}
