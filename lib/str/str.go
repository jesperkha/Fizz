package str

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/util"
)

// Standard string package for common string fuctionality

type i interface{}

var (
	ErrNotNumber = errors.New("string could not be converted to number, line %d")
	ErrNotString = errors.New("expected string value in array, line %d")
)

/*
	Converts value to string.
	func toString(value interface{}) string
*/
func ToString(val i) (str i, err error) {
	return fmt.Sprint(val), err
}

/*
	Formats value to default Fizz print formatting.
	func format(value interface{}) string
*/
func Format(val i) (str i, err error) {
	return util.FormatPrintValue(val), err
}

/*
	Converts all letters in string to lower case.
	func lower(str string) string
*/
func Lower(str string) (val i, err error) {
	return strings.ToLower(str), err
}

/*
	Converts all letters in string to upper case.
	func upper(str string) string
*/
func Upper(str string) (val i, err error) {
	return strings.ToUpper(str), err
}

/*
	Capitalizes the letter at the beginning of each word.
	func capital(str string) string
*/
func Capital(str string) (val i, err error) {
	return strings.Title(str), err
}

/*
	Splits string by substring
	func split(str string, split string) []string
*/
func Split(str string, split string) (val i, err error) {
	splits := []interface{}{}
	for _, s := range strings.Split(str, split) {
		splits = append(splits, s)
	}

	return env.NewArray(splits), err
}

/*
	Joins array of strings into one string with the substring.
	func join(strings []string, sub string) string
*/
func Join(str *env.Array, sub string) (val i, err error) {
	s := []string{}
	for _, i := range str.Values {
		if newStr, ok := i.(string); ok {
			s = append(s, newStr)
			continue
		}

		return val, ErrNotString
	}

	return strings.Join(s, sub), err
}

/*
	Replaces all instances of substring with new string.
	func replace(str string, old string, new string) string
*/
func Replace(str string, old string, new string) (val i, err error) {
	return strings.ReplaceAll(str, old, new), err
}

/*
	Converts string to number.
	func toNumber(str string) float64
*/
func ToNumber(val string) (num i, err error) {
	num, err = strconv.ParseFloat(val, 64)
	if err != nil {
		return num, ErrNotNumber
	}

	return num, err
}
