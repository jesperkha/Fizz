package str

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Standard string package for common string fuctionality

type i interface{}

var (
	Includes = map[string]interface{}{}

	ErrNotNumber = errors.New("string could not be converted to number, line %d")
)

func init() {
	Includes = map[string]interface{}{
		// Convert value to string
		"toString": toString,
		// Convert string to all lower case
		"lower": lower,
		// Convert string to all upper case
		"upper": upper,
		// Capitalize letters at start of each word in string
		"capital": capital,
		// Split string by other string
		"split": split,
		// Replace substring with new string
		"replace": replace,
		// Convert string number representation to string
		"toNumber": toNumber,
	}
}

func toString(val i) (str i, err error) {
	return fmt.Sprint(val), err
}

func lower(str string) (val i, err error) {
	return strings.ToLower(str), err
}

func upper(str string) (val i, err error) {
	return strings.ToUpper(str), err
}

func capital(str string) (val i, err error) {
	return strings.Title(str), err
}

func split(str string, split string) (val i, err error) {
	return strings.Split(str, split), err
}

func replace(str string, old string, new string) (val i, err error) {
	return strings.ReplaceAll(str, old, new), err
}

func toNumber(val string) (num i, err error) {
	num, err = strconv.ParseFloat(val, 64)
	if err != nil {
		return num, ErrNotNumber
	}

	return num, err
}
