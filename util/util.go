package util

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/lexer"
)

// Format error with line numbers for local errors, but ignore for errors passed from
// expression parsing as they are already formatted with line numbers.
func FormatError(err error, line int) error {
	if err == nil {
		return err
	}

	if strings.Contains(err.Error(), "%d") {
		return fmt.Errorf(err.Error(), line)
	}

	return err
}

// Prints red error message to console
func PrintError(err error) {
	if err == nil {
		return
	}

	ct.Foreground(ct.Red, true)
	fmt.Fprintln(os.Stderr, err.Error())
	ct.ResetColor()
}

// Prints error followed by program exit. Exit code 1 is reserved for crashes
func ErrorAndExit(err error) {
	PrintError(err)
	os.Exit(0)
}

// Checks if tokens is in tokenlist
func Contains(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}

	return false
}

func SContains(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}

	return false
}

// Returns index of last token
func SeekClosingBracket(tokens []lexer.Token, start int, beginT, endT int) (endIdx int, eof bool) {
	numParen := 0
	for i := start; i < len(tokens); i++ {
		switch tokens[i].Type {
		case beginT:
			numParen++
		case endT:
			numParen--
		}

		if numParen == 0 {
			return i, false
		}
	}

	return endIdx, true
}

// Returns Fizz name for value
func GetType(value interface{}) string {
	if i, ok := value.(env.FizzObject); ok {
		return i.Type()
	}

	switch value.(type) {
	case float64:
		return "number"
	case nil:
		return "nil"
	}

	return reflect.TypeOf(value).Name()
}

// Adds filename to error message if not already done
func WrapFilename(filename string, err error) error {
	if err == nil {
		return err
	}

	if !strings.Contains(err.Error(), ".fizz") {
		err = fmt.Errorf("%s: %s", filename, err.Error())
	}

	return err
}

type pair struct { a, b string }
type UniquePairs map[pair]bool

// Returns true if pair is already present in map
func (u UniquePairs) Add(a, b string) bool {
	if a > b {
		a, b = b, a
	}

	p := pair{a, b}
	n := u[p]
	u[p] = true
	return n
}

// Removes any path / file extension noise from filename
func GetPlainFilename(path string) string {
	path = strings.TrimSuffix(path, ".fizz")
	split := strings.Split(path, "/")
	return split[len(split)-1]
}