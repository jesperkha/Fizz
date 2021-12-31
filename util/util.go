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

// Converts value to string in proper representation format
func FormatPrintValue(val interface{}) string {
	switch val.(type) {
	case float64, string, bool:
		return fmt.Sprint(val)
	case nil:
		return "nil"
	}

	if e, ok := val.(env.Environment); ok {
		glob := e[0]
		total := ""
		for k, v := range glob {
			total += fmt.Sprintf("%s: %s\n", k, FormatPrintValue(v))
		}

		return total
	}

	if o, ok := val.(*env.Object); ok {
		str := o.Name + ": {\n"
		for key, value := range o.Fields {
			str += fmt.Sprintf("    %s: %v\n", key, FormatPrintValue(value))
		}

		return str + "}"
	}

	if o, ok := val.(*env.Callable); ok {
		return o.Name + "()"
	}

	if a, ok := val.(*env.Array); ok {
		str := "["
		for i, v := range a.Values {
			if i != 0 {
				str += ", "
			}

			str += fmt.Sprintf("%v", FormatPrintValue(v))
		}

		return str + "]"
	}

	return ""
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
	os.Exit(1)
}

// Prints msg and exits with code 0
func PrintAndExit(msg interface{}) {
	fmt.Println(msg)
	os.Exit(1)
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

// Skips token check if in group or array expression. Returns index of ending token or eof.
func SeekBreakPoint(tokens []lexer.Token, verifier func(int, lexer.Token) bool) (targetIdx int, eof bool) {
	parens := 0
	targetIdx = -1
	for idx, token := range tokens {
		// Check before to allow for seeking parens
		if parens == 0 && verifier(idx, token) {
			targetIdx = idx
		}

		switch token.Type {
		case lexer.LEFT_PAREN, lexer.LEFT_SQUARE:
			parens++
		case lexer.RIGHT_PAREN, lexer.RIGHT_SQUARE:
			parens--
		}
	}

	return targetIdx, targetIdx == -1
}

// Splits list of token by split type
func SplitByToken(tokens []lexer.Token, split int) [][]lexer.Token {
	return SplitByTokens(tokens, []int{split})
}

// Splits list of token by multiple split types
func SplitByTokens(tokens []lexer.Token, splits []int) [][]lexer.Token {
	numParen := 0
	start := 0
	result := [][]lexer.Token{}
	for idx, token := range tokens {
		switch token.Type {
		case lexer.LEFT_PAREN, lexer.LEFT_SQUARE:
			numParen++
		case lexer.RIGHT_PAREN, lexer.RIGHT_SQUARE:
			numParen--
		}

		if Contains(splits, token.Type) && numParen == 0 {
			result = append(result, tokens[start:idx])
			start = idx + 1
		}
	}

	// Append last item
	if len(tokens) != 0 {
		result = append(result, tokens[start:])
	}

	return result
}

// Returns Fizz name for value
func GetType(value interface{}) string {
	if i, ok := value.(env.FizzObject); ok {
		return (i).Type()
	}

	switch value.(type) {
	case float64, float32, int:
		return "number"
	case nil:
		return "nil"
	}

	return reflect.TypeOf(value).Name()
}

// Gets fizz type name from reflect name
func GetLibType(typ string) string {
	switch typ {
	case "Object":
		return "object"
	case "Callable":
		return "function"
	case "Array":
		return "array"
	}

	return typ
}

// Adds filename to error message if not already done. Returns nil if err is nil.
func WrapFilename(filename string, err error) error {
	if err == nil || err.Error() == "" {
		return err
	}

	if !strings.Contains(err.Error(), ".fizz") {
		err = fmt.Errorf("%s: %s", filename, err.Error())
	}

	return err
}

type pair struct{ a, b string }
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

func IsInt(value interface{}) (int, bool) {
	if v, ok := value.(float64); ok {
		iv := int(v)
		return iv, v == float64(iv)
	}

	return -1, false
}
