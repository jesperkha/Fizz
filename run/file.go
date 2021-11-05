package run

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Interperates code found in file specified in commandline arguments
func RunFile(filename string) (err error) {
	if !strings.HasSuffix(filename, ".fizz") {
		return ErrNonFizzFile
	}

	if file, err := os.Open(filename); err == nil {
		var buf bytes.Buffer
		bufio.NewReader(file).WriteTo(&buf)

		return Interperate(buf.String())
	}

	// Assumes path error
	return fmt.Errorf(ErrFileNotFound.Error(), filename)
}
