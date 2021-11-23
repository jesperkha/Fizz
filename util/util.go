package util

import (
	"fmt"
	"strings"
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