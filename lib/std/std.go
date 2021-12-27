package std

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Standard library for common io, string, time, and exec fuctionality

type i interface{}

var (
	Includes = map[string]interface{}{}
	scanner  = bufio.NewScanner(os.Stdin)

	ErrNotNumber = errors.New("string could not be converted to number, line %d")
)

func init() {
	Includes = map[string]interface{}{
		"input":    input,
		"toString": toString,
		"toNumber": toNumber,
	}
}

func input(prompt string) (input i, err error) {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text(), nil
}

func toString(val i) (str i, err error) {
	return fmt.Sprint(val), err
}

func toNumber(val string) (num i, err error) {
	num, err = strconv.ParseFloat(val, 64)
	if err != nil {
		return num, ErrNotNumber
	}

	return num, err
}

// read, write, append file
// time, localTime
// lower, upper, capitalize
// copy, run
