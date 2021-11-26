package util

import (
	"fmt"
	"reflect"
	"strings"

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
	if value == nil {
		return "nil"
	}

	if i, ok := value.(env.FizzObject); ok {
		return i.Type()
	}

	return reflect.TypeOf(value).Name()
}
