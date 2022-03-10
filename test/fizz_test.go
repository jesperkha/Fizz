package test

import (
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/jesperkha/Fizz/interp"
)

const (
	validCaseFile   = true
	invalidCaseFile = false
)

var (
	testFiles = []string{
		"stmt",
		"expr",
	}
)

func TestAll(t *testing.T) {
	for i := float64(0); int(i) < len(testFiles); i += 0.5 {
		testType := float64(int(i)) == i
		invalid := ""
		if testType == invalidCaseFile {
			invalid = "in"
		}

		name := testFiles[int(math.Floor(i))]
		filename := fmt.Sprintf("./cases/%svalid_%s.fizz", invalid, name)
		byt, err := os.ReadFile(filename)
		if err != nil {
			t.Error(err)
		}

		cases := strings.Split(string(byt), "\n")
		for idx, c := range cases {
			c = strings.TrimSpace(c) // Removes invisible characters
			if c == "" || strings.HasPrefix(c, "#") {
				continue // Skip comments and whitespace for invalid case files
			}
			_, err := interp.Interperate("", c)

			// Valid case first
			if err != nil && testType == validCaseFile {
				// Case number is line number in case file
				t.Errorf("valid case %d got error: %s", idx+1, err)
			}

			// Invalid cases
			if err == nil && testType == invalidCaseFile {
				t.Errorf("invalid case %d got no error: %s", idx+1, c)
			}
		}
	}
}
